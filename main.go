package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func execWithOutput(cmd *exec.Cmd, flag int) {
	var pipe io.ReadCloser
	if flag == 0 {
		pipe, _ = cmd.StdoutPipe()
	} else {
		pipe, _ = cmd.StderrPipe()
	}
	cmd.Start()
	var reader = bufio.NewReader(pipe)
	for {
		var str, _, err = reader.ReadLine()
		if err != nil || err == io.EOF {
			break
		}
		fmt.Println(string(str))
	}
	cmd.Wait()
}
func main() {
	// Parse pkg name
	var pkg = flag.String("p", "", "")
	flag.Parse()
	fmt.Println(*pkg)
	var lst = strings.SplitN(*pkg, "/", 4)

	// Chdir
	var plst = strings.Split(os.Getenv("GOPATH"), ":")
	os.Chdir(plst[0] + "/src/" + lst[0])
	os.MkdirAll(lst[1], 0700)
	os.Chdir("./" + lst[1])
	fmt.Println(os.Getwd())

	// Download
	var url = strings.Join(lst[:3], "/")
	if lst[0] == "golang.org" {
		url = strings.Replace(url, "golang.org/x", "github.com/golang", -1)
		fmt.Println(url)
	}
	var cmd = exec.Command("wget", "https://"+url+"/archive/master.zip", "-c")
	execWithOutput(cmd, 1)

	// Install
	exec.Command("unzip", "master.zip").Run()
	os.Remove("master.zip")
	os.RemoveAll(lst[2])
	os.Rename(lst[2]+"-master", lst[2])
	os.Chdir(strings.Join(lst[2:], "/"))
	cmd = exec.Command("go", "install")
	execWithOutput(cmd, 1)
}
