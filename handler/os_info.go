package handler

import (
	"encoding/json"
	"github.com/TechBowl-japan/go-stations/xcontext"
	"log"
	"net/http"
)

type OSInfoHandler struct {
}

func NewOSInfoHandler() *OSInfoHandler {
	return &OSInfoHandler{}
}

func (h *OSInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	envInfo := xcontext.GetOSInfo(r.Context())
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err := encoder.Encode(envInfo)
	if err != nil {
		log.Println(err)
	}
}
