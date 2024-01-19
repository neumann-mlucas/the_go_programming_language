package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "root")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir) // Clean up
	fmt.Println("Temporary directory created:", tempDir)

	// Listen on TCP port 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		// Handle the connection
		go handleConnection(conn, tempDir)
	}
}

func handleConnection(conn net.Conn, dir string) {
	defer conn.Close()

	// Read from the connection
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Message received: %s", msg)
		cmd := strings.Fields(msg)
		if len(cmd) == 0 {
			continue
		}

		switch cmd[0] {
		case "ls":
			ls(conn, dir, cmd)
		case "touch":
			touch(conn, dir, cmd)
		case "cd":
			dir = cd(conn, dir, cmd)
		case "mkdir":
			mkdir(conn, dir, cmd)
		}
	}
}

func ls(conn net.Conn, dir string, cmd []string) {
	if !(len(cmd) == 1 || len(cmd) == 2) {
		fmt.Printf("Wrong number of arguments for ls: %s\n", cmd)
		fmt.Fprintln(conn, "Wrong number of arguments for ls")
		return
	}

	if len(cmd) == 2 {
		dir = strings.Join([]string{dir, cmd[1]}, "/")
	} else {
		dir += "/"
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(conn, err.Error())
		return
	}

	for _, file := range files {
		fmt.Fprintln(conn, file.Name())
	}
}

func touch(conn net.Conn, dir string, cmd []string) {
	if !(len(cmd) == 2) {
		fmt.Printf("Wrong number of arguments for touch: %s\n", cmd)
		fmt.Fprintln(conn, "Wrong number of arguments for touch")
		return
	}

	path := strings.Join([]string{dir, cmd[1]}, "/")
	fmt.Println(path)
	_, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(conn, err.Error())
		return
	}
	fmt.Fprintln(conn, path)
}

func cd(conn net.Conn, dir string, cmd []string) string {
	if !(len(cmd) == 2) {
		fmt.Printf("Wrong number of arguments for touch: %s\n", cmd)
		fmt.Fprintln(conn, "Wrong number of arguments for touch")
		return dir
	}
	path := strings.Join([]string{dir, cmd[1]}, "/")
	_, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(conn, err.Error())
		return dir
	}
	fmt.Fprintln(conn, path)
	return path

}

func mkdir(conn net.Conn, dir string, cmd []string) {
	if !(len(cmd) == 2) {
		fmt.Printf("Wrong number of arguments for touch: %s\n", cmd)
		fmt.Fprintln(conn, "Wrong number of arguments for touch")
	}
	path := strings.Join([]string{dir, cmd[1]}, "/")
	err := os.Mkdir(path, 0755)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(conn, err.Error())
		return
	}
	fmt.Fprintln(conn, path)
}
