package things

import (
	"html/template"
	"net/http"
	"path"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/jamesclonk-io/stdlib/env"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/jamesclonk-io/stdlib/web"
)

var (
	data101      *Things
	data101File  string
	data101Mutex *sync.Mutex
	log          *logrus.Logger
)

func init() {
	data101 = &Things{}
	data101File = env.Get("JCIO_101_THINGS_DATA", "https://github.com/jamesclonk-io/101-things/archive/master.zip")
	data101Mutex = &sync.Mutex{}
	log = logger.GetLogger()
}

func ThingsViewHandler(w http.ResponseWriter, req *http.Request) *web.Page {
	vars := mux.Vars(req)
	filename := vars["file"]
	file := path.Join("/", filename)

	// find file
	var html template.HTML
	for _, f := range data101.Files {
		if path.Join("/", f.Path, f.Name) == file {
			html = f.Content
		}
	}

	// wrap into struct
	content := struct {
		Title string
		HTML  template.HTML
	}{
		Title: filename,
		HTML:  html,
	}

	return &web.Page{
		Title:            "jamesclonk.io - 101 Things - " + filename,
		ActiveNavElement: "101",
		Content:          content,
		Template:         "things",
	}
}

func ThingsRefreshHandler(navbar *web.NavBar, thingsIndex int) web.Handler {
	return func(w http.ResponseWriter, req *http.Request) *web.Page {
		data101Mutex.Lock()
		defer data101Mutex.Unlock()

		if err := checkData(navbar, thingsIndex, true); err != nil {
			return web.Error("jamesclonk.io", http.StatusInternalServerError, err)
		}

		return &web.Page{
			Title:            "jamesclonk.io - Refresh",
			ActiveNavElement: "Home",
			Content:          nil,
			Template:         "index",
		}
	}
}

func checkData(navbar *web.NavBar, navIndex int, refresh bool) error {
	// reset data if not set
	if data101 == nil {
		data101 = &Things{}
	}

	// refresh either every 12 hours, or if refresh parameter set to true
	if time.Since(data101.Timestamp).Hours() > 12 || refresh {
		if err := refreshData(data101File); err != nil {
			return err
		}

		// create new navbar elements
		navElements := make([]web.NavElement, 0)
		for _, file := range data101.Files {
			navElements = append(navElements, web.NavElement{
				Name:     path.Join("/", file.Path, file.Basename),
				Link:     path.Join("/101", file.Path, file.Name),
				Dropdown: nil,
			})
		}

		// reset navigation bar element for "101"
		(*navbar)[navIndex].Dropdown = navElements
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
