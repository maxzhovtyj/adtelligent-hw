package app

import (
	"github.com/maxzhovtyj/adtelligent-hw/internal/delivery/http"
	"github.com/maxzhovtyj/adtelligent-hw/internal/services"
	"github.com/maxzhovtyj/adtelligent-hw/internal/storage"
	"github.com/maxzhovtyj/adtelligent-hw/pkg/db/postgres"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func Run() {
	conn, err := postgres.NewConn()
	if err != nil {
		panic(err)
	}

	go func() {
		if err = http.ListenAndServe(":8888", nil); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("initialized db connection")

	appStorage := storage.New(conn)
	appServices := services.New(appStorage)
	appHandler := delivery.New(appServices)

	log.Println("start listening http server")
	if err = http.ListenAndServe(":9999", appHandler.Init()); err != nil {
		log.Fatal(err)
	}
}
