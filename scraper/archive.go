package scraper

import (
  "os"
  "encoding/json"
  //"bufio"
  "fmt"
)

func SaveArchive(pages chan PageData, file string) {
  //delete old archive file (if exists) and create a new one
  if _, err := os.Stat(file); err == nil {
    err = os.Remove(file)
    if err != nil {
      panic(err)
    }
  }

  f, err := os.Create(file)
  if err != nil {
    panic(err)
  }

  encoder := json.NewEncoder(f)
  for {
    page, more := <-pages
    if !more {
      break
    }
    
    _ = "breakpoint"
    err = encoder.Encode(page)
    if err != nil {
      fmt.Println("err")
      panic(err)
    }
  }

  f.Close()
}

/*
func LoadArchive(out_pages chan PageData, file string) {
  f, err := os.Open(file)
  if err != nil {
    panic(err)
  }

  reader := bufio.NewReader(f)
  json_decoder := json.NewDecoder(reader)

  //read open bracket
  _, err = json_decoder.Token()
  if err != nil {
    panic(err)
  }

  var page_data PageData
  for json_decoder.More() {
    err := json_decoder.Decode(&page_data)
    if err != nil {
      panic(err)
    }

    out_pages <- page_data
  }

  //read close bracket
  _, err = json_decoder.Token()
  if err != nil {
    panic(err)
  }
  close(out_pages)
  f.Close()
}
*/
