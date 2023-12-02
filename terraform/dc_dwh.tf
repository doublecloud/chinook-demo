data "aws_caller_identity" "self" {}

data "aws_region" "self" {}

# Prepare BYOC VPC and IAM Role
module "doublecloud_byoc" {
  source  = "doublecloud/doublecloud-byoc/aws"
  version = "1.0.2"
  providers = {
    aws = aws
  }
  ipv4_cidr = var.vpc_cidr_block
}

# Create VPC to peer with
resource "aws_vpc" "peered" {
  cidr_block                       = var.dwh_ipv4_cidr
  provider                         = aws
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
    private_subnets = true
  }
}

# Create VPC Peering from DoubleCloud Network to AWS VPC
resource "doublecloud_network_connection" "example" {
  network_id = doublecloud_network.aws.id
  aws = {
    peering = {
      vpc_id          = aws_vpc.peered.id
      account_id      = data.aws_caller_identity.peered.account_id
      region_id       = var.aws_region
      ipv4_cidr_block = aws_vpc.peered.cidr_block
      ipv6_cidr_block = aws_vpc.peered.ipv6_cidr_block
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
resource "aws_route" "ipv4" {
  provider                  = aws
  route_table_id            = aws_vpc.peered.main_route_table_id
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

