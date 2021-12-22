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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"strconv"
)

type statefulesetServer struct {
	namespace    string
	resourceName string
	volumeName   string
	volumeIndex  int
	log          *logrus.Entry
	kubeclient   *kubernetes.Clientset
	sts          *appsv1.StatefulSet
	pod          *corev1.Pod
}

func NewStatefulesetServer(namespace, resourceName, volumeName string, volumeIndex int, kubeclient *kubernetes.Clientset, log *logrus.Entry) *statefulesetServer {
	return &statefulesetServer{
		namespace:    namespace,
		resourceName: resourceName,
		volumeName:   volumeName,
		volumeIndex:  volumeIndex,
		kubeclient:   kubeclient,
		log:          log,
	}
}

func (s *statefulesetServer) getVolumePod() (volume *corev1.Volume, pod *corev1.Pod, err error) {
	sts, err := s.kubeclient.AppsV1().StatefulSets(s.namespace).Get(context.TODO(), s.resourceName, metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	s.sts = sts
	//get pod from specific statefuleset
	pod, err = s.getPodFromSource()
	if err != nil {
		return nil, nil, err
	}
	s.log.Infof("get pod %s from statefulset %s", pod.Name, sts.Name)

	podPvcName := s.volumeName + "-" + s.resourceName + "-" + strconv.Itoa(s.volumeIndex)
	for _, v := range pod.Spec.Volumes {
		if v.PersistentVolumeClaim.ClaimName == podPvcName {
			tmpVolume := v
			volume = &tmpVolume
			break
		}
	}

	return volume, pod, nil
}

func (s *statefulesetServer) getPodFromSource() (pod *corev1.Pod, err error) {
	label := metav1.LabelSelector{
		MatchLabels: s.sts.Spec.Template.Labels,
	}

	pods, err := s.kubeclient.CoreV1().Pods(s.namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labels.Set(label.MatchLabels).String()})
	if err != nil {
		return nil, err
	}

	if len(pods.Items) < 1 {
		return nil, errors.New(fmt.Sprintf("pods not found from label: %s, in namespace: %s", labels.Set(label.MatchLabels).String(), s.namespace))
	}

	for _, pod := range pods.Items {
		for _, own := range pod.OwnerReferences {
			if *own.Controller && own.Name == s.sts.Name && own.Kind == statefulsetKind &&
				pod.Status.Phase == corev1.PodRunning &&
				pod.Name == s.sts.Name+"-"+strconv.Itoa(s.volumeIndex) {
				return &pod, nil
			}
		}
	}

	return nil, errors.New("pod not found")
}
