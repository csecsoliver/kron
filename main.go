package main

import (
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"

	_ "kron/migrations"

	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {
	app := pocketbase.New()

	// loosely check if it was executed using "go run"
	//isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Dashboard
		// (the isGoRun check is to enable it only during development)
		Automigrate: true,
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/static/{path...}", apis.Static(os.DirFS("./pb_public"), false))
		se.Router.GET("/hello/{name}", hello)

		return se.Next()
	})
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func hello(r *core.RequestEvent) error {
	name := r.Request.PathValue("name")
	return Baselayout("Hello world",
		Div(
			Text("hello "+name),
		),
	).Render(r.Response)
}
