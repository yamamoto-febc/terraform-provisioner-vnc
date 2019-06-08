# Terraform VNC Provisioner

`terraform-provisioner-vnc` provisions Terraform resources via VNC.

[![Build Status](https://travis-ci.org/yamamoto-febc/terraform-provisioner-vnc.svg?branch=master)](https://travis-ci.org/yamamoto-febc/terraform-provisioner-vnc)

## Installation

Download binary from [Releases](https://github.com/yamamoto-febc/terraform-provisioner-vnc/releases/latest), then put it to `~/.terraform.d/plugins/`.

## Usage

```tf
resource your_resource "example" {
  // ...

  provisioner vnc {
    host      = var.vnc_server_host
    port      = var.vnc_server_port
    password  = var.vnc_password
    # timeout   = "5m"
    # boot_wait = "30s" # default: 0s(no wait)

    inline = [
      "<wait5>",
      "user<enter>",      # login prompt: username
      "password<enter>",  # login prompt: password
      "echo /etc/os-release<enter>",  # exec your scripts
    ]
    # script = "path/to/your/script"
    # scripts = [
    #   "path/to/your/script1",
    #   "path/to/your/script2",
    # ]  
  }  
```

## Argument Reference

The following arguments are supported:
   
- `host`(string): The host name or address of VNC Server.
- `port`(int): The port number of VNC Server. 
- `password`(string): The VNC password.
- `boot_wait`(string): The boot_wait to wait for the connection to become available. Should be provided as a string like `30s` or `5m`. 
- `timeout`(string): The timeout to wait for the complete to execute scripts. This defaults to 5 minutes. Should be provided as a string like `30s` or `5m`. 
- `inline`(list of string): This is a list of command strings. They are executed in the order they are provided. This cannot be provided with `script` or `scripts`.
- `script`(string): This is a path (relative or absolute) to a local script that will be copied to the remote resource and then executed. This cannot be provided with `inline` or `scripts`.
- `scripts`(list of string): This is a list of paths (relative or absolute) to local scripts that will be copied to the remote resource and then executed. They are executed in the order they are provided. This cannot be provided with `inline` or `script`.

## License

 `terraform-provisioner-vnc` Copyright (C) 2019 Kazumichi Yamamoto.

  This project is published under [Apache 2.0 License](LICENSE.txt).
  
## Author

  * Kazumichi Yamamoto ([@yamamoto-febc](https://github.com/yamamoto-febc))