# ec2-internal-dns-resolver

Resolve ip-xx-xx-xx-xx.ec2.internal domains, proxy all other requests to 1.1.1.1

### Building
```
make`
```

### Running
1) Set DNS to 127.0.0.1
1) `./bin/ec2-internal-dns-resolver`