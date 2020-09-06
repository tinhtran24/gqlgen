package main

import (
	"fmt"
	"os"

	"github.com/jlightning/gqlgen/cmd"
)

func main() {
	fmt.Fprintf(os.Stderr, "warning: running gqlgen from this binary is deprecated and may be removed in a future release. See https://github.com/jlightning/gqlgen/issues/415\n")
	cmd.Execute()
}
