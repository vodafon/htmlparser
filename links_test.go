package htmlparser

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestLinksHTML(t *testing.T) {
	file, err := ioutil.ReadFile("./testdata/links1.html")
	if err != nil {
		t.Fatal(err)
	}
	links, err := Links(file, "https://google.com")
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 150 {
		t.Errorf("wrong links size %d\n", len(links))
	}
	for _, link := range links {
		if !strings.HasPrefix(link.Abs, "http") {
			t.Errorf("relative URL %s", link.Abs)
		}
	}
}

func TestLinksJSON(t *testing.T) {
	file, err := ioutil.ReadFile("./testdata/json1.json")
	if err != nil {
		t.Fatal(err)
	}
	links, err := Links(file, "https://google.com")
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 0 {
		t.Errorf("wrong links size %d\n", len(links))
	}
}

func TestExternalTLDLinksHTML(t *testing.T) {
	file, err := ioutil.ReadFile("./testdata/links1.html")
	if err != nil {
		t.Fatal(err)
	}
	links, err := ExternalTLDLinks(file, "https://google.com")
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 34 {
		t.Errorf("wrong links size %d\n", len(links))
	}
	for _, link := range links {
		if strings.Contains(link, "google.com") {
			t.Errorf("not external link %s", link)
		}
	}
}
