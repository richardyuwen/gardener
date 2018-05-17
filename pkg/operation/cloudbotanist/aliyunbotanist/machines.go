package aliyunbotanist

import (
	"github.com/gardener/gardener/pkg/operation/terraformer"
	"github.com/gardener/gardener/pkg/operation"
	"github.com/gardener/gardener/pkg/operation/common"
	"fmt"
)

func (c *AliyunBotanist) GetMachineClassInfo() (classKind, classPlural, classChartName string) {
	classKind = "AliyunMachineClass"
	classPlural = "aliyunmachineclasses"
	classChartName = "aliyun-machineclass"
	return
}

// Cosine todo: need to add the machine spec generation logic
func (c *AliyunBotanist) GenerateMachineConfig() ([]map[string]interface{}, []operation.MachineDeployment, error) {
    var (
		machineDeployments = []operation.MachineDeployment{}
		machineClasses = []map[string]interface{}{}

		workers = c.Shoot.Info.Spec.Cloud.Aliyun.Workers
		zones = c.Shoot.Info.Spec.Cloud.Aliyun.Zones
		zoneLen = len(zones)

		vSwitchId = "vSwitchId"
		securityGroupId = "securityGroupId"
		keyPairName = "keyPairName"
		outputVariables = []string{vSwitchId, securityGroupId, keyPairName}
	)

	stateVariable, err := terraformer.NewFromOperation(c.Operation, common.TerraformerPurposeInfra).GetStateOutputVariables(outputVariables...)
	if err != nil {
		return nil, nil, err
	}

	for zoneIndex, zone := range zones {
		for _, worker := range workers {
			cloudConfig, err := c.ComputeDownloaderCloudConfig(worker.Name)
			if err != nil {
				return nil, nil, err
			}

			machineClassSpec := map[string]interface{} {
				"region": c.Shoot.Info.Spec.Cloud.Region,
				"instanceType": worker.MachineType,
				"keyPairName": stateVariable[keyPairName],
				"vSwitchId" : stateVariable[vSwitchId],
				"securityGroupId" : stateVariable[securityGroupId],
				"imageId": c.Shoot.Info.Spec.Cloud.Aliyun.MachineImage.Image,
                "tags": map[string]string {
					"tag1Key": fmt.Sprintf("kubernetes.io/cluster/%s", c.Shoot.SeedNamespace),
					"tag1Value": "1",
					"tag2Key": "kubernetes.io/role/node",
					"tag2Value": "1",
				},
				"secret": map[string]interface{} {
                    "cloudConfig": cloudConfig.FileContent("cloud-config.yaml"),
				},
				"systemDisk": map[string]interface{} {
					"category": worker.VolumeType, //cloud_efficiency, cloud_ssd
					"size": common.DiskSize(worker.VolumeSize),
				},
			}

			if worker.InternetMaxBandwidthIn != nil {
				machineClassSpec["internetMaxBandwidthIn"] = worker.InternetMaxBandwidthIn
			}

			if worker.InternetMaxBandwidthOut != nil {
				machineClassSpec["internetMaxBandwidthOut"] = worker.InternetMaxBandwidthOut
			}

			var (
				machineClassSpecHash = common.MachineClassHash(machineClassSpec, c.Shoot.KubernetesMajorMinorVersion)
				deploymentName = fmt.Sprintf("%s-%s-z%d", c.Shoot.SeedNamespace, worker.Name, zoneIndex+1)
				className = fmt.Sprintf("%s-%s", deploymentName, machineClassSpecHash)
			)

			machineDeployments = append(machineDeployments, operation.MachineDeployment{
				Name: deploymentName,
				ClassName: className,
				Replicas: common.DistributeOverZones(zoneIndex, worker.AutoScalerMax, zoneLen),
			})

			machineClassSpec["name"] = className
			machineClassSpec["secret"].(map[string]interface{})[AliyunAccessKeyId] = string(c.Shoot.Secret.Data[AliyunAccessKeyId])
			machineClassSpec["secret"].(map[string]interface{})[AliyunAccessKeySecret] = string(c.Shoot.Secret.Data[AliyunAccessKeySecret])

			machineClasses = append(machineClasses, machineClassSpec)
		}
	}

	return machineClasses, machineDeployments, nil
}