package main

import (
  "fmt"
  "bytes"
  "encoding/json"
  "net/http"
  "os"
  "go-scraper/scraper"
)

func index_elasticsearch(ch chan scraper.PageData, url string) {
  client := &http.Client{}
  buf := new(bytes.Buffer)

  for {
    page, more := <-ch
    if !more {
      break
    }

    json.NewEncoder(buf).Encode(page)

    req, err := http.NewRequest("POST", url, buf)
    if err != nil {
      fmt.Println(err)
      return
    }

    req.Header.Set("Content-Type", "json")

    _, err = client.Do(req)
    if err != nil {
      fmt.Println(err)
    }
  }
}

const (
  CRAWL = iota
  LOAD_FROM_FILE = iota
)

func main() {
  s := scraper.New()

  pages := make(chan scraper.PageData)

  //archive_ch := make(chan scraper.PageData)
  index_ch := make(chan scraper.PageData)
  root_url := os.Args[1:][0]

  ops := CRAWL

  switch ops {
  case CRAWL:
    //scraper.Multiplex(pages, [2]chan scraper.PageData{archive_ch, index_ch,})

    ignored_urls := []string{
      "http:void(0)",
    }

    /*
    ignore_content_rules := []string {
      "div;class=well,sidebar-nav",
      "div;id=navigation",
      "form;id=commentform",
    }
    */

    s.IgnoreUrls(root_url, ignored_urls)
    go s.Crawl(root_url, pages)
    for {
      p, more := <- pages
      if !more {
        break
      }
      fmt.Println(p)
      _ = "breakpoint"
    }

    return
    //go scraper.SaveArchive(archive_ch, "data/archiver.json")
  case LOAD_FROM_FILE:
    //go scraper.LoadArchive(index_ch, "data/archiver.json")
  }

  go index_elasticsearch(index_ch, "http://localhost:9200/space_name/page")
}
