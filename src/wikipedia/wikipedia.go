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

	creditNames := []string{}
	content.Find("small a").Each(func(idx int, credit *goquery.Selection) {
		creditNames = append(creditNames, credit.Text())
	})
	potd.Credits = strings.Join(creditNames, ", ")

	content.Find("small").Remove()
	content.Find(".noprint").Remove()

	potd.Content = strings.Trim(strings.ReplaceAll(content.Text(), "\n", ""), " ")

	title, exists := img.Attr("alt")
	if !exists {
		return potd, errors.New("image alt not found")
	}
	potd.Title = title

	imgSrc, err := parseBestQualityImage(img)
	if err != nil {
		fmt.Println(err)
		imgSrc, exists = img.Attr("src")
		if !exists {
			return potd, errors.New("image source not found")
		}
	}

	imgUrl, err := url.Parse(imgSrc)
	if err != nil {
		return potd, err
	}
	imgUrl.Scheme = "https"
	potd.Img = imgUrl.String()

	return potd, nil
}

func parseBestQualityImage(img *goquery.Selection) (url string, err error) {
	srcSet, exists := img.Attr("srcset")
	if !exists {
		return "", errors.New("srcset attribure wasn't found")
	}

	bestSize := ""
	sources := strings.Split(srcSet, ",")
	for _, src := range sources {
		srcUrl, srcSize, err := parseSource(src)
		if err != nil {
			return "", err
		}

		if srcSize > bestSize {
			url = srcUrl
			bestSize = srcSize
		}
	}

	return url, nil
}

func parseSource(src string) (url, size string, err error) {
	split := strings.Split(strings.Trim(src, " "), " ")
	if len(split) != 2 {
		err = fmt.Errorf("failed to parse %s into 2 parts", src)
		return
	}

	return split[0], split[1], nil
}
