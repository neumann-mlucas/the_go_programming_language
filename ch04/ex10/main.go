package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

const IssuesURL = "https://api.github.com/search/issues"

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      Issues
}

type Issues []*Issue

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

func (i Issues) Len() int {
	return len(i)
}

func (i Issues) Less(a, b int) bool {
	return i[a].CreatedAt.Sub(i[b].CreatedAt) > 0
}

func (i Issues) Swap(a, b int) {
	i[a], i[b] = i[b], i[a]
}

func (i Issue) String() string {
	dt := i.CreatedAt.Format("01-02-2006")
	return fmt.Sprintf("(%s) #%-5d  %9.9s %.55s\n", dt, i.Number, i.User.Login, i.Title)
}

// SearchIssues queries the GitHub issue tracker.
func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}
	// We must close resp.Body on all execution paths.
	// (Chapter 5 presents 'defer', which makes this simpler.)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}
	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}

func DaysSince(dt time.Time) int64 {
	return int64(time.Since(dt).Hours() / 24)
}

func PrintDateGroup(a, b *Issue) {
	daysA, daysB := DaysSince(a.CreatedAt), DaysSince(b.CreatedAt)

	// first iteration
	if a == b {
		if daysA < 30 {
			fmt.Println("\nLess Then 30 days:")
		} else if daysA < 60 {
			fmt.Println("\nMore Then 30 days:")
		} else {
			fmt.Println("\nMore Then 60 days:")
		}
		return
	}
	if daysA < 30 && daysB > 30 && daysB < 60 {
		fmt.Println("\n\nMore Then 30 days:")
	} else if daysA < 60 && daysB > 60 && daysA > 30 {
		fmt.Println("\n\nMore Then 60 days:")
	}

}
func main() {
	result, err := SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues\n", result.TotalCount)

	sort.Sort(result.Items)
	prev := result.Items[0]
	for _, item := range result.Items {
		PrintDateGroup(prev, item)
		fmt.Print(item)
		prev = item

	}
}
