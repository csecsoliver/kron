package views

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func HomePage() Node {
	return PublicLayout(
		"Home",
		Div(
			P(Text(`This service can be used to check the state of stateless, non-daeomonized and worker type apps in a simple way.`)),
			P(Text(`The app or the user can specify an http request to be made, and Kron will make that request and if it doesn't recieve the expected response, notify the user.`)),
			P(Text(`This can be used as a healthcheck, cron job manager, etc.`)),
		),
	)
}
