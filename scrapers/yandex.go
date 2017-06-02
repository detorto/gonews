package scrapers
import (
    "fmt"
    "net/url"
    "github.com/PuerkitoBio/goquery"
    "time"
    "strings"
    "strconv"

   
)


func getMonthLong(m string) time.Month {
      
    mm  := map[string]time.Month {"января":time.January,
                                  "февраля":time.February,
                                  "марта":time.March,
                                  "апреля":time.April,
                                  "мая":time.May,
                                  "июня":time.June,
                                  "июля":time.February,
                                  "августа":time.August,
                                  "сентября":time.September,
                                  "октября":time.October,
                                  "ноября":time.November,
                                  "декбря":time.December}
    return mm[strings.ToLower(strings.Trim(m,""))]
}


type YandexScraper struct {
    AnyScraper
}


func (scrp YandexScraper) Scrap(query string, page int) ([]Mention, bool) {

    parseTime := func(tm string) time.Time {
        
        var day int
        var month time.Month
        var year int
       
        firstl := strings.Split(tm," ")
    
        if len(firstl) <= 1{
            return time.Now()
            }

        first:= firstl[0]
        if _ ,err := strconv.Atoi(first); err == nil {
            
            d,_  :=  strconv.Atoi(first)
            day = d
            fmt.Println(strings.Split(tm," "))
            month = getMonthLong(strings.Split(tm," ")[1])
            year = time.Now().Year()
            
            fmt.Println("%s %s %s\n",year,month,day)
            return time.Date(year,month,day,23, 0, 0, 0, time.UTC)
        } else {
            if first == "вчера"{
                year = time.Now().Year()
                day = time.Now().Day()
                month = time.Now().Month()

                return time.Date(year,month,day-1,23, 0, 0, 0, time.UTC)
            }

            
        
            layout := "02.01.06"
       
            t, err := time.Parse(layout, first)
            if err != nil {
                panic(err)
            }
            return t

        }
            
        return time.Now()
    }

    var mentions [] Mention

    Url, _ := url.Parse("https://news.yandex.ru/yandsearch")
    
    parameters := url.Values{"numdoc":{"30"},"text":{query},"rpt":{"nnews"},"p":{strconv.Itoa(page)}, "rel":{"tm"}}
    Url.RawQuery = parameters.Encode()
    
    doc := getURL(scrp.bow, Url.String())
    if strings.Contains(scrp.bow.Url().String(),"https://news.yandex.ru/showcaptcha") {
        scrp.banned = true;
        scrp.bantime = time.Now()
        return mentions, true
    }
    
    doc.Find(".search-item").Each(func(i int, s *goquery.Selection) {

        var m Mention
        
        m.Soruce = s.Find(".document__provider-name").Text()
        m.Name = s.Find(".document__title").Text()
            
        m.Time = parseTime(s.Find(".document__time").Text())
       
        link,_ := s.Find(".link").Attr("href")
        m.Link = link
        
        m.Brief = s.Find(".document__snippet").Text()
        m.Aggr = "Yandex"            
        m.Query = query

        mentions = append(mentions,m)
  })
    scrp.banned = false;
    return mentions,false
}