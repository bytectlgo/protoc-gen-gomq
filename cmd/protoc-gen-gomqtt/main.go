package main

import (
	"github.com/bytectlgo/protoc-gen-gomq/cmd/protoc-gen-gomqtt/module"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

func main() {

	pgs.Init(
		pgs.DebugEnv("GOMQ_DEBUG"),
	).
		RegisterModule(module.New()).
		RegisterPostProcessor(
			pgsgo.GoImports(),
			pgsgo.GoFmt()).
		Render()
}
