# Conditional

_Conditional_ is an example/experimental CoreDNS plugin that implements the interfaces
defined by two somewhat independent POC state CoreDNS features, neither of which are part
of CoreDNS proper:

* **CoreDNS Advanced Routing**: with which you can define criteria that control to which server blocks
  queries are routed. (requires: https://github.com/chrisohaver/coredns/tree/views)
* **Conditional forwarding** via pluggable _forward_ plugin policies: with which you
  can define a forward policy based on a user expression that can be used by the forward plugin.
  (requires: https://github.com/chrisohaver/coredns/tree/fwd-poliplug).

## CoreDNS Advanced Routing

This option controls how CoreDNS will route queries to the enclosing server block.
Using this option requires view-capable CoreDNS (https://github.com/chrisohaver/coredns/tree/views).

**Note:** See the [advanced-routing](https://github.com/chrisohaver/conditional/tree/advanced-routing) branch for a version of this plugin that exclusively implements this
feature, and renames the plugin to "serve".


### Syntax
```
conditional {
    view EXPRESSION
}
```

* `view` **EXPRESSION** - CoreDNS will not route incoming queries to the enclosing server block
  if any **EXPRESSION** evaluates to false. See the **Expressions** section below for available variables and functions.
  

### CoreDNS Views Example

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

## Conditional _forward_ Policy

These options define an expression based forward policy that can be used by the policy-pluggable _forward_ plugin.
This requires policy-pluggable _forward_ plugin (https://github.com/chrisohaver/coredns/tree/fwd-poliplug).

### Syntax
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


### Pluggable _forward_ Policy Example

The following (abstract) example defines 3 groups, each containing a single upstream server.
It defines three rules.  When forward uses the `conditional` policy, these rules are
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

### Available Functions

* `incidr(ip,cidr)`: returns true if _ip_ is within _cidr_ 
