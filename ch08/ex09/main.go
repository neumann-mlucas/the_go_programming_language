package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	// "time"
)

const sizeCutOff = 1e8

var vFlag = flag.Bool("v", false, "show verbose progress messages")

type FileInfo struct {
	path   string
	size   int64
	nfiles int
}

func (f FileInfo) inDir(root string) bool {
	return strings.HasPrefix(f.path, root)
}

func printDiskUsageByDir(f []FileInfo) {
	dirs := make(map[string]FileInfo)

	for _, file := range f {
		stat, err := os.Stat(file.path)
		if err != nil {
			continue
		}
		if !stat.IsDir() {
			continue
		}
		dirs[file.path] = file
	}

	for _, file := range f {
		for dirPath, dir := range dirs {
			if file.inDir(dirPath) {
				newDir := FileInfo{dirPath, dir.size + file.size, dir.nfiles + 1}
				dirs[dirPath] = newDir
			}
		}

	}

	for _, dir := range dirs {
		if dir.size > sizeCutOff {
			fmt.Printf("%10d files  %6.1f MB | %s\n", dir.nfiles, float64(dir.size)/1e6, dir.path)
		}
	}
}

func main() {
	flag.Parse()

	// Determine the initial directories.
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Create Container to hold File Information
	var fileInfos []FileInfo

	// Traverse each root of the file tree in parallel.
	fileSizes := make(chan FileInfo)
	var wg sync.WaitGroup
	for _, root := range roots {
		wg.Add(1)
		go walkDir(root, &wg, fileSizes)
	}
	go func() {
		wg.Wait()
		close(fileSizes)
	}()

	var nfiles, nbytes int64
	for info := range fileSizes {
		fileInfos = append(fileInfos, info)
		nfiles++
		nbytes += info.size
	}

	printDiskUsageByDir(fileInfos)
	printDiskUsage(nfiles, nbytes) // final totals
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("\n%10d files  %6.1f GB | TOTAL\n", nfiles, float64(nbytes)/1e9)
}

func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- FileInfo) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			fileSizes <- FileInfo{dir, 0, 0}
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, fileSizes)
		} else {
			fileInfo, err := entry.Info()
			if err == nil {
				fileSizes <- FileInfo{dir, fileInfo.Size(), 1}
			}
		}
	}
}

var sema = make(chan struct{}, 20)

func dirents(dir string) []os.DirEntry {
	sema <- struct{}{}        // acquire token
	defer func() { <-sema }() // release token

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return entries
}
