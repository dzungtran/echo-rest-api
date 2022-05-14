package cmd

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"

	"github.com/dzungtran/echo-rest-api/pkg/utils"
	"github.com/dzungtran/echo-rest-api/tools/mod/cmd/modgen"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type (
	modGenOpts struct {
		Name        string
		Dir         string
		SkipHandler bool
		SkipUseCase bool
	}
	modgenTplVars struct {
		ModuleName string
	}
)

const (
	defaultFlagCodeGen = "// Auto generate"
)

var (
	mgOpts       = &modGenOpts{}
	digFilePaths = map[string]string{
		"usecases/dig.go":              `_ = container.Provide(New{{ .ModuleName }}Usecase)`,
		"repositories/postgres/dig.go": `_ = container.Provide(NewPgsql{{ .ModuleName }}Repository)`,
		"cmd/api/di/params.go":         `{{ .ModuleName }}Usecase usecases.{{ .ModuleName }}Usecase`,
		"cmd/api/di/di.go":             `httpDelivery.New{{ .ModuleName }}Handler(adminGroup, params.MiddlewareManager, params.{{ .ModuleName }}Usecase)`,
	}
	funcsMap = template.FuncMap{
		"ToLower":      strings.ToLower,
		"ToSnake":      utils.ToSnake,
		"ToKebab":      utils.ToKebab,
		"ToCamel":      utils.ToCamel,
		"ToLowerCamel": utils.ToLowerCamel,
	}
)

// modgenCmd represents the modgen command
var modgenCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate module codes from template files",
	Long:  `Generate module codes from template files.`,
	Args: func(cmd *cobra.Command, args []string) error {
		// validate module name
		nameRegex := `^[a-zA-Z0-9_]+$`
		if ok, _ := regexp.MatchString(nameRegex, mgOpts.Name); !ok {
			return errors.New("module name must match ( " + nameRegex + " )")
		}

		mgOpts.SkipHandler = cmd.Flag("skip-handler").Changed
		mgOpts.SkipUseCase = cmd.Flag("skip-usecase").Changed

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		tpls := modgen.GetModGenTemplates()
		tVars := modgenTplVars{ModuleName: mgOpts.Name}
		tplDir := "templates"

		tplFiles, err := tpls.ReadDir(tplDir)
		if err != nil {
			return err
		}

		// parse file templates
		parsedTpls := make([]string, 0)
		for _, f := range tplFiles {
			if f.IsDir() {
				continue
			}

			fname := tplDir + "/" + f.Name()
			fcontent, err := tpls.ReadFile(fname)
			if err != nil {
				return err
			}

			pt, err := template.New(fname).
				Funcs(funcsMap).
				Parse(string(fcontent))
			if err != nil {
				return err
			}

			builder := &bytes.Buffer{}
			err = pt.Execute(builder, tVars)
			if err != nil {
				return err
			}
			parsedTpls = append(parsedTpls, builder.String())
		}

		re := regexp.MustCompile(`(?m)^\/\/\sTarget:(.*)$`)
		parsedFiles := make(map[string]string)
		for _, tpl := range parsedTpls {
			sm := re.FindStringSubmatch(tpl)
			if len(sm) != 2 {
				continue
			}
			fn := strings.Trim(sm[1], " ")

			if mgOpts.SkipHandler && (strings.HasPrefix(fn, "delivery/http")) {
				continue
			}

			if mgOpts.SkipUseCase && (strings.HasPrefix(fn, "delivery/http") || strings.HasPrefix(fn, "delivery/requests") || strings.HasPrefix(fn, "usecases")) {
				continue
			}

			parsedFiles[fn] = tpl
		}

		// write generated file
		curDir, _ := os.Getwd()
		if mgOpts.Dir != "" {
			curDir = strings.Trim(mgOpts.Dir, " ")
		}

		logrus.Debugf("Working Dir: %v", curDir)

		for fname, fcon := range parsedFiles {
			writeGeneratedFile(curDir, fname, fcon)
		}

		// append generated files to dig files for DI
		for fp, str := range digFilePaths {

			if mgOpts.SkipUseCase && strings.HasPrefix(fp, "usecases") {
				continue
			}
			appendGeneratedFileForDependencyInjection(fp, str, tVars)
		}

		runGenerateRoutes()
		logrus.Info("Done!")
		return nil
	},
}

func writeGeneratedFile(curDir, fname, fcon string) error {
	f, err := os.OpenFile(curDir+"/"+fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		logrus.Error(err)
		return err
	}
	_, err = f.Write([]byte(fcon))
	if err != nil {
		logrus.Error(err)
		return err
	}

	if err := f.Close(); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func appendGeneratedFileForDependencyInjection(fp, str string, tVars modgenTplVars) {
	f, err := os.OpenFile(fp, os.O_RDWR, 0644)
	if err != nil {
		logrus.Errorf("failed opening file: %s", err)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadFile(fp)
	if err != nil {
		logrus.Errorf("failed reading data from file: %s", err)
		return
	}

	// parse template
	pt, err := template.New(fp).
		Funcs(funcsMap).
		Parse(str)
	if err != nil {
		logrus.Errorf("failed parse template: %s", err)
		return
	}

	builder := &bytes.Buffer{}
	err = pt.Execute(builder, tVars)
	if err != nil {
		logrus.Errorf("failed parse template: %s", err)
		return
	}

	// Find string existed then ignore append
	if strings.Contains(string(data), builder.String()) {
		return
	}

	parsedStr := builder.String() + "\n\t" + defaultFlagCodeGen
	fcontent := strings.Replace(string(data), defaultFlagCodeGen, parsedStr, 1)
	_, err = f.WriteAt([]byte(fcontent), 0) // Write at 0 beginning
	if err != nil {
		logrus.Errorf("failed writing to file: %s", err)
	}
}

func runGenerateRoutes() {
	cmd := exec.Command("make", "routes")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Run: ", string(out.Bytes()))
}

func init() {
	rootCmd.AddCommand(modgenCmd)
	modgenCmd.Flags().StringVarP(&mgOpts.Name, "name", "n", "", "module name to generate")
	modgenCmd.Flags().BoolVarP(&mgOpts.SkipHandler, "skip-handler", "", false, "skip generate handler")
	modgenCmd.Flags().BoolVarP(&mgOpts.SkipUseCase, "skip-usecase", "", false, "skip generate usecase")
	modgenCmd.Flags().StringVarP(&mgOpts.Dir, "dir", "d", "", "project's root path")
	modgenCmd.MarkFlagRequired("name")
}
