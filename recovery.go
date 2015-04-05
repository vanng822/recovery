// Package recovery implements couple of middleware interfaces
//
//	package main
//
//	import (
//		"github.com/vanng822/r2router"
//		"github.com/vanng822/recovery"
//		"net/http"
//	)
//
//	func main() {
//		seefor := r2router.NewSeeforRouter()
//		rec := recovery.NewRecovery()
//		rec.PrintStack = true
//		seefor.Before(rec.Handler)
//		seefor.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, _ r2router.Params) {
//			panic("Middleware panic")
//		})
//
//		http.ListenAndServe(":8080", seefor)
//	}
package recovery

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type Recovery struct {
	Logger     Logger
	StackAll   bool
	StackSize  int
	PrintStack bool
}

func (rec *Recovery) recovery(w http.ResponseWriter) {
	if err := recover(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		stack := make([]byte, rec.StackSize)
		stack = stack[:runtime.Stack(stack, rec.StackAll)]
		format := "PANIC: %s\n%s"
		rec.Logger.Printf(format, err, stack)

		if rec.PrintStack {
			fmt.Fprintf(w, format, err, stack)
		} else {
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		}
	}
}

func NewRecovery() *Recovery {
	return &Recovery{
		Logger:     log.New(os.Stdout, "[error] ", 0),
		StackAll:   false,
		StackSize:  1024 * 8,
		PrintStack: false,
	}
}

func (rec *Recovery) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer rec.recovery(w)
		next.ServeHTTP(w, req)
	})
}

func (rec *Recovery) HandlerFuncWithNext(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer rec.recovery(w)
	next(w, r)
}
