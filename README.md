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