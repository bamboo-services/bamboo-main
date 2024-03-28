package main

import (
	_ "develop/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"develop/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
