---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "packetfabric_cloud_router_connection_google Resource - terraform-provider-packetfabric"
subcategory: ""
description: |-
  
---

# packetfabric_cloud_router_connection_google (Resource)



A connection from your cloud router to your Google Cloud Platform environment. For more information, see [Cloud Router in the PacketFabric documentation](https://docs.packetfabric.com/cr/).

## Example Usage

```terraform
resource "packetfabric_cloud_router" "cr1" {
  provider     = packetfabric
  asn          = var.pf_cr_asn
  name         = var.pf_cr_name
  account_uuid = var.pf_account_uuid
  capacity     = var.pf_cr_capacity
  regions      = var.pf_cr_regions
}

resource "packetfabric_cloud_router_connection_google" "crc2" {
  provider                    = packetfabric
  description                 = var.pf_crc_description
  circuit_id                  = packetfabric_cloud_router.cr1.id
  account_uuid                = var.pf_account_uuid
  google_pairing_key          = var.pf_crc_google_pairing_key
  google_vlan_attachment_name = var.pf_crc_google_vlan_attachment_name
  pop                         = var.pf_crc_pop
  speed                       = var.pf_crc_speed
  maybe_nat                   = var.pf_crc_maybe_nat
}

output "packetfabric_cloud_router_connection_google" {
  value = packetfabric_cloud_router_connection_google.crc2
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_uuid` (String) The UUID for the billing account that should be billed.
- `circuit_id` (String) Circuit ID of the target cloud router. This starts with "PF-L3-CUST-".
- `description` (String) A brief description of this connection.
- `google_pairing_key` (String) The Google pairing key to use for this connection. This is generated when you create your Google Cloud VLAN attachment.
- `google_vlan_attachment_name` (String) The Google VLAN attachment name.
- `pop` (String) The POP in which you want to provision the connection.
- `speed` (String) The desired speed of the new connection.

	Enum: ["50Mbps" "100Mbps" "200Mbps" "300Mbps" "400Mbps" "500Mbps" "1Gbps" "2Gbps" "5Gbps" "10Gbps"]

### Optional

- `maybe_nat` (Boolean) Set this to true if you intend to use NAT on this connection.Defaults: false
- `published_quote_line_uuid` (String) UUID of the published quote line with which this connection should be associated.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)

