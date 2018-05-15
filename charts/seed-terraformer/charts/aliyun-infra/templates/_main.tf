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
    cidr_block         = "{{ required "zone.cidr.worker is required" $zone.cidr.worker }}"
    availability_zone  = "{{ required "zone.name is required" $zone.name }}"
}

resource "alicloud_key_pair" "publickey" {
    key_name = "{{ required "clusterName is required" .Values.clusterName }}-ssh-publickey"
    public_key = "{{ required "sshPublicKey is required" .Values.sshPublicKey }}"
}

resource "alicloud_security_group" "rule-allow-internal-access" {
    name     = "{{ required "clusterName is required" .Values.clusterName }}-allow-internal-access"
    vpc_id   = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group" "rule-allow-external-access" {
    name     = "{{ required "clusterName is required" .Values.clusterName }}-allow-external-access"
    vpc_id   = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group_rule" "rule-allow-all-incoming-internal" {
    type              = "ingress"
    ip_protocol       = "all"
    port_range        = "1/65535" 
    nic_type          = "intranet"
    security_group_id = "${alicloud_security_group.rule-allow-internal-access.id}"
    cidr_ip           = "10.0.0.0/8"
}

resource "alicloud_security_group_rule" "rule-allow-http-incoming-external" {
    type              = "ingress"
    ip_protocol       = "tcp"
    port_range        = "80"
    nic_type          = "internet"
    security_group_id = "${alicloud_security_group.rule-allow-external-access.id}"
    cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_security_group_rule" "rule-allow-https-incoming-external" {
    type              = "ingress"
    ip_protocol       = "tcp"
    port_range        = "443"
    nic_type          = "internet"
    security_group_id = "${alicloud_security_group.rule-allow-external-access.id}"
    cidr_ip           = "0.0.0.0/0"
}



resource "aws_route_table" "routetable_main" {
  vpc_id = "{{ required "vpc.id is required" .Values.vpc.id }}"

{{ include "aws-infra.common-tags" .Values | indent 2 }}
}

resource "aws_route" "public" {
  route_table_id         = "${aws_route_table.routetable_main.id}"
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = "{{ required "vpc.internetGatewayID is required" .Values.vpc.internetGatewayID }}"
}

resource "aws_security_group" "bastions" {
  name        = "{{ required "clusterName is required" .Values.clusterName }}-bastions"
  description = "Security group for bastions"
  vpc_id      = "{{ required "vpc.id is required" .Values.vpc.id }}"

{{ include "aws-infra.tags-with-suffix" (set $.Values "suffix" "bastions") | indent 2 }}
}

resource "aws_security_group_rule" "bastion_ssh_bastion" {
  type              = "ingress"
  from_port         = 22
  to_port           = 22
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = "${aws_security_group.bastions.id}"
}

resource "aws_security_group_rule" "bastions_egress_all" {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = "${aws_security_group.bastions.id}"
}

resource "aws_security_group" "nodes" {
  name        = "{{ required "clusterName is required" .Values.clusterName }}-nodes"
  description = "Security group for nodes"
  vpc_id      = "{{ required "vpc.id is required" .Values.vpc.id }}"

{{ include "aws-infra.tags-with-suffix" (set $.Values "suffix" "nodes") }}
}

resource "aws_security_group_rule" "nodes_self" {
  type              = "ingress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  self              = true
  security_group_id = "${aws_security_group.nodes.id}"
}

resource "aws_security_group_rule" "nodes_ssh_bastion" {
  type                     = "ingress"
  from_port                = 22
  to_port                  = 22
  protocol                 = "tcp"
  security_group_id        = "${aws_security_group.nodes.id}"
  source_security_group_id = "${aws_security_group.bastions.id}"
}

resource "aws_security_group_rule" "nodes_egress_all" {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = "${aws_security_group.nodes.id}"
}

{{ range $index, $zone := .Values.zones }}
resource "aws_subnet" "nodes_z{{ $index }}" {
  vpc_id            = "{{ required "vpc.id is required" $.Values.vpc.id }}"
  cidr_block        = "{{ required "zone.cidr.worker is required" $zone.cidr.worker }}"
  availability_zone = "{{ required "zone.name is required" $zone.name }}"

{{ include "aws-infra.tags-with-suffix" (set $.Values "suffix" (print "nodes-z" $index)) }}
}

output "subnet_nodes_z{{ $index }}" {
  value = "${aws_subnet.nodes_z{{ $index }}.id}"
}

resource "aws_subnet" "private_utility_z{{ $index }}" {
  vpc_id            = "{{ required "vpc.id is required" $.Values.vpc.id }}"
  cidr_block        = "{{ required "zone.cidr.internal is required" $zone.cidr.internal }}"
  availability_zone = "{{ required "zone.name is required" $zone.name }}"

  tags {
    Name = "{{ required "clusterName is required" $.Values.clusterName }}-private-utility-z{{ $index }}"
    "kubernetes.io/cluster/{{ required "clusterName is required" $.Values.clusterName }}"  = "1"
    "kubernetes.io/role/internal-elb" = "use"
  }
}

resource "aws_security_group_rule" "nodes_tcp_internal_z{{ $index }}" {
  type              = "ingress"
  from_port         = 30000
  to_port           = 32767
  protocol          = "tcp"
  cidr_blocks       = ["{{ required "zone.cidr.internal is required" $zone.cidr.internal }}"]
  security_group_id = "${aws_security_group.nodes.id}"
}

resource "aws_security_group_rule" "nodes_udp_internal_z{{ $index }}" {
  type              = "ingress"
  from_port         = 30000
  to_port           = 32767
  protocol          = "udp"
  cidr_blocks       = ["{{ required "zone.cidr.internal is required" $zone.cidr.internal }}"]
  security_group_id = "${aws_security_group.nodes.id}"
}

resource "aws_subnet" "public_utility_z{{ $index }}" {
  vpc_id            = "{{ required "vpc.id is required" $.Values.vpc.id }}"
  cidr_block        = "{{ required "zone.cidr.public is required" $zone.cidr.public }}"
  availability_zone = "{{ required "zone.name is required" $zone.name }}"

  tags {
    Name = "{{ required "clusterName is required" $.Values.clusterName }}-public-utility-z{{ $index }}"
    "kubernetes.io/cluster/{{ required "clusterName is required" $.Values.clusterName }}"  = "1"
    "kubernetes.io/role/elb" = "use"
  }
}

output "subnet_public_utility_z{{ $index }}" {
  value = "${aws_subnet.public_utility_z{{ $index }}.id}"
}

resource "aws_security_group_rule" "nodes_tcp_public_z{{ $index }}" {
  type              = "ingress"
  from_port         = 30000
  to_port           = 32767
  protocol          = "tcp"
  cidr_blocks       = ["{{ required "zone.cidr.public is required" $zone.cidr.public }}"]
  security_group_id = "${aws_security_group.nodes.id}"
}

resource "aws_security_group_rule" "nodes_udp_public_z{{ $index }}" {
  type              = "ingress"
  from_port         = 30000
  to_port           = 32767
  protocol          = "udp"
  cidr_blocks       = ["{{ required "zone.cidr.public is required" $zone.cidr.public }}"]
  security_group_id = "${aws_security_group.nodes.id}"
}

resource "aws_eip" "eip_natgw_z{{ $index }}" {
  vpc = true
}

resource "aws_nat_gateway" "natgw_z{{ $index }}" {
  allocation_id = "${aws_eip.eip_natgw_z{{ $index }}.id}"
  subnet_id     = "${aws_subnet.public_utility_z{{ $index }}.id}"
}

resource "aws_route_table" "routetable_private_utility_z{{ $index }}" {
  vpc_id = "{{ required "vpc.id is required" $.Values.vpc.id }}"

{{ include "aws-infra.tags-with-suffix" (set $.Values "suffix" (print "private-" $zone.name)) }}
}

resource "aws_route" "private_utility_z{{ $index }}_nat" {
  route_table_id         = "${aws_route_table.routetable_private_utility_z{{ $index }}.id}"
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = "${aws_nat_gateway.natgw_z{{ $index }}.id}"
}

resource "aws_route_table_association" "routetable_private_utility_z{{ $index }}_association_private_utility_z{{ $index }}" {
  subnet_id      = "${aws_subnet.private_utility_z{{ $index }}.id}"
  route_table_id = "${aws_route_table.routetable_private_utility_z{{ $index }}.id}"
}

resource "aws_route_table_association" "routetable_main_association_public_utility_z{{ $index }}" {
  subnet_id      = "${aws_subnet.public_utility_z{{ $index }}.id}"
  route_table_id = "${aws_route_table.routetable_main.id}"
}

resource "aws_route_table_association" "routetable_private_utility_z{{ $index }}_association_nodes_z{{ $index }}" {
  subnet_id      = "${aws_subnet.nodes_z{{ $index }}.id}"
  route_table_id = "${aws_route_table.routetable_private_utility_z{{ $index }}.id}"
}
{{end}}

//=====================================================================
//= IAM instance profiles
//=====================================================================

resource "aws_iam_role" "bastions" {
  name = "{{ required "clusterName is required" .Values.clusterName }}-bastions"
  path = "/"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_instance_profile" "bastions" {
  name = "{{ required "clusterName is required" .Values.clusterName }}-bastions"
  role = "${aws_iam_role.bastions.name}"
}

resource "aws_iam_role_policy" "bastions" {
  name = "{{ required "clusterName is required" .Values.clusterName }}-bastions"
  role = "${aws_iam_role.bastions.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:DescribeRegions"
      ],
      "Resource": [
        "*"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_role" "nodes" {
  name = "{{ required "clusterName is required" .Values.clusterName }}-nodes"
  path = "/"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_instance_profile" "nodes" {
  name = "{{ required "clusterName is required" .Values.clusterName }}-nodes"
  role = "${aws_iam_role.nodes.name}"
}

resource "aws_iam_role_policy" "nodes" {
  name = "{{ required "clusterName is required" .Values.clusterName }}-nodes"
  role = "${aws_iam_role.nodes.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:Describe*"
      ],
      "Resource": [
        "*"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:GetRepositoryPolicy",
        "ecr:DescribeRepositories",
        "ecr:ListImages",
        "ecr:BatchGetImage"
      ],
      "Resource": [
        "*"
      ]
    }
  ]
}
EOF
}

//=====================================================================
//= EC2 Key Pair
//=====================================================================

resource "aws_key_pair" "kubernetes" {
  key_name   = "{{ required "clusterName is required" .Values.clusterName }}-ssh-publickey"
  public_key = "{{ required "sshPublicKey is required" .Values.sshPublicKey }}"
}

//=====================================================================
//= Output variables
//=====================================================================

output "vpc_id" {
  value = "{{ required "vpc.id is required" .Values.vpc.id }}"
}

output "iamInstanceProfileNodes" {
  value = "${aws_iam_instance_profile.nodes.name}"
}

output "keyName" {
  value = "${aws_key_pair.kubernetes.key_name}"
}

output "security_group_nodes" {
  value = "${aws_security_group.nodes.id}"
}

output "nodes_role_arn" {
  value = "${aws_iam_role.nodes.arn}"
}
{{- end -}}


{{- define "aws-infra.common-tags" -}}
tags {
  Name = "{{ required "clusterName is required" .clusterName }}"
  "kubernetes.io/cluster/{{ required "clusterName is required" .clusterName }}" = "1"
}
{{- end -}}
{{- define "aws-infra.tags-with-suffix" -}}
tags {
  Name = "{{ required "clusterName is required" .clusterName }}-{{ required "suffix is required" .suffix }}"
  "kubernetes.io/cluster/{{ required "clusterName is required" .clusterName }}" = "1"
}
{{- end -}}
