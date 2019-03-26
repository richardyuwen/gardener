// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package shoot_test

import (
	gardenv1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	. "github.com/gardener/gardener/pkg/operation/shoot"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("shoot", func() {

	var shoot *Shoot

	BeforeEach(func() {
		shoot = &Shoot{
			Info: &gardenv1beta1.Shoot{},
		}
	})

	Describe("#IPVSEnabled", func() {

		It("should return false when KubeProxy is null", func() {
			shoot.Info.Spec.Kubernetes.KubeProxy = nil

			Expect(shoot.IPVSEnabled()).To(BeFalse())
		})

		It("should return false when KubeProxy.Mode is null", func() {
			shoot.Info.Spec.Kubernetes.KubeProxy = &gardenv1beta1.KubeProxyConfig{}
			Expect(shoot.IPVSEnabled()).To(BeFalse())
		})

		It("should return false when KubeProxy.Mode is not IPVS", func() {
			mode := gardenv1beta1.ProxyModeIPTables
			shoot.Info.Spec.Kubernetes.KubeProxy = &gardenv1beta1.KubeProxyConfig{
				Mode: &mode,
			}
			Expect(shoot.IPVSEnabled()).To(BeFalse())
		})

		It("should return true when KubeProxy.Mode is IPVS", func() {
			mode := gardenv1beta1.ProxyModeIPVS
			shoot.Info.Spec.Kubernetes.KubeProxy = &gardenv1beta1.KubeProxyConfig{
				Mode: &mode,
			}
			Expect(shoot.IPVSEnabled()).To(BeTrue())
		})

	})

})
