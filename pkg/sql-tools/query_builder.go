package sqlTools

import (
	"context"
	"database/sql"
	"reflect"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/dzungtran/echo-rest-api/pkg/constants"
	"github.com/dzungtran/echo-rest-api/pkg/contexts"
	"github.com/dzungtran/echo-rest-api/pkg/utils"
)

type options struct {
	selectFields       []string
	ignoreFields       []string
	autoDateTimeFields []string
}

type CallOptionMapValues struct {
	applyFunc func(*options)
}

func WithMapValuesSelectFields(fields []string) CallOptionMapValues {
	return CallOptionMapValues{
		applyFunc: func(opt *options) {
			if len(opt.selectFields) == 0 {
				opt.selectFields = make([]string, 0)
			}
			if len(fields) > 0 {
				opt.selectFields = append(opt.selectFields, fields...)
			}
		},
	}
}

func WithMapValuesIgnoreFields(fields []string) CallOptionMapValues {
	return CallOptionMapValues{
		applyFunc: func(opt *options) {
			if len(opt.ignoreFields) == 0 {
				opt.ignoreFields = make([]string, 0)
			}
			if len(fields) > 0 {
				opt.ignoreFields = append(opt.ignoreFields, fields...)
			}
		},
	}
}

func WithMapValuesAutoDateTimeFields(fields []string) CallOptionMapValues {
	return CallOptionMapValues{
		applyFunc: func(opt *options) {
			if len(opt.autoDateTimeFields) == 0 {
				opt.autoDateTimeFields = make([]string, 0)
			}
			if len(fields) > 0 {
				opt.autoDateTimeFields = append(opt.autoDateTimeFields, fields...)
			}
		},
	}
}

// GetMapValuesFromStruct For INSERT or UPDATE queries
func GetMapValuesFromStruct(ctx context.Context, st interface{}, callOpts ...CallOptionMapValues) map[string]interface{} {
	rs := make(map[string]interface{})
	t := reflect.TypeOf(st).Elem()
	v := reflect.ValueOf(st)

	opt := appliedFieldSelectOption(callOpts)
	for index := 0; index < t.NumField(); index++ {
		f := t.Field(index)

		// temp fix for embedded struct
		//if f.Type.Kind() == reflect.Struct {
		//	cols, vals := GetColumnsAndValuesFromStruct(ctx, reflect.Indirect(v).FieldByName(f.Name).Interface(), callOpts...)
		//	if len(cols) == 0 {
		//		continue
		//	}
		//	for i := 0; i < len(cols); i++ {
		//		rs[cols[i]] = vals[i]
		//	}
		//	continue
		//}

		dbKey := getKeyFromTag(f.Tag.Get("db"))
		if dbKey == "-" || len(dbKey) == 0 {
			continue
		}

		if len(opt.ignoreFields) > 0 && utils.IsSliceContains(opt.ignoreFields, dbKey) {
			continue
		}

		if len(opt.autoDateTimeFields) > 0 && utils.IsSliceContains(opt.autoDateTimeFields, dbKey) {
			rs[dbKey] = time.Now().UTC()
			continue
		}

		if len(opt.selectFields) > 0 && !utils.IsSliceContains(opt.selectFields, dbKey) {
			continue
		}

		ival := reflect.Indirect(v).FieldByName(f.Name).Interface()
		rs[dbKey] = ival
	}
	return rs
}

func GetColumnsAndValuesFromStruct(ctx context.Context, st interface{}, callOpts ...CallOptionMapValues) ([]string, []interface{}) {
	t := reflect.TypeOf(st).Elem()
	v := reflect.ValueOf(st)

	columns := make([]string, 0)
	values := make([]interface{}, 0)

	opt := appliedFieldSelectOption(callOpts)
	for index := 0; index < t.NumField(); index++ {
		f := t.Field(index)
		dbKey := getKeyFromTag(f.Tag.Get("db"))

		if len(opt.autoDateTimeFields) > 0 && utils.IsSliceContains(opt.autoDateTimeFields, dbKey) {
			columns = append(columns, dbKey)
			values = append(values, time.Now().UTC())
			continue
		}

		// temp fix for embedded struct
		if f.Type.Kind() == reflect.Struct {

			str := reflect.Indirect(v).FieldByName(f.Name)
			pactual := reflect.New(str.Type()).Interface()

			// ignore scan some native struct
			nestedReflect := true
			switch pactual.(type) {
			case *time.Time, *sql.NullTime:
				nestedReflect = false
			}

			if nestedReflect {
				cols, vals := GetColumnsAndValuesFromStruct(ctx, pactual, callOpts...)
				if len(cols) == 0 {
					continue
				}
				columns = append(columns, cols...)
				values = append(values, vals...)
			} else {
				columns = append(columns, dbKey)
				values = append(values, str.Interface())
			}
			continue
		}

		if dbKey == "-" || len(dbKey) == 0 || (len(opt.selectFields) > 0 && !utils.IsSliceContains(opt.selectFields, dbKey)) {
			continue
		}

		if len(opt.ignoreFields) > 0 && utils.IsSliceContains(opt.ignoreFields, dbKey) {
			continue
		}

		columns = append(columns, dbKey)
		values = append(values, reflect.Indirect(v).FieldByName(f.Name).Interface())
	}

	return columns, values
}

func ParseColumnsForSelect(cols []string) []string {
	for i, v := range cols {
		switch v {
		case "_count":
			cols[i] = "count(*) over() as _count"
		}
	}
	return cols
}

func ParseColumnsForSelectWithAlias(cols []string, alias string) []string {
	for i, v := range cols {
		if utils.IsSliceContains([]string{"_count"}, v) {
			switch v {
			case "_count":
				cols[i] = "count(*) over() as _count"
			}
			continue
		}

		cols[i] = alias + "." + v
	}
	return cols
}

func appliedFieldSelectOption(callOptions []CallOptionMapValues) *options {
	if len(callOptions) == 0 {
		return &options{}
	}

	optCopy := &options{}
	for _, f := range callOptions {
		f.applyFunc(optCopy)
	}
	return optCopy
}

// getKeyFromTag -- get db key from struct db tag
func getKeyFromTag(tag string) string {
	pieces := utils.StringSlice(tag, ",")
	if len(pieces) == 0 {
		return ""
	}
	return pieces[0]
}

func BindCommonParamsToSelectBuilder(builder squirrel.SelectBuilder, params contexts.CommonParamsForFetch) squirrel.SelectBuilder {
	limit := constants.DefaultPerPage
	if params.Limit > 0 {
		limit = params.Limit
	}

	if limit > constants.MaximumPerPage {
		limit = constants.DefaultPerPage
	}

	offset := uint64(0)
	if params.Page > 1 {
		offset = limit * (params.Page - 1)
	}

	if !params.NoLimit {
		builder = builder.Limit(limit).Offset(offset)
	}

	return builder
}
