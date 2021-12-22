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

// stsCmd represents the sts command
func newStsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sts",
		Short: "transfer data to StatefulSet kind resource",
		Long: `transfer data to StatefulSet kind resource, you need to specific a sts name.
               In addition, since the STS is stateful, an additional index of the specified instance is required through the "-i" flag.
               instance-index starts from 0.
 For example:
	
	./sync-volume-tool rsync sts web -n my-example -v www-1 -i 0 -p "myPassword" -s=test.file
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
				"kind":   cmd.Use,
				"name": args[0],
			})
			logger.Debug("sts called")
			if cmd.Parent().Use == RsyncTool {
				fmt.Printf("execute rsync deploy %s, volume is %s, namespace is %s, rousce is %v, sshuser: %s, sshpwd:%s, sshport:%s, instanceIdex:%d\n",
					args[0], *volume, *namespace, *source, *sshuser, *sshpwd, *sshPort, *rsyncInstanceIndex)
				s := server.NewServer(RsyncTool, *sshuser, *sshpwd, *sshPort, *namespace, "sts", args[0], *volume, source, *rsyncInstanceIndex, logger)
				s.Run()
			} else if cmd.Parent().Use == ScpTool {
				fmt.Printf("execute scp deploy %s, volume is %s\n", args[0], *volume)
				s := server.NewServer(ScpTool, *sshuser, *sshpwd, *sshPort, *namespace, "sts", args[0], *volume, source, *scpInstanceIndex, logger)
				s.Run()
			}
		},
	}
}

var (
	rsyncInstanceIndex *int
	scpInstanceIndex *int
)

func init() {
	rsyncStsCmd := newStsCmd()
	scpStsCmd := newStsCmd()
	rsyncCmd.AddCommand(rsyncStsCmd)
	scpCmd.AddCommand(scpStsCmd)

	rsyncInstanceIndex = rsyncStsCmd.Flags().IntP("instance-index", "i", -1, "specific instance index when you use statefulset kind resource")
	scpInstanceIndex = scpStsCmd.Flags().IntP("instance-index", "i", -1, "specific instance index when you use statefulset kind resource")
}
