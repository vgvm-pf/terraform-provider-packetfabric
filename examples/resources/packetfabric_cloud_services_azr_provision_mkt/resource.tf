terraform {
  required_providers {
    packetfabric = {
      source  = "packetfabric/packetfabric"
      version = "~> 0.0.1"
    }
  }
}
provider "packetfabric" {
  host  = var.pf_api_server
  token = var.pf_api_key
}

resource "packetfabric_cloud_services_azr_provision_mkt" "new" {
  provider = packetfabric
  description = "my-azure-provisioned-circuit"
  port_circuit_id = "PF-AP-XYZ1-1234"
  vc_request_uuid = "PF-BC-AB1-YZ1-1234567"
  vlan = "25"
}