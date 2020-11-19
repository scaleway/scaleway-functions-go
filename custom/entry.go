package custom

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const defaultPort = 8081

// Request structure sent from core runtime
type runtimeRequest struct {
	Event   interface{}
	Context interface{}
}

// Start - Start the process
func Start(handler interface{}) {
	wrappedHandler := NewHandler(handler)
	StartHandler(wrappedHandler)
}

// StartHandler - Execute a Function handler
func StartHandler(handler Handler) {
	portEnv := os.Getenv("SCW_UPSTREAM_PORT")
	port, err := strconv.Atoi(portEnv)
	if err != nil {
		port = defaultPort
	}

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		ReadTimeout:    3 * time.Second,
		WriteTimeout:   3 * time.Second,
		MaxHeaderBytes: 1 << 20, // Max header of 1MB
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if err := handler.Invoke(writer, request); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte(err.Error()))
		}
	})

	log.Fatal(s.ListenAndServe())
}
