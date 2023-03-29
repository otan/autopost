package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

var flagWatchForIssue = flag.Int("watch", 100000, "what issue to watch for")
var flagOrg = flag.String("org", "cockroachdb", "what org to watch")
var flagRepo = flag.String("repo", "cockroach", "what repo to watch")

func newString(s string) *string {
	return &s
}

func main() {
	flag.Parse()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_API_KEY")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	for {
		fmt.Printf("%s: waiting for issue\n", time.Now().Format(time.RFC3339))
		// list all repositories for the authenticated user
		_, resp, err := client.Issues.Get(ctx, *flagOrg, *flagRepo, *flagWatchForIssue-1)
		if err != nil {
			fmt.Printf("%s\n", err)
			if resp.Response.StatusCode == 404 {
				fmt.Printf("wow, 404\n")
				time.Sleep(time.Second * 5)
				continue
			}
			panic(err)
		}
		break
	}

	fmt.Printf("double check\n")
	_, resp, err := client.Issues.Get(ctx, *flagOrg, *flagRepo, *flagWatchForIssue)
	if err != nil && resp.Response.StatusCode == 404 {
		if _, _, err := client.Issues.Create(ctx, *flagOrg, *flagRepo, &github.IssueRequest{
			Title: newString("WE DID IT! #100,000!"),
			Body:  newString(`![wow](https://media.tenor.com/jBvWbVcN4ioAAAAC/owen-wilson-owen.gif)`),
		}); err != nil {
			panic(err)
		}
		fmt.Printf("issue posted!\n")
	} else {
		fmt.Printf("too late!\n")
	}
	fmt.Printf("done\n")
}
