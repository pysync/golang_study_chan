package main

import (
	"fmt"
	"time"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	// define max worker count and chan buffer length
	const workers, length = 2, 100

	// make slice of strings store request urls
	jobs := []string{url}

	// make map of string store working urls
	working := make(map[string]int)

	// make map of string store done urls
	done := make(map[string]int)

	// create chan using for notify terminate signal
	exit := make(chan bool)

	// create buffer chan for incoming requests
	req := make(chan string, length)

	// create handle function to handle each request
	handle := func(queue chan string, id int) {
		for r := range queue {
			if _, d := done[r]; d {
				continue
			}

			if _, w := working[r]; w {
				continue
			}

			working[r] = 1
			body, urls, err := fetcher.Fetch(r)
			if err != nil {
				fmt.Printf("worker[%d]not found: %s\n", id, url)
				return
			}
			fmt.Printf("worker[%d]found: %s %q\n", id, url, body)
			for _, u := range urls {
				jobs = append(jobs, u)
			}
			delete(working, r)
			done[r] = 1
		}
	}
	// create a serve func
	serve := func(req chan string, exit chan bool) {
		for i := 0; i < workers; i++ {
			go handle(req, i)
		}
		<-exit
	}

	// start request
	go serve(req, exit)

	// create a forever loop listen incoming request
	for {

		// while count of jobs > 0
		// => populate to worker
		for i := 0; i < workers && len(jobs) > 0; i++ {
			job := jobs[0]
			jobs = jobs[1:]
			req <- job
		}
		time.Sleep(time.Microsecond)
	}
}

func main() {
	Crawl("http://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
