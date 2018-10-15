package main

import (
	"github.com/HotelsDotCom/helm-chart-publisher/cmd"
	_ "github.com/HotelsDotCom/helm-chart-publisher/storage/gcs"
	_ "github.com/HotelsDotCom/helm-chart-publisher/storage/s3"
	_ "github.com/HotelsDotCom/helm-chart-publisher/storage/swift"
)

func main() {
	cmd.Execute()
}
