package main

import (
	_ "github.com/fieldkit/cloud/server/api/design"

	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/goagen/codegen"
	genapp "github.com/goadesign/goa/goagen/gen_app"
	genclient "github.com/goadesign/goa/goagen/gen_client"
	// genmain "github.com/goadesign/goa/goagen/gen_main"
	genschema "github.com/goadesign/goa/goagen/gen_schema"
	genswagger "github.com/goadesign/goa/goagen/gen_swagger"
)

func main() {
	codegen.ParseDSL()
	codegen.Run(
		// genmain.NewGenerator(
		// 	genmain.API(design.Design),
		// 	genmain.OutDir("api"),
		// ),
		genapp.NewGenerator(
			genapp.API(design.Design),
			genapp.OutDir("api/app"),
			genapp.Target("app"),
			genapp.NoTest(true),
		),
		genswagger.NewGenerator(
			genswagger.OutDir("api/public"),
			genswagger.API(design.Design),
		),
		genclient.NewGenerator(
			genclient.OutDir("api"),
			genclient.API(design.Design),
		),
		genschema.NewGenerator(
			genschema.OutDir("api/public"),
			genschema.API(design.Design),
		),
	)
}
