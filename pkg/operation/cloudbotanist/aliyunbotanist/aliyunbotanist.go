package aliyunbotanist

import (
	//gardenv1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"github.com/gardener/gardener/pkg/operation"
	//"github.com/gardener/gardener/pkg/operation/common"
)

func New(o *operation.Operation, purpose string) (*AliyunBotanist, error) {

	//Cosine todo: need to enforce the check once apis is updated.
	/*
		var cloudProvider gardenv1beta1.CloudProvider
		switch purpose {
		case common.CloudPurposeShoot:
			cloudProvider = o.Shoot.CloudProvider
		case common.CloudPurposeSeed:
			cloudProvider = o.Seed.CloudProvider
		}


		if cloudProvider != gardenv1beta1.CloudProviderAliyun {
			return nil, errors.New("cannot instantiate an Aliyun botanist if neither Shoot nor Seed cluster specifies Aliyun")
		}
	*/

	return &AliyunBotanist{
		Operation:         o,
		CloudProviderName: "aliyun",
	}, nil
}

func (c *AliyunBotanist) GetCloudProviderName() string {
	return c.CloudProviderName
}
