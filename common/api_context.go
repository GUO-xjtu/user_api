package common

import (
	"context"
	"github.com/ettle/strcase"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"sync"
)

type RequestContext struct {
	context.Context
	GinContext           *gin.Context
	Request              *http.Request
	isFrequenceOverLimit bool
	success              bool
	stable           *bool
	params           interface{}
	rawParams        map[string]string
	antispamFeedback sync.Map
}


func NewRequestContext(c *gin.Context) *RequestContext {
	req := &RequestContext{success: false}

	req.GinContext = c
	req.Request = c.Request
	req.Context = context.Background()
	req.initRawParams()
	return req
}

func (r *RequestContext) GetCtx() context.Context {
	return r.Context
}

func (r *RequestContext) GetGinCtx() *gin.Context {
	return r.GinContext
}

func (r *RequestContext) Params() interface{} {
	return r.params
}

func (r *RequestContext)InitRequestParams(pType reflect.Type) error {
	if r == nil {
		return errors.New("RequestParams Error")
	}

	if pType == nil {
		return nil
	}

	if r.params != nil {
		return nil
	}

	var err error
	r.params, err = parseParams(r, r.rawParams, pType)
	return err
}

func parseParams(ctx context.Context, rawParams map[string]string, pType reflect.Type) (interface{}, error) {
	if pType.Kind() != reflect.Struct {
		log.Fatalf("param type %+v is not struct", pType)
		return nil, errors.Errorf("param type %+v is not struct", pType)
	}

	params := reflect.New(pType)

	for i := 0; i < pType.NumField(); i++ {
		fieldType := pType.Field(i)
		field := params.Elem().FieldByName(fieldType.Name)
		snakeName := strcase.ToSnake(fieldType.Name)

		val, ok := rawParams[snakeName]
		if !ok {
			continue
		}

		var parseRes interface{}
		var parseErr error

		switch field.Kind() {
		case reflect.Struct:
			// TODO
			break
		case reflect.Int:
			parseRes, parseErr = strconv.ParseInt(val, 10, 0)
			break
		case reflect.Int8:
			parseRes, parseErr = strconv.ParseInt(val, 10, 8)
			break
		case reflect.Int16:
			parseRes, parseErr = strconv.ParseInt(val, 10, 16)
			break
		case reflect.Int32:
			parseRes, parseErr = strconv.ParseInt(val, 10, 32)
			break
		case reflect.Int64:
			parseRes, parseErr = strconv.ParseInt(val, 10, 64)
			break
		case reflect.Uint:
			parseRes, parseErr = strconv.ParseUint(val, 10, 0)
			break
		case reflect.Uint8:
			parseRes, parseErr = strconv.ParseUint(val, 10, 8)
			break
		case reflect.Uint16:
			parseRes, parseErr = strconv.ParseUint(val, 10, 16)
			break
		case reflect.Uint32:
			parseRes, parseErr = strconv.ParseUint(val, 10, 32)
			break
		case reflect.Uint64:
			parseRes, parseErr = strconv.ParseUint(val, 10, 64)
			break
		case reflect.Float32:
			parseRes, parseErr = strconv.ParseFloat(val, 32)
			break
		case reflect.Float64:
			parseRes, parseErr = strconv.ParseFloat(val, 64)
			break
		case reflect.Bool:
			parseRes, parseErr = strconv.ParseBool(val)
			break
		case reflect.String:
			parseRes = val
			break
		default:
			continue
		}

		if parseErr == nil && parseRes != nil {
			switch reflect.TypeOf(parseRes).Kind() {
			case reflect.Int64:
				field.SetInt(parseRes.(int64))
				break
			case reflect.Float64:
				field.SetFloat(parseRes.(float64))
				break
			case reflect.Bool:
				field.SetBool(parseRes.(bool))
				break
			case reflect.Uint64:
				field.SetUint(parseRes.(uint64))
				break
			case reflect.String:
				field.SetString(parseRes.(string))
				break
			default:
				continue
			}
		} else {
			return nil, errors.Errorf("unmatched parameter type, param=%s err: %v", snakeName, parseErr)
		}
	}

	return params.Interface(), nil
}

func (r *RequestContext) initRawParams() {
	r.rawParams = map[string]string{}
	for k, v := range r.GinContext.Request.URL.Query() {
		if len(v) > 0 {
			r.rawParams[k] = v[0]
		}
	}
	if r.GinContext.Request.ParseForm() == nil {
		for k, v := range r.GinContext.Request.PostForm {
			if len(v) > 0 {
				r.rawParams[k] = v[0]
			}
		}
	}
	for _, param := range r.GinContext.Params {
		r.rawParams[param.Key] = param.Value
	}
}