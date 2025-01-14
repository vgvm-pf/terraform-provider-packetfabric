---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "packetfabric_cloud_router_connection_port Resource - terraform-provider-packetfabric"
subcategory: ""
description: |-
  
---

# packetfabric_cloud_router_connection_port (Resource)

A connection from your cloud router to one of your PacketFabric access ports. For more information, see [Cloud Router in the PacketFabric documentation](https://docs.packetfabric.com/cr/).

## Example Usage

```terraform
resource "packetfabric_cloud_router" "cr1" {
  provider = packetfabric
  asn      = 4556
  name     = "hello world"
  capacity = "10Gbps"
  regions  = ["US"]
  labels   = ["terraform", "dev"]
}

resource "packetfabric_cloud_router_connection_port" "crc7" {
  provider        = packetfabric
  description     = "hello world"
  circuit_id      = packetfabric_cloud_router.cr1.id
  port_circuit_id = packetfabric_port.port_1.id
  vlan            = 104
  untagged        = false
  speed           = "1Gbps"
  is_public       = false
  maybe_nat       = false
  labels          = ["terraform", "dev"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_uuid` (String) The UUID for the billing account that should be billed. Can also be set with the PF_ACCOUNT_ID environment variable.
- `circuit_id` (String) Circuit ID of the target cloud router. This starts with "PF-L3-CUST-".
- `description` (String) A brief description of this connection.
- `port_circuit_id` (String) The circuit ID of the port to connect to the cloud router. This starts with "PF-AP-".
- `speed` (String) The speed of the new connection.

	Enum: ["50Mbps" "100Mbps" "200Mbps" "300Mbps" "400Mbps" "500Mbps" "1Gbps" "2Gbps" "5Gbps" "10Gbps"]

### Optional

- `is_public` (Boolean) Whether PacketFabric should allocate a public IP address for this connection. Defaults: false
- `labels` (Set of String) Label value linked to an object.
- `maybe_dnat` (Boolean) Set this to true if you intend to use DNAT on this connection. Defaults: false
- `maybe_nat` (Boolean) Set this to true if you intend to use NAT on this connection. Defaults: false
- `po_number` (String) Purchase order number or identifier of a service.
- `published_quote_line_uuid` (String) UUID of the published quote line with which this connection should be associated.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `untagged` (Boolean) Whether the interface should be untagged. Do not specify a VLAN if this is to be an untagged connection. Defaults: false
- `vlan` (Number) Valid VLAN range is from 4-4094, inclusive.

### Read-Only

- `etl` (Number) Early Termination Liability (ETL) fees apply when terminating a service before its term ends. ETL is prorated to the remaining contract days.
- `id` (String) The ID of this resource.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)




## Import

Import a Cloud Router Connection using its corresponding circuit ID and the ID of the Cloud Router it is associated with, in the format `<cloud router ID>:<cloud router connection ID>`.

```bash
terraform import packetfabric_cloud_router_connection_port.crc PF-L3-CUST-1700239:PF-L3-CON-2980512
```