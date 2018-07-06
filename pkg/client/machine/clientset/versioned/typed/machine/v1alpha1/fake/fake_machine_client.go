/*
Copyright 2018 (c) 2018 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/gardener/gardener/pkg/client/machine/clientset/versioned/typed/machine/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeMachineV1alpha1 struct {
	*testing.Fake
}

func (c *FakeMachineV1alpha1) AWSMachineClasses(namespace string) v1alpha1.AWSMachineClassInterface {
	return &FakeAWSMachineClasses{c, namespace}
}

func (c *FakeMachineV1alpha1) AzureMachineClasses(namespace string) v1alpha1.AzureMachineClassInterface {
	return &FakeAzureMachineClasses{c, namespace}
}

func (c *FakeMachineV1alpha1) GCPMachineClasses(namespace string) v1alpha1.GCPMachineClassInterface {
	return &FakeGCPMachineClasses{c, namespace}
}

func (c *FakeMachineV1alpha1) Machines(namespace string) v1alpha1.MachineInterface {
	return &FakeMachines{c, namespace}
}

func (c *FakeMachineV1alpha1) MachineDeployments(namespace string) v1alpha1.MachineDeploymentInterface {
	return &FakeMachineDeployments{c, namespace}
}

func (c *FakeMachineV1alpha1) MachineSets(namespace string) v1alpha1.MachineSetInterface {
	return &FakeMachineSets{c, namespace}
}

func (c *FakeMachineV1alpha1) MachineTemplates(namespace string) v1alpha1.MachineTemplateInterface {
	return &FakeMachineTemplates{c, namespace}
}

func (c *FakeMachineV1alpha1) OpenStackMachineClasses(namespace string) v1alpha1.OpenStackMachineClassInterface {
	return &FakeOpenStackMachineClasses{c, namespace}
}

func (c *FakeMachineV1alpha1) Scales(namespace string) v1alpha1.ScaleInterface {
	return &FakeScales{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeMachineV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
