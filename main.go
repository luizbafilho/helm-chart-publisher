package main

import (
	"log"
	"os"

	"github.com/luizbafilho/chart-server/api"
	"github.com/luizbafilho/chart-server/publisher"
	_ "github.com/luizbafilho/chart-server/storage/s3"
)

func main() {
	publisher, err := publisher.New()
	if err != nil {
		log.Fatal(err)
	}

	a := api.New(publisher)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	a.Serve(":" + port)
}
