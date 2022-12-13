package common

import (
	"github.com/cloudwego/kitex/client"
	"github.com/xiangqin/user_api/kitex_gen/user/user"
)

var (
	UserCli user.Client
)

func Init() {

	UserCli, _ = user.NewClient("xiangqin.user.core", client.WithHostPorts("10.21.21.9:8888"))
}
