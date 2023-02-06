---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "packetfabric_marketplace_service Resource - terraform-provider-packetfabric"
subcategory: ""
description: |-
  
---

# packetfabric_marketplace_service (Resource)

Create a service to make available through the PacketFabric marketplace. For more information, see [Marketplace & IX in the PacketFabric documentation](https://docs.packetfabric.com/eco/). 



## Example Usage

```terraform
# Create a Marketplace Service type quick-connect
resource "packetfabric_marketplace_service" "marketplace_quick_connect" {
  provider                = packetfabric
  name                    = var.pf_name
  description             = var.pf_description
  service_type            = "quick-connect-service"
  sku                     = var.pf_sku
  categories              = var.pf_categories
  published               = var.pf_published
  cloud_router_circuit_id = var.pf_cloud_router_circuit_id
  connection_circuit_ids  = var.pf_connection_circuit_ids
  route_set {
    description = var.pf_route_set_description
    is_private  = var.pf_route_set_is_private
    prefixes {
      prefix     = var.pf_route_set_prefix1
      match_type = var.pf_route_set_match_type1
    }
    prefixes {
      prefix     = var.pf_route_set_prefix2
      match_type = var.pf_route_set_match_type2
    }
  }
}
output "packetfabric_marketplace_service_quick_connect" {
  value = packetfabric_marketplace_service.marketplace_quick_connect
}

# Create a Marketplace Service type port
resource "packetfabric_marketplace_service" "marketplace_port" {
  provider     = packetfabric
  name         = var.pf_name
  description  = var.pf_description
  service_type = "port-service"
  sku          = var.pf_sku
  categories   = var.pf_categories
  published    = var.pf_published
  locations    = var.pf_locations
}
output "packetfabric_marketplace_service_port" {
  value = packetfabric_marketplace_service.marketplace_port
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the service.
- `published` (Boolean) If published, the service appears in your marketplace listing.

### Optional

- `categories` (List of String) Categories in which the service will fit.

	Enum: `"cloud-computing"`, `"content-delivery-network"`, `"edge-computing"`, `"sd-wan"`, `"data-storage"`, `"developer-platform"`, `"internet-service-provider"`, `"security"`, `"video-conferencing"`, `"voice-and-messaging"`, `"web-hosting"`, `"internet-of-things"`, `"private-connectivity"`, `"bare-metal-hosting"`
- `cloud_router_circuit_id` (String) The circuit ID of the Cloud Router this service is associated with (Quick Connect service only).
- `connection_circuit_ids` (List of String) The circuit IDs of the Cloud Router connections that will be included in this service. (Quick Connect service only).
- `description` (String) Brief description of what the service does.
- `locations` (List of String) Locations in which the service will operate (port service only). The location should be a POP, e.g. `NYC5`.
- `route_set` (Block Set) The Cloud Router route set to export (Quick Connect service only). (see [below for nested schema](#nestedblock--route_set))
- `service_type` (String) The service type of this service. Enum: `"port-service"`, `"quick-connect-service"` Defaults: port-service
- `sku` (String) A SKU identifier for the service. This is not shown to the A side user (the requestor).
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The ID of this resource.
- `route_set_circuit_id` (String) The route set circuit ID.
- `service_uuid` (String) The marketplace service UUID
- `state` (String) The marketplace service state.

<a id="nestedblock--route_set"></a>
### Nested Schema for `route_set`

Optional:

- `description` (String) The route set's description.
- `is_private` (Boolean) In a private route set, the return traffic is private. In other words, in a public route set, anyone who imports this route set can also see other clients who are importing the route based on return traffic. Defaults: true
- `prefixes` (Block Set) (see [below for nested schema](#nestedblock--route_set--prefixes))

<a id="nestedblock--route_set--prefixes"></a>
### Nested Schema for `route_set.prefixes`

Required:

- `match_type` (String) The match type for this prefix. Options are: `"exact"` and `"orlonger"`.
- `prefix` (String) A prefix, in CIDR format, to include in this route set.



<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)

