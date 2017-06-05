package scrapers

import (
    
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "time"
    //"strings"
    "math/rand"
    "gopkg.in/headzoo/surf.v1"
    "github.com/headzoo/surf/browser"
    "github.com/headzoo/surf/agent"
   
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
    ua := []string{ agent.AOL(),
                    agent.Chrome(),
                    agent.Firefox(),
                    agent.Konqueror(),
                    agent.MSIE(),
                    agent.Opera(),
                    agent.Safari()}

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



