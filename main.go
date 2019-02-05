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

	var pr github.PullRequest
	if err := json.NewDecoder(f).Decode(&pr); err != nil {
		log.Fatalf("Error decoding GitHub event: %q", err)
	}

	state := *pr.State
	if os.Getenv("DEBUG") != "" {
		log.Printf("PR state: %s\n", state)
		c, err := ioutil.ReadFile("/github/workflow/event.json")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf(string(c))
	}

	if state != "open" {
		log.Printf("Ignore GitHub event state: %q", state)
		return
	}

	u, err := url.Parse("https://api.unsplash.com/photos/random")
	if err != nil {
		log.Fatalf("Error parsing Unsplash URL: %q", err)
	}
	query := os.Getenv("UNSPLASH_QUERY")
	if query != "" {
		u.Query().Set("query", url.QueryEscape(query))
	}
	orientation := os.Getenv("UNSPLASH_ORIENTATION")
	if orientation != "" {
		u.Query().Set("orientation", orientation)
	}

	req, err := http.NewRequest("GET", u.String(), nil)
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

	body := *pr.Body
	body += fmt.Sprintf("\n\n![](%s)\n> Photo by [%s](%s) on [Unsplash](%s)", ur.URLs.Regular, ur.User.Name, ur.User.Links.HTML)
	pr.Body = github.String(body)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	repo := strings.SplitN(os.Getenv("GITHUB_REPOSITORY"), "/", 2)
	if _, _, err := client.PullRequests.Edit(ctx, repo[0], repo[1], *pr.Number, &pr); err != nil {
		log.Fatalf("Error updating the Pull Request with the photo: %q", err)
	}
}
