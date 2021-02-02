package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/vitorarins/todoer/api"
	"github.com/vitorarins/todoer/pb"
	"github.com/vitorarins/todoer/repository"
)

func main() {
	const timeout = 10 * time.Second

	var port int
	var grpcServer bool

	flag.IntVar(&port, "port", 8080, "port where the service will be listening to")
	flag.BoolVar(&grpcServer, "grpc", false, "run todoer service with grpc server")
	flag.Parse()

	repo := repository.NewLocalStorage()

	if grpcServer {
		grpcApi := api.NewGrpcApi(repo)

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		server := grpc.NewServer()
		pb.RegisterTodoerServer(server, grpcApi)
		log.Fatal(server.Serve(lis))
	} else {
		restApi := api.NewApi(repo)
		service := restApi.RegisterRoutes()

		server := &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      service,
			ReadTimeout:  timeout,
			WriteTimeout: timeout,
		}

		log.Fatal(server.ListenAndServe())
	}
}
