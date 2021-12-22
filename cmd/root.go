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
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"path"
	"runtime"
	"strings"
)

var (
	volume *string
	namespace *string
	source *[]string
	sshuser *string
	sshpwd *string
	sshPort *string
)

const (
	RsyncTool = "rsync"
	ScpTool = "scp"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sync-volume-data",
	Short: "sync-volume-data transfers local files/directories to a specified resource kind",
	Long: `sync-volume-data transfers local files/directories to a specified resource kind
           Rsync and SCP are supported. Ensure that the two commands have been installed on the local host.
           And make sure that the local machine has kubeconfig to connect to the K8S cluster, 
           the network of the local machine and the internal IP of the K8S node are communicating.
`,
	Version: "alpha v1.0",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sync-volume-data.yaml)")
	volume = rootCmd.PersistentFlags().StringP("volume", "v", "", "specific volume name in your specific resource")
	namespace = rootCmd.PersistentFlags().StringP("namespace", "n", "", "specific namespace")
	source = rootCmd.PersistentFlags().StringSliceP("source", "s", []string{}, "specific source file/directory which you want to transfer")
	sshuser = rootCmd.PersistentFlags().StringP("ssh-user", "u", "root", "specific user which can ssh to node")
	sshpwd = rootCmd.PersistentFlags().StringP("ssh-password", "p", "", "specific password which can ssh to node")
	sshPort = rootCmd.PersistentFlags().StringP("ssh-port", "P", "22", "specific port which can ssh to node")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// mark Required flag
	rootCmd.MarkPersistentFlagRequired("volume")
	rootCmd.MarkPersistentFlagRequired("namespace")
	rootCmd.MarkPersistentFlagRequired("source")
	rootCmd.MarkPersistentFlagRequired("ssh-password")

}

func newLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{DisableColors: true, FullTimestamp: true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier:  func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File)
			s := strings.Split(frame.Function, ".")
			funcName := s[len(s)-1]
			return fmt.Sprintf("%s()",funcName), fmt.Sprintf("%s:%d", fileName, frame.Line)
		}})
	return logger
}