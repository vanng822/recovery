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

type Options struct {
	Logger     Logger
	StackAll   bool
	StackSize  int
	PrintStack bool
}

type Recovery struct {
	options *Options
	next    http.Handler
}

func (rec *Recovery) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			stack := make([]byte, rec.options.StackSize)
			stack = stack[:runtime.Stack(stack, rec.options.StackAll)]
			format := "PANIC: %s\n%s"
			rec.options.Logger.Printf(format, err, stack)

			if rec.options.PrintStack {
				fmt.Fprintf(w, format, err, stack)
			} else {
				w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			}
		}
	}()
	rec.next.ServeHTTP(w, req)
}

func NewOptions() *Options {
	return &Options{
		Logger:     log.New(os.Stdout, "[error] ", 0),
		StackAll:   false,
		StackSize:  1024 * 8,
		PrintStack: false,
	}
}

func NewRecovery(options *Options, next http.Handler) *Recovery {
	if options == nil {
		options = NewOptions()
	}
	return &Recovery{
		options: options,
		next:    next,
	}
}

func Middleware(options *Options) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return NewRecovery(options, next)
	}
}
