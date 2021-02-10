# Conditional

_Conditional_ is an example/experimental CoreDNS plugin allows a user to define
boolean expressions for use by other CoreDNS plugins/functions.  It currently
only interfaces with two POC state features, neither of which are part of
CoreDNS standard:

* CoreDNS Views (requires: https://github.com/chrisohaver/coredns/tree/views) 
* Conditional forwarding via pluggable _forward_ plugin policies (requires: https://github.com/chrisohaver/coredns/tree/fwd-poliplug).



### Available Variables

* `type`: type of the request (A, AAAA, TXT, ...)
* `name`: name of the request (the domain name requested)
* `class`: class of the request (IN, CH, ...)
* `proto`: protocol used (tcp or udp)
* `client_ip`: client's IP address, for IPv6 addresses these are enclosed in brackets: `[::1]`
* `size`: request size in bytes
* `port`: client's port
* `bufsize`: the EDNS0 buffer size advertised in the query
* `do`: the EDNS0 DO (DNSSEC OK) bit set in the query
* `id`: query ID
* `opcode`: query OPCODE
* `server_ip`: server's IP address; for IPv6 addresses these are enclosed in brackets: `[::1]`
* `server_port` : client's port

### Available Functions 

* `incidr(ip,cidr)`: returns true if _ip_ is within _cidr_ 

## CoreDNS Views Example

Requires view-capable CoreDNS (https://github.com/chrisohaver/coredns/tree/views).

```
.:5399 {
  conditional {
    view incidr(client_ip, '127.0.0.0/24')
  }
  hosts {
    1.2.3.4 test
  }
}

.:5399 {
  conditional {
    view incidr(client_ip, '192.168.0.0/16')
  }
  hosts {
    5.6.7.8 test
  }
}
```

## Pluggable _forward_ Policy Example

Requires policy-pluggable _forward_ plugin (https://github.com/chrisohaver/coredns/tree/fwd-poliplug).

```
.:5399 {
  conditional {
    group a 0
    group b 1
    group c 2
    use c if incidr(client_ip, '127.0.0.0/24') 
    use a if type == 'A'
    use b if type == 'AAAA'
  }
  forward . 127.0.0.1:5390 127.0.0.1:5391  127.0.0.1:5392 {
    policy conditional
  }
}

.:5390 {
  hosts {
    1.2.3.4 a
  }
}

.:5391 {
  hosts {
    0::5:6:7:8 a
  }
}

.:5392 {
  hosts {
    9.9.9.9 a
  }
}

```