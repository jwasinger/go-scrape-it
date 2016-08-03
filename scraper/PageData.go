package scraper

type PageData struct {
  Url string
  Title string
  Content string
  Space string
}

func (page *PageData) Copy() PageData {
  new_page := PageData{ page.Url, page.Title, page.Content, page.Space }
  return new_page
}
