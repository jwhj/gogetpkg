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
	// Parse flags
	var nd = flag.Bool("nd", false, "lsfdkjf")
	var pkg = flag.String("p", "", "")
	flag.Parse()
	if *pkg == "" {
		return
	}
	fmt.Println(*pkg)
	var lst = strings.SplitN(*pkg, "/", 4)

	// Chdir
	var plst = strings.Split(os.Getenv("GOPATH"), ":")
	os.MkdirAll(plst[0]+"/src/"+lst[0]+"/"+lst[1], 0700)
	os.Chdir(plst[0] + "/src/" + lst[0] + "/" + lst[1])
	fmt.Println(os.Getwd())

	// Download
	if !*nd {
		var url = strings.Join(lst[:3], "/")
		for i := 0; i < len(mirrors); i++ {
			if strings.HasPrefix(url, mirrors[i][0]) {
				url = strings.Replace(url, mirrors[i][0], mirrors[i][1], -1)
			}
		}
		fmt.Println(url)
		var fName = lst[2] + ".zip"
		var cmd = exec.Command("wget", "https://"+url+"/archive/master.zip", "-c", "-O", fName)
		execWithOutput(cmd, 1)
		exec.Command("unzip", fName).Run()
		os.Remove(fName)
		os.RemoveAll(lst[2])
		os.Rename(lst[2]+"-master", lst[2])
	}

	// Install
	os.Chdir(strings.Join(lst[2:], "/"))
	var cmd = exec.Command("go", "install")
	execWithOutput(cmd, 1)
}
