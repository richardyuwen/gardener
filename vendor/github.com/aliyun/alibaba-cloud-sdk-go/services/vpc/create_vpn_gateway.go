package vpc

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// CreateVpnGateway invokes the vpc.CreateVpnGateway API synchronously
// api document: https://help.aliyun.com/api/vpc/createvpngateway.html
func (client *Client) CreateVpnGateway(request *CreateVpnGatewayRequest) (response *CreateVpnGatewayResponse, err error) {
	response = CreateCreateVpnGatewayResponse()
	err = client.DoAction(request, response)
	return
}

// CreateVpnGatewayWithChan invokes the vpc.CreateVpnGateway API asynchronously
// api document: https://help.aliyun.com/api/vpc/createvpngateway.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateVpnGatewayWithChan(request *CreateVpnGatewayRequest) (<-chan *CreateVpnGatewayResponse, <-chan error) {
	responseChan := make(chan *CreateVpnGatewayResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateVpnGateway(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// CreateVpnGatewayWithCallback invokes the vpc.CreateVpnGateway API asynchronously
// api document: https://help.aliyun.com/api/vpc/createvpngateway.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateVpnGatewayWithCallback(request *CreateVpnGatewayRequest, callback func(response *CreateVpnGatewayResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateVpnGatewayResponse
		var err error
		defer close(result)
		response, err = client.CreateVpnGateway(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// CreateVpnGatewayRequest is the request struct for api CreateVpnGateway
type CreateVpnGatewayRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	Period               requests.Integer `position:"Query" name:"Period"`
	AutoPay              requests.Boolean `position:"Query" name:"AutoPay"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	Bandwidth            requests.Integer `position:"Query" name:"Bandwidth"`
	EnableIpsec          requests.Boolean `position:"Query" name:"EnableIpsec"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	EnableSsl            requests.Boolean `position:"Query" name:"EnableSsl"`
	SslConnections       requests.Integer `position:"Query" name:"SslConnections"`
	VpcId                string           `position:"Query" name:"VpcId"`
	Name                 string           `position:"Query" name:"Name"`
	InstanceChargeType   string           `position:"Query" name:"InstanceChargeType"`
}

// CreateVpnGatewayResponse is the response struct for api CreateVpnGateway
type CreateVpnGatewayResponse struct {
	*responses.BaseResponse
	RequestId    string `json:"RequestId" xml:"RequestId"`
	VpnGatewayId string `json:"VpnGatewayId" xml:"VpnGatewayId"`
	Name         string `json:"Name" xml:"Name"`
	OrderId      int    `json:"OrderId" xml:"OrderId"`
}

// CreateCreateVpnGatewayRequest creates a request to invoke CreateVpnGateway API
func CreateCreateVpnGatewayRequest() (request *CreateVpnGatewayRequest) {
	request = &CreateVpnGatewayRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Vpc", "2016-04-28", "CreateVpnGateway", "vpc", "openAPI")
	return
}

// CreateCreateVpnGatewayResponse creates a response to parse from CreateVpnGateway response
func CreateCreateVpnGatewayResponse() (response *CreateVpnGatewayResponse) {
	response = &CreateVpnGatewayResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
