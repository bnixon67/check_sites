/*
Copyright 2021 Bill Nixon

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

// check_sites will attempt to connect to each URL provided on stdin
//
// each check will be done via a goroutine to allow multi-tasking since the connection may take some time to process or timeout
func main() {
	// waitGroup for the goroutines
	var waitGroup sync.WaitGroup

	// logger to use for output
	// a Logger can be used simultaneously from multiple goroutines
	// if a fmt.Print* was used, a mutex would be needed
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// loop thru stdin by line
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()

		// trim any leading/trailing spaces
		line = strings.TrimSpace(line)

		// if a valid url, use a goroutine to check if the site is up
		if isValidUrl(line) {
			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()
				logger.Println(checkSite(line))
			}()
		} else {
			logger.Println("invalid URL: ", line)
		}
	}

	// display any errors from reading stdin
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input: ", err)
	}

	// wait for all goroutines to complete
	waitGroup.Wait()
}

// isValidURL determines if urlStr is a valid URL
func isValidUrl(urlStr string) bool {
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}

// checkSite determines if a site is up by doing a simple Get,
// returning the site and status string.
func checkSite(site string) string {
	var status string

	resp, err := http.Head(site)
	if err != nil {
		status = "DOWN " + err.Error()
	} else {
		defer resp.Body.Close()

		status = "UP " + resp.Status

		// read and discard the response body
		if _, err := io.Copy(io.Discard, resp.Body); err != nil {
			log.Fatal(err)
		}
	}

	return fmt.Sprintf("%s %s", site, status)
}
