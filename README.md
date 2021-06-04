# PBJelly
[![Go Report Card](https://goreportcard.com/badge/github.com/the-furry-hubofeverything/PBJelly)](https://goreportcard.com/report/github.com/the-furry-hubofeverything/PBJelly)[![Go Reference](https://pkg.go.dev/badge/github.com/the-furry-hubofeverything/PBJelly.svg)](https://pkg.go.dev/github.com/the-furry-hubofeverything/PBJelly)

Updates DNS records using the Porkbun API 
## Examples 
`PBJelly -l example.com`
```
2021/06/03 22:34:34 Current IP: 203.0.113.32
2021/06/03 22:34:34 {ID:123456789 Name:test.example.com Type:A Content:203.0.113.22 TTL:600 Prio:0 Notes:}
```
`PBJelly -id "123456789" -i "3h" -n "test.example.com" example.com`
```
2021/06/04 01:35:33 Updated test.example.com to 203.0.113.32
```
You can also use config files, note that any commandline flags would override the values in the config file 
 
`PBJelly -c config/example.toml example.com`
## Arguments
```
Usage: PBJelly [options] <domain>
   -c string
        Path to config (default "config.toml")
  -i duration
        Time between updates (default 1h0m0s)
  -id string
        ID of Porkbun DNS entry, use -l to view ID (default "0")
  -l    List DNS records under domain
  -n string
        Name of DNS record
  -o    Execute update only once (For cron or systemd)
  -t string
        Type of DNS record (default "A")
```
