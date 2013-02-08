package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	//"time"
)

type Query struct {
	ItemList []Item `xml:"channel>item"`
}

type Item struct {
	Title   string `xml:"title"`
	Content string `xml:"encoded"`
}

func main() {
	//for {
	stuff, err := http.Get("http://www.rlsbb.com/category/tv-shows/feed/")
	//fmt.Println("stuff: ", stuff)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	defer stuff.Body.Close()
	body, err := ioutil.ReadAll(stuff.Body)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	var q Query
	xml.Unmarshal([]byte(body), &q)
	for _, show := range q.ItemList {
		fmt.Println(show.Title)
	}
	//time.Sleep(20 * time.Minute)
	//}
}
