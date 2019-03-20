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
	"database/sql"
	"fmt"
	"log"
	"os"

	pubConfig "github.com/anthonydenecheau/gopubsub/common/pubsub"
	dogRepository "github.com/anthonydenecheau/gopubsub/common/repository"
	dogService "github.com/anthonydenecheau/gopubsub/common/service"
	pubTask "github.com/anthonydenecheau/gopubsub/common/task"

	"github.com/spf13/cobra"
	"gopkg.in/goracle.v2"
)

var (
	dbOracle                     *sql.DB
	clientVersion, serverVersion goracle.VersionInfo
	testConStr                   string
)

const maxSessions = 64

// publisherCmd represents the publisher command
var publisherCmd = &cobra.Command{
	Use:   "publisher",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("publisher called")
		log.Printf("Publisher %s", configuration.Publisher.Example)
		log.Printf("Publisher username is %s", configuration.Publisher.Database.Username)
		log.Printf("Publisher password is %s", configuration.Publisher.Database.Password)
		log.Printf("Publisher SID is %s", configuration.Publisher.Database.SID)

		log.Printf("Pubsub topic is %s", configuration.PubSub.Topic)
		log.Printf("PubSub project id is %s", configuration.PubSub.GoogleCloudProjectId)

		os.Setenv("PUBSUB_TOPIC", configuration.PubSub.Topic)
		os.Setenv("GOOGLE_CLOUD_PROJECT", configuration.PubSub.GoogleCloudProjectId)

		P := goracle.ConnectionParams{
			Username:    configuration.Publisher.Database.Username,
			Password:    configuration.Publisher.Database.Password,
			SID:         configuration.Publisher.Database.SID,
			MinSessions: 1, MaxSessions: maxSessions, PoolIncrement: 1,
			ConnClass:    "POOLED",
			EnableEvents: true,
		}

		testConStr = P.StringWithPassword()
		var err error
		if dbOracle, err = sql.Open("goracle", testConStr); err != nil {
			fmt.Printf("ERROR: %+v\n", err)
			return
			//panic(err)
		}
		defer dbOracle.Close()

		if dbOracle != nil {
			if clientVersion, err = goracle.ClientVersion(dbOracle); err != nil {
				fmt.Printf("ERROR: %+v\n", err)
				return
			}
			if serverVersion, err = goracle.ServerVersion(dbOracle); err != nil {
				fmt.Printf("ERROR: %+v\n", err)
				return
			}
			fmt.Println("Server:", serverVersion)
			fmt.Println("Client:", clientVersion)
		}

		// dog repository
		dr := dogRepository.NewOraDogRepository(dbOracle)
		// dog service
		ds := dogService.NewDogService(dr)
		// publisher
		pubSubGateway := pubConfig.NewPublisher()
		// initialize all
		pubTask.NewTask(dbOracle, ds, pubSubGateway)

	},
}

func init() {
	rootCmd.AddCommand(publisherCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// publisherCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// publisherCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
