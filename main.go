package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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

type Blog struct {
	Title string
	Url string
	Language string
	Limit int
	Items []string
}

// Channel type for RSS
type Channel struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Desc  string `xml:"description"`
	Language string `xml:"language"`
	Items []Item `xml:"item"`
}

// Rss type for RSS as root
type Rss struct {
	Channel Channel `xml:"channel"`
}

func readFeed(url string) ([]string, string, string) {
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
	return output, rss.Channel.Title, rss.Channel.Language
}

func getFirst(slices []string, limit int) []string {
	if len(slices) > limit {
		return slices[:limit]
	}
	return slices
}

func main() {
	content, err := ioutil.ReadFile("README.md.template")
	if err != nil {
		log.Fatal(err)
	}

	text := string(content)
	blogURLS := os.Getenv("BLOG_URLS")

	urls := strings.Split(blogURLS,",")

	data := []Blog{}

	for _, url := range urls {
		rssItems, title, lang := readFeed(url)
		blog := Blog{
			Title: title,
			Url: url,
			Language: lang,
			Items: getFirst(rssItems, 10),
		}
		data = append(data, blog)
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
