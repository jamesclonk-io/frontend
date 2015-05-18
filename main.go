package main

import (
	"fmt"
	"net/http"

	"github.com/jamesclonk-io/frontend/modules/things"
	"github.com/jamesclonk-io/stdlib/web"
	"github.com/jamesclonk-io/stdlib/web/negroni"
)

func main() {
	frontend := web.NewFrontend()

	// setup navigation
	navbar := web.NavBar{
		web.NavElement{"Home", "/", nil},
		web.NavElement{"101", "/refresh", nil},
		web.NavElement{"Throw 404", "/contact", nil},
		web.NavElement{"Throw Error", "/error", nil},
		web.NavElement{"Menu", "#", []web.NavElement{
			web.NavElement{"Action", "/action", nil},
			web.NavElement{"Something else here", "/something_else", nil},
			web.NavElement{"Link", "/link", nil},
			web.NavElement{"Another Link", "/more_link", nil},
		}},
	}
	frontend.SetNavigation(navbar)

	// setup routes
	frontend.NewRoute("/", index)

	// ThingsRefreshHandler will modify navigation (101 dropdown list)
	frontend.NewRoute("/refresh", things.ThingsRefreshHandler(&frontend.PageMaster.Navbar, 1))
	frontend.NewRoute("/101/{file:.*}", things.ThingsViewHandler)

	frontend.NewRoute("/link", index)
	frontend.NewRoute("/error", createError)

	// setup negroni
	n := negroni.Sbagliato()
	n.UseHandler(frontend.Router)

	// start web server
	server := web.NewServer()
	server.Start(n)
}

func index(w http.ResponseWriter, req *http.Request) *web.Page {
	return &web.Page{
		Title:            "jamesclonk.io",
		ActiveNavElement: "Home",
		Content:          nil,
		Template:         "index",
	}
}

func createError(w http.ResponseWriter, req *http.Request) *web.Page {
	return web.Error("jamesclonk.io", http.StatusInternalServerError, fmt.Errorf("Oops!"))
}
