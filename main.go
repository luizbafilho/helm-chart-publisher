package main

import (
	"log"

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
	a.Serve(":8080")
}
