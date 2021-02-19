# Serve

_Serve_ is an example CoreDNS plugin that allows the routing of queries to be
controlled by user defined expressions.

This plugin requires view-capable CoreDNS (https://github.com/chrisohaver/coredns/tree/views).

## Syntax
```
serve {
    EXPRESSION
}
```

* **EXPRESSION** - CoreDNS will not route incoming queries to the enclosing server block
  if any **EXPRESSION** evaluates to false. Multiple **EXPRESSION** can be given on separate lines.
  See the **Expression Syntax** section below for available variables and functions.
  
CoreDNS will only route a query to a given server block if the query falls within the server block's
zone (per normal routing behavior), and all **EXPRESSION** listed in *serve* evaluate to true.

## Examples

The abstract example below implements CIDR based split DNS routing.  It will return a different
answer for `test.` depending on client's IP address.  It returns ...
* `test. 3600 IN A 1.1.1.1`, for queries with a source address in 127.0.0.0/24
* `test. 3600 IN A 2.2.2.2`, for queries with a source address in 192.168.0.0/16
* `test. 3600 IN A 3.3.3.3`, for all others

```
. {
  serve {
    incidr(client_ip, '127.0.0.0/24')
  }
  hosts {
    1.1.1.1 test
  }
}

. {
  serve {
    incidr(client_ip, '192.168.0.0/16')
  }
  hosts {
    2.2.2.2 test
  }
}

. {
  hosts {
    3.3.3.3 test
  }
}
```

## Expression Syntax

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