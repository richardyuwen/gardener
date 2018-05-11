package aliyunbotanist

import (
	"github.com/gardener/gardener/pkg/operation/common"
	"github.com/gardener/gardener/pkg/operation/terraformer"
)

//Cosine todo: need to hardcode some cloud config things
func (c *AliyunBotanist) GenerateCloudProviderConfig() (string, error) {
	return "cloud_config", nil
}

func (c *AliyunBotanist) GenerateKubeAPIServerConfig() (map[string]interface{}, error) {
	return map[string]interface{}{
		"environment": getAliyunCredentialEnvironment(),
	}, nil
}

func (c *AliyunBotanist) GenerateKubeControllerManagerConfig() (map[string]interface{}, error) {
	return map[string]interface{}{
		"configureRoutes": false,
		"environment":     getAliyunCredentialEnvironment(),
	}, nil
}

func (c *AliyunBotanist) GenerateKubeSchedulerConfig() (map[string]interface{}, error) {
	return nil, nil
}

func getAliyunCredentialEnvironment() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name": "ALIYUN_ACCESS_KEY_ID",
			"valueFrom": map[string]interface{}{
				"secretKeyRef": map[string]interface{}{
					"key":  AliyunAccessKeyId,
					"name": "cloudprovider",
				},
			},
		},
		{
			"name": "ALIYUN_ACCESS_KEY_SECRET",
			"valueFrom": map[string]interface{}{
				"secretKeyRef": map[string]interface{}{
					"key":  AliyunAccessKeySecret,
					"name": "cloudprovider",
				},
			},
		},
	}
}

//Cosine todo: need to use Aliyun OSS for this
func (c *AliyunBotanist) GenerateEtcdBackupConfig() (map[string][]byte, map[string]interface{}, error) {
	ossbucketName := "bucketName"
	stateVariables, err := terraformer.NewFromOperation(c.Operation, common.TerraformerPurposeBackup).GetStateOutputVariables(ossbucketName)
	if err != nil {
		return nil, nil, err
	}
	secretData := map[string][]byte{
		Region:          []byte(c.Seed.Info.Spec.Cloud.Region),
		AccessKeyID:     c.Seed.Secret.Data[AliyunAccessKeyId],
		AccessKeySecret: c.Seed.Secret.Data[AliyunAccessKeySecret],
	}

	backupConfigData := map[string]interface{}{
		"schedule" : c.Shoot.Info.Spec.Backup.Schedule,
		"maxBackups" : c.Shoot.Info.Spec.Backup.Maximum,
		"backupSecret" : common.BackupSecretName,
		"storageContainer" : stateVariables[ossbucketName],
		"env" : []map[string]interface{} {
            {
				"name": "ALIYUN_REGION",
				"valueFrom": map[string]interface{} {
                    "secretKeyRef": map[string]interface{} {
						"name": common.BackupSecretName,
						"key": Region,
					},
				},
			},
		    {
				"name": "ALIYUN_ACCESS_KEY_ID",
				"valueFrom": map[string]interface{} {
                    "secretKeyRef": map[string]interface{} {
						"name": common.BackupSecretName,
						"key": AliyunAccessKeyId,
					},
				},
			},
			{
                "name": "ALIYUN_ACCESS_KEY_SECRET",
				"valueFrom": map[string]interface{} {
                    "secretKeyRef": map[string]interface{} {
						"name": common.BackupSecretName,
						"key": AliyunAccessKeySecret,
					},
				},
			},
		},
		"volumeMount" : []map[string]interface{}{},
	}

	return secretData, backupConfigData, nil
}
