package authz

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"strings"

	"github.com/dzungtran/echo-rest-api/modules/core/domains"
	"github.com/dzungtran/echo-rest-api/pkg/contexts"
	"github.com/dzungtran/echo-rest-api/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage/inmem"
	"github.com/open-policy-agent/opa/util"
	"github.com/tidwall/sjson"
)

var (
	//go:embed data.json
	dataFile []byte
	//go:embed routes.json
	routesFile []byte

	//go:embed rego/*.rego
	regoFs embed.FS
	//go:embed rego/deny/*.rego
	regoDenyFs embed.FS
	//go:embed rego/utils/*.rego
	regoUtilsFs embed.FS

	regoInstance  *rego.Rego
	preparedQuery rego.PreparedEvalQuery
)

type opaInputOpts struct {
	ResourcePermissions map[string][]string
	ExtraData           map[string]interface{}
	RequestEndpoint     string
	RequestMethod       string

	Org *domains.Org
}

type CallOPAInputOption struct {
	applyFunc func(*opaInputOpts)
}

func init() {

	regoContents := map[string]string{}
	readRegoFiles(regoFs, "rego", regoContents)
	readRegoFiles(regoDenyFs, "rego/deny", regoContents)
	readRegoFiles(regoUtilsFs, "rego/utils", regoContents)

	// Compile the module. The keys are used as identifiers in error messages.
	compiler, err := ast.CompileModules(regoContents)
	if err != nil {
		logger.Log().Fatalf("error while init rego compiler, details: %s", err.Error())
	}

	var jsonData map[string]interface{}
	var routesData map[string]interface{}
	err = util.UnmarshalJSON(routesFile, &routesData)
	if err != nil {
		logger.Log().Fatalf("error while init rego instance, details: %s", err.Error())
	}

	dataFile, err = sjson.SetBytes(dataFile, "endpoints_acl", routesData)
	if err != nil {
		logger.Log().Fatalf("error while build data file, details: %s", err.Error())
	}

	err = util.UnmarshalJSON(dataFile, &jsonData)
	if err != nil {
		logger.Log().Fatalf("error while init rego instance, details: %s", err.Error())
	}

	// Manually create the storage layer. inmem.NewFromObject returns an
	// in-memory store containing the supplied data.
	store := inmem.NewFromObject(jsonData)

	/*
		allow = data.authz.allow
		deny = data.authz.deny
		req = data.authz.usr_role
	*/

	// Create new query that returns the value
	regoInstance = rego.New(
		rego.Query(`
			allow = data.authz.allow
			deny = data.authz.deny
		`),
		rego.Store(store),
		rego.Compiler(compiler),
	)

	ctx := context.Background()
	preparedQuery, err = regoInstance.PrepareForEval(ctx)
	if err != nil {
		logger.Log().Fatalf("error while prepared eval opa, details: %v", err.Error())
	}
}

func readRegoFiles(dfs embed.FS, folderName string, filesContent map[string]string) {
	dirs, err := dfs.ReadDir(folderName)
	if err != nil {
		logger.Log().Fatalf("error while init rego compiler, details: %s", err.Error())
	}

	for _, d := range dirs {
		if d.IsDir() {
			continue
		}

		f, err := d.Info()
		if err != nil {
			continue
		}

		fName := folderName + "/" + f.Name()
		fcb, err := dfs.ReadFile(fName)
		if err != nil {
			continue
		}
		filesContent[fName] = string(fcb)
	}
}

func WithInputResourcePermissions(resPerms map[string][]string) CallOPAInputOption {
	return CallOPAInputOption{
		applyFunc: func(oio *opaInputOpts) {
			if oio.ResourcePermissions == nil {
				oio.ResourcePermissions = make(map[string][]string)
			}

			if resPerms != nil {
				currResPerms := oio.ResourcePermissions
				for perm, resIds := range resPerms {
					currResIds := currResPerms[perm]
					if currResIds == nil {
						currResIds = make([]string, 0)
					}

					if len(resIds) > 0 {
						currResIds = append(currResIds, resIds...)
					}
					currResPerms[perm] = currResIds
				}
				oio.ResourcePermissions = currResPerms
			}
		},
	}
}

func WithInputRequestMethod(reqMethod string) CallOPAInputOption {
	return CallOPAInputOption{
		applyFunc: func(oio *opaInputOpts) {
			reqMethod = strings.Trim(reqMethod, " ")
			if reqMethod != "" {
				oio.RequestMethod = reqMethod
			}
		},
	}
}

func WithInputRequestEndpoint(reqEndpoint string) CallOPAInputOption {
	return CallOPAInputOption{
		applyFunc: func(oio *opaInputOpts) {
			reqEndpoint = strings.Trim(reqEndpoint, " ")
			if reqEndpoint != "" {
				oio.RequestEndpoint = reqEndpoint
			}
		},
	}
}

func WithInputExtraData(key string, data interface{}) CallOPAInputOption {
	return CallOPAInputOption{
		applyFunc: func(oio *opaInputOpts) {
			if oio.ExtraData == nil {
				oio.ExtraData = make(map[string]interface{})
			}
			oio.ExtraData[key] = data
		},
	}
}

func WithInputOrg(org *domains.Org) CallOPAInputOption {
	return CallOPAInputOption{
		applyFunc: func(oio *opaInputOpts) {
			oio.Org = org
		},
	}
}

func CheckPolicies(user *domains.UserWithRoles, callOpts ...CallOPAInputOption) (denyMsg []string, err error) {
	ctx := context.Background()
	opts := appliedOPAInputOption(callOpts)
	input := map[string]interface{}{
		"user": user,
	}

	if opts.Org != nil {
		input["org"] = opts.Org
	}

	if len(opts.ExtraData) > 0 {
		for k, v := range opts.ExtraData {
			input[k] = v
		}
	}

	// Apply options to input
	if opts.RequestMethod != "" && opts.RequestEndpoint != "" {
		input["method"] = opts.RequestMethod
		input["endpoint"] = opts.RequestEndpoint
	}

	// Run evaluation.
	rs, err := preparedQuery.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		logger.Log().Errorf("error while eval opa input, details %v", err.Error())
		return
	}

	if len(rs) == 0 {
		err = errors.New("empty decision result")
		return
	}

	var result bool
	result, ok := rs[0].Bindings["allow"].(bool)
	if msg, _ := rs[0].Bindings["deny"].([]interface{}); len(msg) > 0 {
		denyMsg = make([]string, 0)
		for _, v := range msg {
			denyMsg = append(denyMsg, fmt.Sprint(v))
		}
	}

	if !ok {
		logger.Log().Errorf("unexpected decision result, details %v", result)
		err = errors.New("unexpected decision result")
		return
	}

	if !result {
		err = errors.New("forbidden")
		return
	}

	return
}

func CheckPoliciesContext(c echo.Context, callOpts ...CallOPAInputOption) (denyMsg []string, err error) {
	u, err := contexts.GetUserFromContext(c)
	if err != nil {
		return
	}

	callOpts = append(callOpts,
		WithInputRequestMethod(c.Request().Method),
		WithInputRequestEndpoint(c.Path()),
	)

	return CheckPolicies(u, callOpts...)
}

func appliedOPAInputOption(callOptions []CallOPAInputOption) *opaInputOpts {
	if len(callOptions) == 0 {
		return &opaInputOpts{}
	}

	optCopy := &opaInputOpts{}
	for _, f := range callOptions {
		f.applyFunc(optCopy)
	}
	return optCopy
}
