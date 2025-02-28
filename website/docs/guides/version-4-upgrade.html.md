---
subcategory: ""
layout: "aws"
page_title: "Terraform AWS Provider Version 4 Upgrade Guide"
description: |-
  Terraform AWS Provider Version 4 Upgrade Guide
---

# Terraform AWS Provider Version 4 Upgrade Guide

Version 4.0.0 of the AWS provider for Terraform is a major release and includes some changes that you will need to consider when upgrading. This guide is intended to help with that process and focuses only on changes from version 3.X to version 4.0.0. See the [Version 3 Upgrade Guide](/docs/providers/aws/guides/version-3-upgrade.html) for information about upgrading from 1.X to version 3.0.0.

Most of the changes outlined in this guide have been previously marked as deprecated in the Terraform plan/apply output throughout previous provider releases. These changes, such as deprecation notices, can always be found in the [Terraform AWS Provider CHANGELOG](https://github.com/hashicorp/terraform-provider-aws/blob/main/CHANGELOG.md).

~> **NOTE:** Version 4.0.0 of the AWS Provider will be the last major version to support [EC2-Classic resources](#ec2-classic-resource-and-data-source-support) as AWS plans to fully retire EC2-Classic Networking. See the [AWS News Blog](https://aws.amazon.com/blogs/aws/ec2-classic-is-retiring-heres-how-to-prepare/) for additional details.

~> **NOTE:** Version 4.0.0 and 4.x.x versions of the AWS Provider will be the last versions compatible with Terraform 0.12-0.15.

Upgrade topics:

<!-- TOC depthFrom:2 depthTo:2 -->

- [Provider Version Configuration](#provider-version-configuration)
- [Full Resource Lifecycle of Default Resources](#full-resource-lifecycle-of-default-resources)
    - [Resource: aws_default_subnet](#resource-aws_default_subnet)
    - [Resource: aws_default_vpc](#resource-aws_default_vpc)
- [Plural Data Source Behavior](#plural-data-source-behavior)
- [Data Source: aws_cloudwatch_log_group](#data-source-aws_cloudwatch_log_group)
- [Data Source: aws_subnet_ids](#data-source-aws_subnet_ids)
- [Data Source: aws_s3_bucket_object](#data-source-aws_s3_bucket_object)
- [Data Source: aws_s3_bucket_objects](#data-source-aws_s3_bucket_objects)
- [Resource: aws_batch_compute_environment](#resource-aws_batch_compute_environment)
- [Resource: aws_cloudwatch_event_target](#resource-aws_cloudwatch_event_target)
- [Resource: aws_customer_gateway](#resource-aws_customer_gateway)
- [Resource: aws_default_network_acl](#resource-aws_default_network_acl)
- [Resource: aws_elasticache_cluster](#resource-aws_elasticache_cluster)
- [Resource: aws_elasticache_global_replication_group](#resource-aws_elasticache_global_replication_group)
- [Resource: aws_elasticache_replication_group](#resource-aws_elasticache_replication_group)
- [Resource: aws_fsx_ontap_storage_virtual_machine](#resource-aws_fsx_ontap_storage_virtual_machine)
- [Resource: aws_network_acl](#resource-aws_network_acl)
- [Resource: aws_network_interface](#resource-aws_network_interface)
- [Resource: aws_s3_bucket](#resource-aws_s3_bucket)
- [Resource: aws_s3_bucket_object](#resource-aws_s3_bucket_object)
- [Resource: aws_spot_instance_request](#resource-aws_spot_instance_request)

<!-- /TOC -->

Additional Topics:

<!-- TOC depthFrom:2 depthTo:2 -->

- [EC2-Classic resource and data source support](#ec2-classic-resource-and-data-source-support)

<!-- /TOC -->


## Provider Version Configuration

!> **WARNING:** This topic is placeholder documentation until version 4.0.0 is released.

-> Before upgrading to version 4.0.0, it is recommended to upgrade to the most recent 3.X version of the provider and ensure that your environment successfully runs [`terraform plan`](https://www.terraform.io/docs/commands/plan.html) without unexpected changes or deprecation notices.

It is recommended to use [version constraints when configuring Terraform providers](https://www.terraform.io/docs/configuration/providers.html#provider-versions). If you are following that recommendation, update the version constraints in your Terraform configuration and run [`terraform init`](https://www.terraform.io/docs/commands/init.html) to download the new version.

For example, given this previous configuration:

```terraform
provider "aws" {
  # ... other configuration ...

  version = "~> 3.74"
}
```

Update to latest 4.X version:

```terraform
provider "aws" {
  # ... other configuration ...

  version = "~> 4.0"
}
```

## Plural Data Source Behavior

The following plural data sources are now consistent with [Provider Design](https://github.com/hashicorp/terraform-provider-aws/blob/main/docs/contributing/provider-design.md#data-sources)
such that they no longer return an error if zero results are found.

* [aws_cognito_user_pools](/docs/providers/aws/d/cognito_user_pools.html)
* [aws_db_event_categories](/docs/providers/aws/d/db_event_categories.html)
* [aws_ebs_volumes](/docs/providers/aws/d/ebs_volumes.html)
* [aws_ec2_coip_pools](/docs/providers/aws/d/ec2_coip_pools.html)
* [aws_ec2_local_gateway_route_tables](/docs/providers/aws/d/ec2_local_gateway_route_tables.html)
* [aws_ec2_local_gateway_virtual_interface_groups](/docs/providers/aws/d/ec2_local_gateway_virtual_interface_groups.html)
* [aws_ec2_local_gateways](/docs/providers/aws/d/ec2_local_gateways.html)
* [aws_ec2_transit_gateway_route_tables](/docs/providers/aws/d/ec2_transit_gateway_route_tables.html)
* [aws_efs_access_points](/docs/providers/aws/d/efs_access_points.html)
* [aws_emr_release_labels](/docs/providers/aws/d/emr_release_labels.markdown)
* [aws_inspector_rules_packages](/docs/providers/aws/d/inspector_rules_packages.html)
* [aws_ip_ranges](/docs/providers/aws/d/ip_ranges.html)
* [aws_network_acls](/docs/providers/aws/d/network_acls.html)
* [aws_route_tables](/docs/providers/aws/d/route_tables.html)
* [aws_security_groups](/docs/providers/aws/d/security_groups.html)
* [aws_ssoadmin_instances](/docs/providers/aws/d/ssoadmin_instances.html)
* [aws_vpcs](/docs/providers/aws/d/vpcs.html)
* [aws_vpc_peering_connections](/docs/providers/aws/d/vpc_peering_connections.html)

## Full Resource Lifecycle of Default Resources

Default subnets and vpcs can now do full resource lifecycle operations such that resource
creation and deletion are now supported.

### Resource: aws_default_subnet

The `aws_default_subnet` resource behaves differently from normal resources in that if a default subnet exists in the specified Availability Zone, Terraform does not _create_ this resource, but instead "adopts" it into management.
If no default subnet exists, Terraform creates a new default subnet.
By default, `terraform destroy` does not delete the default subnet but does remove the resource from Terraform state.
Set the `force_destroy` argument to `true` to delete the default subnet.

For example, given this previous configuration with no existing default subnet:

```terraform
terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
  required_version = ">= 0.13"
}

provider "aws" {
  region = "eu-west-2"
}

resource "aws_default_subnet" "default" {}
```

The following error was thrown on `terraform apply`:

```
│ Error: Default subnet not found.
│
│   with aws_default_subnet.default,
│   on main.tf line 5, in resource "aws_default_subnet" "default":
│    5: resource "aws_default_subnet" "default" {}
```

Now after upgrading, the above configuration will apply successfully.

To delete the default subnet, the above configuration should be updated to:

```terraform
terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
  required_version = ">= 0.13"
}

resource "aws_default_subnet" "default" {
  force_destroy = true
}
```

### Resource: aws_default_vpc

The `aws_default_vpc` resource behaves differently from normal resources in that if a default VPC exists, Terraform does not _create_ this resource, but instead "adopts" it into management.
If no default VPC exists, Terraform creates a new default VPC, which leads to the implicit creation of [other resources](https://docs.aws.amazon.com/vpc/latest/userguide/default-vpc.html#default-vpc-components).
By default, `terraform destroy` does not delete the default VPC but does remove the resource from Terraform state.
Set the `force_destroy` argument to `true` to delete the default VPC.

For example, given this previous configuration with no existing default VPC:

```terraform
terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
  required_version = ">= 0.13"
}

resource "aws_default_vpc" "default" {}
```

The following error was thrown on `terraform apply`:

```
│ Error: No default VPC found in this region.
│
│   with aws_default_vpc.default,
│   on main.tf line 5, in resource "aws_default_vpc" "default":
│    5: resource "aws_default_vpc" "default" {}
```

Now after upgrading, the above configuration will apply successfully.

To delete the default VPC, the above configuration should be updated to:

```terraform
terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
  required_version = ">= 0.13"
}

resource "aws_default_vpc" "default" {
  force_destroy = true
}
```

## Data Source: aws_cloudwatch_log_group

### Removal of arn Wildcard Suffix

Previously, the data source returned the Amazon Resource Name (ARN) directly from the API, which included a `:*` suffix to denote all CloudWatch Log Streams under the CloudWatch Log Group. Most other AWS resources that return ARNs and many other AWS services do not use the `:*` suffix. The suffix is now automatically removed. For example, the data source previously returned an ARN such as `arn:aws:logs:us-east-1:123456789012:log-group:/example:*` but will now return `arn:aws:logs:us-east-1:123456789012:log-group:/example`.

Workarounds, such as using `replace()` as shown below, should be removed:

```terraform
data "aws_cloudwatch_log_group" "example" {
  name = "example"
}
resource "aws_datasync_task" "example" {
  # ... other configuration ...
  cloudwatch_log_group_arn = replace(data.aws_cloudwatch_log_group.example.arn, ":*", "")
}
```

Removing the `:*` suffix is a breaking change for some configurations. Fix these configurations using string interpolations as demonstrated below. For example, this configuration is now broken:

```terraform
data "aws_iam_policy_document" "ad-log-policy" {
  statement {
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    principals {
      identifiers = ["ds.amazonaws.com"]
      type        = "Service"
    }
    resources = [data.aws_cloudwatch_log_group.example.arn]
    effect = "Allow"
  }
}
```

An updated configuration:

```terraform
data "aws_iam_policy_document" "ad-log-policy" {
  statement {
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    principals {
      identifiers = ["ds.amazonaws.com"]
      type        = "Service"
    }
    resources = ["${data.aws_cloudwatch_log_group.example.arn}:*"]
    effect = "Allow"
  }
}
```

## Data Source: aws_subnet_ids

The `aws_subnet_ids` data source has been deprecated and will be removed removed in a future version. Use the `aws_subnets` data source instead.

For example, change a configuration such as

```hcl
data "aws_subnet_ids" "example" {
  vpc_id = var.vpc_id
}

data "aws_subnet" "example" {
  for_each = data.aws_subnet_ids.example.ids
  id       = each.value
}

output "subnet_cidr_blocks" {
  value = [for s in data.aws_subnet.example : s.cidr_block]
}
```

to

```hcl
data "aws_subnets" "example" {
  filter {
    name   = "vpc-id"
    values = [var.vpc_id]
  }
}

data "aws_subnet" "example" {
  for_each = data.aws_subnets.example.ids
  id       = each.value
}

output "subnet_cidr_blocks" {
  value = [for s in data.aws_subnet.example : s.cidr_block]
}
```

## Data Source: aws_s3_bucket_object

The `aws_s3_bucket_object` data source is deprecated and will be removed in a future version. Use `aws_s3_object` instead, where new features and fixes will be added.

## Data Source: aws_s3_bucket_objects

The `aws_s3_bucket_objects` data source is deprecated and will be removed in a future version. Use `aws_s3_objects` instead, where new features and fixes will be added.

## Resource: aws_batch_compute_environment

No `compute_resources` can be specified when `type` is `UNMANAGED`.

Previously a configuration such as

```hcl
resource "aws_batch_compute_environment" "test" {
  compute_environment_name = "test"

  compute_resources {
    instance_role = aws_iam_instance_profile.ecs_instance.arn
    instance_type = [
      "c4.large",
    ]
    max_vcpus = 16
    min_vcpus = 0
    security_group_ids = [
      aws_security_group.test.id
    ]
    subnets = [
      aws_subnet.test.id
    ]
    type = "EC2"
  }

  service_role = aws_iam_role.batch_service.arn
  type         = "UNMANAGED"
}
```

could be applied and any compute resources were ignored.

Now this configuration is invalid and will result in an error during plan.

To resolve this error simply remove or comment out the `compute_resources` configuration block.

```hcl
resource "aws_batch_compute_environment" "test" {
  compute_environment_name = "test"

  service_role = aws_iam_role.batch_service.arn
  type         = "UNMANAGED"
}
```

## Resource: aws_cloudwatch_event_target

### Removal of `ecs_target` `launch_type` default value

Previously, the `ecs_target` `launch_type` argument defaulted to `EC2` if no value was configured in Terraform.

Workarounds, such as using the empty string `""` as shown below, should be removed:

```terraform
resource "aws_cloudwatch_event_target" "test" {
  arn      = aws_ecs_cluster.test.id
  rule     = aws_cloudwatch_event_rule.test.id
  role_arn = aws_iam_role.test.arn
  ecs_target {
    launch_type         = ""
    task_count          = 1
    task_definition_arn = aws_ecs_task_definition.task.arn
    network_configuration {
      subnets = [aws_subnet.subnet.id]
    }
  }
}
```

An updated configuration:

```terraform
resource "aws_cloudwatch_event_target" "test" {
  arn      = aws_ecs_cluster.test.id
  rule     = aws_cloudwatch_event_rule.test.id
  role_arn = aws_iam_role.test.arn
  ecs_target {
    task_count          = 1
    task_definition_arn = aws_ecs_task_definition.task.arn
    network_configuration {
      subnets = [aws_subnet.subnet.id]
    }
  }
}
```

## Resource: aws_customer_gateway

Previously, `ip_address` could be set to `""`, which would result in an AWS error. However, this value is no longer accepted by the provider.

## Resource: aws_default_network_acl

Previously, `egress.*.cidr_block`, `egress.*.ipv6_cidr_block`, `ingress.*.cidr_block`, or `ingress.*.ipv6_cidr_block` could be set to `""`. However, the value `""` is no longer valid.

For example, previously this type of configuration was valid:

```terraform
resource "aws_default_network_acl" "default" {
  # ...
  egress {
    cidr_block      = "0.0.0.0/0"
    ipv6_cidr_block = ""
    # ...
  }
}
```

Now, set the argument to null (`ipv6_cidr_block = null`) or simply remove the empty-value configuration:

```terraform
resource "aws_default_network_acl" "default" {
  # ...
  egress {
    cidr_block      = "0.0.0.0/0"
    # ...
  }
}
```

## Resource: aws_elasticache_cluster

### Error raised if neither `engine` nor `replication_group_id` is specified

Previously, when neither `engine` nor `replication_group_id` was specified, Terraform would not prevent users from applying the invalid configuration.
Now, this will produce an error similar to the below:

```
Error: Invalid combination of arguments

          with aws_elasticache_cluster.example,
          on terraform_plugin_test.tf line 2, in resource "aws_elasticache_cluster" "example":
           2: resource "aws_elasticache_cluster" "example" {

        "replication_group_id": one of `engine,replication_group_id` must be
        specified

        Error: Invalid combination of arguments

          with aws_elasticache_cluster.example,
          on terraform_plugin_test.tf line 2, in resource "aws_elasticache_cluster" "example":
           2: resource "aws_elasticache_cluster" "example" {

        "engine": one of `engine,replication_group_id` must be specified
```

Configuration that depend on the previous behavior will need to be updated.

## Resource: aws_elasticache_global_replication_group

### actual_engine_version Attribute removal

Switch your Terraform configuration to the `engine_version_actual` attribute instead.

For example, given this previous configuration:

```terraform
output "elasticache_global_replication_group_version_result" {
  value = aws_elasticache_global_replication_group.example.actual_engine_version
}
```

An updated configuration:

```terraform
output "elasticache_global_replication_group_version_result" {
  value = aws_elasticache_global_replication_group.example.engine_version_actual
}
```

## Resource: aws_elasticache_replication_group

!> **WARNING:** This topic is placeholder documentation.

## Resource: aws_fsx_ontap_storage_virtual_machine

We removed the misspelled argument `active_directory_configuration.0.self_managed_active_directory_configuration.0.organizational_unit_distinguidshed_name` that was previously deprecated. Use `active_directory_configuration.0.self_managed_active_directory_configuration.0.organizational_unit_distinguished_name` now instead. Terraform will automatically migrate the state to `active_directory_configuration.0.self_managed_active_directory_configuration.0.organizational_unit_distinguished_name` during planning.

## Resource: aws_network_acl

Previously, `egress.*.cidr_block`, `egress.*.ipv6_cidr_block`, `ingress.*.cidr_block`, or `ingress.*.ipv6_cidr_block` could be set to `""`. However, the value `""` is no longer valid.

For example, previously this type of configuration was valid:

```terraform
resource "aws_network_acl" "default" {
  # ...
  egress {
    cidr_block      = "0.0.0.0/0"
    ipv6_cidr_block = ""
    # ...
  }
}
```

Now, set the argument to null (`ipv6_cidr_block = null`) or simply remove the empty-value configuration:

```terraform
resource "aws_network_acl" "default" {
  # ...
  egress {
    cidr_block      = "0.0.0.0/0"
    # ...
  }
}
```

## Resource: aws_network_interface

!> **WARNING:** This topic is placeholder documentation.

## Resource: aws_s3_bucket

!> **WARNING:** This topic is placeholder documentation.

## Resource: aws_spot_instance_request

### instance_interruption_behaviour Argument removal

Switch your Terraform configuration to the `engine_version_actual` attribute instead.

For example, given this previous configuration:

```terraform
resource "aws_spot_instance_request" "example" {
  # ... other configuration ...
  instance_interruption_behaviour = "hibernate"
}
```

An updated configuration:

```terraform
resource "aws_spot_instance_request" "example" {
  # ... other configuration ...
  instance_interruption_behavior =  "hibernate"
}
```

## Resource: aws_s3_bucket_object

The `aws_s3_bucket_object` resource is deprecated and will be removed in a future version. Use `aws_s3_object` instead, where new features and fixes will be added.

When replacing `aws_s3_bucket_object` with `aws_s3_object` in your configuration, on the next apply, Terraform will recreate the object. If you prefer to not have Terraform recreate the object, import the object using `aws_s3_object`.

For example, the following will import an S3 object into state, assuming the configuration exists, as `aws_s3_object.example`:

```console
% terraform import aws_s3_object.example s3://some-bucket-name/some/key.txt
```

## EC2-Classic Resource and Data Source Support

While an upgrade to this major version will not directly impact EC2-Classic resources configured with Terraform,
it is important to keep in the mind the following AWS Provider resources will eventually no longer
be compatible with EC2-Classic as AWS completes their EC2-Classic networking retirement (expected around August 15, 2022).

* Running or stopped [EC2 instances](/docs/providers/aws/r/instance.html)
* Running or stopped [RDS database instances](/docs/providers/aws/r/db_instance.html)
* [Elastic IP addresses](/docs/providers/aws/r/eip.html)
* [Classic Load Balancers](/docs/providers/aws/r/lb.html)
* [Redshift clusters](/docs/providers/aws/r/redshift_cluster.html)
* [Elastic Beanstalk environments](/docs/providers/aws/r/elastic_beanstalk_environment.html)
* [EMR clusters](/docs/providers/aws/r/emr_cluster.html)
* [AWS Data Pipelines pipelines](/docs/providers/aws/r/datapipeline_pipeline.html)
* [ElastiCache clusters](/docs/providers/aws/r/elasticache_cluster.html)
* [Spot Requests](/docs/providers/aws/r/spot_instance_request.html)
* [Capacity Reservations](/docs/providers/aws/r/ec2_capacity_reservation.html)
