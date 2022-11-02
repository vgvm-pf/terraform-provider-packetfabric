resource "packetfabric_port" "port_1" {
  provider          = packetfabric
  account_uuid      = var.pf_account_uuid
  autoneg           = var.pf_port_autoneg
  description       = var.pf_description
  media             = var.pf_port_media
  nni               = var.pf_port_nni
  pop               = var.pf_port_pop1
  speed             = var.pf_port_speed
  subscription_term = var.pf_port_subterm
  zone              = var.pf_port_avzone1
}

resource "packetfabric_backbone_virtual_circuit_marketplace" "vc_marketplace_conn1" {
  provider    = packetfabric
  description = var.pf_description
  routing_id  = var.pf_routing_id
  market      = var.pf_market
  interface {
    port_circuit_id = packetfabric_port.port_1.id
    untagged        = false
    vlan            = var.pf_vc_vlan1
  }
  bandwidth {
    account_uuid      = var.pf_account_uuid
    longhaul_type     = var.pf_vc_longhaul_type
    speed             = var.pf_vc_speed
    subscription_term = var.pf_vc_subterm
  }
}

output "packetfabric_backbone_virtual_circuit_marketplace" {
  value = packetfabric_backbone_virtual_circuit_marketplace.vc_marketplace_conn1
}