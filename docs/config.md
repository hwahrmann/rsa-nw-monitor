# Monitor configuration

## Format

A config file is a plain text file in [YAML](https://en.wikipedia.org/wiki/YAML) format.

## Configuration Keys
The Monitor configuration supports the following keys. If a key is not specified the default is taken.

|Key                     | Default                        | Description                                      |
|------------------------| -------------------------------|--------------------------------------------------|
|verbose                 | false                          | log output to stdout                             |
|pid-file                | /var/run/rsa-nw-monitor.pid    | file in which server should write its process ID |
|endpointserver          | N/A                            | The address of the RSA Netwitness EP Server      |
|monitorport             | 8080                           | The port to listen for REST calls                |
|user                    | apiuser                        | A read-only user or the deploy_admin user        |
|password                | netwitness                     | The password of the above user                   |

The default configuration path is /etc/rsa-nw-monitor/rsa-nw-monitor.conf but you can change it as below:
```
rsa-nw-monitor -config /usr/local/etc/monitor.conf
```
To show version information use:
```
rsa-nw-monitor -version
```

## Creating a Read-Only user on the EndPoint Server

If a different user than the deploy_admin should be used, this is the way to define a read-only user.
ssh to the EndPointServer and perform the following

```
mongo admin -u deploy_admin
db.createUser({user: "apiuser", pwd: "netwitness", roles: [ { role: "read", db: "endpoint-server"}]})
```

Then change the Config file accordingly
