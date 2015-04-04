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
	options := recovery.NewRecoveryOptions()
	//options.Logger = CustomizeLogger() customized logger
	//options.PrintStack = true printing stacktrace
	seefor.Before(recovery.NewRecovery(options))
	
	http.ListenAndServe(":8080", seefor)
}
```	
