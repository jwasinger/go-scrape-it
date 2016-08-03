package scraper

import (
  "fmt"
  "golang.org/x/net/html"
  "time"
  "net/http"
  "net/url"
  "strings"
  "regexp"
  "errors"
)

type Scraper struct {
  target_urls map[string]bool
  IgnoredUrls map[string]bool
  ScrapedUrls map[string]bool
}

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
  // Iterate over all of the Token's attributes until we find an "href"
  for _, a := range t.Attr {
    if a.Key == "href" {
      href = a.Val
      ok = true
    }
  }

  // "bare" return will return the variables (ok, href) as defined in
  // the function definition
  return
}

//remove all non-alphanumeric-characters
func cleanse_content(content *string) *string {
  reg, err := regexp.Compile("[^A-Za-z0-9]+")
  if err != nil {
    fmt.Println(err)
  }

  safe := reg.ReplaceAllString(*content, " ")
  safe = strings.ToLower(strings.Trim(safe, " "))
  return &safe
}

func extract_content_snippet(raw_bytes []byte, out_page_data *PageData) {
  content_snippet := string(raw_bytes)

  // cleanse the raw bytes of all non-alphanumeric characters
  reg, err := regexp.Compile("[^A-Za-z0-9]+")
  if err != nil {
    fmt.Println(err)
  }

  content_snippet = reg.ReplaceAllString(content_snippet, " ")
  content_snippet = strings.ToLower(strings.Trim(content_snippet, " "))

  if content_snippet != "" {
    out_page_data.Content += " - " + content_snippet
  }
}

func normalize_url(url_str string, href_str string) (string, error) {
  url, err := url.Parse(url_str)
  if err != nil {
    return "", err
  }

  href_url, err := url.Parse(href_str)
  if err != nil {
    return "", err
  }
  
  href_url.Scheme = url.Scheme

  if href_url.Host == "" {
    href_url.Host = url.Host
  } else if href_url.Host != url.Host {
    return "", errors.New("host mismatch")
  }

  href_url.Fragment = ""

  return href_url.String(), nil
}


func (scraper *Scraper) IgnoreUrls(page_url string, urls []string) {
  for _, url := range urls {
    normalized_url, _ := normalize_url(page_url, url)
    exists, _ := scraper.IgnoredUrls[normalized_url]
    if !exists {
      scraper.IgnoredUrls[normalized_url] = true
    }
  }
}

func (scraper *Scraper) shouldIgnore(href string, page_url string) bool {
  normalized_url, err := normalize_url(page_url, href)
  if err != nil{
    return true
  }

  _, exists := scraper.IgnoredUrls[normalized_url]
  if exists {
    return true
  }

  return false
}

func (scraper *Scraper) crawl_page(url string) (*PageData, error) {
  client := &http.Client{}
  page_content := new(PageData)
  page_content.Url = url
  page_content.Space = "space_name"

  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    return nil, err
  }

  req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.84 Safari/537.36")

  resp, err := client.Do(req)

  if err != nil {
    fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
    return nil, err
  }

  b := resp.Body
  defer b.Close() // close Body when the function returns

  z := html.NewTokenizer(b)
  extract_title := false
  level := 0
  ignore := -1

  for {
    tt := z.Next()

    switch {
    case tt == html.ErrorToken:
      //end of the document
      return page_content, nil
    case tt == html.TextToken: 
      if ignore != -1 {
        continue
      }

      if extract_title {
        page_content.Title = string(z.Text())
        extract_title = false
      }

      content_bytes := z.Text()
      extract_content_snippet(content_bytes, page_content)

    case tt == html.StartTagToken:
      level++
      t := z.Token()
      
      switch t.Data {
      case "title":
        extract_title = true
      case "script":
        ignore = level
      case "style":
        ignore = level
      case "a":
        ok, raw_href := getHref(t)
        if !ok {
          continue
        }

        if scraper.shouldIgnore(raw_href, url) {
          continue
        }

        normalized_url, _ := normalize_url(url, raw_href)

        if normalized_url == url {
          continue
        }

        //don't re-scrape pages
        _, exists := scraper.ScrapedUrls[normalized_url]
        if exists {
          continue
        }

        //don't add duplicates of same url to target_urls
        _, exists = scraper.target_urls[normalized_url]
        if exists {
          continue
        }
        scraper.target_urls[normalized_url] = true

      default:
      }
    case tt == html.EndTagToken:
      if ignore == level {
        ignore = -1 //the tag that was being ignored is no longer being parsed
      }
      level--
    }
  }
}

func pop_next_key(m map[string]bool) string {
  k := ""

  if len(m) == 0 {
    return ""
  }

  for key := range m {
    k = key
    break
  }

  delete(m, k)
  return k
}

func (scraper *Scraper) Crawl(root_url string, out_content chan PageData) {
  scraper.target_urls[root_url] = true
  for len(scraper.target_urls) != 0 {
    next_url := pop_next_key(scraper.target_urls)

    //scrape the page
    page_data, err := scraper.crawl_page(next_url)
    if err != nil {
      fmt.Println(err)
      continue
    }
    
    fmt.Println(len(scraper.ScrapedUrls))
    scraper.ScrapedUrls[next_url] = true

    //output the result
    out_content <- *page_data

    //wait two seconds between crawling pages (to reduce strain on whatever site is being crawled)
    time.Sleep(time.Second * 2)
  }


  close(out_content)
}

/*
func (scraper *Scraper) IgnoreContentRules(rules []string) {
  for _, rule := range rules {
    
  }
}
*/

func New() *Scraper{
  scraper := new(Scraper)
  scraper.IgnoredUrls = make(map[string]bool)
  scraper.ScrapedUrls = make(map[string]bool)
  scraper.target_urls = make(map[string]bool)

  return scraper
}
