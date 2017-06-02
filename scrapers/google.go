package scrapers
import (
    "fmt"
    "net/url"
    "github.com/PuerkitoBio/goquery"
    "time"
    "strings"
    "strconv"
)

type GoogleScraper struct {
    AnyScraper
}



func (scrp GoogleScraper) Scrap(query string, page int) ([]Mention, bool) {

    parseTime := func(tm string) time.Time {
        
        day := strings.Split(tm," ")[0]
        if _, err := strconv.Atoi(day); err == nil {
           return time.Now()
        }

        lay := "Jan 2, 2006"
        t, err := time.Parse(lay, strings.Trim(tm," "))
        if err != nil {
            fmt.Println(err)
        }
     
        return t
    }

    var mentions [] Mention

    Url, _ := url.Parse("https://www.google.com/search")

    parameters := url.Values{"hl":{"en"},"q":{query},"tbm":{"nws"},"start":{strconv.Itoa(10*page)}}
    Url.RawQuery = parameters.Encode()
    
    doc := getURL(scrp.bow, Url.String())
    
    if strings.Contains(scrp.bow.Url().String(),"google.com/sorry") {
        scrp.banned = true
        scrp.bantime = time.Now()
        return mentions, true
    }
    
    doc.Find(".g").Each(func(i int, s *goquery.Selection) {

        var m Mention
        
        nametg := s.Find("a")
        m.Name = nametg.Text()
        
        link,_ := nametg.Attr("href")//("href")
        
        m.Link = link
        m.Brief = s.Find(".st").Text()

        src_time := s.Find("span").Text()        
        
        src_time_sl := strings.Split(src_time,"-")
        
        src := strings.Join(src_time_sl[0 : len(src_time_sl)-1],"")
        
        time := src_time_sl[len(src_time_sl)-1]
        
        m.Time = parseTime(time)
        m.Soruce = src
        m.Aggr = "Google"
        m.Query = query
        
        mentions = append(mentions,m)
  })
    
    scrp.banned = false;
    return mentions, false
}
