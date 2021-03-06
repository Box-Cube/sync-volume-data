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

package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"os/exec"
	"strings"
	remote "sync-volume-data/remote_execute"
	"sync-volume-data/utils"
	"syscall"
)

type TransferAction string

const (
	TransferTo   = "to"
	TransferFrom = "from"
)

type Server struct {
	kubeclient    *kubernetes.Clientset
	sshuser       string
	sshpwd        string
	sshPort       string
	tool          string
	namespace     string
	resourceKind  string
	resourceName  string
	volume        string
	instanceIndex int
	sourceDir     *[]string
	errMsg        []error
	action        string
	log           *logrus.Entry
}

func NewServer(tool, sshuser, sshpwd, sshPort, namespace, resourceKind, resourceName, volume string, sourceDir *[]string,
	instanceIndex int, logger *logrus.Entry, action string) *Server {
	errMsg := new([]error)

	kubeclient := utils.NewClientset()

	return &Server{
		kubeclient:    kubeclient,
		tool:          tool,
		namespace:     namespace,
		resourceKind:  resourceKind,
		resourceName:  resourceName,
		volume:        volume,
		instanceIndex: instanceIndex,
		sourceDir:     sourceDir,
		sshuser:       sshuser,
		sshpwd:        sshpwd,
		sshPort:       sshPort,
		errMsg:        *errMsg,
		log:           logger,
		action:        action,
	}
}

// 获取 volume directory 全路径
///var/lib/kubelet/pods/<podUID>/volumes/<volume-plugin-name>/<volume-dir>
//<volume-dir> == {PV-NAME} + /mount (CSI 用到就有/mount)
type resourceInfoer interface {
	getVolumePod() (volume *corev1.Volume, pod *corev1.Pod, err error)
}

const (
	deployKind      = "Deployment"
	statefulsetKind = "StatefulSet"
	daemonsetKind   = "DaemonSet"
	replicaSetKind  = "ReplicaSet"
	podKind         = "Pod"
)

func (s *Server) Run() {
	s.validateParameter()
	var pod *corev1.Pod
	var err error
	var volume *corev1.Volume
	var nodeIP string
	var sourceExec resourceInfoer
	defaultRootDir := "/var/lib/kubelet/pods/"

	if s.resourceKind == deployKind {
		sourceExec = NewDeployServer(s.namespace, s.resourceName, s.volume, s.kubeclient, s.log)
		//volume, pod, err = deployRun.getVolumeInfo()
	} else if s.resourceKind == daemonsetKind {
		sourceExec = NewDaemonsetServer(s.namespace, s.resourceName, s.volume, s.kubeclient, s.log)
	} else if s.resourceKind == statefulsetKind {
		sourceExec = NewStatefulesetServer(s.namespace, s.resourceName, s.volume, s.instanceIndex, s.kubeclient, s.log)
	} else if s.resourceKind == podKind {
		sourceExec = NewPodServer(s.namespace, s.resourceName, s.volume, s.kubeclient)
	}

	volume, pod, err = sourceExec.getVolumePod()
	if err != nil {
		s.log.Fatal(err)
	}

	nodeIP, err = s.getNodeIPFromPod(pod)
	if err != nil {
		s.log.Fatal(err)
	}
	s.log.Infof("get node ip %s from pod %s", nodeIP, pod.Name)

	volumeDir, err := s.GetVolumeDirectory(volume)
	if err != nil {
		s.log.Fatal(err)
	}

	volumePath := defaultRootDir + string(pod.UID) + "/volumes/*/" + volumeDir
	s.log.Infof("get volume path: %s", volumePath)

	//for debug
	//nodeIP = "180.184.65.175"
	//nodeIP = "180.184.64.139"
	//TODO, now only support password method,key method will be supported later
	sshcli := remote.NewCli(s.sshuser, s.sshpwd, fmt.Sprintf("%s:%s", nodeIP, s.sshPort), remote.SshPassword, "")

	//get only a row as expected
	actualVolumePath, err := sshcli.Run(fmt.Sprintf("ls -d %s | awk 'NR=1{printf $NF}'", volumePath))
	if err != nil {
		s.log.Fatal("get err from remote node %s : %s", nodeIP, err.Error())
	}

	s.log.Infof("get volume path from remote node: %s", actualVolumePath)

	var args []string
	command := s.tool

	if s.tool == "scp" {
		args = []string{
			"-rp",
			"-P",
			s.sshPort,
		}
		if s.action == TransferTo {
			for _, file := range *s.sourceDir {
				args = append(args, file)
			}

			args = append(args, fmt.Sprintf("%s@%s:%s", s.sshuser, nodeIP, actualVolumePath))
		} else if s.action == TransferFrom {
			if len(*s.sourceDir) == 1 {
				args = append(args, fmt.Sprintf(`%s@%s:%s/%s`, s.sshuser, nodeIP, actualVolumePath, strings.Join(*s.sourceDir, "")))
			} else if len(*s.sourceDir) > 1 {
				args = append(args, fmt.Sprintf(`%s@%s:%s/{%s}`, s.sshuser, nodeIP, actualVolumePath, strings.Join(*s.sourceDir, ",")))

			}
			//TODO support destination
			args = append(args, ".")
		}

	} else if s.tool == "rsync" {
		args = []string{
			"-av",
			"--progress",
			"-e", fmt.Sprintf("ssh -p %s", s.sshPort),
		}
		if s.action == TransferTo {
			for _, file := range *s.sourceDir {
				args = append(args, file)
			}

			args = append(args, fmt.Sprintf("%s@%s:%s", s.sshuser, nodeIP, actualVolumePath))
		} else if s.action == TransferFrom {
			if len(*s.sourceDir) == 1 {
				filePath := actualVolumePath + "/" + strings.Join(*s.sourceDir, "")
				args = append(args, fmt.Sprintf("%s@%s:%s", s.sshuser, nodeIP, filePath))
			} else if len(*s.sourceDir) > 1 {
				args = append(args, fmt.Sprintf("%s@%s:%s/{%s} ", s.sshuser, nodeIP, actualVolumePath, strings.Join(*s.sourceDir, ",")))
			}
			//TODO support destination
			args = append(args, ".")
		}
	}
	s.log.Infof("execute command: %s args: %s", command, args)

	cmd := exec.Command(command, args...)
	// 命令的错误输出和标准输出都连接到同一个管道
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		s.log.Fatal(err.Error())
	}

	if err = cmd.Start(); err != nil {
		s.log.Fatal(err.Error())
	}
	// Get the output from the pipe in real time and print it to the terminal
	for {
		tmp := make([]byte, 2048)
		_, err := stdout.Read(tmp)
		fmt.Print(string(tmp))
		if err != nil {
			break
		}
	}

	if err = cmd.Wait(); err != nil {
		if ex, ok := err.(*exec.ExitError); ok {
			res := ex.Sys().(syscall.WaitStatus).ExitStatus() //获取命令执行返回状态，相当于shell: echo $?
			fmt.Println("#####################################################################################")
			s.log.Errorf("sync data failed, exit code is %d, err: %s", res, err)
			return
		}
	}

	fmt.Println("#####################################################################################")
	s.log.Infof("sync data to pod volume succeed !!")
}

func (s *Server) validateParameter() {
	//s.ValidateTool()
	s.ValidateSshPwd()
	s.ValidateNamespace()
	s.ValidateSourceKind()
	s.ValidateSourceName()
	s.ValidateVolume()
	s.ValidateInstanceIndex()
	s.ValidateSourceDir()

	if len(s.errMsg) > 0 {
		for _, err := range s.errMsg {
			s.log.Errorf(err.Error())
		}
		os.Exit(1)
	}
}

func (s *Server) getNodeIPFromPod(pod *corev1.Pod) (nodeIP string, err error) {
	//get node where pod is running
	node, err := s.kubeclient.CoreV1().Nodes().Get(context.TODO(), pod.Spec.NodeName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	for _, address := range node.Status.Addresses {
		switch address.Type {
		case corev1.NodeInternalIP:
			nodeIP = address.Address
		default:
		}
	}

	if nodeIP == "" {
		return "", errors.New("node where pod is running, can't get Addresses")
	}

	return nodeIP, nil
}

func (s *Server) GetVolumeDirectory(volume *corev1.Volume) (string, error) {
	// This case implies the administrator created the PV and attached it directly, without PVC.
	// Note that only one VolumeSource can be populated per Volume on a pod
	if volume.VolumeSource.PersistentVolumeClaim == nil {
		if volume.VolumeSource.CSI != nil {
			return volume.Name + "/mount", nil
		}
		return volume.Name, nil
	}

	// Most common case is that we have a PVC VolumeSource, and we need to check the PV it points to for a CSI source.
	pvc, err := s.kubeclient.CoreV1().PersistentVolumeClaims(s.namespace).Get(context.TODO(), volume.VolumeSource.PersistentVolumeClaim.ClaimName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	pv, err := s.kubeclient.CoreV1().PersistentVolumes().Get(context.TODO(), pvc.Spec.VolumeName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// PV's been created with a CSI source.
	if pv.Spec.CSI != nil {
		return pvc.Spec.VolumeName + "/mount", nil
	}

	return pvc.Spec.VolumeName, nil
}
