---
subcategory: "Security Group"
layout: "twcc"
page_title: "TWCC: twcc_security_group"
description: |-
  Provides a security group.
---

# Resource: twcc_security_group

Provides a security group resource. Security groups control network traffic to and from VCS instances and other resources by managing security rules.

## Example Usage

### Creating a standalone security group with rules

```hcl
data "twcc_project" "testProject" {
    name = "ENT108079"
    platform = "openstack-taichung-default-2"
}

# Create a standalone security group
resource "twcc_security_group" "my_sg" {
    name = "my-security-group"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    description = "Security group for my application"
}

# Add SSH rule to the security group
resource "twcc_security_group_rule" "allow_ssh" {
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    security_group = twcc_security_group.my_sg.id
    direction = "ingress"
    protocol = "tcp"
    port_range_min = 22
    port_range_max = 22
    remote_ip_prefix = "0.0.0.0/0"
}

# Add HTTP rule to the same security group
resource "twcc_security_group_rule" "allow_http" {
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    security_group = twcc_security_group.my_sg.id
    direction = "ingress"
    protocol = "tcp"
    port_range_min = 80
    port_range_max = 80
    remote_ip_prefix = "0.0.0.0/0"
}
```

### Using security group with VCS instance

```hcl
data "twcc_project" "testProject" {
    name = "ENT108079"
    platform = "openstack-taichung-default-2"
}

# Create security group first
resource "twcc_security_group" "web_sg" {
    name = "web-server-sg"
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    description = "Security group for web servers"
}

# Add rules to security group
resource "twcc_security_group_rule" "web_https" {
    platform = data.twcc_project.testProject.platform
    project = data.twcc_project.testProject.id
    security_group = twcc_security_group.web_sg.id
    direction = "ingress"
    protocol = "tcp"
    port_range_min = 443
    port_range_max = 443
    remote_ip_prefix = "0.0.0.0/0"
}

# VCS instance can reference the security group ID
# (exact usage depends on VCS resource configuration)
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the security group.

* `platform` - (Required) The name of the platform where security group is created.

* `project` - (Required) The ID of the project where security group is created.

* `description` - (Optional) The description of the security group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the security group.

* `create_time` - The create time (UTC) of the security group.

* `type` - The type of the security group.

* `user` - The user information who created the security group.

* `security_group_rules` - The list of security group rules associated with this security group.

## Import

Security groups can be imported using the `id`, e.g.

```
$ terraform import twcc_security_group.my_sg <security_group_id>
```
