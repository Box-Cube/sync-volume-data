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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/api/apps/v1"
)

type daemonsetServer struct {
	namespace string
	resourceName string
	volumeName string
	kubeclient *kubernetes.Clientset
	daemonset *appsv1.DaemonSet
	pod *corev1.Pod
}

func NewDaemonsetServer(namespace, resourceName, volumeName string, kubeclient *kubernetes.Clientset) *daemonsetServer {
	return &daemonsetServer{
		namespace: namespace,
		resourceName: resourceName,
		volumeName: volumeName,
		kubeclient: kubeclient,
	}
}

func (d *daemonsetServer) getVolumePod() (volume *corev1.Volume, pod *corev1.Pod, err error) {
	ds, err := d.kubeclient.AppsV1().DaemonSets(d.namespace).Get(context.TODO(), d.resourceName, metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	d.daemonset = ds
	pod, err = d.getPodFromSource()
	if err != nil {
		return nil, nil, err
	}

	for _, v := range ds.Spec.Template.Spec.Volumes {
		if v.Name == d.volumeName {
			tmpVolume := v
			volume = &tmpVolume
			break
		}
	}

	return volume, pod, nil
}

func (d *daemonsetServer) getPodFromSource() (pod *corev1.Pod, err error) {
	pods, err := d.kubeclient.CoreV1().Pods(d.namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labels.Set(d.daemonset.Spec.Template.Labels).String()})
	if err != nil {
		return nil, err
	}

	if len(pods.Items) < 1 {
		return nil, errors.New("pods not found")
	}

	for _, pod := range pods.Items {
		for _, own := range pod.OwnerReferences {
			if own.Name == d.resourceName && *own.Controller && pod.Status.Phase == corev1.PodRunning &&
				own.Kind == "DaemonSet" {
				return &pod, nil
			}
		}
	}

	return nil, errors.New("pod not found")
}