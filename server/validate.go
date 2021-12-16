package server

import (
	"context"
	"errors"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
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

		err := errors.New("statefulset kind is not currently supported, so stay tuned")
		s.errMsg = append(s.errMsg, err)
		return err
	case "ds", "daemonset":
		s.resourceKind = daemonsetKind

		err := errors.New("daemonset kind is not currently supported, so stay tuned")
		s.errMsg = append(s.errMsg, err)
		return err
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
	} else {
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
	}

	if exist {
		return true, nil
	} else {
		err = errors.New(fmt.Sprintf("volume %s not exist", s.volume))
		s.errMsg = append(s.errMsg, err)
		return exist, err
	}
}

func (s *Server) ValidateSourceDir() (exist bool, err error) {
	if s.sourceDir == "" {
		err = errors.New("source file/directory cannot be empty")
		s.errMsg = append(s.errMsg, err)
		return false, err
	}

	_, err = os.Stat(s.sourceDir)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		s.errMsg = append(s.errMsg, err)
		return false, err
	}

	s.errMsg = append(s.errMsg, err)
	return false, err
}

// DeploymentComplete considers a deployment to be complete once all of its desired replicas
// are updated and available, and no old pods are running.
func DeploymentComplete(deployment *appsv1.Deployment, newStatus *appsv1.DeploymentStatus) bool {
	return newStatus.UpdatedReplicas == *(deployment.Spec.Replicas) &&
		newStatus.Replicas == *(deployment.Spec.Replicas) &&
		newStatus.AvailableReplicas == *(deployment.Spec.Replicas) &&
		newStatus.ObservedGeneration >= deployment.Generation
}
