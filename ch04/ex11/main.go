package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func invalidArgument() {
	var usage = "%s\n usage:\t[read|close|open] OWNER REPO ISSUE_NUMBER\n"
	log.Fatalf(usage, os.Args)
}

func main() {
	if len(os.Args) != 5 {
		invalidArgument()
	}
	cmd, args := os.Args[1], os.Args[2:]
	url := buildURL(args)

	var issue *Issue
	var err error

	switch cmd {
	case "read":
		issue, err = GetIssue(url)
	case "close":
		issue, err = EditIssue(url, map[string]string{"state": "closed"})
	case "open":
		issue, err = EditIssue(url, map[string]string{"state": "open"})
	default:
		invalidArgument()
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(issue)
}

func buildURL(args []string) string {
	owner, repo, number := args[0], args[1], args[2]
	return strings.Join([]string{APIURL, "repos", owner, repo, "issues", number}, "/")
}

const IssuesURL = "https://api.github.com/search/issues"
const APIURL = "https://api.github.com"

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

func (i Issue) String() string {
	return fmt.Sprintf("\nUSER: %s\n -> [%05d] [%s] %s\n%s",
		i.User.Login, i.Number, i.State, i.Title, i.Body)
}

func GetIssue(url string) (*Issue, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("can't get %s: %s", url, resp.Status)
	}

	var issue Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return nil, err
	}
	return &issue, nil
}

func EditIssue(url string, fields map[string]string) (*Issue, error) {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(fields)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", url, buf)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(os.Getenv("GITHUB_USER"), os.Getenv("GITHUB_PASS"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to edit issue: %s", resp.Status)
	}

	var issue Issue
	if err = json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return nil, err
	}
	return &issue, nil
}
