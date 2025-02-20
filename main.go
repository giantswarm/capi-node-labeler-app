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
	// master will already have the label set, worker wont
	if hasLabel(node.Labels, MasterNodeRoleLabel) || hasLabel(node.Labels, ControlPlaneNodeRoleLabel) {
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
