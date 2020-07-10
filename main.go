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

type Item struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	Desc    string `xml:"description"`
	Guid    string `xml:"guid"`
	PubDate string `xml:"pubDate"`
}

type Channel struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Desc  string `xml:"description"`
	Items []Item `xml:"item"`
}

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
	rssItems := readFeed("https://www.prskavec.net/post/index.xml")
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
