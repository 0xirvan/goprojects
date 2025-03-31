package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sync"

	"golang.org/x/net/html"
)

const baseUrl string = "https://scrape-me.dreamsofcode.io"

func main() {
	linksCh := make(chan string)
	htmlChan := make(chan *html.Node)

	client := &http.Client{}

	var wg sync.WaitGroup
	visited := make(map[string]bool)
	var mu sync.Mutex
	var fileMutex sync.Mutex

	wg.Add(1)
	go func() {
		defer wg.Done()
		fetchPage(baseUrl, client, htmlChan)
	}()

	go func() {
		var extractWg sync.WaitGroup
		for doc := range htmlChan {
			extractWg.Add(1)
			go func(doc *html.Node) {
				defer extractWg.Done()
				extractLinks(doc, linksCh, &fileMutex)
			}(doc)
		}
		extractWg.Wait()
		close(linksCh)
	}()

	go func() {
		for link := range linksCh {
			mu.Lock()
			if visited[link] {
				mu.Unlock()
				continue
			}
			visited[link] = true
			mu.Unlock()

			wg.Add(1)
			go func(link string) {
				defer wg.Done()
				fetchPage(link, client, htmlChan)
			}(link)
		}
	}()

	wg.Wait()
	close(htmlChan)
}

func extractLinks(n *html.Node, linksCh chan<- string, fileMutex *sync.Mutex) {

	regex, err := regexp.Compile(`^(#|/)`)
	if err != nil {
		log.Printf("Error compiling regex: %v", err)
		return
	}

	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					var link string
					if regex.MatchString(attr.Val) {
						link = fmt.Sprintf("%s%s", baseUrl, attr.Val)

						parsedUrl, err := url.Parse(link)
						if err != nil {
							log.Printf("Error parsing URL %s: %v", link, err)
							continue
						}
						parsedUrl.RawQuery = ""
						normalizedLink := parsedUrl.String()

						linksCh <- normalizedLink
					} else {
						link = attr.Val
					}

					fileMutex.Lock()
					file, err := os.OpenFile("links.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						log.Printf("Error opening/creating file: %v", err)
						fileMutex.Unlock()
						continue
					}

					if _, err := file.WriteString(link + "\n"); err != nil {
						log.Printf("Error writing to file: %v", err)
					}
					file.Close()
					fileMutex.Unlock()

				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
}

func fetchPage(baseUrl string, client *http.Client, htmlChan chan<- *html.Node) {

	parsedUrl, err := url.Parse(baseUrl)
	if err != nil {
		log.Printf("Error parsing URL %s: %v", baseUrl, err)
		return
	}
	parsedUrl.RawQuery = ""
	normalizedURL := parsedUrl.String()

	res, err := client.Get(normalizedURL)
	if err != nil {
		log.Println("Error fetching URL:", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Println("Error: ", res.StatusCode)
		return
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		log.Println("Error parsing HTML:", err)
		return
	}

	htmlChan <- doc
	log.Println("Fetched: ", normalizedURL)
}
