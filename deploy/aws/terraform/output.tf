# vpc
output "vpc_cidr_block" {
  value = aws_vpc.vpc.cidr_block
}

output "subnet_cidr_block" {
  value = aws_subnet.nodes_subnet.cidr_block
}

# vm 1
output "vm_01_instance_id" {
  value = aws_instance.vm_01_instance.id
}

output "vm_01_private_ip" {
  value = aws_eip.vm_01_eip.private_ip
}

output "vm_01_private_dns" {
  value = aws_eip.vm_01_eip.private_dns
}

output "vm_01_public_ip" {
  value = aws_eip.vm_01_eip.public_ip
}

output "vm_01_public_dns" {
  value = aws_eip.vm_01_eip.public_dns
}