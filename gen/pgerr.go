package main

import (
	"fmt"
)

type pgErrType struct {
	sqlstate  string
	code      string // This could be typed, but it doesn't worth the effort.
	macroName string
	spec_name string
	name      string
}

func (pe pgErrType) Name() string {
	return pe.name
}

func (pe pgErrType) String() string {
	return fmt.Sprintf("%s: %s", pe.sqlstate, pe.name)
}

func (pge pgErrType) Error() string {
	return pge.String()
}
