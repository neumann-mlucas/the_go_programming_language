package main

import (
	"fmt"
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

func main() {
	mt := MultiTierSort{tracks, []string{"Year", "Title"}}
	sort.Sort(mt)
	printTracks(tracks)
}
