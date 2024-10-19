package handler

import (
	"html/template"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/config"
)

type HomeData struct {
	LoginURL             string
	SwaggerURL           string
	LoginInternalUserURL string
	AppVersion           string
	AppBuildTime         string
}

func HomeHandler(w http.ResponseWriter, _ *http.Request) {
	t, _ := template.ParseFiles("internal/webserver/templates/home.html")

	data := HomeData{
		SwaggerURL:           "docs/index.html",
		LoginInternalUserURL: "auth/internal",
		LoginURL:             config.GetApplicationURL(),
		AppVersion:           config.GetBuildInfo().Version,
		AppBuildTime:         config.GetBuildInfo().BuildTime,
	}
	_ = t.Execute(w, data)
}
