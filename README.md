# Summary

The RSA Netwitness Monitor Service provides a simple REST interface to be able
to retrieve the Risc Score of a Machine being queried.

A call is performed using:
http://apihost:8080/machine/server2016

A JSON Object is returned with the Machine Name, IP Address and Risc Score along
with a Status code. e.g.:

{"Machine":"server2016","IP":"192.168.1.132","Score":"0","Status":"200"}

## Documentation
- [Configuration](/docs/config.md).

## Supported platform
- Linux

## Build
Given that the Go Language compiler (version 1.11 or above preferred) is installed, you can build it with:
```
go get github.com/hwahrmann/rsa-nw-monitor
cd $GOPATH/src/github.com/hwahrmann/rsa-nw-monitor
make build
```
The binary is then in the subfolder named monitor.

To check the version:
```
./rsa-nw-monitor -version
```

Parameters are specified in a config file, like described in [Configuration](/docs/config.md).
A sample is in the conf folder, and automatically installed in /etc/rsa-nw-monitor:
```
./rsa-nw-monitor -config myconfig.conf
```

## Installation
You can download and install a pre-built rpm package as below ([RPM](https://github.com/hwahrmann/rsa-nw-monitor/releases)).

```
yum localinstall rsa-nw-monitor-1.0.0-1.x86_64.rpm
```

Once you installed you need to configure some basic parameters, for more information check [[Configuration](/docs/config.md):
```
/etc/rsa-nw-monitor/rsa-nw-monitor.conf
```
You can start the service by the below:
```
service rsa-nw-monitor start
```

## License
Licensed under the Apache License, Version 2.0 (the "License")

## Contribute
Welcomes any kind of contribution, please follow the next steps:

- Fork the project on github.com.
- Create a new branch.
- Commit changes to the new branch.
- Send a pull request.
