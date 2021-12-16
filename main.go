/*
Copyright 2021 Box-Cube

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

package main

import (
	"errors"
	"flag"
	"github.com/google/martian/log"
	"os"
	"strings"
	"sync-volume-data/server"
	"sync-volume-data/utils"
)

/*
	用法：
	sync-volume-data --kubeconfig={kubeconfig} --namespace={namespace} --resource={kind/source-name} \
	--volume={volume-name} --tool={rsync/scp} --source-dir={dir-or-file-path} --sshuser={ssh-user} \
	--sshpwd={ssh-password} --ssh-port={ssh-port}
*/

func main() {

	namespace := flag.String("namespace", "", "specific namespace")
	resource := flag.String("resource", "", "specific resource. exam: deploy/web")
	//container := flag.String("container", "", "specific container name")
	volume := flag.String("volume", "", "specific volume name in your specific resource")
	tool := flag.String("tool", "", "specific sync tool, now only support rsync/scp")
	sourceDir := flag.String("source-dir", "", "specific source directory where you want to sync")
	sshuser := flag.String("sshuser", "root", "specific user which can ssh to node")
	sshpwd := flag.String("sshpwd", "", "specific user which can ssh to node")
	sshPort := flag.String("ssh-port", "22", "specific port which can ssh to node")
	clientset := utils.NewClientset()

	sourceData := strings.Split(*resource, "/")
	if len(sourceData) < 2 {
		log.Errorf(errors.New("source need to specific, try -h get useage").Error())
		os.Exit(1)
	}
	resourceKind := sourceData[0]
	resourceName := sourceData[1]
	server :=server.NewServer(clientset, sshuser, sshpwd, sshPort, tool, namespace, &resourceKind, &resourceName, volume, sourceDir)

	server.Run()
}
