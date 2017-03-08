package main

import (
	"log"
	"os"

	"github.com/luizbafilho/helm-chart-publisher/api"
	"github.com/luizbafilho/helm-chart-publisher/publisher"
	_ "github.com/luizbafilho/helm-chart-publisher/storage/s3"
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
