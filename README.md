# Experience-Gap-Indentifier
This tool takes a resume file and a list of job posting files or URLs as input, and identifies any experience gaps in the resume based on the requirements of the job postings. Jobalytics is an example of a tool that takes a single resume and a single job posting, then points out buzz words that you can use to spruce up your resume. **This tool**, however, is designed for people who are looking for new projects, new ideas, new things to learn, but just don't know what to prioritize. Gather a few of you dream roles, test them against your list of capabilities, and see what you could stand to get some more experience with! Ideal for people looking to pivot careers or choose a new skill to learn.

### Usage
To run this program, save your resume as 'resume.txt' or change the name of the stored constant in the program. Then, supply one or multiple job listings, either as files or as a url.
```
go run main.go [-v] [job posting file/URL]...
```
-v prints out in verbose. Not very suitable for urls, but may add some info for files.

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
CERTIFICATIONS:
cissp (1 mentions)

CYBER BUZZ WORDS:
strategy (4 mentions)
model (2 mentions)
compliance (2 mentions)
risk (3 mentions)
verification (4 mentions)
governance (1 mentions)
framework (1 mentions)
privacy (2 mentions)
awareness (3 mentions)
identity (3 mentions)
access (1 mentions)
regulation (2 mentions)
intrusion (1 mentions)
policy (6 mentions)
detection (2 mentions)
design (6 mentions)

PROGRAMMING LANGUAGES:
ruby (2 mentions)

ACRONYMS:
api (6 mentions)
it (7 mentions)
```
### Note

This program reads the job postings from the provided files and URLs. If a file is provided, it is read as a string. If a URL is provided, the program performs a GET request to the URL and reads the response body as a string.
