package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	targetFileSuffixes = []string{".txt", ".go"}
)

func main() {
	start := time.Now().UnixNano()
	lineCounts, _ := Traverse(os.Args[1])
	for k, v := range lineCounts {
		fmt.Printf("%s : %v\n", k, v)
	}
	end := time.Now().UnixNano()
	fmt.Println("time cost", end-start)
}

func merge(dest map[string]int64, source map[string]int64) error {
	for k, v := range source {
		dest[k] = v
	}
	return nil
}

func isTargetFile(filePath string) bool {
	for _, targetSuffix := range targetFileSuffixes {
		if strings.HasSuffix(filePath, targetSuffix) {
			return true
		}
	}
	return false
}

func countLines(filePath string) int64 {
	cmd := exec.Command("wc", []string{"-l", filePath}...)
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	reader := bufio.NewReader(stdout)
	line, _ := reader.ReadString('\n')
	cmd.Wait()
	n, _ := strconv.Atoi(strings.Split(line, " ")[0])
	return int64(n)
}

func Traverse(dirPth string) (map[string]int64, error) {
	result := make(map[string]int64)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return result, err
	}
	PthSep := string(os.PathSeparator)
	files := []string{}
	for _, fi := range dir {
		fullPath := dirPth + PthSep + fi.Name()
		if fi.IsDir() {
			subResult, err := Traverse(fullPath)
			if err == nil {
				merge(result, subResult)
			}
			continue
		}
		if isTargetFile(fullPath) {
			files = append(files, fullPath)
			result[fullPath] = countLines(fullPath)
		}
	}
	return result, nil
}
