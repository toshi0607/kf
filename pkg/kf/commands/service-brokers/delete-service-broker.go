// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package servicebrokers

import (
	"context"
	"fmt"

	servicecatalogclient "github.com/google/kf/pkg/client/servicecatalog/clientset/versioned"
	"github.com/google/kf/pkg/kf/commands/config"
	installutil "github.com/google/kf/pkg/kf/commands/install/util"
	"github.com/google/kf/pkg/kf/commands/utils"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewDeleteServiceBrokerCommand adds a cluster service broker to the service catalog.
// TODO (juliaguo): Add flag to allow namespaced service broker
func NewDeleteServiceBrokerCommand(p *config.KfParams, client servicecatalogclient.Interface) *cobra.Command {
	var (
		spaceScoped bool
		force       bool
	)
	deleteCmd := &cobra.Command{
		Use:     "delete-service-broker BROKER_NAME",
		Aliases: []string{"dsb"},
		Short:   "Remove a cluster service broker from service catalog",
		Example: `  kf delete-service-broker mybroker`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceBrokerName := args[0]

			cmd.SilenceUsage = true

			if err := utils.ValidateNamespace(p); err != nil {
				return err
			}

			var toDelete func() error
			if !spaceScoped {
				toDelete = func() error {
					err := client.ServicecatalogV1beta1().ClusterServiceBrokers().Delete(serviceBrokerName, &metav1.DeleteOptions{})
					return err
				}
			} else {
				toDelete = func() error {
					err := client.ServicecatalogV1beta1().ServiceBrokers(p.Namespace).Delete(serviceBrokerName, &metav1.DeleteOptions{})
					return err
				}
			}

			shouldDelete := true
			if !force {
				var err error
				shouldDelete, err = installutil.SelectYesNo(context.Background(), fmt.Sprintf("Really delete service-broker %s?", serviceBrokerName))
				if err != nil {
					return err
				}
			}
			if shouldDelete {
				return toDelete()
			} else {
				return nil
			}
		},
	}

	deleteCmd.Flags().BoolVar(
		&spaceScoped,
		"space-scoped",
		false,
		"Set to delete a space scoped service broker.")

	deleteCmd.Flags().BoolVar(
		&force,
		"force",
		false,
		"Set to force deletion without a confirmation prompt.")

	return deleteCmd
}
