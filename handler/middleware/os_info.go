package middleware

import (
	"context"
	"github.com/TechBowl-japan/go-stations/model"
	"github.com/mssola/user_agent"
	"net/http"
)

func getenvFromUserAgent(userAgent string) model.EnvInfo {
	ua := user_agent.New(userAgent)
	osName := ua.OS()
	browserName, _ := ua.Browser()
	envInfo := model.EnvInfo{
		OS:      osName,
		Browser: browserName,
	}

	return envInfo
}

func OSInfo(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.UserAgent()
		envInfo := getenvFromUserAgent(userAgent)

		ctx := context.WithValue(r.Context(), model.EnvinfoKey{}, envInfo)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
