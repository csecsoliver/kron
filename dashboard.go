package main

import (
	//"github.com/pocketbase/pocketbase"
	//"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	//"github.com/pocketbase/pocketbase/plugins/migratecmd"
	//. "maragu.dev/gomponents"
	//. "maragu.dev/gomponents/html"

	_ "kron/migrations"
	"kron/views"
)
func gDashboard(r *core.RequestEvent) error {
	html := views.DashboardHome()
	return html.Render(r.Response)
}