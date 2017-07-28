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
	countMap := countAll(os.Args[1])
	//fmt.Println("counts of files", len(countMap))
	for k, v := range countMap {
		fmt.Printf("%s : %v\n", k, v)
	}
	end := time.Now().UnixNano()
	fmt.Println("time cost", end-start)
}

func countAll(dir string) map[string]int64 {
	//children, err := ioutil.ReadDir(dir)
	//if err != nil {
	//	return result
	//}
	resultChannel := make(chan map[string]int64)
	return Traverse(dir, resultChannel, true)
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
	resultChannel <- line
}

func Traverse(dirPth string, resultChannel chan map[string]int64, shouldReturn bool) map[string]int64 {
	writeChannel := func(content map[string]int64) {
		if shouldReturn {
			return
		}
		resultChannel <- content
	}
	result := make(map[string]int64)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		fmt.Printf("read path %v failed %v\n", dirPth, err)
		writeChannel(result)
		return result
	}
	PthSep := string(os.PathSeparator)
	files := []string{}
	dirs := []string{}
	for _, fi := range dir {
		fullPath := dirPth + PthSep + fi.Name()
		if fi.IsDir() {
			//subResult, err := Traverse(fullPath)
			//if err == nil {
			//	merge(result, subResult)
			//}
			dirs = append(dirs, fullPath)
			continue
		}
		if isTargetFile(fullPath) {
			files = append(files, fullPath)
		}
	}
	if len(files) == 0 && len(dirs) == 0 {
		writeChannel(result)
		return result
	}
	fileChannel := make(chan string, len(files))
	defer close(fileChannel)
	if len(files) > 0 {
		for _, file := range files {
			go countLines(file, fileChannel)
		}
	}
	dirChannel := make(chan map[string]int64, len(dirs))
	defer close(dirChannel)
	if len(dirs) > 0 {
		for _, dir := range dirs {
			go Traverse(dir, dirChannel, false)
		}
	}
	fileReceived := 0
	dirReceived := 0
	for {
		select {
		case line := <-fileChannel:
			splitedLine := strings.Split(line, " ")
			countStr, fileName := splitedLine[0], splitedLine[1]
			fileName = strings.TrimRight(fileName, "\n")
			n, _ := strconv.Atoi(countStr)
			result[fileName] = int64(n)
			fileReceived += 1
			if fileReceived == len(files) && dirReceived == len(dirs) {
				writeChannel(result)
				return result
			}
		case countMap := <-dirChannel:
			merge(result, countMap)
			dirReceived += 1
			if fileReceived == len(files) && dirReceived == len(dirs) {
				writeChannel(result)
				return result
			}
		}
	}
	return result
}
