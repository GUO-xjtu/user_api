package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xiangqin/user_api/biz/handler"
)

func register(r *gin.Engine) {
	xiangQinUser(r)
	r.GET("/ping", handler.Ping)
}


