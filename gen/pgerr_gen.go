package main

import (
	"regexp"
	"strings"
)

// # The format of this file is one error code per line, with the following
// # whitespace-separated fields:
// #
// #      sqlstate    E/W/S    errcode_macro_name    spec_name
// #
// # where sqlstate is a five-character string following the SQLSTATE conventions,
// # the second field indicates if the code means an error, a warning or success,
// # errcode_macro_name is the C macro name starting with ERRCODE that will be put
// # in errcodes.h, and spec_name is a lowercase, underscore-separated name that
// # will be used as the PL/pgSQL condition name and will also be included in the
// # SGML list. The last field is optional, if not present the PL/pgSQL condition
// # and the SGML entry will not be generated.
var lineRe = regexp.MustCompile(`^(\w{5})\s+(.)\s+(\w+)\s*(\w+)?`)

func newPgError(line string) (pgErrType, bool) {
	m := lineRe.FindStringSubmatch(line)
	if m == nil {
		return pgErrType{}, false
	}

	pe := pgErrType{
		sqlstate:  m[1],
		code:      m[2],
		macroName: m[3],
		spec_name: m[4],
	}

	if pe.spec_name == "" {
		pe.spec_name = strings.ToLower(strings.TrimPrefix(pe.macroName, "ERRCODE_"))
	}

	switch pe.code {
	case "E":
		pe.spec_name += " error"
	case "W":
		pe.spec_name += " warning"
	}

	pe.name = strings.ReplaceAll(strings.Title(strings.ReplaceAll(pe.spec_name, "_", " ")), " ", "")

	return pe, true
}
