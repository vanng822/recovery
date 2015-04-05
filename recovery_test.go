package recovery

import (
	"github.com/stretchr/testify/assert"
	"github.com/vanng822/r2router"
	"github.com/codegangsta/negroni"
	//"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSeeforRecovery(t *testing.T) {
	router := r2router.NewSeeforRouter()
	rec := NewRecovery()
	router.Before(rec.Handler)

	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		panic("This shouldn't crash Seefor")
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	// get
	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, string(content), "Internal Server Error")
}

func TestSeeforRecoveryMiddlewarePanic(t *testing.T) {
	router := r2router.NewSeeforRouter()
	rec := NewRecovery()
	rec.PrintStack = true
	router.Before(rec.Handler)

	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		panic("This shouldn't crash Seefor")
	})

	router.Before(r2router.WrapBeforeHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Middleware panic")
	})))

	ts := httptest.NewServer(router)
	defer ts.Close()

	// get
	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
	assert.Contains(t, string(content), "Middleware panic")
	assert.NotContains(t, string(content), "This shouldn't crash Seefor")
}

func TestSeeforRecoveryPrintStack(t *testing.T) {
	router := r2router.NewRouter()
	n := negroni.New()
	rec := NewRecovery()
	rec.PrintStack = true
	n.UseFunc(rec.HandlerFuncWithNext)
	
	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		panic("This shouldn't crash Seefor")
	})
	n.UseHandler(router)
	ts := httptest.NewServer(n)
	defer ts.Close()

	// get
	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
	assert.Contains(t, string(content), "This shouldn't crash Seefor")
}
