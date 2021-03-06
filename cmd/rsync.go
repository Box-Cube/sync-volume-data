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

// rsyncCmd represents the rsync command
var rsyncCmd = &cobra.Command{
	Use:   "rsync",
	Short: "use rsync tool to trans your data",
	Long: `command will use rsync tool to trans data from local to remote volume
	please make sure you have rsync command on your machine
			you can use rsync to trans your local File/Directory
 For example:
	
	sync-volume-data rsync to deploy nginx -n my-web -v web -u root -p "myPassword" -s=test.file
`,
}

func init() {
	rootCmd.AddCommand(rsyncCmd)
}
