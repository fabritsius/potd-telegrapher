package wikipedia

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type POTD struct {
	Title   string `json:"title"`
	Img     string `json:"img"`
	Content string `json:"content"`
	Credits string `json:"credits"`
}

func ParsePOTD(date string) (*POTD, error) {
	res, err := http.Get(fmt.Sprintf("https://en.wikipedia.org/wiki/Template:POTD/%s", date))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
		return nil, err
	}

	return fillInPOTD(res.Body)
}

func fillInPOTD(body io.ReadCloser) (*POTD, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	potd := new(POTD)

	contentBox := doc.Find("#mw-content-text")
	img := contentBox.Find("a.image img")
	content := contentBox.Find(".mw-parser-output div:nth-child(2) div:nth-child(3)")

	potd.Credits = content.Find("small a").Text()

	content.Find("small").Remove()
	content.Find(".noprint").Remove()

	potd.Content = strings.Trim(strings.ReplaceAll(content.Text(), "\n", ""), " ")

	title, exists := img.Attr("alt")
	if !exists {
		return potd, errors.New("image alt not found")
	}
	potd.Title = title

	imgSrc, exists := img.Attr("src")
	if !exists {
		return potd, errors.New("image src not found")
	}
	imgUrl, err := url.Parse(imgSrc)
	if err != nil {
		return potd, err
	}
	imgUrl.Scheme = "https"
	potd.Img = imgUrl.String()

	return potd, nil
}
