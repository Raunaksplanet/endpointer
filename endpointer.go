package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

func main() {
	// Check if the user provided a file with URLs
	if len(os.Args) < 2 {
		fmt.Println("Usage: endpointer <url_file.txt>")
		return
	}

	// Read the file containing the URLs
	urlFile := os.Args[1]
	urls, err := readLines(urlFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Regular expression to match JavaScript variable names
	varRegex := regexp.MustCompile(`\b(var|let|const)\s+([a-zA-Z_$][0-9a-zA-Z_$]*)`)

	for _, url := range urls {
		fmt.Printf("Processing URL: %s\n", url)

		// Fetch the JavaScript content
		jsContent, err := fetchJSContent(url)
		if err != nil {
			fmt.Printf("Error fetching JS content: %v\n", err)
			continue
		}

		// Extract variable names
		variables := extractVariableNames(jsContent, varRegex)

		// Generate potential test URLs (ensuring unique variable names)
		testURLs := generateTestURLs(url, variables)
		for _, testURL := range testURLs {
			fmt.Println(testURL)
		}
	}
}

// readLines reads a file line by line and returns a slice of strings
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// fetchJSContent fetches the JavaScript content from a URL
func fetchJSContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// extractVariableNames extracts variable names from JavaScript content
func extractVariableNames(jsContent string, varRegex *regexp.Regexp) []string {
	matches := varRegex.FindAllStringSubmatch(jsContent, -1)
	variables := make([]string, 0, len(matches))

	// Use a map to track unique variable names
	uniqueVars := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 2 {
			variable := match[2]
			if !uniqueVars[variable] {
				uniqueVars[variable] = true
				variables = append(variables, variable)
			}
		}
	}

	return variables
}

// generateTestURLs generates potential test URLs with unique variable names as GET parameters
func generateTestURLs(baseURL string, variables []string) []string {
	var testURLs []string

	for _, variable := range variables {
		testURL := fmt.Sprintf("%s?%s=", baseURL, variable)
		testURLs = append(testURLs, testURL)
	}

	return testURLs
}
