/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	projectID string = ""

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "pubsub-poller [FLAGS] SUBSCRIPTION...",
		Args:  cobra.MinimumNArgs(1),
		Short: "Polls gcloud pubsub subscriptions and outputs messages to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			if gcpProject, ok := os.LookupEnv("GCP_PROJECT"); ok {
				projectID = gcpProject
			}

			ctx := context.Background()
			client, err := pubsub.NewClient(ctx, projectID)
			if err != nil {
				fmt.Printf("pubsub.NewClient: %v", err)
				return
			}

			subscriptions := args[1:]
			var wg sync.WaitGroup

			for _, sub := range subscriptions {
				wg.Add(1)
				go pullMsgs(ctx, client, "seng-6285-spring21", sub, &wg)
			}

			wg.Wait()

		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&projectID, "project-id", "p", "", "Project ID")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".pubsub-poller" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".pubsub-poller")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func pullMsgs(ctx context.Context, client *pubsub.Client, projectID, subID string, wg *sync.WaitGroup) error {
	var mu sync.Mutex
	// received := 0
	sub := client.Subscription(subID)
	flog := log.WithFields(log.Fields{
		"project-id":   projectID,
		"subscription": subID,
	})
	flog.Debug("Starting poller")
	cctx := context.WithValue(ctx, "subscriptionID", subID)
	err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		flog.WithField("message-id", msg.ID).Info(string(msg.Data))
		msg.Ack()
		// received++
		// if received == 1000 {
		// 	// fmt.Printf("Received 1000 messages on %s subscription.", subID)
		// 	wg.Done()
		// 	cancel()
		// }
	})
	if err != nil {
		return fmt.Errorf("Receive: %v", err)
	}
	return nil
}
