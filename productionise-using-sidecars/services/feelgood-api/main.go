package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	log := newLogger()
	router := newRouter(log)
	addr := fmt.Sprintf(":%d", 8080)

	srv := http.Server{
		Addr:    addr,
		Handler: router,
	}

	log.Info("starting server...")
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	log.Info("server has started...")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		log.Warn("received SIGINT...")
	case syscall.SIGTERM:
		log.Warn("received SIGTERM...")
	}

	log.Info("server is shutting down...")
	srv.Shutdown(context.Background())
	log.Info("server shutdown complete...")
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "not found"}`))
}

func goodvibes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("content-type", "application/json")

	randomSource := rand.NewSource(time.Now().UnixNano())
	randomNum := rand.New(randomSource).Intn(100)

	if randomNum%3 == 0 {
		fail(w, randomNum)
		return
	}

	ok(w)
}

func fail(w http.ResponseWriter, num int) {
	body := fmt.Sprintf(`{"message": "random number %d divisible by 3"}`, num)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(body))
}

func ok(w http.ResponseWriter) {
	body, err := getMessage()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.Write([]byte(body))
}

func getMessage() ([]byte, error) {
	rand.Seed(time.Now().Unix())
	msg := messages[rand.Intn(len(messages))]
	return json.Marshal(msg)
}

func newRouter(logger *logrus.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", logMiddleware(logger, notFound))
	mux.Handle("/goodvibes", logMiddleware(logger, goodvibes))

	return mux
}

func newLogger() *logrus.Logger {
	logger := logrus.New()
	file, err := os.OpenFile("/var/log/feelgood/server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.Out = file
	} else {
		logger.Warn("Failed to log to file, using default stderr")
	}

	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	return logger
}

func logMiddleware(logger *logrus.Logger, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wi := &responseRecorder{
			statusCode:     http.StatusOK,
			ResponseWriter: w,
		}

		defer func() {
			path := r.URL.RequestURI()
			method := r.Method
			proto := r.Proto

			logger.WithFields(
				logrus.Fields{
					"path":     path,
					"method":   method,
					"protocol": proto,
					"status":   wi.statusCode,
					"response": wi.body,
				},
			).Info()
		}()

		next(wi, r)
	})
}

// statusCodeRecorder implements WriteHeader in order to intercept status codes
// from the response. The WriteHeader may not be called if the request was
// successful so it should be initialised to http.StatusOK
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       string
}

// WriteHeader pull the status code from the response to the interceptor
func (rr *responseRecorder) WriteHeader(statusCode int) {
	rr.statusCode = statusCode
	rr.ResponseWriter.WriteHeader(statusCode)
}

func (rr *responseRecorder) Write(body []byte) (int, error) {
	rr.body = string(body)
	return rr.ResponseWriter.Write(body)
}
