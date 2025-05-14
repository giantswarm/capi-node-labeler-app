/*
Copyright 2021.

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
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/utils/env"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	//+kubebuilder:scaffold:imports
)

const (
	ControlPlaneNodeRoleLabel = "node-role.kubernetes.io/control-plane"
	MasterNodeRoleLabel       = "node-role.kubernetes.io/master"
	LegacyMasterNodeLabel     = "node.kubernetes.io/master"

	WorkerNodeRoleLabel   = "node-role.kubernetes.io/worker"
	LegacyWorkerNodeLabel = "node.kubernetes.io/worker"

	LegacyRoleLabel = "kubernetes.io/role"

	ControlPlaneNodeTaint = "node-role.kubernetes.io/control-plane"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	//+kubebuilder:scaffold:scheme
}

func main() {
	nodeName := env.GetString("NODE_NAME", "")
	if nodeName == "" {
		fmt.Printf("ERROR: NODE_NAME env cannot be empty\n")
		os.Exit(1)
	}

	config, err := ctrl.GetConfig()
	if err != nil {
		fmt.Printf("ERROR: failed to get config for controlelr runtime client\n")
		panic(err)
	}

	ctrlClient, err := client.New(config, client.Options{})
	if err != nil {
		fmt.Printf("ERROR: failed to create controller runtime client\n")
		panic(err)
	}

	ctx := context.TODO()
	var node v1.Node

	err = ctrlClient.Get(ctx, client.ObjectKey{Name: nodeName}, &node)
	if err != nil {
		fmt.Printf("ERROR: failed to get node %s\n", nodeName)
		panic(err)
	}

	shouldUpdate := false
	// check if the node is worker or master
	if isControlPlaneNode(node, ctrlClient) {
		// master node
		if !hasLabel(node.Labels, MasterNodeRoleLabel) {
			node.Labels[MasterNodeRoleLabel] = ""
			fmt.Printf("adding label %s=''\n", MasterNodeRoleLabel)
			shouldUpdate = true
		}
		if !hasLabel(node.Labels, ControlPlaneNodeRoleLabel) {
			node.Labels[ControlPlaneNodeRoleLabel] = ""
			fmt.Printf("adding label %s=''\n", ControlPlaneNodeRoleLabel)
			shouldUpdate = true
		}
		if !hasLabel(node.Labels, LegacyMasterNodeLabel) {
			node.Labels[LegacyMasterNodeLabel] = ""
			fmt.Printf("adding label %s=''\n", LegacyMasterNodeLabel)
			shouldUpdate = true
		}
		if !hasLabel(node.Labels, LegacyRoleLabel) {
			node.Labels[LegacyRoleLabel] = "master"
			fmt.Printf("adding label %s='master'\n", LegacyRoleLabel)
			shouldUpdate = true
		}
		if hasLabel(node.Labels, WorkerNodeRoleLabel) {
			delete(node.Labels, WorkerNodeRoleLabel)
			fmt.Printf("removing label %s\n", WorkerNodeRoleLabel)
			shouldUpdate = true
		}
		if hasLabel(node.Labels, LegacyWorkerNodeLabel) {
			delete(node.Labels, LegacyWorkerNodeLabel)
			fmt.Printf("removing label %s\n", LegacyWorkerNodeLabel)
			shouldUpdate = true
		}

		if !hasTaint(node.Spec.Taints, ControlPlaneNodeTaint) {
			node.Spec.Taints = append(node.Spec.Taints, v1.Taint{
				Key:    ControlPlaneNodeTaint,
				Effect: v1.TaintEffectNoSchedule,
			})
			fmt.Printf("adding taint %s=''\n", ControlPlaneNodeTaint)
			shouldUpdate = true
		}
	} else {
		// worker node
		if !hasLabel(node.Labels, WorkerNodeRoleLabel) {
			node.Labels[WorkerNodeRoleLabel] = ""
			fmt.Printf("adding label %s=''\n", WorkerNodeRoleLabel)
			shouldUpdate = true
		}
		if !hasLabel(node.Labels, LegacyWorkerNodeLabel) {
			node.Labels[LegacyWorkerNodeLabel] = ""
			fmt.Printf("adding label %s=''\n", LegacyWorkerNodeLabel)
			shouldUpdate = true
		}
		if !hasLabel(node.Labels, LegacyRoleLabel) {
			node.Labels[LegacyRoleLabel] = "worker"
			fmt.Printf("adding label %s='worker'\n", LegacyRoleLabel)
			shouldUpdate = true
		}
	}

	if shouldUpdate {
		err = ctrlClient.Update(ctx, &node)
		if err != nil {
			fmt.Printf("ERROR: failed to apply new labels\n")
			panic(err)
		}
		fmt.Printf("new labels applied to node\n")
	} else {
		fmt.Printf("required labels are already applied\n")
	}

	fmt.Printf("capi-node-labeler finished successfully\n")
	fmt.Printf("sleeping forever\n")

	// Don't try select {} here, it will cause the following error:
	//
	// fatal error: all goroutines are asleep - deadlock!
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
}

func hasLabel(labels map[string]string, labelName string) bool {
	_, ok := labels[labelName]
	return ok
}

func hasTaint(taints []v1.Taint, taintKey string) bool {
	for _, taint := range taints {
		if taint.Key == taintKey {
			return true
		}
	}
	return false
}

func isControlPlaneNode(node v1.Node, ctrlClient client.Client) bool {
	if hasLabel(node.Labels, MasterNodeRoleLabel) || hasLabel(node.Labels, ControlPlaneNodeRoleLabel) || hasLabel(node.Labels, LegacyMasterNodeLabel) {
		return true
	}

	// if node doesn't have any of the control plane labels, check if it has control plane Pods running as fallback.
	// During DR scenarios, the node may not have any of the control plane labels but may still be a control plane node.
	var apiPods v1.PodList
	err := ctrlClient.List(context.TODO(), &apiPods, client.InNamespace("kube-system"), client.MatchingLabels{"component": "kube-apiserver", "tier": "control-plane"})
	if err != nil {
		fmt.Printf("ERROR: failed to list api-server Pods in kube-system namespace\n")
		panic(err)
	}

	var etcdPods v1.PodList
	err = ctrlClient.List(context.TODO(), &etcdPods, client.InNamespace("kube-system"), client.MatchingLabels{"component": "etcd", "tier": "control-plane"})
	if err != nil {
		fmt.Printf("ERROR: failed to list etcd Pods in kube-system namespace\n")
		panic(err)
	}

	var cpPods v1.PodList
	cpPods.Items = append(apiPods.Items, etcdPods.Items...)

	for _, pod := range cpPods.Items {
		if pod.Spec.NodeName == node.Name {
			return true
		}
	}

	return false
}
