package secretsnitch

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

var (
	// directory module
	directory *string

	// url module
	url     *string
	urlList *string

	//github module
	github *bool
	from   *string
	to     *string

	//gitlab module
	gitlab *bool

	//phishtank module
	phishtank *bool

	// output file name
	outputFile *string
)

func customUsage() {
	fmt.Println("\nSecretsnitch\nhttps://github.com/0x4f53/secretsnitch\n")
	fmt.Println("A lightning-fast secret scanner in Golang!\n")
	fmt.Fprintf(os.Stderr, "Usage:\n%s [input options] [output options]\n", os.Args[0])
	fmt.Println("")
	fmt.Println("Input options:")
	fmt.Println("")
	fmt.Println("  --github           Scan public GitHub commits from the past hour")
	fmt.Println("    --from           (optional) Timestamp to start from (format: 2006-01-02-15)")
	fmt.Println("    --to             (optional) Timestamp to stop at (format: 2006-01-02-15)")
	fmt.Println("")
	fmt.Println("  --gitlab           Scan the last 100 public GitLab commits")
	fmt.Println("")
	fmt.Println("  --phishtank        Scan reported phishtank.org URLs from the past day")
	fmt.Println("")
	fmt.Println("  --url              Single URL to scan")
	fmt.Println("  --urlList          A file containing a list of URLs to scan for secrets")
	fmt.Println("")
	fmt.Println("  --directory        Scan an entire directory")
	fmt.Println("")
	fmt.Println("Output options:")
	fmt.Println("")
	fmt.Println("  --output           Save scan output to file")
	fmt.Println("")
}

func setFlags() {
	github = pflag.Bool("github", false, "")
	from = pflag.String("from", "", "")
	to = pflag.String("to", "", "")

	url = pflag.String("url", "", "")
	urlList = pflag.String("urlList", "", "")
	directory = pflag.String("directory", "", "")

	gitlab = pflag.Bool("gitlab", false, "")

	phishtank = pflag.Bool("phishtank", false, "")

	outputFile = pflag.String("output", defaultOutputDir, "")

	pflag.Usage = customUsage

	pflag.Parse()

	if !*github && !*gitlab && !*phishtank && *url == "" && *urlList == "" && *directory == "" {
		pflag.Usage()
		fmt.Println("Come on, you've got to pick something!")
		os.Exit(-1)
	}

}
