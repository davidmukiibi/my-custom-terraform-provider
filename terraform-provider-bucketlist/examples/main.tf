terraform {
  required_providers {
    bucketlist = {
      version = "0.2"
      source  = "hashicorp.com/edu/bucketlist"
    }
  }
}

provider "bucketlist" {}

variable "user_name" {
  type    = string
  default = "david"
}

data "bucketlist_users" "all" {}

# Returns all users
output "all_users" {
  value = data.bucketlist_users.all.users
}

# Only returns packer spiced latte
output "chosen" {
  value = {
    for user in data.bucketlist_users.all.users :
    user.surname => user
    if user.first_name == var.user_name
  }
}

resource "bucketlist_user" "david" {
    first_name = "David"
    surname = "Mukiibi"
    email = "david.mukiibi@gmail.com"
    password = "1234567890"
}

output "new_user" {
  value = bucketlist_user.david
}