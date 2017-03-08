// Copyright © 2017 Luiz Filho
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/luizbafilho/helm-chart-publisher/api"
	"github.com/luizbafilho/helm-chart-publisher/config"
	"github.com/luizbafilho/helm-chart-publisher/publisher"
	"github.com/spf13/cobra"
)

var (
	configFile string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "helm-chart-publisher",
	Short: "publishes helm chart",
	Long:  `helm-chart-publisher publishes helm chart to a configured storage`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.ReadConfigFile(configFile); err != nil {
			log.Fatal(err)
		}

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
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")
}
