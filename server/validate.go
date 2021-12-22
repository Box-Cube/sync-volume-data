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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

func (s *Server) ValidateTool() error {
	if s.tool == "rsync" || s.tool == "scp" {
		return nil
	} else if s.tool == "" {
		err := errors.New("sync tool command cannot be empty")
		s.errMsg = append(s.errMsg, err)
		return err
	} else {
		err := errors.New(fmt.Sprintf("not support tool command %s, pleas try \"-h\" to get useage", s.tool))
		s.errMsg = append(s.errMsg, err)
		return err
	}
}

func (s *Server) ValidateSshPwd() error {
	if s.sshpwd == "" {
		err := errors.New("ssh password cannot be empty")
		s.errMsg = append(s.errMsg, err)
		return err
	}
	return nil
}

func (s *Server) ValidateNamespace() error {
	if s.namespace == "" {
		err := errors.New("namespace cannot be empty")
		s.errMsg = append(s.errMsg, err)
		return err
	}

	_, err := s.kubeclient.CoreV1().Namespaces().Get(context.TODO(), s.namespace, metav1.GetOptions{})
	if err != nil {
		s.errMsg = append(s.errMsg, err)
		return err
	}

	return nil
}

func (s *Server) ValidateSourceKind() error {
	switch s.resourceKind {
	case "":
		err := errors.New("sourceKind cannot be empty")
		s.errMsg = append(s.errMsg, err)
		return err
	case "deploy", "deployment":
		s.resourceKind = deployKind
		return nil
	case "sts", "statefulset":
		s.resourceKind = statefulsetKind
		return nil
	case "ds", "daemonset":
		s.resourceKind = daemonsetKind
		return nil
	case "pod":
		s.resourceKind = podKind
		return nil
	default:
		err := errors.New(fmt.Sprintf("sourceKind %s not supported, pleas try \"-h\" to get useage", s.resourceKind))
		s.errMsg = append(s.errMsg, err)
		return err
	}
}

func (s *Server) ValidateSourceName() error {
	if s.resourceName == "" {
		err := errors.New("sourceName cannot be empty")
		s.errMsg = append(s.errMsg, err)
		return err
	}

	if s.resourceKind == deployKind {
		deploy, err := s.kubeclient.AppsV1().Deployments(s.namespace).Get(context.TODO(), s.resourceName, metav1.GetOptions{})
		if err != nil {
			s.errMsg = append(s.errMsg, err)
			return err
		}

		//Check whether the deploy status is ready
		expect := DeploymentComplete(deploy, &deploy.Status)
		if !expect {
			err = errors.New(fmt.Sprintf("deploy %s satuts not Completed", deploy.Name))
			s.errMsg = append(s.errMsg, err)
			return err
		}

		return nil
	} else if s.resourceKind == statefulsetKind {
		sts, err := s.kubeclient.AppsV1().StatefulSets(s.namespace).Get(context.TODO(), s.resourceName, metav1.GetOptions{})
		if err != nil {
			s.errMsg = append(s.errMsg, err)
			return err
		}

		expect := StatefulsetComplete(sts, &sts.Status)
		if !expect {
			err = errors.New(fmt.Sprintf("statefulset %s satuts not Completed", sts.Name))
			s.errMsg = append(s.errMsg, err)
			return err
		}

		return nil
	} else if s.resourceKind == daemonsetKind {
		ds, err := s.kubeclient.AppsV1().DaemonSets(s.namespace).Get(context.TODO(), s.resourceName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		expect := DaemonsetComplete(ds, &ds.Status)
		if !expect {
			err = errors.New(fmt.Sprintf("daemonset %s satuts not Completed", ds.Name))
			s.errMsg = append(s.errMsg, err)
			return err
		}

		return nil
	} else if s.resourceKind == podKind {
		pod, err := s.kubeclient.CoreV1().Pods(s.namespace).Get(context.TODO(), s.resourceName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if pod.Status.Phase != corev1.PodRunning {
			return errors.New(fmt.Sprintf("pod %s not running...", s.resourceName))
		}

		return nil
	} else {
		//TODO support statefulset and daemonset
		err := errors.New(fmt.Sprintf("sourceKind %s not supported, pleas try \"-h\" to get useage", s.resourceKind))
		s.errMsg = append(s.errMsg, err)
		return err
	}
}

func (s *Server) ValidateVolume() (exist bool, err error) {
	if s.volume == "" {
		err = errors.New("volume name cannot be empty")
		s.errMsg = append(s.errMsg, err)
		return false, err
	}

	//can execute ValidateVolume(), we think ValidateSourceName and ValidateSourceKind is ok
	if s.resourceKind == deployKind {
		deploy, err := s.kubeclient.AppsV1().Deployments(s.namespace).Get(context.TODO(), s.resourceName, metav1.GetOptions{})
		if err != nil {
			s.errMsg = append(s.errMsg, err)
			return false, err
		}

		for _, v := range deploy.Spec.Template.Spec.Volumes {
			if v.Name == s.volume {
				exist = true
			}
		}
	} else if s.resourceKind == statefulsetKind {
		sts, err := s.kubeclient.AppsV1().StatefulSets(s.namespace).Get(context.TODO(), s.resourceName, metav1.GetOptions{})
		if err != nil {
			s.errMsg = append(s.errMsg, err)
			return false, err
		}

		for _, v := range sts.Spec.VolumeClaimTemplates {
			if v.Name == s.volume {
				exist = true
			}
		}
	} else if s.resourceKind == daemonsetKind {
		ds, err := s.kubeclient.AppsV1().DaemonSets(s.namespace).Get(context.TODO(), s.resourceName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		for _, v := range ds.Spec.Template.Spec.Volumes {
			if v.Name == s.volume {
				exist = true
			}
		}
	} else if s.resourceKind == podKind {
		pod, err := s.kubeclient.CoreV1().Pods(s.namespace).Get(context.TODO(), s.resourceName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		for _, v := range pod.Spec.Volumes {
			if v.Name == s.volume {
				exist = true
			}
		}
	}

	if exist {
		return true, nil
	} else {
		err = errors.New(fmt.Sprintf("volume %s not exist", s.volume))
		s.errMsg = append(s.errMsg, err)
		return exist, err
	}
}

func (s *Server) ValidateInstanceIndex() (err error) {

	if s.resourceKind == statefulsetKind {
		if s.instanceIndex < 0 {
			err = errors.New(fmt.Sprintf("you need to specific volume-index when you use resource %s", s.resourceKind))
			s.errMsg = append(s.errMsg, err)
			return err
		}

		sts, err := s.kubeclient.AppsV1().StatefulSets(s.namespace).Get(context.TODO(), s.resourceName, metav1.GetOptions{})
		if err != nil {
			s.errMsg = append(s.errMsg, err)
			return err
		}

		if s.instanceIndex+1 > int(*sts.Spec.Replicas) {
			err = errors.New(fmt.Sprintf("The value of instance-index is greater than the Replicas of statefulset %s", s.resourceKind))
			s.errMsg = append(s.errMsg, err)
			return err
		}

		return nil
	} else {
		if s.instanceIndex >= 0 {
			err = errors.New(fmt.Sprintf("you don't need to specific volume-index when you use resource %s", s.resourceKind))
			s.errMsg = append(s.errMsg, err)
			return err
		}
		return nil
	}

}

func (s *Server) ValidateSourceDir() (exist bool, err error) {

	if len(*s.sourceDir) < 1 {
		err = errors.New("source file/directory cannot be empty")
		s.errMsg = append(s.errMsg, err)
		return false, err
	}

	for _, file := range *s.sourceDir {
		_, err = os.Stat(file)
		if err == nil {
			exist = true
		} else if os.IsNotExist(err) {
			s.errMsg = append(s.errMsg, err)
			exist = false
		}
	}

	//s.errMsg = append(s.errMsg, err)
	if exist {
		return true, nil
	} else {
		return false, err
	}
}

// DeploymentComplete considers a deployment to be complete once all of its desired replicas
// are updated and available, and no old pods are running.
func DeploymentComplete(deployment *appsv1.Deployment, newStatus *appsv1.DeploymentStatus) bool {
	return newStatus.UpdatedReplicas == *(deployment.Spec.Replicas) &&
		newStatus.Replicas == *(deployment.Spec.Replicas) &&
		newStatus.AvailableReplicas == *(deployment.Spec.Replicas) &&
		newStatus.ReadyReplicas == *(deployment.Spec.Replicas) &&
		newStatus.ObservedGeneration >= deployment.Generation
}

func StatefulsetComplete(sts *appsv1.StatefulSet, newStatus *appsv1.StatefulSetStatus) bool {
	return newStatus.ReadyReplicas == *(sts.Spec.Replicas) &&
		newStatus.Replicas == *(sts.Spec.Replicas) &&
		newStatus.CurrentReplicas == *(sts.Spec.Replicas) &&
		newStatus.UpdatedReplicas == *(sts.Spec.Replicas) &&
		newStatus.ObservedGeneration >= sts.Generation
	//newStatus.AvailableReplicas == *(sts.Spec.Replicas)
}

func DaemonsetComplete(ds *appsv1.DaemonSet, newStatus *appsv1.DaemonSetStatus) bool {
	return newStatus.CurrentNumberScheduled == newStatus.DesiredNumberScheduled &&
		newStatus.NumberAvailable == newStatus.DesiredNumberScheduled &&
		newStatus.NumberMisscheduled == 0 && newStatus.NumberReady == newStatus.DesiredNumberScheduled &&
		newStatus.ObservedGeneration >= ds.Generation
}
