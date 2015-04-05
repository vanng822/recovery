## Recovery

Recovery is a middleware with implementing func(next http.Handler) http.Handler

## Example

```go	
package main

import (
    "github.com/vanng822/r2router"
    "github.com/vanng822/recovery"
    "net/http"
)

func main() {
    seefor := r2router.NewSeeforRouter()
    rec := recovery.NewRecovery()
    rec.PrintStack = true
    seefor.Before(rec.Handler)
	seefor.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, _ r2router.Params) {
		panic("This shouldn't crash Seefor")
	})
	
    http.ListenAndServe(":8080", seefor)
}
```	
