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
	"errors"
	"flag"
	"fmt"
	"path/filepath"

	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	horizoncomponents "github.com/blackducksoftware/horizon/pkg/components"
	"github.com/blackducksoftware/horizon/pkg/deployer"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

var secretType horizonapi.SecretType

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Deploys the synopsys operator onto your cluster",
	Args: func(cmd *cobra.Command, args []string) error {
		// Check the Secret Type
		switch init_secretType {
		case "Opaque":
			secretType = horizonapi.SecretTypeOpaque
		case "ServiceAccountToken":
			secretType = horizonapi.SecretTypeServiceAccountToken
		case "Dockercfg":
			secretType = horizonapi.SecretTypeDockercfg
		case "DockerConfigJSON":
			secretType = horizonapi.SecretTypeDockerConfigJSON
		case "BasicAuth":
			secretType = horizonapi.SecretTypeBasicAuth
		case "SSHAuth":
			secretType = horizonapi.SecretTypeSSHAuth
		case "TypeTLS":
			secretType = horizonapi.SecretTypeTLS
		default:
			fmt.Printf("Invalid Secret Type: %s\n", init_secretType)
			return errors.New("Bad Secret Type")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("at this point we would call kube/install.sh -i %s -p %s -k %s -d %s\n", init_synopsysOperatorImage, init_promethiusImage, init_blackduckRegistrationKey, init_dockerConfigPath)

		// check if operator is already installed
		out, err := RunKubeCmd("get", "clusterrolebindings", "synopsys-operator-admin", "-o", "go-template='{{range .subjects}}{{.namespace}}{{end}}'")
		if err == nil {
			fmt.Printf("You have already installed the operator in namespace %s.\n", out)
			fmt.Printf("To delete the operator run: synopsysctl stop --namespace %s\n", out)
			fmt.Printf("Nothing to do...\n")
			return
		}

		// Start Horizon
		var kubeconfig *string
		if home := homeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// Use the current context in kubeconfig
		rc, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}

		// Create Horizon Deployer with current kube context
		horizonDeployer, err := deployer.NewDeployer(rc)

		// create a new namespace
		// out, err = RunKubeCmd("create", "namespace", namespace)
		ns := horizoncomponents.NewNamespace(horizonapi.NamespaceConfig{
			// APIVersion:  "string",
			// ClusterName: "string",
			Name:      namespace,
			Namespace: namespace,
		})

		horizonDeployer.AddNamespace(ns)

		// create a secret
		secret := horizoncomponents.NewSecret(horizonapi.SecretConfig{
			APIVersion: "v1",
			// ClusterName : "cluster",
			Name:      init_secretName,
			Namespace: namespace,
			Type:      secretType,
		})

		secret.AddData(map[string][]byte{
			"ADMIN_PASSWORD":    []byte(init_secretAdminPassword),
			"POSTGRES_PASSWORD": []byte(init_secretPostgresPassword),
			"USER_PASSWORD":     []byte(init_secretUserPassword),
			"HUB_PASSWORD":      []byte(init_secretBlackduckPassword),
		})

		horizonDeployer.AddSecret(secret)

		// Deploy Kubernetes Resources from Horizon
		err = horizonDeployer.Run()
		if err != nil {
			fmt.Printf("Error deploying Kubernetes Resources with Horizon : %s\n", err)
			return
		}

		// create -f the operator

		// secret link stuff

		// expose the routes

	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&init_synopsysOperatorImage, "synopsys-operator-image", "i", init_synopsysOperatorImage, "synopsys operator image URL")
	initCmd.Flags().StringVarP(&init_promethiusImage, "promethius-image", "p", init_promethiusImage, "promethius image URL")
	initCmd.Flags().StringVarP(&init_blackduckRegistrationKey, "blackduck-registration-key", "k", init_blackduckRegistrationKey, "key to register with KnowledgeBase")
	initCmd.Flags().StringVarP(&init_dockerConfigPath, "docker-config", "d", init_dockerConfigPath, "path to docker config (image pull secrets etc)")

	initCmd.Flags().StringVar(&init_secretName, "secret-name", init_secretName, "name of kubernetes secret for postgres and blackduck")
	initCmd.Flags().StringVar(&init_secretType, "secret-type", init_secretType, "type of kubernetes secret for postgres and blackduck")
	initCmd.Flags().StringVar(&init_secretAdminPassword, "admin-password", init_secretAdminPassword, "postgres admin password")
	initCmd.Flags().StringVar(&init_secretPostgresPassword, "postgres-password", init_secretPostgresPassword, "postgres password")
	initCmd.Flags().StringVar(&init_secretUserPassword, "user-password", init_secretUserPassword, "postgres user password")
	initCmd.Flags().StringVar(&init_secretBlackduckPassword, "blackduck-password", init_secretBlackduckPassword, "blackduck password for 'sysadmin' account")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
