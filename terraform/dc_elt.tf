resource "doublecloud_transfer_endpoint" "pg-source" {
  name       = "chinook-pg-source"
  project_id = var.dc_project_id
  settings {
    postgres_source {
      connection {
        on_premise {
          tls_mode {
            ca_certificate = file("global-bundle.pem")
          }
          hosts = [
            aws_db_instance.tutorial_database.address
          ]
          port = 5432
        }
      }
      database = aws_db_instance.tutorial_database.db_name
      user     = aws_db_instance.tutorial_database.username
      password = var.db_password
    }
  }
}

data "doublecloud_clickhouse" "dwh" {
  name       = doublecloud_clickhouse_cluster.alpha-clickhouse.name
  project_id = var.dc_project_id
}

resource "doublecloud_transfer_endpoint" "dwh-target" {
  name       = "alpha-clickhouse-target"
  project_id = var.dc_project_id
  settings {
    clickhouse_target {
      connection {
        address {
          cluster_id = doublecloud_clickhouse_cluster.alpha-clickhouse.id
        }
        database = "default"
        user     = data.doublecloud_clickhouse.dwh.connection_info.user
        password = data.doublecloud_clickhouse.dwh.connection_info.password
      }
    }
  }
}

resource "doublecloud_transfer" "pg2ch" {
  name       = "postgres-to-clickhouse-snapshot"
  project_id = var.dc_project_id
  source     = doublecloud_transfer_endpoint.pg-source.id
  target     = doublecloud_transfer_endpoint.dwh-target.id
  type       = "SNAPSHOT_ONLY"
  activated  = false
  runtime = {
    dedicated = {
      flavor = "TINY"
      vpc_id = doublecloud_network.aws.id
    }
  }
}
