// Code generated by goa v3.3.1, DO NOT EDIT.
//
// Code Generator
//
// Command:
// goa

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/fieldkit/cloud/server/api/design"
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/generator"
	"goa.design/goa/v3/eval"
	goa "goa.design/goa/v3/pkg"
)

func main() {
	var (
		out     = flag.String("output", "", "")
		version = flag.String("version", "", "")
		cmdl    = flag.String("cmd", "", "")
		ver     int
	)
	{
		flag.Parse()
		if *out == "" {
			fail("missing output flag")
		}
		if *version == "" {
			fail("missing version flag")
		}
		if *cmdl == "" {
			fail("missing cmd flag")
		}
		v, err := strconv.Atoi(*version)
		if err != nil {
			fail("invalid version %s", *version)
		}
		ver = v
	}

	if ver > goa.Major {
		fail("cannot run goa %s on design using goa v%s\n", goa.Version(), *version)
	}
	if err := eval.Context.Errors; err != nil {
		fail(err.Error())
	}
	if err := eval.RunDSL(); err != nil {
		fail(err.Error())
	}
	if err := os.RemoveAll("gen/firmware"); err != nil {
		fail(err.Error())
	}
	if err := os.RemoveAll("gen/tasks"); err != nil {
		fail(err.Error())
	}
	if err := os.RemoveAll("gen/ingestion"); err != nil {
		fail(err.Error())
	}
	if err := os.RemoveAll("gen/records"); err != nil {
		fail(err.Error())
	}
	if err := os.RemoveAll("gen/information"); err != nil {
		fail(err.Error())
	}
	if err := os.RemoveAll("gen/test"); err != nil {
		fail(err.Error())
	}
	if err := os.RemoveAll("gen/activity"); err != nil {
		fail(err.Error())
	}
	if err := os.RemoveAll("gen/oidc"); err != nil {
		fail(err.Error())
	}
	if err := os.RemoveAll("gen/user"); err != nil {
		fail(err.Error())
	}
	if err := os.RemoveAll("gen/project"); err != nil {
		fail(err.Error())
	}
	if err := os.RemoveAll("gen/notes"); err != nil {
		fail(err.Error())
	}
	if err := os.RemoveAll("gen/http"); err != nil {
		fail(err.Error())
	}
	if err := os.RemoveAll("gen/csv"); err != nil {
		fail(err.Error())
	}
	if err := os.RemoveAll("gen/sensor"); err != nil {
		fail(err.Error())
	}
	if err := os.RemoveAll("gen/station"); err != nil {
		fail(err.Error())
	}
	if err := os.RemoveAll("gen/discourse"); err != nil {
		fail(err.Error())
	}
	if err := os.RemoveAll("gen/following"); err != nil {
		fail(err.Error())
	}
	if err := os.RemoveAll("gen/export"); err != nil {
		fail(err.Error())
	}
	if err := os.RemoveAll("gen/discussion"); err != nil {
		fail(err.Error())
	}
	if err := os.RemoveAll("gen/modules"); err != nil {
		fail(err.Error())
	}
	if err := os.RemoveAll("gen/data"); err != nil {
		fail(err.Error())
	}
	codegen.DesignVersion = ver
	outputs, err := generator.Generate(*out, "gen")
	if err != nil {
		fail(err.Error())
	}

	fmt.Println(strings.Join(outputs, "\n"))
}

func fail(msg string, vals ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, vals...)
	os.Exit(1)
}