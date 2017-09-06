package main

import (
	"github.com/luizbafilho/helm-chart-publisher/cmd"
	_ "github.com/luizbafilho/helm-chart-publisher/storage/gcs"
	_ "github.com/luizbafilho/helm-chart-publisher/storage/s3"
	_ "github.com/luizbafilho/helm-chart-publisher/storage/swift"
)

func main() {
	cmd.Execute()
}
