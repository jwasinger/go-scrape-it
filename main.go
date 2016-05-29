package main

import (
  "fmt"
  "io"
  //"io/ioutil"
  "net/http"
  "net/url"
  "golang.org/x/net/html"
  "sync"
  //"go_scraper/stack"
  "github.com/andybalholm/cascadia"
)

type PageData struct {
  url, data string;
}

type ScrapeState struct {
  
}

func extract_text(html_root *html.Node) {
  for {
    
  }
}

func match_all(selector string, html_root *html.Node) []*html.Node {
  sel, err := Compile(string)
  if err {
    panic err
  }

  match_nodes, err := sel.MatchAll(html_root)
  if err {
    panic err
  }

  return match_nodes
}

func parse_data(body io.Reader, pages chan string) {
  //output_links := make(
  tokenizer := html.NewTokenizer(body)
  for {
    token := tokenizer.Next()

    switch {
    case token == html.ErrorToken:
      return
    case token == html.StartTagToken:
      token := tokenizer.Token()
      for _, attr := range token.Attr {
        if attr.Key == "href" {
          pages <- attr.Val
          continue
        }
      }
    case token == html.EndTagToken:
      continue
    }
  }
}

func scrape(url string, out_content chan string, wg *sync.WaitGroup) {
  resp, _ := http.Get(url)

  defer wg.Done();
  
  parse_data(resp.Body, out_content)

  /*
  c <- PageData{
    url,
    string(bytes))
  }
  */

  resp.Body.Close()
}

func main () {
  links := []string{}

  urls := []string{
    "http://reddit.com",
    "https://huffingtonpost.com",
  }

  data := make(chan string)
  var wg sync.WaitGroup

  wg.Add(len(urls));
  fmt.Println(len(urls))

  for _, url := range urls {
    url, err := url.Parse(url)
    if err {
      panic(err)
    }

    go scrape(url, data, &wg)
  }

  done := false

  go func() {
    for !done {
      link := <-data
      links = append(links, link)
    }
  }()

  wg.Wait()
  done = true

  for _, link := range links{
    fmt.Println(link)
  }
}
