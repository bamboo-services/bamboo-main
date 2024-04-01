package main

import (
	_ "develop/internal/packed"
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"

	"github.com/gogf/gf/v2/os/gctx"

	"develop/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
