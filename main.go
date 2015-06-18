package main

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/jcio-frontend/modules/newsfeed"
	"github.com/jamesclonk-io/jcio-frontend/modules/quotes"
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
	// setup http handler
	n := setup()

	// start web server
	server := web.NewServer()
	server.Start(n)
}

func setup() *negroni.Negroni {
	frontend := frontend()

	// setup negroni
	n := negroni.Sbagliato()
	n.UseHandler(quotes.NewQuoteMiddleware(frontend))
	n.UseHandler(frontend.Router)

	return n
}

func frontend() *web.Frontend {
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
	newsfeed.UpdateFeeds(news)

	// setup routes
	frontend.NewRoute("/", index)
	frontend.NewRoute("/example", example)

	frontend.NewRoute("/refresh", c.RefreshHandler)

	frontend.NewRoute("/news", news.ViewHandler)

	frontend.NewRoute("/101/{.*}", c.ViewHandler)
	frontend.NewRoute("/101/{.*}/{.*}", c.ViewHandler)
	frontend.NewRoute("/goty/{.*}", c.ViewHandler)
	frontend.NewRoute("/static/{.*}", c.ViewHandler)

	frontend.NewRoute("/error/{.*}", createError)

	return frontend
}

func index(w http.ResponseWriter, req *http.Request) *web.Page {
	return &web.Page{
		ActiveLink: "/",
		Content:    nil,
		Template:   "index",
	}
}

func example(w http.ResponseWriter, req *http.Request) *web.Page {
	return &web.Page{
		ActiveLink: "/",
		Content:    nil,
		Template:   "example",
	}
}

func createError(w http.ResponseWriter, req *http.Request) *web.Page {
	return web.Error("jamesclonk.io - Error", http.StatusInternalServerError, fmt.Errorf("Error!"))
}
