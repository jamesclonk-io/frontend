package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/jamesclonk-io/stdlib/web/negroni"
	"github.com/stretchr/testify/assert"
)

var (
	m *negroni.Negroni
)

func init() {
	os.Setenv("PORT", "3003")
	os.Setenv("JCIO_CMS_DATA", "https://github.com/jamesclonk-io/content/archive/master.zip")
	logrus.SetOutput(ioutil.Discard)
	logger.GetLogger().Out = ioutil.Discard
	m = setup()
}

func Test_Main_Setup(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3003/", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>jamesclonk.io</title>`)
	assert.Contains(t, body, `<img src="/images/welcome.png" class="welcome-picture" alt="jamesclonk.io"/>`)
}

func Test_Main_404(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3003/something", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusNotFound, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>jamesclonk.io</title>`)
	assert.Contains(t, body, `<div class="alert alert-warning">This is not the page you are looking for..</div>`)
}

func Test_Main_500(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3003/error/something", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusInternalServerError, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>jamesclonk.io - Error</title>`)
	assert.Contains(t, body, `<div class="alert alert-danger">Error: Error!</div>`)
}

func Test_Main_Index(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>jamesclonk.io</title>`)
	assert.Contains(t, body, `<img src="/images/welcome.png" class="welcome-picture" alt="jamesclonk.io"/>`)
}

func Test_Main_Example(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/example", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>jamesclonk.io</title>`)
	assert.Contains(t, body, `<button type="button" class="btn btn-lg btn-danger">Danger</button>`)
}

func Test_Main_News(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3003/news", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>jamesclonk.io - News</title>`)
	assert.Contains(t, body, `Hacker News`)
	assert.Contains(t, body, `<a href="https://news.ycombinator.com/" target="_new" class="list-group-item active">`)
	assert.Contains(t, body, `<p class="list-group-item-text"><i class="fa fa-external-link fa-fw"></i> https://news.ycombinator.com/</p>`)
	assert.Contains(t, body, `<a href="https://news.ycombinator.com/item?id=`)

	assert.Contains(t, body, `Ars Technica`)
	assert.Contains(t, body, `<a href="https://arstechnica.com" target="_new" class="list-group-item active">`)
	assert.Contains(t, body, `<p class="list-group-item-text"><i class="fa fa-external-link fa-fw"></i> https://arstechnica.com</p>`)
	
	//assert.Contains(t, body, `<a href="http://www.heise.de/newsticker/meldung/`)
	//assert.Contains(t, body, `<a href="https://www.reddit.com/r/technology`)
}

func Test_Main_101(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3003/101/Links", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>jamesclonk.io - Links</title>`)
	assert.Contains(t, body, `<li><a href="https://github.com/JamesClonk">https://github.com/JamesClonk</a></li>`)
	assert.Contains(t, body, `<li><a href="http://golang.org/doc/effective_go.html">http://golang.org/doc/effective_go.html</a></li>`)

	// 404 response expected
	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "http://localhost:3003/101/does_not_exist", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusNotFound, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `<title>jamesclonk.io - does_not_exist - Not Found</title>`)
	assert.Contains(t, body, `<div class="alert alert-warning">This is not the page you are looking for..</div>`)
}

func Test_Main_GOTY(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3003/goty/2009", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>jamesclonk.io - 2009</title>`)
	assert.Contains(t, body, `Game of the Year`)
	assert.Contains(t, body, `<a href="http://en.wikipedia.org/wiki/Uncharted_2:_Among_Thieves">Uncharted 2: Among Thieves</a>`)
	assert.Contains(t, body, `<a href="http://store.steampowered.com/app/47810/">Dragon Age: Origins</a>`)
	assert.Contains(t, body, `Fortune favours the bold..`)
	assert.Contains(t, body, `<img src="/images/goty/fortune.jpg" alt="Uncharted 2: Among Thieves" title="Uncharted 2: Among Thieves" />`)

	// 404 response expected
	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "http://localhost:3003/goty/2006", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusNotFound, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `<title>jamesclonk.io - 2006 - Not Found</title>`)
	assert.Contains(t, body, `<div class="alert alert-warning">This is not the page you are looking for..</div>`)
}

func Test_Main_MyMovies(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3003/static/Movies", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>jamesclonk.io - Movies</title>`)
	assert.Contains(t, body, `Hier meine krassen Filmchen!`)
	assert.Contains(t, body, `<a href="http://files.jamesclonk.ch/movies/klickschtduhier.avi">Klickscht du hier!</a>`)
	assert.Contains(t, body, `<img src="/images/icon_catch.gif" alt=":catch:" title=":catch:" />`)
}

func Test_Main_Quake3(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3003/static/Quake", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>jamesclonk.io - Quake</title>`)
	assert.Contains(t, body, `Quake 3 CPMA`)
	assert.Contains(t, body, `<a href="http://www.playmorepromode.org/">http://www.playmorepromode.org/</a>`)
	assert.Contains(t, body, `<a href="http://ioquake3.org/">http://ioquake3.org/</a>`)
	assert.Contains(t, body, `<img src="/images/cpma1.jpg" alt="Quake 3 CPMA" title="Quake 3 CPMA" />`)
}

func Test_Main_Gallery(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3003/static/Gallery", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>jamesclonk.io - Gallery</title>`)
	assert.Contains(t, body, `Japan / Korea`)
	assert.Contains(t, body, `mp4 - AVC / AAC - High Resolution / 720p`)
	assert.Contains(t, body, `<a href="http://files.jamesclonk.ch/movies/japan_day_six.mp4">Day Six</a>`)
	assert.Contains(t, body, `<a href="http://www.jamesclonk.ch/gallery/">http://www.jamesclonk.ch/gallery/</a>`)
}
