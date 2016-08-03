package scraper

func Multiplex(src chan PageData, output [2]chan PageData) {
  go func(src chan PageData, output [2]chan PageData) {
    for {
      p, more := <-src
      if !more {
        for _, ch := range output {
          close(ch)
        }
        break
      }

      for _, ch := range output {
        ch <- p.Copy()
      }
    }
  }(src, output)
}
