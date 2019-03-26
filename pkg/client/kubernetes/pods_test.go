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

package kubernetes_test

import (
	. "github.com/gardener/gardener/pkg/client/kubernetes"
	mockcorev1 "github.com/gardener/gardener/pkg/mock/client-go/core/v1"
	mockrest "github.com/gardener/gardener/pkg/mock/client-go/rest"
	mockio "github.com/gardener/gardener/pkg/mock/go/io"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"net/http"
	"net/url"
)

var _ = Describe("Pods", func() {
	var (
		ctrl *gomock.Controller
		pods *mockcorev1.MockPodInterface
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		pods = mockcorev1.NewMockPodInterface(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("#GetPodLogs", func() {
		It("should read all pod logs and close the stream", func() {
			const name = "name"
			var (
				options = &corev1.PodLogOptions{}
				logs    = []byte("logs")
				client  = mockrest.NewMockHTTPClient(ctrl)
				body    = mockio.NewMockReadCloser(ctrl)
			)

			response := &http.Response{StatusCode: http.StatusOK, Body: body}

			gomock.InOrder(
				pods.EXPECT().GetLogs(name, options).Return(rest.NewRequest(client, http.MethodGet, &url.URL{}, "", rest.ContentConfig{}, rest.Serializers{}, nil, nil, 0)),
				client.EXPECT().Do(gomock.Any()).Return(response, nil),
				body.EXPECT().Read(gomock.Any()).DoAndReturn(func(data []byte) (int, error) {
					copy(data, logs)
					return len(logs), io.EOF
				}),
				body.EXPECT().Close(),
			)

			actual, err := GetPodLogs(pods, name, options.DeepCopy())
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(Equal(logs))
		})
	})
})
