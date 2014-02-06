package main

import (
	. "github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"
)

const (
	CNY = iota
	BTC
	LTC
	MAX
)

type ticker struct {
	symbol string
	from   uint
	to     uint
}

/* global market matrix */
var M [MAX][MAX]float64
var T [MAX]ticker

func init() {
	T = [MAX]ticker{
		{
			symbol: "btc_cny",
			from:   0,
			to:     1,
		},
		{
			symbol: "ltc_cny",
			from:   0,
			to:     2,
		},
		{
			symbol: "ltc_btc",
			from:   2,
			to:     1,
		},
	}
}

func main() {
	runtime.GOMAXPROCS(3)
	log.Println("FXBTC Starts")
	changed := make(chan string)

	go func() {
		var who string
		for {
			who = <-changed
			log.Println(who, " changed")
		}
	}()
	//for {
	go func() {
		log.Println("Goroutine")
		for i := 0; ; i = (i + 1) % 3 {
			select {
			case <-time.Tick(time.Second * 10):
				{
					log.Println("Request")
					resp, _ := http.Get("https://data.fxbtc.com/api?op=query_ticker&symbol=" + T[i].symbol)
					defer resp.Body.Close()

					body, _ := ioutil.ReadAll(resp.Body)
					js, _ := NewJson(body)

					ask, _ := js.Get("ticker").Get("ask").Float64()
					bid, _ := js.Get("ticker").Get("bid").Float64()

					log.Println(T[i].symbol, ask, bid)

					if M[T[i].from][T[i].to] != ask || M[T[i].to][T[i].from] != 1.0/bid {
						M[T[i].from][T[i].to] = ask
						M[T[i].to][T[i].from] = 1.0 / bid
						changed <- T[i].symbol
					}
				}
			}
		}

	}()

	//}
	for {

	}
}
