package api

import (
	"fmt"
	"net/http"
)

type SseWriterInterceptor struct {
	http.ResponseWriter
}

func (w *SseWriterInterceptor) Write(data []byte) (int, error) {
	formatted := fmt.Sprintf("data: %s\n\n", data)
	return w.ResponseWriter.Write([]byte(formatted))
}
