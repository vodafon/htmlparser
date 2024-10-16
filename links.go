package htmlparser

import (
	"bytes"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

var (
	invalidSymbols = []string{"'", "&", "$", "^", "::", "\"", ";"}
)

type Link struct {
	Raw         string
	Abs         string
	TLD         string
	External    bool
	ExternalTLD bool
}

func ExternalTLDLinks(page []byte, domain string) ([]string, error) {
	links, err := Links(page, domain)
	if err != nil {
		return nil, err
	}

	res := []string{}
	for _, link := range links {
		if link.ExternalTLD {
			res = append(res, link.Abs)
		}
	}
	return res, nil
}

func TLD(host string) (string, error) {
	tld, icann := publicsuffix.PublicSuffix(host)
	if icann {
		return publicsuffix.EffectiveTLDPlusOne(host)
	}
	return tld, nil
}

func Links(page []byte, domain string) ([]Link, error) {
	if !strings.HasPrefix(domain, "http") {
		domain = "https://" + domain
	}
	uri, err := url.Parse(domain)
	if err != nil {
		return nil, err
	}
	tld, err := TLD(uri.Host)
	if err != nil {
		return nil, err
	}

	z := html.NewTokenizer(bytes.NewReader(page))
	links := []Link{}
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return links, nil
		case tt == html.StartTagToken:
			linkRaw := parseToken(z.Token())
			if linkRaw == "" {
				continue
			}

			linkAbs := urlABS(linkRaw, uri)
			linkURI, err := url.Parse(linkAbs)
			if err != nil {
				return nil, err
			}

			if IsContain(linkURI.Host, invalidSymbols) {
				continue
			}

			linkTLD, err := TLD(linkURI.Host)
			if err != nil {
				return nil, err
			}

			link := Link{
				Raw: linkRaw,
				Abs: linkAbs,
				TLD: linkTLD,
			}
			link.External = linkURI.Host != uri.Host
			link.ExternalTLD = link.TLD != tld
			links = append(links, link)
		}
	}
	return links, nil
}

func urlABS(link string, host *url.URL) string {
	if strings.HasPrefix(link, "http") {
		return link
	}
	if strings.HasPrefix(link, "//") {
		return host.Scheme + ":" + link
	}
	if strings.HasPrefix(link, "/") {
		return host.String() + link
	}
	return host.String() + "/" + link
}

func parseToken(t html.Token) string {
	switch t.Data {
	case "a", "link":
		return getAttr(t, "href")
	case "img", "script", "iframe":
		return getAttr(t, "src")
	case "form":
		// TODO: method
		return ""
	// case "b", "div", "input", "h2", "span", "meta", "title", "html", "head", "style", "body":
	// 	return ""
	default:
		// fmt.Printf("%s %+v\n", t.Data, t)
	}
	return ""
}

func getAttr(t html.Token, name string) string {
	for _, a := range t.Attr {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}

func IsContain(host string, sl []string) bool {
	for _, s := range sl {
		if strings.Contains(host, s) {
			return true
		}
	}
	return false
}
