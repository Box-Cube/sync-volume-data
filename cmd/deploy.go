/*
Copyright © 2021 Box-Cube

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
	"errors"
	"github.com/sirupsen/logrus"
	//cmd2 "sync-volume-data/cmd"
	"sync-volume-data/server"

	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
func newDeployCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "deploy",
		Short: "transfer data from/to Deployment kind resource",
		Long: `transfer data from/to Deployment kind resource, you need to specific a deploy name.
 For example:
	
	sync-volume-data rsync to deploy nginx -n my-web -v web -u root -p "myPassword" -s=test.file
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("you need specific a deploy name")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger := newLogger().WithFields(logrus.Fields{
				"namespace": *namespace,
				"kind":      cmd.Use,
				"name":      args[0],
			})

			logger.Info("deploy called")

			if cmd.Parent().Parent().Use == RsyncTool {
				logger.Infof("execute rsync deploy %s, volume is %s, namespace is %s, rousce is %v, sshuser: %s, sshpwd:%s, sshport:%s, action:%s\n",
					args[0], *volume, *namespace, *source, *sshuser, *sshpwd, *sshPort, cmd.Parent().Use)

				s := server.NewServer(RsyncTool, *sshuser, *sshpwd, *sshPort, *namespace, "deploy",
					args[0], *volume, source, -1, logger, cmd.Parent().Use)
				s.Run()
			} else if cmd.Parent().Parent().Use == ScpTool {
				logger.Infof("execute rsync deploy %s, volume is %s, namespace is %s, rousce is %v, sshuser: %s, sshpwd:%s, sshport:%s, action:%s\n",
					args[0], *volume, *namespace, *source, *sshuser, *sshpwd, *sshPort, cmd.Parent().Use)

				s := server.NewServer(ScpTool, *sshuser, *sshpwd, *sshPort, *namespace, "deploy",
					args[0], *volume, source, -1, logger, cmd.Parent().Use)
				s.Run()
			}
		},
	}
}

func init() {
	rsyncFromCmd.AddCommand(newDeployCmd())
	rsyncToCmd.AddCommand(newDeployCmd())
	scpFromCmd.AddCommand(newDeployCmd())
	scpToCmd.AddCommand(newDeployCmd())
}
