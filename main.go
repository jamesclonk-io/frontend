package main

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/jamesclonk-io/stdlib/web"
	"github.com/jamesclonk-io/stdlib/web/cms"
	"github.com/jamesclonk-io/stdlib/web/negroni"
)

var (
	log *logrus.Logger
)

func init() {
	log = logger.GetLogger()
}

func main() {
	frontend := web.NewFrontend()

	// setup navigation
	navbar, err := cms.GetNavBar()
	if err != nil {
		log.Fatal(err)
	}
	frontend.SetNavigation(navbar)

	// setup routes
	frontend.NewRoute("/", index)

	// ThingsRefreshHandler will modify navigation (101 dropdown list)
	frontend.NewRoute("/refresh", cms.ThingsRefreshHandler(&frontend.PageMaster.Navbar, 1))
	frontend.NewRoute("/101/{.*}", cms.ViewHandler)

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
		Title:      "jamesclonk.io",
		ActiveLink: "/",
		Content:    nil,
		Template:   "index",
	}
}

func createError(w http.ResponseWriter, req *http.Request) *web.Page {
	return web.Error("jamesclonk.io", http.StatusInternalServerError, fmt.Errorf("Oops!"))
}
