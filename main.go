package main

import (
	_ "xiaoMain/internal/packed"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"

	_ "xiaoMain/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	"xiaoMain/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
