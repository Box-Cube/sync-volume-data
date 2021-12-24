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
	"github.com/spf13/cobra"
)

// toCmd represents the to command
func newToCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "to",
		Short: "transfer data from local machine to remote deploy/sts/ds/pod kind resource",
		Long: `transfer data from local machine to remote deploy/sts/ds/pod kind resource
	For example:
		./sync-volume-tool rsync to pod web-1-789cb6ff95-wfhk2  -n my-example -v mypd  -p 'password' -s  file-test,dir-test/local
`,
	}
}

var (
	rsyncToCmd = newToCmd()
	scpToCmd   = newToCmd()
)

func init() {
	rsyncCmd.AddCommand(rsyncToCmd)
	scpCmd.AddCommand(scpToCmd)
}
