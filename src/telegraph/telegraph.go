package telegraph

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/fabritsius/potd-telegrapher/src/wikipedia"
)

func MakeArticle(date string) (string, error) {
	telegraph, err := NewTelegraphClient()
	if err != nil {
		return "", err
	}

	wikiPage, err := wikipedia.ParsePOTD(date)
	if err != nil {
		return "", err
	}

	page, err := fillPage(wikiPage)
	if err != nil {
		return "", err
	}

	resp, err := telegraph.createPage(page)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func fillPage(wikiPage *wikipedia.POTD) (TelegraphPage, error) {
	page := TelegraphPage{
		Title: wikiPage.Title,
	}

	page.AddImg(wikiPage.Img)
	page.AddText(wikiPage.Content)
	page.AddText(fmt.Sprintf("Credits: %s", wikiPage.Credits))

	return page, nil
}

type TelegraphPage struct {
	Title      string `json:"title"`
	AuthorName string `json:"authorName"`
	AuthorURL  string `json:"authorURL"`
	Content    []Node `json:"content"`
}

func (page *TelegraphPage) AddText(text string) {
	page.addContent(Node{
		Tag:      "p",
		Attrs:    []Attribute{},
		Children: fmt.Sprintf(`["%s"]`, text),
	})
}

func (page *TelegraphPage) AddImg(src string) {
	node := Node{
		Tag:      "img",
		Children: `[""]`,
	}
	node.AddAttr(Attribute{
		Name:  "src",
		Value: src,
	})
	page.addContent(node)
}

func (page *TelegraphPage) addContent(node Node) {
	page.Content = append(page.Content, node)
}

func (page TelegraphPage) StringContent() string {
	result, err := json.Marshal(page.Content)
	if err != nil {
		return ""
	}

	return string(result)
}

type Node struct {
	Tag      string      `json:"tag"`
	Attrs    []Attribute `json:"attrs"`
	Children string      `json:"children"`
}

func (n *Node) AddAttr(attribute Attribute) {
	n.Attrs = append(n.Attrs, attribute)
}

func (n Node) StringAttrs() string {
	result, err := json.Marshal(n.Attrs)
	if err != nil {
		return ""
	}

	return string(result)
}

type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (t *TelegraphClient) createPage(page TelegraphPage) (*http.Response, error) {
	requestURL, err := t.getRequestURL("createPage")
	if err != nil {
		return nil, err
	}

	values := requestURL.Query()
	values.Add("title", page.Title)
	values.Add("author_name", page.AuthorName)
	values.Add("content", page.StringContent())

	res, _ := json.MarshalIndent(page, "", "  ")
	fmt.Println(string(res))

	requestURL.RawQuery = values.Encode()
	fmt.Println(requestURL.String())
	return http.Get(requestURL.String())
}

func (t *TelegraphClient) getRequestURL(method string) (*url.URL, error) {
	result, err := url.Parse("https://api.telegra.ph")
	if err != nil {
		return nil, err
	}
	result = result.JoinPath(method)
	values := result.Query()
	values.Add("access_token", t.token)

	result.RawQuery = values.Encode()
	return result, nil
}

type TelegraphClient struct {
	token string
}

func NewTelegraphClient() (*TelegraphClient, error) {
	token, found := getTelegraphToken()
	if !found {
		return nil, errors.New("please set TELEGRAPH_TOKEN")
	}

	client := new(TelegraphClient)
	client.token = token
	return client, nil
}

func getTelegraphToken() (string, bool) {
	return os.LookupEnv("TELEGRAPH_TOKEN")
}
