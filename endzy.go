package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sync"
	"time"
)

func main() {

	var concurrency int
	var domains string
	flag.IntVar(&concurrency, "concurrency", 30, "The concurrency or speed")
	flag.StringVar(&domains, "domains", "", "The domain in the format example.com with out (http or https)")
	flag.Parse()

	// Check to see if the parameters argument is not empty
	if domains != "" {
		// Create our goroutines
		var wg sync.WaitGroup
		for i := 0; i <= concurrency; i++ {
			wg.Add(1)
			go func() {
				// Call our function to do stuff
				checkdb(domains)
				wg.Done()
			}()
			wg.Wait()
		}
	} else {
		flag.PrintDefaults()
	}
}

// It will compare the endpoints discovered to the ones in a database for new matches
func checkdb(newd string) {

	// Sleep in order for goroutines to work
	time.Sleep(time.Millisecond * 10)

	endpoints := getEndpoints(newd)

	// Read all the files in the directory "db"
	files, err := ioutil.ReadDir("db")
	if err != nil {
		log.Fatal(err)
	}
	// If there are no files save them first
	if len(files) == 0 {
		for _, e := range endpoints {
			fmt.Println(e)
			writeEndpoint(e)
			return
		}
	}

	// Open all the files in the "db" dir
	for _, file := range files {
		f, err := os.Open("db/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		// Reads the file in the db and compares it with the endpoint discovered
		fScanner := bufio.NewScanner(f)
		for fScanner.Scan() {
			for _, e := range endpoints {
				if e != fScanner.Text() {
					fmt.Println(e)
					//writeEndpoint(e)
				}
			}
		}
	}
}

// Write the endpoint to a file
func writeEndpoint(e string) {

	if fileExists("db/data") == false {
		file, err := os.Create("db/data")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		data := []byte(e)
		_, err2 := file.Write(data)
		if err2 != nil {
			log.Fatal(err)
		}

	} else {
		file, err := os.Open("db/data")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		data := []byte(e)
		_, err2 := file.Write(data)
		if err2 != nil {
			log.Fatal(err)
		}
	}

}

// Check if a file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Generate a wordlist
func getEndpoints(d string) []string {

	// These are our slices / array's
	var endpoints = make([]string, 0)
	var links = make([]string, 0)
	var javascript = make([]string, 0)

	// Use our scanners for Stdin
	jsScanner := bufio.NewScanner(os.Stdin)
	for jsScanner.Scan() {
		javascript = append(javascript, jsScanner.Text())
	}

	// Iterate over the slices
	for _, js := range javascript {

		// The javascript link
		link := js

		// Make a new Request
		req, err := http.NewRequest("GET", link, nil)
		if err != nil {
			return nil
		}

		// Send the request to the server fetch the response
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil
		}

		// Read the response & convert it to a string
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil
		}
		bodyString := string(bodyBytes)

		// Match all the links in the js
		re := regexp.MustCompile(`(https?|ftp|file)://[-A-Za-z0-9\+&@#/%?=~_|!:,.;]*[-A-Za-z0-9\+&@#/%=~_|]`)
		matches := re.FindStringSubmatch(bodyString)
		if matches != nil {
			match := matches[0]
			links = append(links, match)
		}
	}

	// Iterate over the links & parameters
	for _, _link := range links {
		// The link is going to be in the form http://domain.com/endpoint/blahblah
		// We need the path which /endpoint/blahblah
		u, err := url.Parse(_link)
		if err != nil {
			return nil
		}

		// This is the final result which should look like
		// /endpoint/blahblah/?param=FUZZ
		endpoint := d + u.Path
		endpoints = append(endpoints, endpoint)

		return endpoints
	}
	return endpoints
}
