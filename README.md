# Scaleway Functions Go

Library authored By [Scaleway's Serverless Team](https://scaleway.com).

This repository contains a runtime wrapper necessary to develop with Golang on `Scaleway Functions` (Scaleway's Function As A Service Product). It allows users  to deploy Golang Function Handlers in the cloud by adding the transport/event formatting layer on top of the developer's codebase.

Basically, this project will stand as the Function's Gateway, starting an HTTP server to handle incoming traffic by transforming requests in Golang Event Structures and execute handlers defined by the developer with these structures.

**Disclaimer**: This library is heavily inspired by `aws lambda go` as we target compatibility with AWS Lambda Code (Allow users to develop for lambda and deploy on Scaleway Functions with minimal changes in their codebase). This way, you may already be familiar with API Gateway Event, Context and Response structures.

`Scaleway Functions Go` is intended to be used by developers, who want to use `Golang` runtime to run FAAS Applications on Scaleway Functions.

## Requirements

In order to use this project, you will need:
- Golang
- Package manager for go such as [dep](https://github.com/golang/dep).

## Install

In order to start development on your Serverless Application with Go, you have to install this project as a dependency, for example with `dep`:
```bash
dep ensure -add github.com/scaleway/scaleway-functions-go
```

**Please Note that you will have to package your vendors when uploading your codebase to Scaleway Functions, as we will take care of building your code**. This is why we use `dep` in this example.

## Getting Started

In order to run a Function in the cloud, you must specify a `handler` function in your code, which will execute your business logic.

Here is some sample code to run a basic function on Scaleway Functions.

You may take a look at [this directory for other examples](./example).

``` go
package main

import (
	"encoding/json"

	"github.com/scaleway/scaleway-functions-go/events"
	"github.com/scaleway/scaleway-functions-go/lambda"
)

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := map[string]interface{}{
		"message": "We're all good",
	}

	responseB, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(responseB),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
```

There are multiple important things to note inside above code snippet:
- **package main** is `required` (so we can execute the main function inside your package).
- import both `scaleway-functions-go/events` and `scaleway-functions-go/lambda` (library developed inside this repository).
- a `handler` function is defined, with `events.APIGatewayProxyRequest` as a parameter (contains informations about HTTP event triggered), and must return an `APIGatewayProxyResponse` structure, and an `error`.
- a `main` function, calling `lambda.Start` with the handler, described by the developer.

All the `wrapping` logic is contained inside `lambda.Start` and makes sure your code runs properly on our platform. As we will execute your `main package`, this main function will be executed to bootstrap the environment when a function instance is triggered.


**Please Note** that a function instance will be created when an event is triggered, but will stay up for a little bit of time. Thus, the same function instance may be executed for multiple events. 
This way, you may configure some `initialization logic` (for example, opening a connection to a Database), outside of the `main` function so it only gets executed once at startup.
Here is an example starting up a connection to a MySQL Database:

```go

package main

import (
	"fmt"
	"os"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/scaleway/scaleway-functions-go/events"
	"github.com/scaleway/scaleway-functions-go/lambda"
)

// global instance of database
var db *sql.DB

// Initialization logic, executed at function startup
func init() {
	var err error

	// Get Database configuration from environment variables
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	connectionString := fmt.Sprintf("%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatalf("error received while opening connection to Database: %v", err)
	}
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Do something with db instance
	//...

	return events.APIGatewayProxyResponse{
		Body:       "How to initialize stuff in a Scaleway Function",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
```

## Documentation and Useful links

As `Scaleway Functions` is in early access phase, developers invited to use our product will receive a link to the documentation of the platform.

You may use the [Serverless Framework](https://serverless.com) to deploy your Golang functions, with our plugin [for Scaleway Functions platform](https://github.com/scaleway/serverless-scaleway-functions).

## Contributing

As said above, we are only in `early access phase`, so this plugin is mainly developed and maintained by `Scaleway Serverless Team`. When the platform will reach a stable release, contributions via Pull Requests will be open.

Until then, you are free to open issues or discuss with us on our [Community Slack Channels](https://slack.online.net/).

## License

This project is MIT licensed.