{{- define "aliyun-infra.variables" -}}
variable "ALIYUN_ACCESS_KEY_ID" {
  description = "Aliyun Access Key ID of technical user"
  type        = "string"
}

variable "ALIYUN_ACCESS_KEY_SECRET" {
  description = "Aliyun Access Key Secret of technical user"
  type        = "string"
}
{{- end -}}
