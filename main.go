package main

import (
  "fmt"
  "io"
  //"io/ioutil"
  "net/http"
  "golang.org/x/net/html"
  "sync"
  //"go_scraper/stack"
)

type PageData struct {
  url, data string;
}

func parse_links(body io.Reader, links chan string) (int, int) {
  //output_links := make(
  tokenizer := html.NewTokenizer(body)
  for {
    token := tokenizer.Next()

    switch {
    case token == html.ErrorToken:
      return 0, 1
    case token == html.StartTagToken:
      token := tokenizer.Token()
      for _, attr := range token.Attr {
        if attr.Key == "href" {
          links <- attr.Val
          continue
        }
      }
    case token == html.EndTagToken:
      continue
    }
  }

  return 0, 1
}

func scrape(url string, c chan string, wg *sync.WaitGroup) {
  resp, _ := http.Get(url)

  defer wg.Done();
  
  //links, content := parse_links(resp.Body)
  parse_links(resp.Body, c)

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
  }

  data := make(chan string);
  var wg sync.WaitGroup

  wg.Add(len(urls));
  fmt.Println(len(urls))

  for i := range urls {
    go scrape(urls[i], data, &wg);
  }

  done := false

  go func() {
    for {
      if done {
        break
      }

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
