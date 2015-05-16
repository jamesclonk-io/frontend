package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/stdlib/env"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/jamesclonk-io/stdlib/web"
	"github.com/jamesclonk-io/stdlib/web/negroni"
)

var (
	data101     *ZipFile
	data101File string
	log         *logrus.Logger
)

func init() {
	data101File = env.Get("JCIO_101_THINGS_DATA_FILE", "https://github.com/jamesclonk-io/101-things/archive/content.zip")
	log = logger.GetLogger()
}

func main() {
	// setup routes
	router := web.NewRouter()
	router.NewRoute("/", "jamesclonk.io", "Home", index)
	router.NewRoute("/about", "jamesclonk.io", "About", things)
	router.NewRoute("/link", "jamesclonk.io", "Dropdown", index)
	router.Handle("/error", web.ErrorHandler("jamesclonk.io", "Error", router.Render, fmt.Errorf("Crap!")))

	// setup negroni
	n := negroni.Sbagliato()
	n.UseHandler(router)

	// start web server
	server := web.NewServer()
	server.Start(n)
}

func index(w http.ResponseWriter, req *http.Request) *web.Page {
	return &web.Page{
		Content:  nil,
		Template: "index",
	}
}

func things(w http.ResponseWriter, req *http.Request) *web.Page {
	// reset data if not set
	if data101 == nil {
		data101 = &ZipFile{}
	}

	// refresh either every 12 hours, or if refresh parameter set to true
	if time.Since(data101.Timestamp).Hours() > 12 ||
		req.URL.Query().Get("refresh") == "true" {
		if err := refreshData(data101File); err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
				"file":  data101File,
			}).Error("Could not refresh data")

			return &web.Page{
				StatusCode: http.StatusInternalServerError,
				Error:      err,
			}
		}
		for _, folder := range data101.Folders {
			fmt.Println(folder)
			for _, file := range data101.Data[folder].Files {
				fmt.Printf("%#v\n", file.Name)
				fmt.Printf("%#v\n", file.Basename)
				fmt.Printf("%#v\n", file.Path)
			}
		}
	}

	return &web.Page{
		Content:  nil,
		Template: "index",
	}
}

func refreshData(input string) (err error) {
	data101, err = getData(input)
	if err != nil {
		return err
	}
	return nil
}
