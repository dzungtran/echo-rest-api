package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/dzungtran/echo-rest-api/pkg/logger"
	"github.com/dzungtran/echo-rest-api/pkg/utils"
	"github.com/gertd/go-pluralize"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type (
	modGenOpts struct {
		Name string
		Dir  string
	}
	modgenTplVars struct {
		ModuleName   string
		SingularName string
		PluralName   string
		RootPackage  string
	}
	modConfigs struct {
		RootPackage string `yaml:"RootPackage"`
	}
)

const (
	configFile = ".modgen.yaml"
)

var (
	mgOpts   = &modGenOpts{}
	cnf      = &modConfigs{}
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
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		parsedFiles := make(map[string]string)
		folderStruct := make(map[string]bool)

		pluralize := pluralize.NewClient()
		yamlFile, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(yamlFile, cnf)
		if err != nil {
			return err
		}

		tVars := modgenTplVars{
			ModuleName:   pluralize.Plural(mgOpts.Name),
			SingularName: pluralize.Singular(mgOpts.Name),
			PluralName:   pluralize.Plural(mgOpts.Name),
			RootPackage:  cnf.RootPackage,
		}

		tplDir := "tools/mod/cmd/modgen/template"
		err = filepath.Walk(tplDir,
			func(path string, finfo os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if finfo.IsDir() {
					return nil
				}

				fname := strings.ReplaceAll(path, tplDir, "modules/"+utils.ToKebab(tVars.ModuleName))
				folderStruct[strings.ReplaceAll(fname, finfo.Name(), "")] = true

				fcontent, err := os.ReadFile(path)
				if err != nil {
					logger.Log().Errorf("Error while read template, details: %v", err)
					return err
				}

				pt, err := template.New(fname).
					Funcs(funcsMap).
					Parse(string(fcontent))
				if err != nil {
					logger.Log().Errorf("Error while parse template, details: %v", err)
					return err
				}

				builder := &bytes.Buffer{}
				err = pt.Execute(builder, tVars)
				if err != nil {
					logger.Log().Errorf("Error while parse template, details: %v", err)
					return err
				}

				parsedFiles[fname] = builder.String()
				return nil
			})
		if err != nil {
			logger.Log().Errorf("Error while parse template, details: %v", err)
			return err
		}

		if len(folderStruct) > 0 {
			for fd := range folderStruct {
				os.MkdirAll(fd, os.ModePerm)
			}
		}

		for fn, c := range parsedFiles {
			fn = strings.ReplaceAll(fn, ".gotpl", ".go")
			fn = strings.ReplaceAll(fn, "placeholder", utils.ToSnake(tVars.SingularName))
			err = ioutil.WriteFile(fn, []byte(c), 0644)
			if err != nil {
				fmt.Println("failed writing to file: ", fn, "\n", err.Error())
			}
		}

		runGenerateRoutes()
		logger.Log().Info("Done!")
		return nil
	},
}

func runGenerateRoutes() {
	cmd := exec.Command("make", "routes")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		logger.Log().Fatal(err)
	}
	logger.Log().Info("Run: ", out.String())
}

func init() {
	rootCmd.AddCommand(modgenCmd)
	modgenCmd.Flags().StringVarP(&mgOpts.Name, "name", "n", "", "module name to generate")
	modgenCmd.Flags().StringVarP(&mgOpts.Dir, "dir", "d", "", "project's root path")
	modgenCmd.MarkFlagRequired("name")
}
