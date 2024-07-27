package routes

import (
	"net/http"

	components "github.com/sfluor/test-htmx/views/components/buldan"
	"github.com/sfluor/test-htmx/views/model/buldan"
)

func RegisterBuldan(engine buldan.Engine) {
    home := components.BuldanHome()
    http.HandleFunc("GET /buldan", func(w http.ResponseWriter, r *http.Request) {
		home.Render(r.Context(), w)
    })
}
