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
	//fmt.Println("counts of files", len(lineCounts))
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

func countLines(filePath string, resultChannel chan string) {
	cmd := exec.Command("wc", []string{"-l", filePath}...)
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	reader := bufio.NewReader(stdout)
	line, _ := reader.ReadString('\n')
	cmd.Wait()
	//n, _ := strconv.Atoi(strings.Split(line, " ")[0])
	resultChannel <- line
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
			//result[fullPath] = countLines(fullPath)
		}
	}
	if len(files) == 0 {
		return result, nil
	}
	resultChannel := make(chan string, len(files))
	for _, file := range files {
		go countLines(file, resultChannel)
	}
	received := 0
	for {
		select {
		case line := <-resultChannel:
			splitedLine := strings.Split(line, " ")
			countStr, fileName := splitedLine[0], splitedLine[1]
			fileName = strings.TrimRight(fileName, "\n")
			n, _ := strconv.Atoi(countStr)
			result[fileName] = int64(n)
			received += 1
			if received == len(files) {
				return result, nil
			}
		}
	}
	return result, nil
}
