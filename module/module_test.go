package module

import (
	"os"
	"testing"

	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

func TestModule(t *testing.T) {
	req, err := os.Open("../example/code_generator_request.pb.bin")
	if err != nil {
		t.Fatal(err)
	}

	// fs := afero.NewMemMapFs()
	// res := &bytes.Buffer{}
	pgs.Init(
		pgs.ProtocInput(req), // use the pre-generated request
		// pgs.ProtocOutput(res), // capture CodeGeneratorResponse
		// pgs.FileSystem(fs),    // capture any custom files written directly to disk
	).
		RegisterModule(
			//	ASTPrinter(),
			New(),
		).
		RegisterPostProcessor(
			pgsgo.GoImports(),
			pgsgo.GoFmt(),
		).
		Render()

}
