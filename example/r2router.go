package main

import (
	"github.com/vanng822/r2router"
	"github.com/vanng822/recovery"
	"net/http"
	"fmt"
)

func main() {
	seefor := r2router.NewSeeforRouter()
	rec := recovery.NewRecovery()
	rec.PrintStack = true
	seefor.Before(rec.Handler)
	
	seefor.Get("/hello/:name", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		fmt.Fprintf(w, "Hello %s!", p.Get("name"))
	})
	
	seefor.Get("/panic", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		panic("Me Panic")
	})

	http.ListenAndServe(":8080", seefor)
}
