---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "packetfabric_cs_google_hosted_connection Data Source - terraform-provider-packetfabric"
subcategory: ""
description: |-
  
---

# packetfabric_cs_google_hosted_connection (Data Source)


## Example Usage

```terraform
{
  "hosted_connections" = tolist([
    {
      "account_uuid" = "bbbfb3fe-cdd1-48a9-90ea-9fc59ea41234"
      "cloud_circuit_id" = "PF-AE-PDX1-1739482"
      "cloud_provider" = toset([
        {
          "pop" = "PDX1"
          "site" = "Pittock Building"
        },
      ])
      "customer_uuid" = "58c80946-5fbc-400e-8060-95b5dfbf1234"
      "deleted" = false
      "description" = "Google Hosted connection for Foo update"
      "is_cloud_router_connection" = false
      "pop" = "PDX1"
      "port_type" = "dedicated"
      "service_class" = "metro"
      "service_provider" = "aws"
      "settings" = toset([
        {
          "autoneg" = false
          "aws_region" = "us-west-2"
          "zone_dest" = "B"
        },
      ])
      "settings_aws_region" = ""
      "site" = "Pittock Building"
      "speed" = "10Gbps"
      "state" = "active"
      "time_created" = "2022-06-16T23:13:21.126145-0700"
      "time_updated" = "2022-06-16T23:15:09.089127-0700"
      "user_uuid" = "4e3bb859-9f64-4d12-ae9c-be3a0231234"
      "uuid" = "3adadf96-3c27-4598-baf8-f4d993401234"
    },
  ])
  "id" = "f9a76a87-e7d0-44b9-a35e-9a89b2241234"
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cloud_circuit_id` (String) The unique PF circuit ID for this connection
		Example: PF-AP-LAX1-1002

### Optional

- `cloud_provider_pop` (String) Point of Presence for the cloud provider location.
		Example: DAL1
- `cloud_provider_region` (String) Region short name.
		Example: us-west-1
- `customer_uuid` (String) The UUID for the customer this connection belongs to.
- `description` (String) The description of this connection.
		Example: AWS connection for Foo Corp.
- `pop` (String) Point of Presence for the connection.
		Example: LAS1
- `port_type` (String) The port type for the given port.
		Enum: [ "hosted", "dedicated" ]
- `service_class` (String) The service class for the given port, either long haul or metro.
		Enum: [ "longhaul", "metro" ]
- `service_provider` (String) The service provider of the connection
		Enum: [ "aws", "azure", "packet", "google", "ibm", "salesforce", "webex" ]
- `site` (String) Site name
		Example: SwitchNAP Las Vegas 7
- `speed` (String) The desired speed of the connection.
		Enum: [ "50Mbps", "100Mbps", "200Mbps", "300Mbps", "400Mbps", "500Mbps", "1Gbps", "2Gbps", "5Gbps", "10Gbps" ]
- `state` (String) The state of the connection.
		Enum: [ "active", "deleting", "inactive", "pending", "requested" ]
- `time_created` (String) Date and time of connection creation
- `time_updated` (String) Date and time connection was last updated
- `user_uuid` (String) The UUID for the user this connection belongs to.

### Read-Only

- `id` (String) The ID of this resource.

