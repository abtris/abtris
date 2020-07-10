package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

// Item type for RSS
type Item struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	Desc    string `xml:"description"`
	GUID    string `xml:"guid"`
	PubDate string `xml:"pubDate"`
}

// Channel type for RSS
type Channel struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Desc  string `xml:"description"`
	Items []Item `xml:"item"`
}

// Rss type for RSS as root
type Rss struct {
	Channel Channel `xml:"channel"`
}

func readFeed(url string) []string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	rss := Rss{}

	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&rss)
	if err != nil {
		log.Fatal(err)
	}

	output := []string{}
	for _, item := range rss.Channel.Items {
		output = append(output, fmt.Sprintf("[%s](%s)\n", item.Title, item.Link))
	}
	return output
}

func main() {
	content, err := ioutil.ReadFile("README.md.template")
	if err != nil {
		log.Fatal(err)
	}

	text := string(content)
	blogURL := os.Getenv("BLOG_URL")
	rssItems := readFeed(blogURL)
	data := struct {
		Title string
		Items []string
	}{
		Title: "Last blog posts",
		Items: rssItems,
	}

	t, err := template.New("readme").Parse(text)
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(os.Stdout, data)
	if err != nil {
		log.Fatal(err)
	}
}
