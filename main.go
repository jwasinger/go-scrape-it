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
  links []string
  url, data string;
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

func scrape(url string, out_content chan *PageData, wg *sync.WaitGroup) {
  resp, _ := http.Get(url)

  defer wg.Done();
  defer resp.Body.Close()
  
  var page_data PageData

  tokenizer := html.NewTokenizer(body)
  for {
    token := tokenizer.Next()

    switch {
    case token == html.ErrorToken:
      out_data <- &page_data
      return
    case token == html.TextToken:
      page_data.url.WriteString(token.Text())
    case token == html.StartTagToken:
      token := tokenizer.Token()
      for _, attr := range token.Attr {
        if attr.Key == "href" {
          page_data.links = append(page_data.links, attr.Val)
          continue
        }
      }
    case token == html.EndTagToken:
      continue
    }
  }

  out_content <- &page_data
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

  var delta_time int64
  var last_time int64

  for _, url := range urls {
    url, err := url.Parse(url)
    if err {
      panic(err)
    }

    last_time = time.Time.Now().Unix()

    go scrape(url, data, &wg)
    
    if time.Time.Now().Unix() - last_time < request_delay {

    }
  }

  done := false

  go func() {
    for !done {
      link := <-data
    }
  }()

  wg.Wait()
  done = true

  for _, link := range links{
    fmt.Println(link)
  }
}
