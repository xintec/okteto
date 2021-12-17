// Copyright 2021 The Okteto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package namespace

import (
	"context"
	"fmt"

	contextCMD "github.com/okteto/okteto/cmd/context"
	"github.com/okteto/okteto/cmd/utils"
	"github.com/okteto/okteto/pkg/analytics"
	"github.com/okteto/okteto/pkg/errors"
	"github.com/okteto/okteto/pkg/log"
	"github.com/okteto/okteto/pkg/okteto"
	"github.com/spf13/cobra"
)

// Delete deletes a namespace
func Delete(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a namespace",
		RunE: func(cmd *cobra.Command, args []string) error {

			if err := contextCMD.NewContextCommand().Run(ctx, &contextCMD.ContextOptions{}); err != nil {
				return err
			}

			if !okteto.IsOkteto() {
				return errors.ErrContextIsNotOktetoCluster
			}

			nsCmd, err := newNamespaceCommand()
			if err != nil {
				return err
			}
			err = nsCmd.executeDeleteNamespace(ctx, args[0])
			analytics.TrackDeleteNamespace(err == nil)
			return err
		},
		Args: utils.ExactArgsAccepted(1, ""),
	}
	return cmd
}

func (nc *namespaceCommand) executeDeleteNamespace(ctx context.Context, namespace string) error {

	if err := nc.okClient.Namespaces().Delete(ctx, namespace); err != nil {
		return fmt.Errorf("failed to delete namespace: %s", err)
	}
	log.Success("Namespace '%s' deleted", namespace)

	if okteto.Context().Namespace == namespace {
		personalNamespace := okteto.Context().PersonalNamespace
		if personalNamespace == "" {
			personalNamespace = okteto.GetSanitizedUsername()
		}
		ctxOptions := &contextCMD.ContextOptions{
			Namespace: personalNamespace,
			Context:   okteto.Context().Name,
			Show:      true,
			Save:      true,
		}
		return nc.ctxCmd.Run(ctx, ctxOptions)
	}
	return nil
}
