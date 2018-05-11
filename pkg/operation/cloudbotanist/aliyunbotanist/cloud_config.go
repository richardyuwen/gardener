package aliyunbotanist

import "github.com/gardener/gardener/pkg/operation/common"

func (c *AliyunBotanist) GenerateCloudConfigUserDataConfig() *common.CloudConfigUserDataConfig {
	return &common.CloudConfigUserDataConfig{
		WorkerNames: c.Shoot.GetWorkerNames(),
		HostnameOverride: true, //Cosine todo: what's this
		ProvisionCloudProviderConfig: true, //Cosine todo: what's this
	}
}