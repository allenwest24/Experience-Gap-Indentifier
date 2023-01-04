package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Read the resume from a file
	resume, err := readResume("resume.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Read the job postings from the provided files and URLs
	jobPostings, err := readJobPostings(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		return
	}

	// Compare the set of words in the resume to the set of words in the job postings
	// to identify words not found in the resume that are in a job posting
	missingWords := identifyMissingWords(resume, jobPostings)

	// Write the missing words to an output file
	if err := writeMissingWords("output.txt", missingWords); err != nil {
		fmt.Println(err)
		return
	}
}

func readResume(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var resume []string
	for scanner.Scan() {
		resume = append(resume, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return strings.Join(resume, " "), nil
}
func readJobPostings(filenames []string) ([]string, error) {
	var jobPostings []string
	for _, filename := range filenames {
		if strings.HasPrefix(filename, "http://") || strings.HasPrefix(filename, "https://") {
			// Read the job posting from a URL
			resp, err := http.Get(filename)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			bytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			jobPostings = append(jobPostings, string(bytes))
		} else {
			// Read the job posting from a file
			bytes, err := ioutil.ReadFile(filename)
			if err != nil {
				return nil, err
			}
			jobPostings = append(jobPostings, string(bytes))
		}
	}
	return jobPostings, nil
}

func identifyMissingWords(resume string, jobPostings []string) map[string]int {
	// Split the resume and job postings into a slice of words
	resumeWords := strings.Fields(resume)
	var jobPostingWords [][]string
	for _, posting := range jobPostings {
		jobPostingWords = append(jobPostingWords, strings.Fields(posting))
	}

	// Create a map of the missing words and their frequency in the job postings
	missingWords := make(map[string]int)
	for _, postingWords := range jobPostingWords {
		for _, word := range postingWords {
			if !contains(resumeWords, word) {
				missingWords[word]++
			}
		}
	}
	return missingWords
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func writeMissingWords(filename string, missingWords map[string]int) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for word, count := range missingWords {
		if _, err := file.WriteString(fmt.Sprintf("%s (%d job postings)\n", word, count)); err != nil {
			return err
		}
	}
	return nil
}
