package main

import (
	"io/fs"
	"log"
	"net/http"

	//"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	// . "maragu.dev/gomponents"
	// . "maragu.dev/gomponents/html"

	"embed"
	_ "kron/migrations"
	"kron/views"
)

//go:embed pb_public/*
var pb_public embed.FS

func main() {
	app := pocketbase.New()
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: true,
	})
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		err := reloadJobs(app)
		if err != nil {
			println("Failed to reload jobs: " + err.Error())
		}
		se.Router.BindFunc(func(e *core.RequestEvent) error {
			cookie, err := e.Request.Cookie("pb_auth")
			if err == nil {
				record, err := e.App.FindAuthRecordByToken(cookie.Value)
				if err == nil {
					e.Auth = record
				}
			}
			return e.Next()
		})
		var staticFS, _ = fs.Sub(pb_public, "pb_public")
		se.Router.GET("/static/{path...}", apis.Static(staticFS, false))
		se.Router.GET("/", func(r *core.RequestEvent) error { return views.HomePage().Render(r.Response) })
		se.Router.GET("/login", gLogin)
		se.Router.POST("/login", pLogin)
		se.Router.GET("/dash/jobs", gJobs)
		se.Router.POST("/dash/jobs", pJob)
		se.Router.GET("/dash/jobs/{id}", gJob)
		//se.Router.DELETE("/dash/jobs/{id}", dJob)

		se.Router.GET("/dash", gDashboard)
		
		se.Router.POST("/test", func(r *core.RequestEvent) error {
			return r.String(200, "test")
		})
		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

//	func hello(r *core.RequestEvent) error {
//		name := r.Request.PathValue("name")
//		return views.BaseLayout("Hello world",
//			Div(
//				Text("hello "+name),
//			),
//		).Render(r.Response)
//	}
func gLogin(r *core.RequestEvent) error {
	if (r.Auth != nil) {
		r.Redirect(302, "/dash")
	}
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
func reloadJobs(app core.App) error {

	records, err := app.FindAllRecords("jobs")
	if err != nil {
		return err
	}
	for _, record := range records {
		job, err := jobRecordToStruct(record)
		if err != nil {
			return err
		}
		err = job.RegisterCron(app)
		if err != nil {
			return err
		}
	}
	return nil
}
