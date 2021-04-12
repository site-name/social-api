package routes

// import (
// 	"net/http"
// 	"time"

// 	"github.com/sitename/sitename/modules/log"
// )

// // LoggerHandler is a handler that will log the routing to the default gitea log
// func LoggerHandler(level log.Level) func(next http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			start := time.Now()
// 			_ = log.GetLogger("router").Log(0, level, "Started %s %s for %s", log.ColoredMethod(r.Method), r.URL.RequestURI(), r.RemoteAddr)
// 			next.ServeHTTP(w, r)
// 			var status int

// 		})
// 	}
// }
