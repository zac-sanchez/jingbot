package web

import (
	"embed"
	"html/template"
	"net/http"
)

type FrontendData interface {
	GetData() interface{}
}

type HomepageData struct {
	PageTitle string
	Endpoints []string
}

func (hpd *HomepageData) GetData() interface{} {
	return hpd
}

//go:embed static
var static embed.FS

func Frontend(data FrontendData, path string) http.HandlerFunc {
	t, err := template.ParseFS(static, path)
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/html")
		t.Execute(w, data.GetData())
	}
}
