package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

// TestRenderHeadings runs the real README.md.template through render and
// verifies each feed gets the right section heading: writing and talks get
// their custom labels (matched on the feed URL), everything else falls back
// to "<Title> in <Language>".
func TestRenderHeadings(t *testing.T) {
	tmpl, err := os.ReadFile("README.md.template")
	if err != nil {
		t.Fatalf("reading template: %v", err)
	}

	data := []Blog{
		{Title: "Prskavčí blog", Url: "https://blog.prskavec.net/index.xml", Language: "cs-CZ", Items: []string{"[cs post](https://blog.prskavec.net/a)\n"}},
		{Title: "Posts on Prskavec.Net", Url: "https://www.prskavec.net/post/index.xml", Language: "en", Items: []string{"[a post](https://www.prskavec.net/post/a)\n"}},
		{Title: "Talks on Prskavec.Net", Url: "https://www.prskavec.net/talk/index.xml", Language: "en", Items: []string{"[a talk](https://www.prskavec.net/talk/a)\n"}},
	}

	var buf strings.Builder
	if err := render(&buf, string(tmpl), data); err != nil {
		t.Fatalf("render: %v", err)
	}
	out := buf.String()

	want := []string{
		"### ✍️ Writing",
		"### 🎤 Talks",
		"### Prskavčí blog in cs-CZ",
		"- [a post](https://www.prskavec.net/post/a)",
		"- [a talk](https://www.prskavec.net/talk/a)",
	}
	for _, w := range want {
		if !strings.Contains(out, w) {
			t.Errorf("rendered output missing %q\n--- output ---\n%s", w, out)
		}
	}

	// The generic label must not leak onto the writing/talks feeds.
	for _, bad := range []string{"### Posts on Prskavec.Net in en", "### Talks on Prskavec.Net in en"} {
		if strings.Contains(out, bad) {
			t.Errorf("rendered output should not contain generic heading %q", bad)
		}
	}
}

func TestParseLimits(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		n    int
		want []int
	}{
		{"empty defaults", "", 3, []int{defaultLimit, defaultLimit, defaultLimit}},
		{"all set", "5,2,7", 3, []int{5, 2, 7}},
		{"partial fills rest with default", "5", 3, []int{5, defaultLimit, defaultLimit}},
		{"whitespace trimmed", " 5 , 2 ", 2, []int{5, 2}},
		{"invalid keeps default", "5,x,3", 3, []int{5, defaultLimit, 3}},
		{"non-positive keeps default", "5,0,-1", 3, []int{5, defaultLimit, defaultLimit}},
		{"extra parts ignored", "5,2,7,9", 2, []int{5, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseLimits(tt.raw, tt.n)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseLimits(%q, %d) = %v, want %v", tt.raw, tt.n, got, tt.want)
			}
		})
	}
}

func TestGetFirst(t *testing.T) {
	in := []string{"a", "b", "c"}
	if got := getFirst(in, 2); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Errorf("getFirst limit<len = %v, want [a b]", got)
	}
	if got := getFirst(in, 5); !reflect.DeepEqual(got, in) {
		t.Errorf("getFirst limit>len = %v, want %v", got, in)
	}
}
