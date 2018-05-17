{{- define "aliyun-infra.main" -}}
provider "alicloud" {
  access_key = "${var.ALIYUN_ACCESS_KEY_ID}"
  secret_key = "${var.ALIYUN_ACCESS_KEY_SECRET}"
  region     = "{{ required "alicloud.region is required" .Values.alicloud.region }}"
}

//=====================================================================
//= VPC, Subnets, Security Groups, Gateways
//=====================================================================

{{ if .Values.create.vpc -}}
resource "alicloud_vpc" "vpc" {
  name                 = "{{ required "clusterName is required" .Values.clusterName }}"
  cidr_block           = "{{ required "vpc.cidr is required" .Values.vpc.cidr }}"
}
{{- end}}

resource "alicloud_vswitch" "vsw" {
    name               = " {{ required "clusterName is required" .Values.clusterName }} "
    vpc_id             = "${alicloud_vpc.vpc.id}"
    cidr_block         = "{{ required "zone.cidr is required" $zone.cidr }}"
    availability_zone  = "{{ required "zone.name is required" $zone.name }}"
}

resource "alicloud_key_pair" "publickey" {
    key_name = "{{ required "clusterName is required" .Values.clusterName }}-ssh-publickey"
    public_key = "{{ required "sshPublicKey is required" .Values.sshPublicKey }}"
}

resource "alicloud_security_group" "rule-allow-various-access" {
    name     = "{{ required "clusterName is required" .Values.clusterName }}-allow-internal-access"
    vpc_id   = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group_rule" "rule-allow-all-incoming-internal" {
    type              = "ingress"
    ip_protocol       = "all"
    port_range        = "-1/-1" 
    nic_type          = "intranet"
    security_group_id = "${alicloud_security_group.rule-allow-various-access.id}"
    cidr_ip           = "{{ required "vpc.cidr is required" .Values.vpc.cidr }}"
}

resource "alicloud_security_group_rule" "rule-allow-http-incoming-external" {
    type              = "ingress"
    ip_protocol       = "tcp"
    port_range        = "80/80"
    security_group_id = "${alicloud_security_group.rule-allow-various-access.id}"
    cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_security_group_rule" "rule-allow-https-incoming-external" {
    type              = "ingress"
    ip_protocol       = "tcp"
    port_range        = "443/443"
    security_group_id = "${alicloud_security_group.rule-allow-various-access.id}"
    cidr_ip           = "0.0.0.0/0"
}


output "vSwitchId" {
    value = "${alicloud_vswitch.vsw.id}"
}

output "securityGroupId" {
    value = "${alicloud_security_group.rule-allow-various-access.id}"
}

output "keyPairName" {
    value = "${alicloud_key_pair.publickey.key_name}"
}