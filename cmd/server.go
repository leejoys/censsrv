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

	// Запускаем веб-сервер на порту 8083 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов.
	log.Println("HTTP server is started on localhost:8083")
	defer log.Println("HTTP server has been stopped")
	log.Fatal(http.ListenAndServe("localhost:8083", srv.api.Router()))
}
