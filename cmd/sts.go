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
	"github.com/spf13/cobra"
	"sync-volume-data/server"
)

// stsCmd represents the sts command
func newStsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sts",
		Short: "transfer data from/to StatefulSet kind resource",
		Long: `transfer data from/to StatefulSet kind resource, you need to specific a sts name.
               In addition, since the STS is stateful, an additional index of the specified instance is required through the "-i" flag.
               instance-index starts from 0.
 For example:
	
	./sync-volume-tool rsync to sts web -n my-example -v www-1 -i 0 -p "myPassword" -s=test.file
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("you need specific a statefulset name")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger := newLogger().WithFields(logrus.Fields{
				"namespace": *namespace,
				"kind":      cmd.Use,
				"name":      args[0],
			})
			logger.Debug("sts called")
			switch cmd.Parent().Parent().Use {
			case RsyncTool:
				if cmd.Parent().Use == "from" {
					instanceIndex = rsyncFromInstanceIndex
				} else if cmd.Parent().Use == "to" {
					instanceIndex = rsyncToInstanceIndex
				}
			case ScpTool:
				if cmd.Parent().Use == "from" {
					instanceIndex = scpFromInstanceIndex
				} else if cmd.Parent().Use == "to" {
					instanceIndex = scpToInstanceIndex
				}
			}

			if cmd.Parent().Parent().Use == RsyncTool {
				logger.Printf("execute rsync deploy %s, volume is %s, namespace is %s, rousce is %v, sshuser: %s, sshpwd:%s, sshport:%s, instanceIdex:%d\n",
					args[0], *volume, *namespace, *source, *sshuser, *sshpwd, *sshPort, *instanceIndex)

				s := server.NewServer(RsyncTool, *sshuser, *sshpwd, *sshPort, *namespace, "sts",
					args[0], *volume, source, *instanceIndex, logger, cmd.Parent().Use)
				s.Run()
			} else if cmd.Parent().Parent().Use == ScpTool {
				logger.Printf("execute rsync deploy %s, volume is %s, namespace is %s, rousce is %v, sshuser: %s, sshpwd:%s, sshport:%s, instanceIdex:%d\n",
					args[0], *volume, *namespace, *source, *sshuser, *sshpwd, *sshPort, *instanceIndex)

				s := server.NewServer(ScpTool, *sshuser, *sshpwd, *sshPort, *namespace, "sts",
					args[0], *volume, source, *instanceIndex, logger, cmd.Parent().Use)
				s.Run()
			}
		},
	}
}

var (
	rsyncFromInstanceIndex *int
	rsyncToInstanceIndex   *int
	scpFromInstanceIndex   *int
	scpToInstanceIndex     *int
	instanceIndex          *int
)

func init() {
	rsyncFromStsCmd := newStsCmd()
	rsyncToStsCmd := newStsCmd()
	scpFromStsCmd := newStsCmd()
	scpToStsCmd := newStsCmd()

	rsyncFromCmd.AddCommand(rsyncFromStsCmd)
	rsyncToCmd.AddCommand(rsyncToStsCmd)
	scpFromCmd.AddCommand(scpFromStsCmd)
	scpToCmd.AddCommand(scpToStsCmd)

	rsyncFromInstanceIndex = rsyncFromStsCmd.Flags().IntP("instance-index", "i", -1, "specific instance index when you use statefulset kind resource")
	rsyncToInstanceIndex = rsyncToStsCmd.Flags().IntP("instance-index", "i", -1, "specific instance index when you use statefulset kind resource")
	scpFromInstanceIndex = scpFromStsCmd.Flags().IntP("instance-index", "i", -1, "specific instance index when you use statefulset kind resource")
	scpToInstanceIndex = scpToStsCmd.Flags().IntP("instance-index", "i", -1, "specific instance index when you use statefulset kind resource")
}
