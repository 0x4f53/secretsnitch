# Secretsnitch

A lightning-fast secret scanner in Golang!

```
Usage:
./secretsnitch [input options] [output options]

Input options:

  --github           Scan public GitHub commits from the past hour
    --from           (optional) Timestamp to start from (format: 2006-01-02-15)
    --to             (optional) Timestamp to stop at (format: 2006-01-02-15)

  --gitlab           Scan the last 100 public GitLab commits

  --phishtank        Scan reported phishtank.org URLs from the past day

  --url              Single URL to scan
  --urlList          A file containing a list of URLs to scan for secrets

  --directory        Scan an entire directory

Output options:

  --output           Save scan output to file
```