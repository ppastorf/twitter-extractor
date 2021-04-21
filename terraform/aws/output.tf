# vpc
output "vpc_cidr_block" {
  value = "${aws_vpc.cluster_vpc.cidr_block}"
}

output "subnet_cidr_block" {
  value = "${aws_subnet.cluster_subnet.cidr_block}"
}

# cluster 01
output "cluster_01_private_ip" {
  value = "${aws_eip.cluster_01_eip.private_ip}"
}

output "cluster_01_private_dns" {
  value = "${aws_eip.cluster_01_eip.private_dns}"
}

output "cluster_01_public_ip" {
  value = "${aws_eip.cluster_01_eip.public_ip}"
}

output "cluster_01_public_dns" {
  value = "${aws_eip.cluster_01_eip.public_dns}"
}