# Experience-Gap-Indentifier
This tool takes a resume file and a list of job posting files or URLs as input, and identifies any experience gaps in the resume based on the requirements of the job postings. Ideal for people looking to pivot careers or choose a new skill to learn.

### Usage
To run this program, save your resume as 'resume.txt' or change the name of the stored constant in the program. Then, supply one or multiple job listings, either as files or as a url.
```
go run main.go [job posting file/URL]...
```

### Output

The program writes the identified experience gaps to an output file called "output.txt". Each line of the output file lists an experience gap and the number of job postings that required that experience.
Example

Given the following resume file:
```
Name: John Smith
Job Title: Software Developer
Skills: Go, Python, JavaScript
Experience:
- Job Title: Software Developer
  Company: ABC Company
  Description: Developed software applications using Go and Python
- Job Title: Web Developer
  Company: XYZ Company
  Description: Developed web applications using JavaScript
```
And the following job posting file:
```
Software Engineer | Google | Go, Python, Java
Full Stack Developer | Facebook | JavaScript, Ruby, Python
```
Running the program as follows:
```
go run main.go job_posting.txt
```
Will produce the following output file:
```
Java (1 job postings)
Ruby (1 job postings)
```
### Note

This program reads the job postings from the provided files and URLs. If a file is provided, it is read as a string. If a URL is provided, the program performs a GET request to the URL and reads the response body as a string.
