---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "packetfabric_cloud_router_quick_connect Resource - terraform-provider-packetfabric"
subcategory: ""
description: |-
  
---

# packetfabric_cloud_router_quick_connect (Resource)

Use Cloud Router Quick Connect to import a service provider's routes into one of your Cloud Router connections. For more information, see [Quick Connect in the PacketFabric documentation](https://docs.packetfabric.com/cr/qc/).

## Example Usage

```terraform
resource "packetfabric_cloud_router_quick_connect" "cr_quick_connect" {
  provider              = packetfabric
  cr_circuit_id         = packetfabric_cloud_router.cr1.id
  connection_circuit_id = packetfabric_cloud_router_connection_aws.crc1.id
  service_uuid          = var.pf_service_uuid
  return_filters {
    prefix     = "185.56.153.165/32"
    match_type = "exact"
  }
  return_filters {
    prefix     = "185.56.153.166/32"
    match_type = "exact"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `connection_circuit_id` (String) The circuit ID of the Cloud Router connection that will be importing the routes.
- `cr_circuit_id` (String) The circuit ID of your Cloud Router.
- `return_filters` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--return_filters))
- `service_uuid` (String) The service UUID associated with the service provider's Quick Connect.

### Optional

- `import_filters` (Block Set) This is set by the service provider. (see [below for nested schema](#nestedblock--import_filters))
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `connection_speed` (String) The speed of the target cloud router connection.
- `id` (String) The Circuit ID of this Cloud Router Import.
- `is_defunct` (Boolean) Whether the Quick Connect is defunct. This typically happens when the provider removes the service.
- `route_set_circuit_id` (String) The Circuit ID of the Route Set selected for this Cloud Router Import.
- `state` (String) Shows the state of this import.

	Enum: `"pending"` `"rejected"` `"provisioning"`  `"active"`  `"deleting"`  `"inactive"`

<a id="nestedblock--return_filters"></a>
### Nested Schema for `return_filters`

Required:

- `prefix` (String) The prefix to export to the service provider that they will use for return traffic.

Optional:

- `as_prepend` (Number) The AS prepend to apply to the exported/returned prefix.

	Available range is 1 through 5.
- `match_type` (String) The match type of this prefix.

	Enum: `"exact"` `"orlonger"` Defaults: exact
- `med` (Number) The MED to apply to the exported/returned prefix.

	Available range is 1 through 4294967295.


<a id="nestedblock--import_filters"></a>
### Nested Schema for `import_filters`

Optional:

- `local_preference` (Number) The local preference to apply to the prefix.

	Available range is 1 through 4294967295.
- `match_type` (String) The match type for the imported prefix. This is set by the service provider.

	Enum: `"exact"` `"orlonger"`
- `prefix` (String) The route prefix that you will be importing from the Quick Connect. This is set by the service provider.


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)




## Import

Import a Cloud Router Quick Connec using its corresponding Import Circuit ID, circuit ID and the ID of the Cloud Router it is associated with, in the format `<cloud router ID>:<cloud router connection ID>:<import circuit ID>`.

```bash
terraform import packetfabric_cloud_router_quick_connect.cr_quick_connect PF-L3-CUST-1700239:PF-L3-CON-2980512:PF-L3-IMP-12345
```