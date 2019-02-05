package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type prWrap struct {
	Action      string             `json:"action"`
	Number      int                `json:"number"`
	PullRequest github.PullRequest `json:"pull_request"`
}

type unsplashResponse struct {
	URLs struct {
		Regular string `json:"regular"`
	} `json:"urls"`
	User struct {
		Name  string `json:"name"`
		Links struct {
			HTML string `json:"html"`
		} `json:"links"`
	} `json:"user"`
}

func main() {
	eventName := os.Getenv("GITHUB_EVENT_NAME")
	if eventName != "pull_request" {
		log.Printf("Ignore GitHub event: %q", eventName)
		return
	}

	eventPath := os.Getenv("GITHUB_EVENT_PATH")
	f, err := os.Open(eventPath)
	if err != nil {
		log.Fatalf("Error opening GitHub event file: %q", err)
	}

	var pr prWrap
	if err := json.NewDecoder(f).Decode(&pr); err != nil {
		log.Fatalf("Error decoding GitHub event: %q", err)
	}

	var debug bool
	if os.Getenv("DEBUG") != "" {
		debug = true
		c, err := ioutil.ReadFile("/github/workflow/event.json")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf(string(c))

		log.Printf("%q\n", pr)
		log.Printf("%q\n", pr.PullRequest.State)
	}

	if pr.Action != "open" && !debug {
		log.Printf("Ignore GitHub event state: %q", pr.Action)
		return
	}

	u, err := url.Parse("https://api.unsplash.com/photos/random")
	if err != nil {
		log.Fatalf("Error parsing Unsplash URL: %q", err)
	}

	q := u.Query()
	query := os.Getenv("UNSPLASH_QUERY")
	if query != "" {
		q.Set("query", query)
	}
	orientation := os.Getenv("UNSPLASH_ORIENTATION")
	if orientation != "" {
		q.Set("orientation", orientation)
	}

	u.RawQuery = q.Encode()
	us := u.String()
	if debug {
		log.Printf("Unsplash URL: %s", us)
	}

	req, err := http.NewRequest("GET", us, nil)
	if err != nil {
		log.Fatalf("Error fetching photo from Unsplash: %q", err)
	}
	req.Header.Set("Authorization", "Client-ID "+os.Getenv("UNSPLASH_CLIENT_ID"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error fetching photo from Unsplash: %q", err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("Error fetching photo from Unsplash: status code %q", resp.StatusCode)
	}

	var ur unsplashResponse
	if err := json.NewDecoder(resp.Body).Decode(&ur); err != nil {
		log.Fatalf("Error decoding photo from Unsplash: %q", err)
	}

	body := *pr.PullRequest.Body
	body += fmt.Sprintf("\n\n![](%s)\n> Photo by [%s](%s?utm_source=splashed_pull_requests&utm_medium=referral) on [Unsplash](https://unsplash.com?utm_source=splashed_pull_requests&utm_medium=referral)", ur.URLs.Regular, ur.User.Name, ur.User.Links.HTML)
	pr.PullRequest.Body = github.String(body)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	repo := strings.SplitN(os.Getenv("GITHUB_REPOSITORY"), "/", 2)
	if _, _, err := client.PullRequests.Edit(ctx, repo[0], repo[1], pr.Number, &pr.PullRequest); err != nil {
		log.Fatalf("Error updating the Pull Request with the photo: %q", err)
	}
}
