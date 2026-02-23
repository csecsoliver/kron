package main

import (
	"log"
	"net/http"
	"os"
	//"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"

	_ "kron/migrations"
	"kron/views"
)

func main() {
	app := pocketbase.New()

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: true,
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/static/{path...}", apis.Static(os.DirFS("./pb_public"), false))
		se.Router.GET("/", func(r *core.RequestEvent) error { return views.HomePage().Render(r.Response) })
		se.Router.GET("/login", gLogin)
		se.Router.POST("/login", pLogin)
		se.Router.GET("/dash/jobs", gJobs)
		//se.Router.POST("/dash/jobs", pJobs)
		//se.Router.GET("/dash/jobs/{id}", gJob)
		//se.Router.DELETE("/dash/jobs/{id}", dJob)
		
		
		se.Router.GET("/dash", gDashboard)
		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func hello(r *core.RequestEvent) error {
	name := r.Request.PathValue("name")
	return views.BaseLayout("Hello world",
		Div(
			Text("hello "+name),
		),
	).Render(r.Response)
}
func gLogin(r *core.RequestEvent) error {

	return views.LoginPage("").Render(r.Response)
}
func pLogin(r *core.RequestEvent) error {
	email := r.Request.FormValue("email")
	password := r.Request.FormValue("password")
	record, err := r.App.FindAuthRecordByEmail("users", email)
	if err != nil { // if the user does not exist
		collection, err := r.App.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		record = core.NewRecord(collection)
		record.SetEmail(email)
		record.SetPassword(password)
		if err := r.App.Save(record); err != nil {
			return views.LoginPage("Failed to create user: " + err.Error()).Render(r.Response)
		}
	} else { // when the user exists
		if !record.ValidatePassword(password) {
			return views.LoginPage("Invalid Email or Password").Render(r.Response)
		}

	}
	token, err := record.NewAuthToken()
	if err != nil {
		return views.LoginPage("Failed to create auth token: " + err.Error()).Render(r.Response)
	}

	http.SetCookie(r.Response, &http.Cookie{
		Name:     "pb_auth",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		// Expires:  time.Now().Add(time.Hour * 10),
	})
	return r.Redirect(302, "/dash")
}