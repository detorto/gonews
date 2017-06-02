package scrapers

import (
    
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "time"
    //"strings"
    "math/rand"
    "gopkg.in/headzoo/surf.v1"
    "github.com/headzoo/surf/browser"
   
)

type Mention struct {
  Query string
  Time time.Time
  Brief string
  Link string
  Name string
  Soruce string
  Aggr string
}


func randomUserAgent() string {
    ua := []string{ "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
                        "Mozilla/5.0 (Windows NT 6.3; WOW64; rv:53.0) Gecko/20100101 Firefox/53.0",
                        "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
                        "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:53.0) Gecko/20100101 Firefox/53.0",
                        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.96 Safari/537.36",
                        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/603.1.30 (KHTML, like Gecko) Version/10.1 Safari/603.1.30",
                        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36",
                    }

    return ua[rand.Intn(len(ua))]
}


func makeNewBOW() *browser.Browser {    
	    bow := surf.NewBrowser()
    	bow.SetUserAgent(randomUserAgent())
    	return bow
}

type Scraper interface {
	Scrap(query string, page int) ([]Mention, bool)
	Banned() bool
	BanTime() time.Time
	Name() string
}

type AnyScraper struct {
	bow *browser.Browser
	banned bool
	bantime time.Time
	name string
}

func (scrp AnyScraper) Banned() bool {
    return scrp.banned
}

func (scrp AnyScraper) BanTime() time.Time {
    return scrp.bantime
}

func (scrp AnyScraper) Name() string {
    return scrp.name
}

var all_scrapers map[string]Scraper

func init() {
  	fmt.Println("init!")
  	gscrap := &GoogleScraper{AnyScraper{makeNewBOW(), false, time.Now(), "Google"}}
  	yscrap := &YandexScraper{AnyScraper{makeNewBOW(), false, time.Now(), "Yandex"}}

  	fmt.Println("G: ",gscrap)
  	fmt.Println("Y: ",yscrap)

    all_scrapers = map[string]Scraper {"Google":gscrap,
    								   "Yandex":yscrap}
}

func getURL(bow  *browser.Browser, url string) (*goquery.Selection) {
    
    err := bow.Open(url)
    
    if err != nil {
        panic(err)
    }

    return bow.Dom()
}
    
func GetAllScrapers() *map[string]Scraper{
	return &all_scrapers
}



