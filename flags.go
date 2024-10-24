package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

var (
	// directory module
	directory *string

	// file module
	file *string

	// url module
	url     *string
	urlList *string

	//github module
	github      *bool
	from        *string
	to          *string
	githubGists *bool

	//gitlab module
	gitlab *bool

	//phishtank module
	phishtank *bool

	//recursively crawl URLs on page.
	recurse *bool

	// output file name
	outputFile *string
)

func customUsage() {
	fmt.Println("\nSecretsnitch\nhttps://github.com/0x4f53/secretsnitch\n")
	fmt.Println("A lightning-fast secret scanner in Golang!\n")
	fmt.Fprintf(os.Stderr, "Usage:\n%s [input options] [output options]\n", os.Args[0])
	fmt.Println("")
	fmt.Println("Input (pick at least one):")
	fmt.Println("")
	fmt.Println("  --github                             Scan public GitHub commits from the past hour")
	fmt.Println("    --from                             (optional) Timestamp to start from (format: 2006-01-02-15)")
	fmt.Println("    --to                               (optional) Timestamp to stop at (format: 2006-01-02-15)")
	fmt.Println("")
	fmt.Println("  --github-gists                       Scan the last 100 public GitHub Gists")
	fmt.Println("")
	fmt.Println("  --gitlab                             Scan the last 100 public GitLab commits")
	fmt.Println("")
	fmt.Println("  --phishtank                          Scan reported phishtank.org URLs from the past day")
	fmt.Println("")
	fmt.Println("  --url=<http://url>                   A single URL to scan")
	fmt.Println("  --urlList=<file containing URLs>     A line-separated file containing a list of URLs to scan for secrets")
	fmt.Println("")
	fmt.Println("  --directory=<directory/with/files/>  Scan an entire directory")
	fmt.Println("  --file=<file.js>  					Scan a file")
	fmt.Println("")
	fmt.Println("Optional arguments:")
	fmt.Println("")
	fmt.Println("  --output                             Save scan output to a custom location")
	fmt.Println("")
	fmt.Println("  --recurse                   Crawl URLs and hyperlinks inside page to 1 level")
	fmt.Println("")
}

func setFlags() {
	github = pflag.Bool("github", false, "")
	from = pflag.String("from", "", "")
	to = pflag.String("to", "", "")

	githubGists = pflag.Bool("github-gists", false, "")

	url = pflag.String("url", "", "")
	urlList = pflag.String("urlList", "", "")
	directory = pflag.String("directory", "", "")
	file = pflag.String("file", "", "")

	gitlab = pflag.Bool("gitlab", false, "")

	phishtank = pflag.Bool("phishtank", false, "")
	recurse = pflag.Bool("recurse", false, "")

	outputFile = pflag.String("output", defaultOutputDir, "")

	pflag.Usage = customUsage

	pflag.Parse()

	if !*github && !*gitlab && !*phishtank && *url == "" && *urlList == "" && *directory == "" && *file == "" && !*githubGists {
		pflag.Usage()
		fmt.Println("Come on, you'll have to pick some option!")
		os.Exit(-1)
	}

}
