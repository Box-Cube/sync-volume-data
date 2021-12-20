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
	Long: `you can use rsync to trans your local File/Directory
 For example:
	
	sync-volume-data rsync deploy nginx -n my-web -v web -sshuser root -sshpwd "myPassword" 
`,
	//Args: func(cmd *cobra.Command, args []string) error {
	//	if len(args) < 1 {
	//		return errors.New("source file/directory cannot be empty")
	//	}
	//	_, err := os.Stat(args[0])
	//	if err != nil {
	//		return err
	//	}
	//
	//	return nil
	//},
	//Run: func(cmd *cobra.Command, args []string) {
	//	//fmt.Println("you need to specific a resource kind ")
	//	log.Errorf("you need to specific a resource kind ")
	//	*tool = "rsync"
	//},
}

func init() {
	rootCmd.AddCommand(rsyncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rsyncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rsyncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
