package api

import (
	"fmt"
	"net/http"
)

func (s *Server) LoadIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	}
}

func (s *Server) StartSSE() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		}

		sseW := &SseWriterInterceptor{w}

		for {
			select {
			case <-r.Context().Done():
				return
			case track := <-s.service.Tracks:
				s.log.Info("New track", "track", track)
				if _, err := fmt.Fprintf(sseW, track); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					break
				}
				flusher.Flush()
			}
		}
	}
}
