package main

import (
	"bufio"
	"fmt"
	gitignore "github.com/sabhiram/go-gitignore"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/tabwriter"
)

var ignoreMatcher *gitignore.GitIgnore

func init() {
	gitignoreContents, err := os.ReadFile(".gitignore")
	if err != nil {
		log.Fatal(err)
	}

	ignoreMatcher = gitignore.CompileIgnoreLines(strings.Split(string(gitignoreContents), "\n")...)
}

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
		content, err := os.ReadFile(f)
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
	err = os.WriteFile(filename, []byte(sb.String()), 0644)
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

// isBinaryFile checks if a file is a binary file
func isBinaryFile(filename string) (bool, error) {
	f, err := os.Open(filename)
	if err != nil {
		// If the file cannot be opened, treat it as a non-binary file
		return false, nil
	}
	defer f.Close()

	buf := make([]byte, 512)
	_, err = f.Read(buf)
	if err != nil {
		// If the file cannot be read, treat it as a non-binary file
		return false, nil
	}

	contentType := http.DetectContentType(buf)
	return !strings.HasPrefix(contentType, "text/"), nil
}

func shouldIgnore(path string) bool {
	// Ignore specific files and directories
	base := filepath.Base(path)
	if base == ".git" || base == "favicon.ico" {
		return true
	}

	// Ignore files based on their extension
	ext := filepath.Ext(path)
	ignoredExtensions := []string{".gitignore", ".idx", ".woff", ".png", ".jpg", ".jpeg", ".gif", ".ico", ".svg", ".txt", ".html", ".properties", ".scss", ".sh", ".sample", ".prettierignore", ".prettierrc"}
	for _, ignoredExt := range ignoredExtensions {
		if ext == ignoredExt {
			return true
		}
	}

	// Check if the path is in the .gitignore rules
	relPath, err := filepath.Rel(".", path)
	if err != nil {
		log.Fatal(err)
	}
	if ignoreMatcher.MatchesPath(relPath) {
		return true
	}

	// Check if the file is a binary file
	isBinary, err := isBinaryFile(path)
	if err != nil {
		log.Fatal(err)
	}
	if isBinary {
		return true
	}

	return false
}

func getFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if shouldIgnore(path) {
				return filepath.SkipDir
			}
			return nil
		}

		if !shouldIgnore(path) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
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
