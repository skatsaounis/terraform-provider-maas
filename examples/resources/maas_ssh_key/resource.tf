resource "tls_private_key" "tf_key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "maas_ssh_key" "tf_key" {
  key = tls_private_key.tf_key.public_key_openssh
}
