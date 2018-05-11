package aliyunbotanist

import (
	"github.com/gardener/gardener/pkg/operation"
)

func (c *AliyunBotanist) GetMachineClassInfo() (classKind, classPlural, classChartName string) {
	classKind = "AliyunMachineClass"
	classPlural = "aliyunmachineclasses"
	classChartName = "aliyun-machineclass"
	return
}

// Cosine todo: need to add the machine spec generation logic
func (c *AliyunBotanist) GenerateMachineConfig() ([]map[string]interface{}, []operation.MachineDeployment, error) {
	return nil, nil, nil
}