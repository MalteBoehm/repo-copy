package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"text/tabwriter"
)

func main() {
	if len(os.Args) != 1 {
		fmt.Println("Usage: copy-repo")
		return
	}

	files, err := getFiles(".")
	if err != nil {
		fmt.Println(err)
		return
	}

	var sb strings.Builder
	var totalFiles, totalLines, totalWords, totalChars int
	for _, f := range files {
		if !isCodeFile(f) {
			continue
		}

		content, err := ioutil.ReadFile(f)
		if err != nil {
			fmt.Println(err)
			continue
		}

		totalFiles++

		sb.WriteString("// " + f + "\n")
		sb.WriteString("// " + filepath.Base(f) + "\n\n")
		sb.WriteString(string(content) + "\n\n")

		scanner := bufio.NewScanner(strings.NewReader(string(content)))
		for scanner.Scan() {
			totalLines++
			totalWords += len(strings.Fields(scanner.Text()))
			totalChars += len(scanner.Text())
		}
	}

	filename := "codebase.txt"
	err = ioutil.WriteFile(filename, []byte(sb.String()), 0644)
	if err != nil {
		log.Fatal(err)
	}

	cmd := "pbcopy" // for macOS
	if runtime.GOOS == "windows" {
		cmd = "clip"
	}
	c := exec.Command(cmd)
	stdin, err := c.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stdin.Close()

	err = c.Start()
	if err != nil {
		log.Fatal(err)
	}

	_, err = stdin.Write([]byte(sb.String()))
	if err != nil {
		log.Fatal(err)
	}

	err = stdin.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Codebase written to", filename, "and copied to clipboard")

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(tw, "Statistics:")
	fmt.Fprintln(tw, "-----------")
	fmt.Fprintf(tw, "Total Files:\t%d\n", totalFiles)
	fmt.Fprintf(tw, "Total Lines:\t%d\n", totalLines)
	fmt.Fprintf(tw, "Total Words:\t%d\n", totalWords)
	fmt.Fprintf(tw, "Total Chars:\t%d\n", totalChars)
	fmt.Fprint(tw, "\nLanguages:\n-----------\n")

	languages := make(map[string]int)
	for _, f := range files {
		if !isCodeFile(f) {
			continue
		}
		ext := filepath.Ext(f)
		languages[ext]++
	}
	for ext, count := range languages {
		fmt.Fprintf(tw, "%s:\t%d\n", ext, count)
	}

	fmt.Fprintln(tw, "\nToken Count:")
	tokenCount := countTokens(sb.String())
	fmt.Fprintf(tw, "Total Tokens:\t%d\n", tokenCount)

	fmt.Fprintln(tw, "\nEstimated LLaMA 3 Requirements:")
	fmt.Fprintf(tw, "Tokens per Request:\t4,096\n")
	fmt.Fprintf(tw, "Estimated Requests:\t%.2f\n", float64(tokenCount)/4096)

	tw.Flush()
}

func getFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if shouldIgnoreDir(path) {
				return filepath.SkipDir
			}
			return nil
		}
		if !shouldIgnoreFile(path) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func shouldIgnoreDir(path string) bool {
	if filepath.Base(path) == "vendor" {
		return true
	}

	gitignore, err := ioutil.ReadFile(".gitignore")
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(gitignore), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && strings.HasSuffix(line, "/") && strings.HasPrefix(path, line) {
			return true
		}
	}
	return false
}

func shouldIgnoreFile(path string) bool {
	gitignore, err := ioutil.ReadFile(".gitignore")
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(gitignore), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasSuffix(line, "/") {
			match, err := filepath.Match(line, path)
			if err != nil {
				continue
			}
			if match {
				return true
			}
		}
	}
	return false
}

var codeFileExtensions = regexp.MustCompile(`\.(go|java|py|c|cpp|h|hpp|js|ts|php|rb|swift|kt|scala|rs|cs)$`)

func isCodeFile(path string) bool {
	return codeFileExtensions.MatchString(path)
}

func countTokens(text string) int {
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanWords)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count
}
