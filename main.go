package main

import (
    "fmt"
    "time"
    "net/http"
    "html/template"    
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "log"
    "math/rand"
    "./scrapers"
    "encoding/json"
)

var db *mgo.Session

type Query struct {
  SumbitTime time.Time
  Text string
  Status string
  MentionsCount int
  LastMention time.Time
  ParseStatus map[string]int
  ID        bson.ObjectId `bson:"_id,omitempty"`
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

func randomSleepMS() int {
    return (rand.Intn(55)+5)*100
}


func saveMentions(men []scrapers.Mention) {
    session, err := mgo.Dial("localhost")
    b := session
    defer session.Close()
    c := b.DB("db").C("mentions")
    for _,m := range men{
        err = c.Insert(&m)
        if err != nil {
            log.Fatal(err)
        }
    }    
}

func do_scrap(scraper scrapers.Scraper, maxpage int) {

    sname := scraper.Name()

    session, err := mgo.Dial("localhost")
    defer session.Close()
    if err != nil {
        panic(err)
    }
   

    c := session.DB("db").C("queries")
    fmt.Println("Runed scraper: ", sname)
    updateStatus := func (ide  bson.ObjectId , status string) {
    	
        colQuerier := bson.M{"_id": ide}
        change := bson.M{"$set": bson.M{"status":status}}
        err := c.Update(colQuerier, change)
        if err != nil {
            panic(err)
        }
    }

    for {
		
      	query := Query{}

        iter:=c.Find(nil).Sort("-SumbitTime").Iter()
      	
        for iter.Next(&query) {

            full_processes := true;
			for _, v := range query.ParseStatus { 
				if v != maxpage {
					full_processes = false
				}
			}

			if full_processes {
					updateStatus(query.ID, "COMPLETE")
			}

            if query.ParseStatus[sname] == maxpage {
                continue
            } else {
                
                fmt.Printf("Goin scrap %s [%d] - %s\n",sname,query.ParseStatus[sname],query.Text)
                updateStatus(query.ID, fmt.Sprintf("Сrawling..."))
                mentions,banned  := scraper.Scrap(query.Text,query.ParseStatus[sname] )
                
                if banned {
                    fmt.Printf("Banned in %s [%d]\n",sname,query.ParseStatus[sname])
                    time.Sleep(2 * time.Minute)
                    continue
                }
                
                saveMentions(mentions)
                nps := fmt.Sprintf("parsestatus.%s",sname)
                mpg := query.ParseStatus[sname]+1
                colQuerier := bson.M{"_id": query.ID}
                change := bson.M{"$set": bson.M{nps:mpg,"mentionscount":query.MentionsCount+len(mentions)}}
                err := c.Update(colQuerier, change)
                if err != nil {
                    panic(err)
                }
            
                time.Sleep(time.Duration(randomSleepMS())*time.Millisecond)


            }
        }  
   }
        
}


func submitQueryForScrap(query string) {
    
    c := db.DB("db").C("queries")
            
    err := c.Insert(&Query{  SumbitTime: time.Now(),
    						 Text: query,
  							 Status: "Pending..." })
    if err != nil {
            log.Fatal(err)
    }
}

func getEnginesStatus() map[string]bool {

    return map[string]bool {}
}


func mainpage(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("./templates/main.html")
    t.Execute(w, getEnginesStatus())
}

func result(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func submit(w http.ResponseWriter, r *http.Request) {
    query := r.FormValue("query")
    if query != "" {
        submitQueryForScrap(query)
    }
    fmt.Println("Subminted", query)
    http.Redirect(w, r, "/", http.StatusFound) 
}

func getQueries() []Query {
    
    c := db.DB("db").C("queries")

    result := []Query{}
    
    iter:=c.Find(nil).Sort("-sumbittime").Iter()
    err := iter.All(&result)


    if err != nil {
            log.Fatal(err)
    }

    return result
}

func queries(w http.ResponseWriter, r *http.Request) {
    
    queries := getQueries()
    
   	b, err := json.Marshal(struct {Queries []Query} {Queries: queries})
	if err != nil {
		fmt.Println(err)
		return
	}
    	
    fmt.Fprintf(w, string(b[:]))
    
}


func getMentions(qwery string) []scrapers.Mention {
    
    c := db.DB("db").C("mentions")

    result := []scrapers.Mention{}
    
    iter := c.Find(bson.M{"query":qwery}).Iter()
    err := iter.All(&result)

    if err != nil {
            log.Fatal(err)
    }

    return result
}

func resultp(w http.ResponseWriter, r *http.Request) {
    
    mentions := getMentions(r.URL.Query()["q"][0])
    for _,m:= range mentions {
        fmt.Fprintf(w, "Aggr: %s<p>Source: %s<p>Name: %s<p>Time: %s <p>Link: %s<p>Brief: %s<p><p>",m.Aggr, m.Soruce, m.Name, m.Time, m.Link, m.Brief)
    }
}

func main() {

    session, err := mgo.Dial("localhost")
    
    if err != nil {
            panic(err)
    }
    db = session
    defer session.Close()

  
    // Optional. Switch the session to a monotonic behavior.
    session.SetMode(mgo.Monotonic, true)

   // fmt.Println((*scrapers.GetAllScrapers())["Google"].Scrap("hello",0))

    fmt.Println("Mongo client initialized")
    http.HandleFunc("/", mainpage)
    http.HandleFunc("/submit/", submit)
    http.HandleFunc("/result", resultp)
    http.HandleFunc("/data/queries/", queries)

    
    
    
    
    for k, v := range (*scrapers.GetAllScrapers()) { 
    	fmt.Println("GO ",k)
    	go do_scrap(v,5)
	}
    
    //go scrapper_google()
    //go scrapper_yandex()
   
    http.ListenAndServe(":8080", nil)
}