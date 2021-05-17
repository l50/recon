# recon
## Install
```
make init && make
```

## ASNs
Get the ASNs and associated IP ranges from a list of domains:
```
recon asn -t recon_targets.txt 
```

`recon_targets.txt`:
```
target1.com
target2.com
```