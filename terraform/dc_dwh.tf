data "aws_caller_identity" "self" {}

data "aws_region" "self" {}

# Prepare BYOC VPC and IAM Role
module "doublecloud_byoc" {
  source  = "doublecloud/doublecloud-byoc/aws"
  version = "1.0.2"
  providers = {
    aws = aws
  }
  ipv4_cidr = var.dwh_ipv4_cidr
}

# Get account ID to peer with
data "aws_caller_identity" "peered" {
  provider = aws
}

# Create DoubleCloud BYOC Network
resource "doublecloud_network" "aws" {
  project_id = var.dc_project_id
  name       = "alpha-network"
  region_id  = module.doublecloud_byoc.region_id
  cloud_type = "aws"
  aws = {
    vpc_id          = module.doublecloud_byoc.vpc_id
    account_id      = module.doublecloud_byoc.account_id
    iam_role_arn    = module.doublecloud_byoc.iam_role_arn
    private_subnets = false
  }
}

# Create VPC Peering from DoubleCloud Network to AWS VPC
resource "doublecloud_network_connection" "example" {
  network_id = doublecloud_network.aws.id
  aws = {
    peering = {
      vpc_id          = aws_vpc.tutorial_vpc.id
      account_id      = data.aws_caller_identity.peered.account_id
      region_id       = var.aws_region
      ipv4_cidr_block = aws_vpc.tutorial_vpc.cidr_block
      ipv6_cidr_block = aws_vpc.tutorial_vpc.ipv6_cidr_block
    }
  }
}

# Accept Peering Request on AWS side
resource "aws_vpc_peering_connection_accepter" "own" {
  provider                  = aws
  vpc_peering_connection_id = time_sleep.avoid_aws_race.triggers["peering_connection_id"]
  auto_accept               = true
}

# Confirm Peering creation
resource "doublecloud_network_connection_accepter" "accept" {
  id = doublecloud_network_connection.example.id

  depends_on = [
    aws_vpc_peering_connection_accepter.own,
  ]
}

# Create ipv4 routes to DoubleCloud Network
resource "aws_route" "ipv4_private_route" {
  provider                  = aws
  route_table_id            = aws_route_table.tutorial_private_rt.id
  destination_cidr_block    = doublecloud_network_connection.example.aws.peering.managed_ipv4_cidr_block
  vpc_peering_connection_id = time_sleep.avoid_aws_race.triggers["peering_connection_id"]
}
resource "aws_route" "ipv4_public_route" {
  provider                  = aws
  route_table_id            = aws_route_table.tutorial_public_rt.id
  destination_cidr_block    = doublecloud_network_connection.example.aws.peering.managed_ipv4_cidr_block
  vpc_peering_connection_id = time_sleep.avoid_aws_race.triggers["peering_connection_id"]
}

# Sleep to avoid AWS InvalidVpcPeeringConnectionID.NotFound error
resource "time_sleep" "avoid_aws_race" {
  create_duration = "30s"

  triggers = {
    peering_connection_id = doublecloud_network_connection.example.aws.peering.peering_connection_id
  }
}

## Actual Clickhouse DB

resource "doublecloud_clickhouse_cluster" "alpha-clickhouse" {
  project_id = var.dc_project_id
  name       = "alpha-clickhouse"
  region_id  = var.aws_region
  cloud_type = "aws"
  network_id = doublecloud_network.aws.id

  resources {
    clickhouse {
      resource_preset_id = "s1-c2-m4"
      disk_size          = 34359738368
      replica_count      = 1
    }
  }

  config {
    log_level       = "LOG_LEVEL_TRACE"
    max_connections = 120
  }

  access {
    data_services    = ["visualization"]
    ipv4_cidr_blocks = [
      {
        value       = doublecloud_network.aws.ipv4_cidr_block
        description = "DC Network interconnection"
      },
      {
        value       = aws_vpc.tutorial_vpc.cidr_block
        description = "Peered VPC"
      },
      {
        value       = "${var.my_ip}/32"
        description = "My IP"
      }
    ]
    ipv6_cidr_blocks = [
      {
        value       = "${var.my_ipv6}/128"
        description = "My IPv6"
      }
    ]
  }
}
