package main

import (
	"censsrv/pkg/api"
	"log"
	"net/http"
)

type server struct {
	api *api.API
}

func main() {
	srv := server{}

	srv.api = api.New()
	log.Fatal(http.ListenAndServe("localhost:8083", srv.api.Router()))
}
