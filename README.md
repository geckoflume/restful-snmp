# restful-smnp

`restful-snmp` provides a handy way to query SNMP via a HTTP.

### Install and run

To install:

```sh
$ go get github.com/unprofession-al/restful-snmp
```

To run:

```sh
$ RS_LISTEN=0.0.0.0 RS_PORT=3030 /path/to/restful-snmp
```

### Query a OID of a device
To issue a SNMP GET, fetch the API like so:
```sh
$ curl -X GET http://server:port/[device_to_query]/[oid_to_fetch]
```
Where `[device_to_query]` is an IP address or an FQDN of a device (eg. a switch, router etc.) that runs an SNMP agent and `[oid_to_fetch]` is the OID that should be queried.

The response in case of a successful SNMP GET query comes as a HTTP 200 containing JSON of the following schema:
```json
{
    Name: ".1.3.6.1.4.1.6.3.16.1.1.1.1.4.49.51.51.52",
    Type: 4,
    Value: "1334"
}
```
If you only care about the value just append `only_value=true` to the value string.

### Specify a SNMP community other that `public`
The community that is used by default is `public`. To specify an custom community specify the name via query string:
```sh
$ curl -X GET http://server:port/[device_to_query]/[oid_to_fetch]?community=[community_name]
```
