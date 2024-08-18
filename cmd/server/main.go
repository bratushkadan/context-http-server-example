package main

import (
	"log"
	"net/http"
	"time"

	"github.com/bratushkadan/context-http-server-example/internal"
)

func main() {
	requestIdMiddleware := internal.CreateRequestIdMiddleware()
	authMiddleware := internal.CreateAuthMiddleware()
	timeoutMiddleware := internal.CreateTimeoutMiddleware(1 * time.Second)

	http.HandleFunc("/hello-world", internal.HelloWorldHandler)
	http.Handle("/bar", requestIdMiddleware(http.HandlerFunc(internal.BarHandler)))
	http.Handle("/uuid", requestIdMiddleware(http.HandlerFunc(internal.BarHandler)))

	// curl localhost:8079/private-endpoint                                                         | no "X-Auth-Token" provided : authenticated
	// curl -H "X-Auth-Token: foo"  localhost:8079/private-endpoint                                 | invalid "X-Auth-Token" provided : authenticated
	// curl -H "X-Auth-Token: 84214dc4-73d8-4b4a-af62-c919252ab8b3" localhost:8079/private-endpoint |
	http.Handle("/private-endpoint", requestIdMiddleware(authMiddleware(http.HandlerFunc(internal.PrivateHandler))))

	http.Handle("/slow", requestIdMiddleware(timeoutMiddleware(http.HandlerFunc(internal.SlowHandler))))

	log.Fatal(http.ListenAndServe("0.0.0.0:8079", nil))
}
