package router

import (
	"database/sql"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	mux.Handle("/healthz", middleware.OSInfo(handler.NewHealthzHandler()))
	svc := service.NewTODOService(todoDB)
	mux.Handle("/todos", handler.NewTODOHandler(svc))
	mux.Handle("/do_panic", middleware.Recovery(handler.NewPanicHandler()))
	mux.Handle("/os_info", middleware.CheckAuth(middleware.OSInfo(middleware.Logging(handler.NewOSInfoHandler()))))
	mux.Handle("/heavy", middleware.OSInfo(middleware.Logging(handler.NewHeavyHandler())))

	return mux
}
