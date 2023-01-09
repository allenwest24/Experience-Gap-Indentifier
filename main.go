package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"log"
	"flag"
)

func main() {
	// Define the flags
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "Verbose mode")

	// Parse the flags
	flag.Parse()

	// Read the resume from a file
	resume, err := readResume("resume.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	resume = sanitizeString(resume)
	fmt.Println(resume)

	// Read the job postings from the provided files and URLs
	//jobPostings, err := readJobPostings(os.Args[1:])
	jobPostings, err := readJobPostings(flag.Args())
	if err != nil {
		fmt.Println(err)
		return
	}
	for i, s := range jobPostings {
		jobPostings[i] = sanitizeString(s)
	}
	fmt.Println(jobPostings)

	// Compare the set of words in the resume to the set of words in the job postings
	// to identify words not found in the resume that are in a job posting
	missingWords := identifyMissingWords(resume, jobPostings)

	// Write the missing words to an output file
	if err := writeMissingWords("output.txt", missingWords, verbose); err != nil {
		fmt.Println(err)
		return
	}
}

func sanitizeString(s string) string {
	// Make the string lowercase
	s = strings.ToLower(s)

	// Remove all extra characters
	s = strings.ReplaceAll(s, ")", "")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "`", "")

	// Remove all newlines
	s = strings.ReplaceAll(s, "\n", " ")

	// Separate based combinators
	s = strings.ReplaceAll(s, "/", " ")

	return s
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

func readFile(filename string) []string {
	// Open the file using the os.Open function
	file, err := os.Open(filename)
	if err != nil {
	    log.Fatal(err)
	}
	defer file.Close()

	// Create a new Scanner to read the file
	scanner := bufio.NewScanner(file)

	// Use the Scan function to read the file line by line
	lines := []string{}
	for scanner.Scan() {
	    lines = append(lines, sanitizeString(scanner.Text()))
	}
	
	if err := scanner.Err(); err != nil {
	    log.Fatal(err)
	}
	return lines
}

func writeMissingWords(filename string, missingWords map[string]int, verbose bool) error {
	knownCerts := readFile(strings.Join([]string{"Skill Categories", "certifications.txt"}, "/"))
	knownCBW := readFile(strings.Join([]string{"Skill Categories", "cyber_buzz_words.txt"}, "/"))
	knownPL := readFile(strings.Join([]string{"Skill Categories", "programming_languages.txt"}, "/"))
	knownAcronyms := readFile(strings.Join([]string{"Skill Categories", "tech_acronyms.txt"}, "/"))

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var certList []string
	var CBWList []string
	var PLList []string
	var acronymList []string
	var other []string

	for word := range missingWords {
		if contains(knownCerts, word) {
			certList = append(certList, word)
		} else if contains(knownCBW, word) {
			CBWList = append(CBWList, word)
		} else if contains(knownPL, word) {
			PLList = append(PLList, word)
		} else if contains(knownAcronyms, word) {
			acronymList = append(acronymList, word)
		} else {
			other = append(other, word)
		}
	}

	// Write all certifications to output.txt.
	if _, err := file.WriteString("\nCERTIFICATIONS:\n"); err != nil {
		return err
	}
	for _, str := range certList {
		if _, err := file.WriteString(fmt.Sprintf("%s (%d mentions)\n", str, missingWords[str])); err != nil {
			return err
		}
	}

	// Write all cyber buzz words to output.txt.
        if _, err := file.WriteString("\nCYBER BUZZ WORDS:\n"); err != nil {
                return err
        }
        for _, str := range CBWList {
                if _, err := file.WriteString(fmt.Sprintf("%s (%d mentions)\n", str, missingWords[str])); err != nil {
                        return err
                }
        }

	// Write all programming languages to output.txt.
        if _, err := file.WriteString("\nPROGRAMMING LANGUAGES:\n"); err != nil {
                return err
        }
        for _, str := range PLList {
                if _, err := file.WriteString(fmt.Sprintf("%s (%d mentions)\n", str, missingWords[str])); err != nil {
                        return err
                }
        }

	// Write all acronyms to output.txt.
        if _, err := file.WriteString("\nACRONYMS:\n"); err != nil {
                return err
        }
        for _, str := range acronymList {
                if _, err := file.WriteString(fmt.Sprintf("%s (%d mentions)\n", str, missingWords[str])); err != nil {
                        return err
                }
        }

	// Write all else  to output.txt if verbose is enabled
	if verbose {
        	if _, err := file.WriteString("\nEVERYTHING ELSE:\n"); err != nil {
                	return err
        	}
        	for _, str := range other {
                	if _, err := file.WriteString(fmt.Sprintf("%s (%d mentions)\n", str, missingWords[str])); err != nil {
                        	return err
                	}
        	}
	}

	return nil
}
