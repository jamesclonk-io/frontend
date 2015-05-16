package things

import (
	"net/http"
	"path"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/jamesclonk-io/stdlib/env"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/jamesclonk-io/stdlib/web"
)

var (
	data101     *Things
	data101File string
	log         *logrus.Logger
)

func init() {
	data101File = env.Get("JCIO_101_THINGS_DATA", "https://github.com/jamesclonk-io/101-things/archive/master.zip")
	log = logger.GetLogger()
}

func ThingsHandler(navbar *web.NavBar, thingsIndex int) web.Handler {
	return func(w http.ResponseWriter, req *http.Request) *web.Page {
		if err := checkData(navbar, thingsIndex, req.URL.Query().Get("refresh") == "true"); err != nil {
			return web.Error("jamesclonk.io", http.StatusInternalServerError, err)
		}

		vars := mux.Vars(req)
		file := path.Join("/", vars["file"])

		// find file
		var content string
		for _, f := range data101.Files {
			if path.Join("/", f.Path, f.Name) == file {
				content = f.Content
			}
		}

		return &web.Page{
			Title:    "jamesclonk.io - 101 Things - " + vars["file"],
			Content:  content,
			Template: "index",
		}
	}
}

func checkData(navbar *web.NavBar, thingsIndex int, refresh bool) error {
	// reset data if not set
	if data101 == nil {
		data101 = &Things{}
	}

	// refresh either every 12 hours, or if refresh parameter set to true
	if time.Since(data101.Timestamp).Hours() > 12 || refresh {
		if err := refreshData(data101File); err != nil {
			return err
		}

		navElements := make([]web.NavElement, 0)
		for _, file := range data101.Files {
			// create new navbar element
			navElements = append(navElements, web.NavElement{
				Name:     path.Join("/", file.Path, file.Basename),
				Link:     path.Join("/101", file.Path, file.Name),
				Dropdown: nil,
			})
		}

		// reset navigation bar element
		(*navbar)[thingsIndex].Dropdown = navElements
	}
	return nil
}

func refreshData(input string) (err error) {
	data101, err = getData(input)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
			"file":  data101File,
		}).Error("Could not refresh data")
		return err
	}
	return nil
}
