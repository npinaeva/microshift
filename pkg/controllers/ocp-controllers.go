/*
Copyright © 2021 Microshift Contributors

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
package controllers

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	openshift_apiserver "github.com/openshift/openshift-apiserver/pkg/cmd/openshift-apiserver"
	openshift_controller_manager "github.com/openshift/openshift-controller-manager/pkg/cmd/openshift-controller-manager"

	"github.com/openshift/microshift/pkg/config"
)

func newOpenshiftApiServerCommand(stopCh <-chan struct{}) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "openshift-apiserver",
		Short: "Command for the OpenShift API Server",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(1)
		},
	}
	start := openshift_apiserver.NewOpenShiftAPIServerCommand("start", os.Stdout, os.Stderr, stopCh)
	cmd.AddCommand(start)

	return cmd
}
func OCPAPIServer(cfg *config.MicroshiftConfig) error {
	stopCh := make(chan struct{})
	command := newOpenshiftApiServerCommand(stopCh)
	args := []string{
		"start",
		"--config=" + cfg.DataDir + "/resources/openshift-apiserver/config/config.yaml",
		"--authorization-kubeconfig=" + cfg.DataDir + "/resources/kubeadmin/kubeconfig",
		"--authentication-kubeconfig=" + cfg.DataDir + "/resources/kubeadmin/kubeconfig",
		"--requestheader-client-ca-file=" + cfg.DataDir + "/certs/ca-bundle/ca-bundle.crt",
		"--requestheader-allowed-names=kube-apiserver-proxy,system:kube-apiserver-proxy,system:openshift-aggregator",
		"--requestheader-username-headers=X-Remote-User",
		"--requestheader-group-headers=X-Remote-Group",
		"--requestheader-extra-headers-prefix=X-Remote-Extra-",
		"--client-ca-file=" + cfg.DataDir + "/certs/ca-bundle/ca-bundle.crt",
	}
	command.SetArgs(args)
	logrus.Infof("starting openshift-apiserver, args: %v", args)
	go func() {
		logrus.Fatalf("ocp apiserver exited: %v", command.Execute())
	}()

	return nil
}

func newOpenShiftControllerManagerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "openshift-controller-manager",
		Short: "Command for the OpenShift Controllers",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(1)
		},
	}
	start := openshift_controller_manager.NewOpenShiftControllerManagerCommand("start", os.Stdout, os.Stderr)
	cmd.AddCommand(start)
	return cmd
}

func OCPControllerManager(cfg *config.MicroshiftConfig) {
	command := newOpenShiftControllerManagerCommand()
	args := []string{
		"--config=" + cfg.DataDir + "/resources/openshift-controller-manager/config/config.yaml",
	}
	startArgs := append(args, "start")
	command.SetArgs(startArgs)

	go func() {
		logrus.Fatalf("ocp controller-manager exited: %v", command.Execute())
	}()
}
