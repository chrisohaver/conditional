# conditional

_conditional_ defines an expression based forwarding policy that the _fwdpolicy_ plugin can use.

## Syntax

These options define an expression based forward policy that can be used by the policy-pluggable _fwdpolicy_ plugin.
(https://github.com/infobloxopen/fwdpolicy).

```
conditional {
    group GROUP-NAME UPSTREAM-INDEX ...
    use GROUP-NAME if EXPRESSION
}
```

* `group` - assigns a **GROUP-NAME** to a set of **UPSTREAM-INDEX**. **UPSTREAM-INDEX** are the
  integer index of the upstream defined in the forward plugin. e.g. If there are three upstreams
  defined by forward, then the index values are 0, 1, and 2.
* `use` - if the **EXPRESSION** evaluates to true for the incoming query, the forward policy will return upstream
  servers assigned to the **GROUP-NAME**. The forward plugin will then attempt to route the query to those upstreams.
  See the **Expressions** section below for available variables and functions.


## Example

The following (abstract) example defines 3 groups, each containing a single upstream server.
It defines three rules.  When _fwdpolicy_ uses the `conditional` policy, these rules are
evaluated...
* If the client IP address is local (in 127.0.0.0/24), it will forward to group `c` (127.0.0.1:5392)
* If the query type is `A`, it will forward to group `a` (127.0.0.1:5390)
* If the query type is `AAAA`, it will forward to group `b` (127.0.0.1:5391)

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
  fwdpolicy . 127.0.0.1:5390 127.0.0.1:5391  127.0.0.1:5392 {
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

## Expressions

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

#### Metadata

When you enable the _metadata_ plugin, metadata are available as a variables in expressions.
When including a metadata value in an expression, wrap the metadata name in square brackets e.g.`[example/metadata]`.

### Available Functions

* `incidr(ip,cidr)`: returns true if _ip_ is within _cidr_ 


## CoreDNS Advanced Routing

_conditional_ also defines an expression based view, with which you can define criteria that control to which server
blocks queries are routed.  Views are not supported in https://github.com/coredns/coredns. So to use this feature
you need apply changes to coredns code as applied in the branch: https://github.com/chrisohaver/coredns/tree/views.

### CoreDNS Advanced Routing - Syntax
```
conditional {
    view EXPRESSION
}
```

* `view` **EXPRESSION** - CoreDNS will not route incoming queries to the enclosing server block
  if any **EXPRESSION** evaluates to false. See the **Expressions** section for available variables and functions.
  Note metadata variables are not supported in view expressions.


### CoreDNS Advanced Routing - Example

The abstract example below implements CIDR based split DNS routing.  It will return a different
answer for `test.` depending on client's IP address.  It returns ...
* `test. 3600 IN A 1.1.1.1`, for queries with a source address in 127.0.0.0/24
* `test. 3600 IN A 2.2.2.2`, for queries with a source address in 192.168.0.0/16
* `test. 3600 IN A 3.3.3.3`, for all others

```
. {
  conditional {
    view incidr(client_ip, '127.0.0.0/24')
  }
  hosts {
    1.1.1.1 test
  }
}

. {
  conditional {
    view incidr(client_ip, '192.168.0.0/16')
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