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
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"os"
	"strings"
)

type deployServer struct {
	namespace    string
	resourceName string
	volumeName   string
	kubeclient   *kubernetes.Clientset
	deploy       *appsv1.Deployment
	pod          *corev1.Pod
	log          *logrus.Entry
}

func NewDeployServer(namespace, resourName, volumeName string, kubeclient *kubernetes.Clientset, logger *logrus.Entry) *deployServer {
	return &deployServer{
		namespace:    namespace,
		resourceName: resourName,
		volumeName:   volumeName,
		kubeclient:   kubeclient,
		log:          logger,
	}
}

func (d *deployServer) getVolumePod() (volume *corev1.Volume, pod *corev1.Pod, err error) {
	deploy, err := d.kubeclient.AppsV1().Deployments(d.namespace).Get(context.TODO(), d.resourceName, metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	d.deploy = deploy
	//get pod from specific deploy
	pod, err = d.getPodFromSource()
	if err != nil {
		d.log.Errorf("%s", err.Error())
		os.Exit(1)
	}
	d.log.Infof("get pod %s from deployment %s\n", pod.Name, deploy.Name)

	for _, v := range deploy.Spec.Template.Spec.Volumes {
		if v.Name == d.volumeName {
			tmpVolume := v
			volume = &tmpVolume
			break
		}
	}

	return volume, pod, nil
}

func (d *deployServer) getPodFromSource() (pod *corev1.Pod, err error) {

	lebel := metav1.LabelSelector{
		MatchLabels:      d.deploy.Spec.Template.Labels,
		MatchExpressions: nil,
	}
	pods, err := d.kubeclient.CoreV1().Pods(d.namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labels.Set(lebel.MatchLabels).String()})
	if err != nil {
		return nil, err
	}

	if len(pods.Items) < 1 {
		return nil, errors.New("pods not found")
	}

	//Since deploy is working as expected, pick a pod at random
	for _, pod := range pods.Items {
		for _, own := range pod.OwnerReferences {
			deployName := own.Name[0:strings.LastIndex(own.Name, "-")]
			//log.Infof("get deploy name: %s from pod %s", deployName, pod.Name)
			if deployName == d.resourceName && *own.Controller &&
				own.Kind == replicaSetKind && pod.Status.Phase == corev1.PodRunning {
				return &pod, nil
				// ??????????????????????????????...
				//rs, err := s.kubeclient.AppsV1().ReplicaSets(s.namespace).Get(context.TODO(), own.Name, metav1.GetOptions{})
				//if err != nil {
				//	return nil, err
				//}
				//
				//for _, rsOwn := range rs.OwnerReferences {
				//	if deployName == rsOwn.Name {
				//		return &pod, nil
				//	}
				//}
			}
		}
	}

	return nil, errors.New("pod not found")
}
