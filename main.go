package main

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/frontend/modules/quotes"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/jamesclonk-io/stdlib/web"
	"github.com/jamesclonk-io/stdlib/web/cms"
	"github.com/jamesclonk-io/stdlib/web/negroni"
	"github.com/jamesclonk-io/stdlib/web/newsreader"
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

	// newsreader
	news, err := newsreader.NewReader(frontend, c.GetConfiguration())
	if err != nil {
		log.Fatal(err)
	}
	news.GetFeeds()

	// setup routes
	frontend.NewRoute("/", index)
	frontend.NewRoute("/refresh", c.RefreshHandler)

	frontend.NewRoute("/news", news.ViewHandler)

	frontend.NewRoute("/101/{.*}", c.ViewHandler)
	frontend.NewRoute("/101/{.*}/{.*}", c.ViewHandler)
	frontend.NewRoute("/goty/{.*}", c.ViewHandler)
	frontend.NewRoute("/static/{.*}", c.ViewHandler)

	// setup negroni
	n := negroni.Sbagliato()
	n.UseHandler(quotes.NewQuoteMiddleware(frontend))
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
