package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xiangqin/user_api/common"
)

func main() {
	r := gin.Default()
	common.Init()
	defer r.Run(":8080")
	register(r)
}
