package decomposition_tools

import (
	"golang.org/x/net/html"
	"regexp"
)

func GetElHref(d []html.Attribute) string {
	cl := ""
	for _, v := range d {
		if v.Key == "href" {
			cl = v.Val
		}
	}

	return cl
}

func GetElAttribute(key string, d []html.Attribute) string {
	a := ""
	for _, v := range d {
		if v.Key == key {
			a = v.Val
		}
	}

	return a
}

func GetElClass(d []html.Attribute) string {
	cl := ""
	for _, v := range d {
		if v.Key == "class" {
			cl = v.Val
		}
	}

	return cl
}

func MatchElClassByRegExp(exp string, d []html.Attribute) bool {
	elementClass := GetElClass(d)
	if elementClass != "" {
		cnt, err := regexp.MatchString(exp, elementClass)
		if err != nil {
			return false
		}

		return cnt
	}
	return false
}
