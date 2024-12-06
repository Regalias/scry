package scrape

import "github.com/PuerkitoBio/goquery"

func FindAttr(s *goquery.Selection, attrName string) (string, bool) {
	for _, node := range s.Nodes {
		for _, attr := range node.Attr {
			if attr.Key == attrName {
				return attr.Val, true
			}
		}
	}
	return "", false
}

func FindChildAttr(s *goquery.Selection, selector string, attrName string) (string, bool) {
	return FindAttr(s.Find(selector), attrName)
}
