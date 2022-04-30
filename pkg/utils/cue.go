package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/format"
	cueJsonEncoding "cuelang.org/go/encoding/json"
)

// CueEvaluatingJsonValue - Simple func to build json from a CueLang definition
// Parameters:
// - schemaName: Main cue schema name to eval
// - cueDef: Cue document, contains Cue definitions
// - jsonValue: Use to as Cue values for evaluating
func CueEvaluatingJsonValue(schemaName string, cueDef string, jsonValue []byte) ([]byte, error) {
	out := &bytes.Buffer{}

	// convert
	d := cueJsonEncoding.NewDecoder(nil, "", bytes.NewReader(jsonValue))
	for {
		e, err := d.Extract()
		if err == io.EOF {
			break
		}
		appendAstExprToBuffer(out, e, err)
		if err != nil {
			break
		}
	}

	if out.Len() == 0 {
		return nil, errors.New("cannot parse json to cue")
	}

	str := fmt.Sprintf("#%s & ", schemaName) + out.String()

	// create a context
	ctx := cuecontext.New()

	// TODO: needs some cache
	// compile our schema first
	s := ctx.CompileString(cueDef, cue.Filename("schema.cue"))
	if err := s.Err(); err != nil {
		return nil, err
	}

	// compile our value
	v := ctx.CompileString(str, cue.Scope(s), cue.Filename("values.cue"))
	if err := v.Err(); err != nil {
		return nil, err
	}

	rs, err := v.MarshalJSON()

	return rs, err
}

// CueValidateJson - Simple func to validate json by cue definitions
// Parameters:
// - schemaName: Main cue schema name to eval
// - cueDef: Cue document, contains Cue definitions
// - jsonBytes: Use to as Cue values for evaluating
func CueValidateJson(schemaName string, cueDef string, jsonBytes []byte) error {
	ctx := cuecontext.New()

	// TODO: needs some cache
	s := ctx.CompileString(cueDef, cue.Filename("schema.cue"))
	if err := s.Err(); err != nil {
		return err
	}

	v := s.LookupPath(cue.ParsePath("#" + schemaName))
	if err := v.Err(); err != nil {
		return err
	}

	data := ctx.CompileBytes(jsonBytes, cue.Filename("content.json"))
	if err := data.Err(); err != nil {
		return err
	}

	unified := data.Unify(v)
	if err := unified.Err(); err != nil {
		return err
	}

	opts := []cue.Option{
		cue.Attributes(true),
		cue.Definitions(true),
		cue.Hidden(true),
	}

	err := unified.Validate(opts...)
	return err
}

// Becareful with json ommitempty options.
// Prefer use CueValidateJson for raw payload before marshal to struct
func CueValidateObject(schemaName string, cueDef string, obj interface{}) error {
	jsonB, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	return CueValidateJson(schemaName, cueDef, jsonB)
}

func appendAstExprToBuffer(w *bytes.Buffer, e ast.Expr, err error) {
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	b, err := format.Node(e)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprint(w, string(b))
}
