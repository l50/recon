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

### Output ASNs and IP ranges to files using anew and bash
```
recon asn -t recon_targets.txt | while read -r line; do ([[ $line =~ ^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\/[0-9]{2}$ ]] && echo $line | anew ip_ranges.txt) || ([[ $line =~ ^[0-9]{5}$ ]] && echo $line | anew asns.txt); done
```

`recon_targets.txt`:
```
target1.com
target2.com
```