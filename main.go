package main

import (
	_ "bamboo-main/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"bamboo-main/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
