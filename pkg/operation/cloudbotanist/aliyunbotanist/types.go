package aliyunbotanist

import (
	"github.com/gardener/gardener/pkg/operation"
)

type AliyunBotanist struct {
	*operation.Operation
	CloudProviderName string
}

const (
	AliyunAccessKeyId = "accessKeyId"
    AliyunAccessKeySecret = "accessKeySecret"
	Region = "region"
)