// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"time"

	subConfig "github.com/anthonydenecheau/gopubsub/common/config"
	dogRepository "github.com/anthonydenecheau/gopubsub/common/repository"
	dogService "github.com/anthonydenecheau/gopubsub/common/service"
	subTask "github.com/anthonydenecheau/gopubsub/common/task"

	"github.com/go-pg/pg"
	"github.com/spf13/cobra"
)

// subscriberCmd represents the subscriber command
var subscriberCmd = &cobra.Command{
	Use:   "subscriber",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("subscriber called")
		log.Printf("Subscriber %s", configuration.Subscriber.Example)
		log.Printf("Subscriber database name is %s", configuration.Subscriber.Database.Name)

		os.Setenv("PUBSUB_TOPIC", configuration.PubSub.Topic)
		os.Setenv("GOOGLE_CLOUD_PROJECT", configuration.PubSub.GoogleCloudProjectId)
		os.Setenv("PUBSUB_SUBSRCIPTION", configuration.PubSub.Subscription)

		dbPg := pg.Connect(&pg.Options{
			User:                  configuration.Subscriber.Database.Username,
			Password:              configuration.Subscriber.Database.Password,
			Database:              configuration.Subscriber.Database.Name,
			Addr:                  configuration.Subscriber.Database.Host + ":" + configuration.Subscriber.Database.Port,
			RetryStatementTimeout: true,
			MaxRetries:            4,
			MinRetryBackoff:       250 * time.Millisecond,
		})

		var n int
		_, err := dbPg.QueryOne(pg.Scan(&n), "SELECT 1")
		if err != nil {
			fmt.Printf("ERROR: %+v\n", err)
			return
		}
		defer dbPg.Close()

		// dog repository
		dr := dogRepository.NewPgDogRepository(dbPg)
		// dog service
		ds := dogService.NewDogService(dr)
		// publisher
		pubSubGateway := subConfig.NewSubscriber()
		// initialize all
		subTask.NewSubTask(dbPg, ds, pubSubGateway)

	},
}

func init() {
	rootCmd.AddCommand(subscriberCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// subscriberCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// subscriberCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
