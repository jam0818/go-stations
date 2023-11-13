package handler

import (
	"log"
	"net/http"
	"time"
)

type HeavyHandler struct{}

func NewHeavyHandler() *HeavyHandler {
	return &HeavyHandler{}
}

func (h HeavyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("heavy process starts")
	time.Sleep(5 * time.Second)
	log.Println("done")
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("hello\n"))
	if err != nil {
		return
	}
}
