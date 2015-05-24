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
	// frontend
	frontend := web.NewFrontend("jamesclonk.io")

	// cms
	c, err := cms.NewCMS(frontend)
	if err != nil {
		log.Fatal(err)
	}

	// setup routes
	frontend.NewRoute("/", index)

	frontend.NewRoute("/refresh", c.RefreshHandler)

	frontend.NewRoute("/101/{.*}", c.ViewHandler)
	frontend.NewRoute("/101/{.*}/{.*}", c.ViewHandler)
	frontend.NewRoute("/goty/{.*}", c.ViewHandler)

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
		ActiveLink: "/",
		Content:    nil,
		Template:   "index",
	}
}

func createError(w http.ResponseWriter, req *http.Request) *web.Page {
	return web.Error("jamesclonk.io", http.StatusInternalServerError, fmt.Errorf("Oops!"))
}
