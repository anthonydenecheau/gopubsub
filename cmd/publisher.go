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
	"os"

	logConfig "github.com/anthonydenecheau/gopubsub/common/config"
	pubConfig "github.com/anthonydenecheau/gopubsub/common/config"
	breederRepository "github.com/anthonydenecheau/gopubsub/common/repository"
	dogRepository "github.com/anthonydenecheau/gopubsub/common/repository"
	ownerRepository "github.com/anthonydenecheau/gopubsub/common/repository"
	parentRepository "github.com/anthonydenecheau/gopubsub/common/repository"
	pedigreeRepository "github.com/anthonydenecheau/gopubsub/common/repository"
	titleRepository "github.com/anthonydenecheau/gopubsub/common/repository"
	breederService "github.com/anthonydenecheau/gopubsub/common/service"
	constantService "github.com/anthonydenecheau/gopubsub/common/service"
	dogService "github.com/anthonydenecheau/gopubsub/common/service"
	ownerService "github.com/anthonydenecheau/gopubsub/common/service"
	parentService "github.com/anthonydenecheau/gopubsub/common/service"
	pedigreeService "github.com/anthonydenecheau/gopubsub/common/service"
	personService "github.com/anthonydenecheau/gopubsub/common/service"
	titleService "github.com/anthonydenecheau/gopubsub/common/service"
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

		os.Setenv("PUBSUB_TOPIC", configuration.PubSub.Topic)
		os.Setenv("GOOGLE_CLOUD_PROJECT", configuration.PubSub.GoogleCloudProjectId)

		var err error
		log, file := logConfig.NewLogger(configuration.Logger, configuration.App)
		if file != nil {
			defer file.Close()
		}

		log.Info("Start Application!")
		log.Infof("Publisher %s", configuration.Publisher.Example)
		log.Infof("Publisher username is %s", configuration.Publisher.Database.Username)
		log.Infof("Publisher password is %s", configuration.Publisher.Database.Password)
		log.Infof("Publisher SID is %s", configuration.Publisher.Database.SID)

		log.Infof("Pubsub topic is %s", configuration.PubSub.Topic)
		log.Infof("PubSub project id is %s", configuration.PubSub.GoogleCloudProjectId)

		P := goracle.ConnectionParams{
			Username:    configuration.Publisher.Database.Username,
			Password:    configuration.Publisher.Database.Password,
			SID:         configuration.Publisher.Database.SID,
			MinSessions: 1, MaxSessions: maxSessions, PoolIncrement: 1,
			ConnClass:    "POOLED",
			EnableEvents: true,
		}

		testConStr = P.StringWithPassword()
		if dbOracle, err = sql.Open("goracle", testConStr); err != nil {
			log.Fatalf("Oracle Connection FAILED : %s", err.Error())
			return
		}
		defer dbOracle.Close()

		if dbOracle != nil {
			if clientVersion, err = goracle.ClientVersion(dbOracle); err != nil {
				log.Fatalf("clientVersion FAILED : %s", err.Error())
				return
			}
			if serverVersion, err = goracle.ServerVersion(dbOracle); err != nil {
				log.Fatalf("serverVersion FAILED : %s", err.Error())
				return
			}
			log.Info("Server :", serverVersion)
			log.Info("Client : ", clientVersion)
		}

		// dog repository, service
		dr := dogRepository.NewOraDogRepository(dbOracle)
		ds := dogService.NewDogService(dr)

		// breeder repository, service
		br := breederRepository.NewOraBreederRepository(dbOracle)
		bs := breederService.NewBreederService(br)

		// owner repository, service
		or := ownerRepository.NewOraOwnerRepository(dbOracle)
		os := ownerService.NewOwnerService(or)

		// parent repository, service
		pr := parentRepository.NewOraParentRepository(dbOracle)
		ps := parentService.NewParentService(pr)

		// pedigree repository, service
		lr := pedigreeRepository.NewOraPedigreeRepository(dbOracle)
		ls := pedigreeService.NewPedigreeService(lr)

		// titre francais et etranger repository, service
		tr := titleRepository.NewOraTitleRepository(dbOracle)
		tsf := titleService.NewTitleService(tr, constantService.TitleDomaineFr)
		tse := titleService.NewTitleService(tr, constantService.TitleDomaineEtr)

		// personne repository, service
		ns := personService.NewPersonService(br, or)

		// publisher
		pubSubGateway := pubConfig.NewPublisher()
		// initialize all
		pubTask.NewTask(
			dbOracle,
			ds,
			bs,
			os,
			ps,
			ls,
			tsf,
			tse,
			ns,
			pubSubGateway,
			log)

		log.Info("End Application!")
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
