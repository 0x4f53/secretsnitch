# Secretsnitch

A lightning-fast secret scanner in Golang!

## Features

### Fast and efficient

Secretsnitch is a rapid scanner for secrets, written in Golang. It utilizes concurrency via thousands of goroutines and downloads
and processes files, runs regular expression checks on them for secrets, grabs associated URLs and domains, tags them and performs tons of other operations concurrently.

Features like caching make sure files aren't re-downloaded. This speeds up the tool significantly while keeping network resource
consumption low.

### Modular

It supports tons of modules for popular online sources that may potentially contain thousands of secrets. This includes GitHub commits,
GitHub Gists, GitLab, Phishtank, random webpage scraping etc. More modules can be added this way, making the tool extremely dynamic.

### Smart and accurate

Secretsnitch doesn't just use regular expressions, it also scans files for common cloud provider strings, performs entropy checks on
captured strings, and gives you the variable/filename associated with a captured secret and gives you a precise indication on whether something may be an actual secret or not.

### Easy to use

Secretsnitch was designed to be easy to use, whether you are a pentester, bounty hunter or want to deploy it across your organization. The
tool can be run in singular commands as shown in the examples below.

### Community-driven

The signatures list is completely community-driven and is a combination of trial-and-error, Google searching, ChatGPT and from existing lists like that of GitGuardian. Pull requests for signature additions and corrections are welcome, and feel free to use these signatures in other cybersecurity tools you build.

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