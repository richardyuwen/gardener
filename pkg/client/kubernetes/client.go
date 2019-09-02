// Copyright (c) 2018 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubernetes

import (
	"context"
	"errors"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	gardencoreclientset "github.com/gardener/gardener/pkg/client/core/clientset/versioned"
	gardenclientset "github.com/gardener/gardener/pkg/client/garden/clientset/versioned"
	machineclientset "github.com/gardener/gardener/pkg/client/machine/clientset/versioned"
	"github.com/gardener/gardener/pkg/utils"
	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"
	corev1 "k8s.io/api/core/v1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kubernetesclientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	apiserviceclientset "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
)

// KubeConfig is the key to the kubeconfig
const KubeConfig = "kubeconfig"

// NewRuntimeClientFromSecret creates a new controller runtime Client struct for a given secret.
func NewRuntimeClientFromSecret(secret *corev1.Secret, fns ...ConfigFunc) (client.Client, error) {
	if kubeconfig, ok := secret.Data[KubeConfig]; ok {
		return NewRuntimeClientFromBytes(kubeconfig, fns...)
	}
	return nil, errors.New("no valid kubeconfig found")

}

// NewRuntimeClientFromBytes creates a new controller runtime Client struct for a given kubeconfig byte slice.
func NewRuntimeClientFromBytes(kubeconfig []byte, fns ...ConfigFunc) (client.Client, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	opts := append([]ConfigFunc{WithRESTConfig(config)}, fns...)
	return NewRuntimeClientForConfig(opts...)
}

// NewRuntimeClientForConfig returns a new controller runtime client from a config.
func NewRuntimeClientForConfig(fns ...ConfigFunc) (client.Client, error) {
	conf := &config{}
	for _, f := range fns {
		if err := f(conf); err != nil {
			return nil, err
		}
	}
	return client.New(conf.restConfig, conf.clientOptions)
}

// NewClientFromFile creates a new Client struct for a given kubeconfig. The kubeconfig will be
// read from the filesystem at location <kubeconfigPath>. If given, <masterURL> overrides the
// master URL in the kubeconfig.
// If no filepath is given, the in-cluster configuration will be taken into account.
func NewClientFromFile(masterURL, kubeconfigPath string, fns ...ConfigFunc) (Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		return nil, err
	}

	opts := append([]ConfigFunc{WithRESTConfig(config)}, fns...)
	return NewWithConfig(opts...)
}

// NewClientFromBytes creates a new Client struct for a given kubeconfig byte slice.
func NewClientFromBytes(kubeconfig []byte, fns ...ConfigFunc) (Interface, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	opts := append([]ConfigFunc{WithRESTConfig(config)}, fns...)
	return NewWithConfig(opts...)
}

// NewClientFromSecret creates a new Client struct for a given kubeconfig stored as a
// Secret in an existing Kubernetes cluster. This cluster will be accessed by the <k8sClient>. It will
// read the Secret <secretName> in <namespace>. The Secret must contain a field "kubeconfig" which will
// be used.
func NewClientFromSecret(k8sClient Interface, namespace, secretName string, fns ...ConfigFunc) (Interface, error) {
	secret := &corev1.Secret{}
	if err := k8sClient.Client().Get(context.TODO(), kutil.Key(namespace, secretName), secret); err != nil {
		return nil, err
	}
	return NewClientFromSecretObject(secret, fns...)
}

// NewClientFromSecretObject creates a new Client struct for a given Kubernetes Secret object. The Secret must
// contain a field "kubeconfig" which will be used.
func NewClientFromSecretObject(secret *corev1.Secret, fns ...ConfigFunc) (Interface, error) {
	if kubeconfig, ok := secret.Data[KubeConfig]; ok {
		return NewClientFromBytes(kubeconfig, fns...)
	}
	return nil, errors.New("the secret does not contain a field with name 'kubeconfig'")
}

var supportedKubernetesVersions = []string{
	"1.10",
	"1.11",
	"1.12",
	"1.13",
	"1.14",
	"1.15",
}

func checkIfSupportedKubernetesVersion(gitVersion string) error {
	for _, supportedVersion := range supportedKubernetesVersions {
		ok, err := utils.CompareVersions(gitVersion, "~", supportedVersion)
		if err != nil {
			return err
		}

		if ok {
			return nil
		}
	}
	return fmt.Errorf("unsupported kubernetes version %q", gitVersion)
}

// NewWithConfig returns a new Kubernetes base client.
func NewWithConfig(fns ...ConfigFunc) (Interface, error) {
	conf := &config{}

	for _, f := range fns {
		if err := f(conf); err != nil {
			return nil, err
		}
	}

	return new(conf)
}

func new(conf *config) (Interface, error) {
	c, err := client.New(conf.restConfig, conf.clientOptions)
	if err != nil {
		return nil, err
	}

	applier, err := NewApplierForConfig(conf.restConfig)
	if err != nil {
		return nil, err
	}

	kubernetes, err := kubernetesclientset.NewForConfig(conf.restConfig)
	if err != nil {
		return nil, err
	}

	garden, err := gardenclientset.NewForConfig(conf.restConfig)
	if err != nil {
		return nil, err
	}

	gardenCore, err := gardencoreclientset.NewForConfig(conf.restConfig)
	if err != nil {
		return nil, err
	}

	machine, err := machineclientset.NewForConfig(conf.restConfig)
	if err != nil {
		return nil, err
	}

	apiRegistration, err := apiserviceclientset.NewForConfig(conf.restConfig)
	if err != nil {
		return nil, err
	}

	apiExtension, err := apiextensionsclientset.NewForConfig(conf.restConfig)
	if err != nil {
		return nil, err
	}

	clientSet := &Clientset{
		config:     conf.restConfig,
		restMapper: conf.clientOptions.Mapper,
		restClient: kubernetes.Discovery().RESTClient(),

		applier: applier,

		client: c,

		kubernetes:      kubernetes,
		garden:          garden,
		gardenCore:      gardenCore,
		machine:         machine,
		apiregistration: apiRegistration,
		apiextension:    apiExtension,
	}

	serverVersion, err := clientSet.kubernetes.Discovery().ServerVersion()
	if err != nil {
		return nil, err
	}

	if err := checkIfSupportedKubernetesVersion(serverVersion.GitVersion); err != nil {
		return nil, err
	}
	clientSet.version = serverVersion.GitVersion

	return clientSet, nil
}
