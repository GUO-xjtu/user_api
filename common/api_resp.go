package common

import (
	"encoding/json"
	"fmt"
	"github.com/GUO-xjtu/pb_gen/gen/data"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
	"time"
)

// IsUsePBResp 判断请求返回数据是否 PB 序列化
func IsUsePBResp(ctx *RequestContext) bool {
	content := ctx.Request.Header.Get("response-format")
	if len(content) == 0 || content == "json" { // json or protobuf
		return false
	}
	return true
}

func genErrResp(ctx *RequestContext, err error) {
	logger.CtxNoticef(ctx.Context, "err_code", err)
	ctx.GinContext.Set("status", err)
	if IsUsePBResp(ctx) {
		genPBErrResp(ctx, err)
	} else {
		genJsonErrResp(ctx, err)
	}
}

func genJsonErrResp(ctx *RequestContext, err error) {
	data := map[string]interface{}{}
	if err.Error() != "" {
		data["message"] = err.Error()
	}

	ctx.GinContext.JSON(200, gin.H{
		"data":        data,
	})
}

func genPBErrResp(ctx *RequestContext, err error) {
	pbResp := data.PBResponse{}
	extra, _ := json.Marshal(err.Error())
	pbResp.Header = &data.PBResponse_Header{
		Now:        GetMSTimeStamp(time.Now()),

		Extra: string(extra),
	}

	pbResp.Body = []byte{}
	ctx.GinContext.Header("net_3", fmt.Sprintf("%d", GetMSTimeStamp(time.Now())))
	ctx.GinContext.ProtoBuf(200, &pbResp)
}
