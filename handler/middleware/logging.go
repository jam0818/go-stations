package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/TechBowl-japan/go-stations/model"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		accessTime := time.Now()
		next.ServeHTTP(w, r)
		endTime := time.Now()
		latency := endTime.Sub(accessTime).Milliseconds()
		envInfo, ok := r.Context().Value(model.EnvinfoKey{}).(model.EnvInfo)
		if !ok {
			envInfo = model.EnvInfo{
				OS:      "",
				Browser: "",
			}
		}
		accessLog := struct {
			Timestamp time.Time
			Latency   int64
			Path      string
			OS        string
		}{
			Timestamp: accessTime,
			Latency:   latency,
			Path:      r.URL.String(),
			OS:        envInfo.OS,
		}
		logJSON, _ := json.Marshal(accessLog)
		fmt.Println(string(logJSON))
	}
	return http.HandlerFunc(fn)
}
