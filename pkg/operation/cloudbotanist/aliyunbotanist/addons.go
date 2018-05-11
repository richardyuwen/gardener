package aliyunbotanist

import "github.com/gardener/gardener/pkg/operation/common"

func (c *AliyunBotanist) DeployKube2IAMResources() error {
	return nil
}

func (c *AliyunBotanist) DestroyKube2IAMResources() error {
	return nil
}

func (c *AliyunBotanist) GenerateKube2IAMConfig() (map[string]interface{}, error) {
	return common.GenerateAddonConfig(nil, false), nil
}

func (c *AliyunBotanist) GenerateNginxIngressConfig() (map[string]interface{}, error) {
	return common.GenerateAddonConfig(nil, c.Shoot.NginxIngressEnabled()), nil
}

// Cosine todo: need to add some pv logic here
func (c *AliyunBotanist) GenerateAdmissionControlConfig() (map[string]interface{}, error) {
	return nil, nil
}