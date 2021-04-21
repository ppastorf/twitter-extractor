# base ami for instances
data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/*/ubuntu-focal-20.04-amd64-server-*"]
  }

  # https://docs.aws.amazon.com/pt_br/AWSEC2/latest/UserGuide/virtualization_types.html
  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

# cluster virtual network
resource "aws_vpc" "cluster_vpc" {
  cidr_block       = "10.0.0.0/16"
  instance_tenancy = "default"
  enable_dns_hostnames = true

  tags = {
    Name = "scc0158_cluster_vpc"
  }
}

resource "aws_security_group" "cluster_sg" {
  name   = "scc0158_cluster_sg"
  vpc_id = "${aws_vpc.cluster_vpc.id}"

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
    Name = "scc0158_cluster_sg"
  }
}

resource "aws_subnet" "cluster_subnet" {
  vpc_id            = "${aws_vpc.cluster_vpc.id}"
  cidr_block        = "10.0.1.0/24"
  map_public_ip_on_launch = true

  tags = {
    Name = "scc0158_cluster_subnet"
  }
}

resource "aws_internet_gateway" "cluster_igw" {
  vpc_id = "${aws_vpc.cluster_vpc.id}"

  tags = {
    Name = "scc0158_cluster_igw"
  }
}

resource "aws_route_table" "cluster_route_table" {
  vpc_id = "${aws_vpc.cluster_vpc.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.cluster_igw.id
  }

  route {
    ipv6_cidr_block = "::/0"
    gateway_id = aws_internet_gateway.cluster_igw.id
  }

  tags = {
    Name = "scc0158_cluster_route_table"
  }
}

resource "aws_key_pair" "cluster_keypair" {
  key_name   = "scc0158_cluster_keypair"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCt7leg0VkmuFOzvAJNFy7xCFLOAujJl1bcHd6yP3t4mG91h0ea/Ezr2FVGktc62I44rBWflZ6Lw+E0RmeHILUGDV49sPs7TMi43Dw3JHeNg4ek6Xleu4oGY6Nh+76vG+5jgZFbNIYTNX28mc1Q8G+ZUYXzDX4DyzSpNCmPv1owUv+7KdPjxQr9Op0fWBWkKK/ilsI/4YlnJW1KeiubwAxjq92hgNy6jAFmHPQ79OzqGx3ccF3iN9/f0Wfcxy+je1n3j9Df60WRaND1nj3m4/vok2HSIQRWUhcCjUaCbTYd4gsIL2y1ch1CQ9EymBUB55HS6iFWP1eCo9pExD2krEm9"

  tags = {
    Name = "scc0158_cluster_keypair"
  }
}

# cluster member 01
resource "aws_eip" "cluster_01_eip" {
  vpc                       = true
  network_interface         = "${aws_network_interface.cluster_01_ni.id}"
  depends_on                = [aws_internet_gateway.cluster_igw]
  associate_with_private_ip = "10.0.1.100"

  tags = {
    Name = "scc0158_cluster_01_eip"
  }
}

resource "aws_network_interface" "cluster_01_ni" {
  subnet_id   = "${aws_subnet.cluster_subnet.id}"
  private_ips = ["10.0.1.100"]
  security_groups = ["${aws_security_group.cluster_sg.id}"]

  tags = {
    Name = "scc0158_cluster_01_ni"
  }
}

resource "aws_instance" "cluster_01_instance" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"
  key_name      = aws_key_pair.cluster_keypair.id

  network_interface {
    network_interface_id = "${aws_network_interface.cluster_01_ni.id}"
    device_index         = 0
  }

  tags = {
    Name = "scc0158_cluster_01_instance"
  }
}

