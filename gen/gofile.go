package main

import (
	"fmt"
	"os"
	"os/exec"
)

type goFile struct {
	*os.File
}

// createGoFile is an expirement in code generation. It wraps os.Create, with a header,
// and adds a `go fmt` to the Close function.
func createGoFile(pkg, name string) (*goFile, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(f, "package %s\n\n", pkg)
	fmt.Fprintf(f, "\n")
	fmt.Fprintf(f, "//////////////////////////////\n")
	fmt.Fprintf(f, "// THIS IS A GENERATED FILE //\n")
	fmt.Fprintf(f, "//////////////////////////////\n")
	fmt.Fprintf(f, "\n")

	return &goFile{f}, nil
}

func (gf *goFile) Close() error {
	if err := gf.File.Close(); err != nil {
		return err
	}

	cmd := exec.Command("gofmt", "-w", gf.File.Name())
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error formatting %s: %s", gf.File.Name(), stdoutStderr)
	}

	return nil
}
