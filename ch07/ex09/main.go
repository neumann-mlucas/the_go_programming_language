package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
)

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

func printTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	tw.Flush() // calculate column widths and print table
}

type MultiTierSort struct {
	items    []*Track
	priority []string
}

func (m MultiTierSort) Len() int      { return len(m.items) }
func (m MultiTierSort) Swap(i, j int) { m.items[i], m.items[j] = m.items[j], m.items[i] }
func (m MultiTierSort) Less(i, j int) bool {
	for idx := range m.priority {
		field := m.priority[len(m.priority)-1-idx]
		comp := compField(m.items[i], m.items[j], field)
		if comp > 0 {
			return true
		} else if comp < 0 {
			return false
		}
	}
	return true
}

func compField(a, b *Track, field string) int {
	switch field {
	case "Title":
		return strings.Compare(a.Title, b.Title)
	case "Artist":
		return strings.Compare(a.Artist, b.Artist)
	case "Album":
		return strings.Compare(a.Album, b.Album)
	case "Year":
		return a.Year - b.Year
	case "Length":
		return int(a.Length.Microseconds() - b.Length.Microseconds())
	default:
		return 0
	}
}

func (m MultiTierSort) Select(col string) []string {
	m.priority = append(m.priority, col)
	if len(m.priority) > 5 {
		m.priority = m.priority[len(m.priority)-6:]
	}
	return m.priority
}

var html = template.Must(template.New("tracks").Parse(`
<html>
<body>

<table>
	<tr>
		<th><a href="?sort=Title">Title</a></th>
		<th><a href="?sort=Artist">Artist</a></th>
		<th><a href="?sort=Album">Album</a></th>
		<th><a href="?sort=Year">Year</a></th>
		<th><a href="?sort=Length">Length</a></th>
	</tr>
{{range .}}
	<tr>
		<td>{{.Title}}</td>
		<td>{{.Artist}}</td>
		<td>{{.Album}}</td>
		<td>{{.Year}}</td>
		<td>{{.Length}}</td>
	</td>
{{end}}
</body>
</html>
`))

func main() {
	mt := MultiTierSort{tracks, []string{}}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sortCol := r.FormValue("sort")
		p := mt.Select(sortCol)
		// no ideia why I can only set the slice on the level of the handler
		mt.priority = p
		sort.Sort(mt)
		err := html.Execute(w, tracks)
		if err != nil {
			log.Printf("template error: %s", err)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
