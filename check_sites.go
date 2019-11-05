package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

func main() {
	var waitGroup sync.WaitGroup
	var m sync.Mutex

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
				m.Lock()
				defer m.Unlock()
				fmt.Println(checkSite(line))
			}()
		} else {
			fmt.Fprintln(os.Stderr, "invalid URL: ", line)
		}
	}
	// display any errors from reading stdin
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input: ", err)
	}

	// wait for all goroutines to complete
	waitGroup.Wait()
}

// check if urlStr is a valid URL
func isValidUrl(urlStr string) bool {
	_, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return false
	}
	return true
}

// check if site is up
func checkSite(site string) string {
	var status string

	resp, err := http.Get(site)
	if err != nil {
		status = "DOWN " + err.Error()
	} else {
		defer resp.Body.Close()

		status = "UP " + resp.Status

		io.Copy(ioutil.Discard, resp.Body)
	}

	return fmt.Sprintf("%s %s", site, status)
}
