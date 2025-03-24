package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
	//"database/sql"
	//"github.com/joho/godotenv"
)

func workerId(id int, urlChannel <-chan string, resultChannel chan<- string) {
	for url := range urlChannel { //taking a url from the urlchannel
		fmt.Printf("Worker %d Scraping %s\n", id, url)

		//Simulating scraping by using a GET method
		//resp -> pointer for HTTP response
		//err -> error value that will be returned if there is an error
		resp, err := http.Get(url)
		if err != nil {
			resultChannel <- fmt.Sprintf("Worker %d: failed to fetch %s, Error: %s", id, url, err)
			continue
		}
		defer resp.Body.Close()

		//sending the result to the resultChannel
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			resultChannel <- fmt.Sprintf("Worker %d: failed to read response body from %s, Error: %s", id, url, err)
			continue
		}
		//Here I am trying to extract the title of the page using regex
		//regex is not the best way to parse HTML, but for this example, it will work (maybe)
		// a better option would be to use a library like goquery but I'm not sure how it will work with the current setup
		titleRegex := regexp.MustCompile("<title>(.*?)</title>")
		matches := titleRegex.FindStringSubmatch(string(body))

		//sending the result to the resultChannel
		if len(matches) > 1 {
			resultChannel <- fmt.Sprintf("Worker %d title of %s : %s", id, url, matches[1])
		} else {
			resultChannel <- fmt.Sprintf("Worker %d: failed to find the title in %s", id, url)
		}
		time.Sleep(1 * time.Second) //making it sleep for 1 second to simulate the scraping process because it's too fast
	}
}

func main() {
	//I will give some lists of urls to work with
	urls := []string{
		"https://golang.org",
		"https://www.google.com",
		"https://github.com",
		"https://github.com/VaibhavSingh79/Netflix",
		"https://www.stackoverflow.com",
		"https://news.ycombinator.com",
	}

	//creating channels
	urlChannel := make(chan string)    //channels for urls to scrape
	resultChannel := make(chan string) //channels for receiving results

	//number of worker
	numWorkers := 3 //let's say there are 3 workers

	for i := 1; i <= numWorkers; i++ {
		go workerId(i, urlChannel, resultChannel)
	}
	go func() {
		for _, url := range urls {
			urlChannel <- url
		}
		close(urlChannel)
	}()

	var results []string
	for i := 0; i < len(urls); i++ {
		result := <-resultChannel // Read from resultChannel
		fmt.Println(result)       // Print to console
		results = append(results, result)
	}

	//Here I will create a file and write the results to it
	file, err := os.Create("scraped_results.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Here i am writing the results to the file
	for _, result := range results {
		file.WriteString(result + "\n")
	}
	fmt.Printf("All scraping tasks completed!")
}
