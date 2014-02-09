package main

import (
	. "github.com/bitly/go-simplejson"
	"github.com/hoisie/web"
	"io/ioutil"
	"log"
	"net/http"
	//"runtime"
	"strconv"
	"time"
)

const TICK = 30

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
			from:   1,
			to:     2,
		},
	}
}

func MakeAbitrage() {

	/* CNY -> BTC -> LTC */
	log.Println(M[0][1], M[1][2], M[2][0])
	factor := M[0][1] * M[1][2] * M[2][0]

	log.Println(factor)
	if factor < 1 && factor != 0 {
		log.Panicln("haha CNY -> BTC -> LTC")
	}
	/* CNY -> LTC -> BTC */
	factor = M[0][2] * M[2][1] * M[1][0]

	log.Println(factor)
	if factor < 1 && factor != 0 {
		log.Panicln("haha CNY -> LTC -> BTC")
	}
}

func main() {
	log.Println("FXBTC Starts")
	changed := make(chan int)

	go func() {
		var tid int
		for {
			tid = <-changed
			log.Println(T[tid].symbol, " changed")
			MakeAbitrage()
		}
	}()
	//for {
	go func() {
		log.Println("Goroutine")
		for i := 0; ; i = (i + 1) % 3 {
			select {
			case <-time.Tick(time.Second * TICK):
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
						changed <- i
					}
				}
			}
		}

	}()

	web.Get("/(.*)", hello)
	web.Run("0.0.0.0:7777")
}

func hello(ctx *web.Context, val string) string {
	abitrage1 := strconv.FormatFloat(M[0][1]*M[1][2]*M[2][0], 'f', 6, 64)
	abitrage2 := strconv.FormatFloat(M[0][2]*M[2][1]*M[1][0], 'f', 6, 64)
	return "hello: A1: " + abitrage1 + " A2: " + abitrage2 + val
}
