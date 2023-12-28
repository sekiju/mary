package request

import "github.com/PuerkitoBio/goquery"

func GetDocument(url string, args *Config) (*goquery.Document, error) {
	if args == nil {
		args = &Config{}
	}

	req, err := newRequest("GET", url, nil, *args)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
