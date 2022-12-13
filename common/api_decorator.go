package common

import (
	"encoding/json"
	"fmt"
	"github.com/GUO-xjtu/pb_gen/gen/data"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"net"
	"net/http"
	"os"
	"reflect"
	"runtime/debug"
	"strings"
	"syscall"
	"time"
)

type RequestHandler func(context *RequestContext) (proto.Message, error)

const (
	optionsContextKey = "xiangqin_options_context_key"
	requestContextKey = "request_context"


)


func API(handler RequestHandler, paramsMsg proto.Message, module, view string, options ...APIOption) func(*gin.Context) {
	w := baseAPI(handler, paramsMsg, module, view, options...)
	return w
}

func baseAPI(handler RequestHandler, paramsMsg proto.Message, module, view string, options ...APIOption) func(*gin.Context) {
	h := func(c *gin.Context) {
		v1, exists := c.Get(optionsContextKey)
		if !exists {
			fmt.Sprintf("is not")
		}
		fmt.Sprintf("v1:%+v", v1)
		Opts := v1.(*APIOptions)

		v2, _ := c.Get(requestContextKey)
		ctx := v2.(*RequestContext)

		if paramsMsg != nil {
			err := ctx.InitRequestParams(reflect.TypeOf(paramsMsg).Elem())
			if err != nil {
				logger.CtxWarnf(ctx, "request params error: %s", err.Error())
				genErrResp(ctx, err)
				return
			}
		}

		respMsg, bizErr := handler(ctx)
		if bizErr != nil {
			genErrResp(ctx, bizErr)
			return
		}
		if respMsg == nil {
			respMsg = &empty.Empty{}
		}

		// 1. JSON 格式数据处理
		if !IsUsePBResp(ctx) {
			b, err := MarshalJsonPb(respMsg)
			if err != nil {
				logger.CtxErrorf(ctx.GetCtx(), "Marshal schema result to json failed, err: %s stack: %s", err, debug.Stack())
				genErrResp(ctx, err)
				return
			}

			var (
				respData  interface{}
				respExtra map[string]interface{}
			)
			if Opts.PureJSONResp {
				r := struct {
					Data  interface{}            `json:"data"`
					Extra map[string]interface{} `json:"extra"`
				}{}
				err = json.Unmarshal(b, &r)
				if err != nil {
					logger.CtxErrorf(ctx.GetCtx(), "Unmarshal schema json result failed, err: %s stack: %s", err, debug.Stack())
					genErrResp(ctx, err)
					return
				}

				respData, respExtra = r.Data, r.Extra
			} else {
				r := struct {
					Data  json.RawMessage     `json:"data"`
					Extra map[string]interface{} `json:"extra"`
				}{}
				err = json.Unmarshal(b, &r)
				if err != nil {
					logger.CtxErrorf(ctx.GetCtx(), "Unmarshal schema json result failed, err: %s stack: %s", err, debug.Stack())
					genErrResp(ctx, err)
					return
				}

				if len(r.Data) == 0 {
					r.Data = json.RawMessage("{}")
				}

				respData, respExtra = r.Data, r.Extra
			}

			if respExtra == nil {
				respExtra = map[string]interface{}{}
			}

			logger.CtxNoticef(ctx.GetCtx(), "err_code", 0)

			respExtra["now"] = GetMSTimeStamp(time.Now())
			c.Header("net_3", fmt.Sprintf("%d", GetMSTimeStamp(time.Now())))
			if Opts.PureJSONResp {
				c.PureJSON(200, gin.H{
					"status_code": 0,
					"data":        respData,
					"extra":       respExtra,
				})
			} else {
				c.JSON(200, gin.H{
					"status_code": 0,
					"data":        respData,
					"extra":       respExtra,
				})
			}
		} else {
			pbResp := data.PBResponse{}
			pbResp.Header = &data.PBResponse_Header{
				StatusCode: 0,
				Now:        GetMSTimeStamp(time.Now()),
			}

			if v, ok := respMsg.(*wrappers.BytesValue); ok {
				pbResp.Body = v.Value
			} else {
				pbResult, err := proto.Marshal(respMsg)
				if err != nil {
					logger.CtxErrorf(ctx.GetCtx(), "pb marshal failed error:%v", err.Error())
					pbResp.Body = []byte{}
				} else {
					pbResp.Body = pbResult
				}
			}

			logger.CtxNoticef(ctx.GetCtx(), "err_code", 0)
			c.Header("net_3", fmt.Sprintf("%d", GetMSTimeStamp(time.Now())))
			c.ProtoBuf(200, &pbResp)
		}
	}
	return handlerMW(h, module, view, options...)
}

func ensureT2Set(c *gin.Context) {
	if len(c.Writer.Header().Get("net_2")) == 0 {
		c.Header("net_2", fmt.Sprintf("%d", GetMSTimeStamp(time.Now())))
	}
}

func handlerMW(handlerFn func(c *gin.Context), module, view string, options ...APIOption) func(*gin.Context) {
	return func(c *gin.Context) {
		start := time.Now()
		var ctx *RequestContext
		ctx = NewRequestContext(c)
		c.Set("request_context", ctx)
		ensureT2Set(c)
		defer func() {
			if e := recover(); e != nil {
				var brokenPipe bool
				if ne, ok := e.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						// ugly method used by gin.RecoveryWithWriter()
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") ||
							se.Err == syscall.EPIPE ||
							se.Err == syscall.ECONNRESET {
							brokenPipe = true
						}
					}
				}

				if brokenPipe {
					logger.CtxWarnf(ctx, "recover from broken pipe")
					ctx.GinContext.Abort()
				} else {
					logger.CtxErrorf(ctx, "recover from panic: err=%s stack=%s", e, debug.Stack())
					ctx.GinContext.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()
		Opts := NewAPIOptions(options...)
		c.Set(optionsContextKey, Opts)

		logger.CtxNoticef(ctx.Context, "method:%s", ctx.Request.Method)

		logger.CtxNoticef(ctx.Context, "view:%s", view)
		logger.CtxNoticef(ctx.Context, "module:%s", module)
		logger.CtxNoticef(ctx.Context, "remote_addr:%s", ctx.Request.RemoteAddr)

		c.Set("view_name", view)
		c.Set("module", module)


		handlerFn(c)

		logger.CtxNoticef(ctx.Context, "total_cost(us):%v", time.Since(start).Nanoseconds()/1000)
	}
}


