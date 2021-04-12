package context

// import (
// 	"bytes"
// 	"html/template"
// 	"net/http"
// 	"time"

// 	"code.gitea.io/gitea/modules/log"
// 	"code.gitea.io/gitea/modules/setting"
// )

// type routerLoggerOptions struct {
// 	req            *http.Request
// 	Identity       *string
// 	Start          *time.Time
// 	ResponseWriter http.ResponseWriter
// 	Ctx            map[string]interface{}
// }

// // AccessLogger returns a middleware to log access logger
// func AccessLogger() func(http.Handler) http.Handler {
// 	logger := log.GetLogger("access")
// 	logTemplate, _ := template.New("log").Parse(setting.AccessLogTemplate)
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 			start := time.Now()
// 			next.ServeHTTP(w, req)
// 			identity := "-"
// 			if val := SignedUserName(req); val != "" {
// 				identity = val
// 			}
// 			rw := w.(ResponseWriter)

// 			buf := bytes.NewBuffer([]byte{})
// 			err := logTemplate.Execute(buf, routerLoggerOptions{
// 				req:            req,
// 				Identity:       &identity,
// 				Start:          &start,
// 				ResponseWriter: rw,
// 				Ctx: map[string]interface{}{
// 					"RemoteAddr": req.RemoteAddr,
// 					"Req":        req,
// 				},
// 			})
// 			if err != nil {
// 				log.Error("Could not set up chi access logger: %v", err.Error())
// 			}

// 			err = logger.SendLog(log.INFO, "", "", 0, buf.String(), "")
// 			if err != nil {
// 				log.Error("Could not set up chi access logger: %v", err.Error())
// 			}
// 		})
// 	}
// }
