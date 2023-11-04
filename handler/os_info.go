package handler

import (
	"encoding/json"
	"github.com/TechBowl-japan/go-stations/model"
	"log"
	"net/http"
)

type OSInfoHandler struct {
}

func NewOSInfoHandler() *OSInfoHandler {
	return &OSInfoHandler{}
}

func (h *OSInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	envInfo, ok := r.Context().Value(model.EnvinfoKey{}).(model.EnvInfo)
	if !ok {
		envInfo = model.EnvInfo{
			OS:      "",
			Browser: "",
		}
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err := encoder.Encode(envInfo)
	if err != nil {
		log.Println(err)
	}
}
