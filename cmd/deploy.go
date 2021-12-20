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
	//cmd2 "sync-volume-data/cmd"
	"sync-volume-data/server"

	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
func newDeployCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "deploy",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("you need specific a deploy name")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("deploy called")

			if cmd.Parent().Use == "rsync" {
				fmt.Printf("execute rsync deploy %s, volume is %s, namespace is %s, rousce is %v, sshuser: %s, sshpwd:%s, sshport:%s\n",
					args[0], *volume, *namespace, *source, *sshuser, *sshpwd, *sshPort)
				s := server.NewServer("rsync", *sshuser, *sshpwd, *sshPort, *namespace, "deploy", args[0], *volume, source, -1)
				s.Run()
			} else if cmd.Parent().Use == "scp" {
				fmt.Printf("execute scp deploy %s, volume is %s\n", args[0], *volume)
				s := server.NewServer("scp", *sshuser, *sshpwd, *sshPort, *namespace, "deploy", args[0], *volume, source, -1)
				s.Run()
			}
		},
	}
}

func init() {
	rsyncCmd.AddCommand(newDeployCmd())
	scpCmd.AddCommand(newDeployCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
