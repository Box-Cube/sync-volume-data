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
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sync-volume-data/server"
)

// dsCmd represents the ds command
func newDsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ds",
		Short: "transfer data from/to DaemonSet kind resource",
		Long: `transfer data from/to DaemonSet kind resource, you need to specific a daemonset name.
 For example:
	
	sync-volume-data rsync from ds nginx -n my-web -v web -u root -p "myPassword" -s=test.file
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("you need specific a daemonset name")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger := newLogger().WithFields(logrus.Fields{
				"namespace": *namespace,
				"kind":      cmd.Use,
				"name":      args[0],
			})
			logger.Debug("ds called")
			if cmd.Parent().Parent().Use == RsyncTool {
				fmt.Printf("execute rsync daemonset %s, volume is %s, namespace is %s, rousce is %v, sshuser: %s, sshpwd:%s, sshport:%s\n",
					args[0], *volume, *namespace, *source, *sshuser, *sshpwd, *sshPort)

				s := server.NewServer(RsyncTool, *sshuser, *sshpwd, *sshPort, *namespace, "ds",
					args[0], *volume, source, -1, logger, cmd.Parent().Use)
				s.Run()
			} else if cmd.Parent().Parent().Use == ScpTool {
				fmt.Printf("execute scp daemonset %s, volume is %s\n", args[0], *volume)

				s := server.NewServer(ScpTool, *sshuser, *sshpwd, *sshPort, *namespace, "ds",
					args[0], *volume, source, -1, logger, cmd.Parent().Use)
				s.Run()
			}
		},
	}
}

func init() {
	rsyncFromCmd.AddCommand(newDsCmd())
	rsyncToCmd.AddCommand(newDsCmd())
	scpFromCmd.AddCommand(newDsCmd())
	scpToCmd.AddCommand(newDsCmd())
}
