package main

import (
	"github.com/luizbafilho/helm-chart-publisher/cmd"
	_ "github.com/luizbafilho/helm-chart-publisher/storage/s3"
)

func main() {
	cmd.Execute()
}
