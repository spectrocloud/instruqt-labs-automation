resource "local_file" "kubeconfig" {
  content              = local.cluster_kubeconfig
  filename             = "kubeconfig_${var.sc_cluster_name}"
  file_permission      = "0644"
  directory_permission = "0755"
}
