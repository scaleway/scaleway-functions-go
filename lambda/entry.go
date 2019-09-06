package lambda

import (
	"fmt"
	"log"
	"os"

	"github.com/scaleway/scaleway-functions-go/events"
)

const defaultPort = 8080

// FunctionHandler - Handler for Event
type FunctionHandler func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

// Start - Start the process
func Start(handler interface{}) {
	wrappedHandler := NewHandler(handler)
	StartHandler(wrappedHandler)
}

// StartHandler - Execute Handler
func StartHandler(handler Handler) {
	// 1: Parse arguments
	event := os.Args[1]
	// context := os.Args[4]

	// 2:
	response, err := handler.Invoke(nil, []byte(event))
	if err != nil {
		errorMessage := fmt.Sprintf("[SCW_ERROR] %s", err.Error())
		log.Print(errorMessage)
		return
	}

	successMessage := fmt.Sprintf("[SCW_END] %s", string(response))
	fmt.Print(successMessage)
}
