## Integration Application

Sample terraform application for client-facing analytics.

# Stage 1: Exist infra 

Basic application infrastructure looks as follows:

![stage1.png](diagrams%2Fstage1.png)

Everything is described in `terraform`-folder.

To run this in your own AWS space do follows:

```shell
cd terraform
terraform init
terraform apply -var-file=env.tfvars -var="my_ip=$(curl -4 ifconfig.me)"
```

Inside `env.tfvars` consist main argument that used inside module, default one:

```terraform
db_password = "Password"
db_username = "chinook_admin"
```

Final goal is extract data from private Postgresql, to feed postgresql with test data see readme in [here](./loadgen)

# Stage 2: Create and connect with Data VPC

![stage2.png](diagrams/stage2.png)

First thing first - connectivity. For that we need to connect double-cloud via [BYOA](https://double.cloud/docs/en/vpc/connect-dc-to-aws)

# Stage 3: Add Clickhouse in Data VPC

Once we establish connectivity bridge we can introduce private DWH instance ([Clickhouse](https://double.cloud/services/managed-clickhouse/))

![stage3.png](diagrams/stage3.png)

# Stage 4: Transfer from PostgreSQL to Clickhouse

Once we have a cluster, we need to get some data. To make it possible we will create a transfer from exist postgres to this clickhouse.

![stage4.png](diagrams/stage4.png)

# Stage 5: Expose Clickhouse via Dashboard

Finally, we have all data in clickhouse, last step is here: Expose it via Dashboard

![stage5.png](diagrams/stage5.png)
