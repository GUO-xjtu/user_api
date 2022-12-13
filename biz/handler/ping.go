package handler

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	logger.CtxInfof(c, "a sample app log")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

