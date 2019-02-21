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

package synopsysctl

import (
	"fmt"

	alertclientset "github.com/blackducksoftware/synopsys-operator/pkg/alert/client/clientset/versioned"
	blackduckclientset "github.com/blackducksoftware/synopsys-operator/pkg/blackduck/client/clientset/versioned"
	opssightclientset "github.com/blackducksoftware/synopsys-operator/pkg/opssight/client/clientset/versioned"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Remove a Synopsys Resource from your cluster",
	Args: func(cmd *cobra.Command, args []string) error {
		numArgs := 1
		if len(args) < numArgs {
			return fmt.Errorf("Not enough arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("Deleting a Non-Synopsys Resource")
		kubeCmdArgs := append([]string{"delete"}, args...)
		out, err := RunKubeCmd(kubeCmdArgs...)
		if err != nil {
			fmt.Printf("Error Deleting the Resource: %s", out)
		} else {
			fmt.Printf("%+v", out)
		}
	},
}

var deleteBlackduckCmd = &cobra.Command{
	Use:   "blackduck",
	Short: "Delete a Blackduck",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("Deleting a Blackduck\n")
		// Read Commandline Parameters
		blackduckNamespace := args[0]

		// Get Kubernetes Rest Config
		restconfig := getKubeRestConfig()
		blackduckClient, err := blackduckclientset.NewForConfig(restconfig)
		if err != nil {
			fmt.Printf("Error creating the Blackduck Clientset: %s\n", err)
			return
		}
		blackduckClient.SynopsysV1().Blackducks(blackduckNamespace).Delete(blackduckNamespace, &metav1.DeleteOptions{})

	},
}

var deleteOpsSightCmd = &cobra.Command{
	Use:   "opssight",
	Short: "Delete an OpsSight",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("Deleting an OpsSight\n")
		// Read Commandline Parameters
		opsSightNamespace := args[0]

		// Get Kubernetes Rest Config
		restconfig := getKubeRestConfig()
		opssightClient, err := opssightclientset.NewForConfig(restconfig)
		if err != nil {
			fmt.Printf("Error creating the OpsSight Clientset: %s\n", err)
			return
		}
		opssightClient.SynopsysV1().OpsSights(opsSightNamespace).Delete(opsSightNamespace, &metav1.DeleteOptions{})
	},
}

var deleteAlertCmd = &cobra.Command{
	Use:   "alert",
	Short: "Delete an Alert",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("Deleting an Alert\n")
		// Read Commandline Parameters
		alertNamespace := args[0]

		// Get Kubernetes Rest Config
		restconfig := getKubeRestConfig()
		alertClient, err := alertclientset.NewForConfig(restconfig)
		if err != nil {
			fmt.Printf("Error creating the Alert Clientset: %s\n", err)
			return
		}
		alertClient.SynopsysV1().Alerts(alertNamespace).Delete(alertNamespace, &metav1.DeleteOptions{})
	},
}

func init() {
	deleteCmd.DisableFlagParsing = true
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.AddCommand(deleteBlackduckCmd)
	deleteCmd.AddCommand(deleteOpsSightCmd)
	deleteCmd.AddCommand(deleteAlertCmd)
}
