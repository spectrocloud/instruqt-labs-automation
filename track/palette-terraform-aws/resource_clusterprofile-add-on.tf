# If looking up a cluster profile instead of creating a new one
# data "spectrocloud_cluster_profile" "profile" {
#   # id = <uid>
#   name = "eks-basic"
# }
#

resource "spectrocloud_cluster_profile" "add-on" {
  name        = "tf-add-on-profile"
  description = "OPA + Hipster App"
  type        = "add-on"
  tags        = ["owner:ops"]
  version     = "1.0.0"

  pack {
    name   = data.spectrocloud_pack.opa.name
    tag    = local.opa_version
    uid    = data.spectrocloud_pack.opa.id
    values = data.spectrocloud_pack.opa.values
  }

  pack {
    name   = data.spectrocloud_pack.hipster.name
    tag    = local.hipster_version
    uid    = data.spectrocloud_pack.hipster.id
    values = data.spectrocloud_pack.hipster.values
  }
}
  

