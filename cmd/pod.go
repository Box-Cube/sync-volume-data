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
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sync-volume-data/server"
)

// podCmd represents the pod command
func newPodCmd() *cobra.Command {
	return  &cobra.Command{
		Use:   "pod",
		Short: "transfer data to Pod kind resource",
		Long: `transfer data to Pod kind resource, you need to specific a pod name.
 For example:
	
	sync-volume-data rsync pod web-1-789cb6ff95-wfhk2 -n my-web -v web -u root -p "myPassword" -s=test.file
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("you need specific a pod name")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger := newLogger().WithFields(logrus.Fields{
				"namespace": *namespace,
				"kind":   cmd.Use,
				"name": args[0],
			})
			logger.Debug("pod called")
			if cmd.Parent().Use == RsyncTool {
				fmt.Printf("execute rsync pod %s, volume is %s, namespace is %s, rousce is %v, sshuser: %s, sshpwd:%s, sshport:%s\n",
					args[0], *volume, *namespace, *source, *sshuser, *sshpwd, *sshPort)
				s := server.NewServer(RsyncTool, *sshuser, *sshpwd, *sshPort, *namespace, "pod", args[0], *volume, source, -1, logger)
				s.Run()
			} else if cmd.Parent().Use == ScpTool {
				fmt.Printf("execute scp pod %s, volume is %s\n", args[0], *volume)
				s := server.NewServer(ScpTool, *sshuser, *sshpwd, *sshPort, *namespace, "pod", args[0], *volume, source, -1, logger)
				s.Run()
			}
		},
	}
}

func init() {
	rsyncCmd.AddCommand(newPodCmd())
	scpCmd.AddCommand(newPodCmd())
}
