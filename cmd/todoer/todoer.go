package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vitorarins/todoer/api"
	"github.com/vitorarins/todoer/repository"
)

func main() {
	const timeout = 10 * time.Second

	var port int

	flag.IntVar(&port, "port", 8080, "port where the service will be listening to")
	flag.Parse()

	repo := repository.NewLocalStorage()
	restApi := api.NewApi(repo)
	service := restApi.RegisterRoutes()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      service,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	log.Infof("running todoer service, listening on port %d", port)
	log.Fatal(server.ListenAndServe())
}
