package logger

import "net/http"

type responseRecorder struct {
	http.ResponseWriter
	status       int
	bytesWritten int
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.bytesWritten += n
	return n, err
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
