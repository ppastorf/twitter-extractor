data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/*/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

# vpc
resource "aws_vpc" "vpc" {
  cidr_block           = "10.0.0.0/16"
  instance_tenancy     = "default"
  enable_dns_hostnames = true

  tags = {
    Name = "gcloud13_vpc"
  }
}

resource "aws_security_group" "nodes_sg" {
  name   = "gcloud13_nodes_sg"
  vpc_id = aws_vpc.vpc.id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "TCP"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "gcloud13_nodes_sg"
  }
}

resource "aws_subnet" "nodes_subnet" {
  vpc_id                  = aws_vpc.vpc.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true

  tags = {
    Name = "gcloud13_nodes_subnet"
  }
}

resource "aws_internet_gateway" "vpc_igw" {
  vpc_id = aws_vpc.vpc.id

  tags = {
    Name = "gcloud13_vpc_igw"
  }
}

resource "aws_route_table" "vpc_route_table" {
  vpc_id = aws_vpc.vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.vpc_igw.id
  }

  tags = {
    Name = "gcloud13_vpc_route_table"
  }
}

resource "aws_route_table_association" "nodes_subnet_route_table_association" {
  subnet_id      = aws_subnet.nodes_subnet.id
  route_table_id = aws_route_table.vpc_route_table.id
}

resource "aws_key_pair" "ssh_key" {
  key_name   = "gcloud13_ssh_key"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCt7leg0VkmuFOzvAJNFy7xCFLOAujJl1bcHd6yP3t4mG91h0ea/Ezr2FVGktc62I44rBWflZ6Lw+E0RmeHILUGDV49sPs7TMi43Dw3JHeNg4ek6Xleu4oGY6Nh+76vG+5jgZFbNIYTNX28mc1Q8G+ZUYXzDX4DyzSpNCmPv1owUv+7KdPjxQr9Op0fWBWkKK/ilsI/4YlnJW1KeiubwAxjq92hgNy6jAFmHPQ79OzqGx3ccF3iN9/f0Wfcxy+je1n3j9Df60WRaND1nj3m4/vok2HSIQRWUhcCjUaCbTYd4gsIL2y1ch1CQ9EymBUB55HS6iFWP1eCo9pExD2krEm9"

  tags = {
    Name = "gcloud13_ssh_key"
  }
}

# vm 1
resource "aws_network_interface" "vm_01_ni" {
  subnet_id       = aws_subnet.nodes_subnet.id
  private_ips     = ["10.0.1.100"]
  security_groups = [aws_security_group.nodes_sg.id]

  tags = {
    Name = "gcloud13_vm_01_ni"
  }
}

resource "aws_eip" "vm_01_eip" {
  vpc                       = true
  network_interface         = aws_network_interface.vm_01_ni.id
  depends_on                = [aws_internet_gateway.vpc_igw, aws_instance.vm_01_instance]
  associate_with_private_ip = "10.0.1.100"

  tags = {
    Name = "gcloud13_vm_01_eip"
  }
}

resource "aws_instance" "vm_01_instance" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.small"
  key_name      = aws_key_pair.ssh_key.id

  network_interface {
    network_interface_id = aws_network_interface.vm_01_ni.id
    device_index         = 0
  }

  tags = {
    Name = "gcloud13_vm_01_instance"
  }
}

