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

/* DASHAOZI 0.0.1 */
const TICK = 30
const RATE = 1.002

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
var Log string

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
	factor := M[CNY][BTC] * RATE * M[BTC][LTC] * RATE * M[LTC][CNY] * RATE

	//log.Println(factor)
	if factor < 1 && factor != 0 {
		log.Println("MakeAbitrage CNY -> BTC -> LTC")
		log.Println(M[CNY][BTC], M[BTC][LTC], M[LTC][CNY])
		Log = strconv.FormatFloat(factor, 'f', 6, 64)

		MakeTrade(true)
	}
	/* CNY -> LTC -> BTC */
	factor = M[CNY][LTC] * RATE * M[LTC][BTC] * RATE * M[BTC][CNY] * RATE

	//log.Println(factor)
	if factor < 1 && factor != 0 {
		log.Println("MakeAbitrage CNY -> LTC -> BTC")
		log.Println(M[CNY][LTC], M[LTC][BTC], M[BTC][CNY])
		Log = strconv.FormatFloat(factor, 'f', 6, 64)

		MakeTrade(false)
	}
}

func MakeTrade(dir bool) {

	for i := 0; i < 3; i++ {
		Log = Log + " Trading"
		cny, btc, ltc, _ := GetAccount()
		log.Println(i, cny, btc, ltc)
		if dir == true {
			if cny > 100 {
				Buy(M[CNY][BTC], cny/M[CNY][BTC]-0.0001, 0)
				Log = Log + " cny buy btc:" + strconv.FormatFloat((cny/M[CNY][BTC]-0.0001), 'f', 4, 64)
			}
			if btc > 0.02 {
				Buy(M[BTC][LTC], btc/M[BTC][LTC]-0.0001, 2)
				Log = Log + " btc buy ltc:" + strconv.FormatFloat((btc/M[BTC][LTC]-0.0001), 'f', 4, 64)
			}
			if ltc > 1 {
				Sell(M[LTC][CNY], ltc, 1)
				Log = Log + " ltc sell cny:" + strconv.FormatFloat((ltc), 'f', 4, 64)
			}
		} else {
			if cny > 100 {
				Buy(M[CNY][LTC], cny/M[CNY][LTC]-0.0001, 1)
				Log = Log + " cny buy ltc:" + strconv.FormatFloat((cny/M[CNY][LTC]-0.0001), 'f', 4, 64)
			}
			if ltc > 1 {
				Sell(M[LTC][BTC], ltc, 2)
				Log = Log + " ltc sell btc:" + strconv.FormatFloat((ltc), 'f', 4, 64)
			}
			if btc > 0.02 {
				Sell(M[BTC][CNY], btc, 0)
				Log = Log + " btc sell cny:" + strconv.FormatFloat((btc), 'f', 4, 64)
			}

		}
		time.Sleep(1)
	}
	CancelAllOrders()
}

func main() {
	log.Println("FXBTC Starts")
	changed := make(chan int)

	/* Initial Configuration */
	var cfg Config
	var err error
	err = LoadConfig(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	USERNAME = cfg.Email
	PASSWORD = cfg.Password

	GetToken()
	go func() {
		for {
			select {
			case <-time.Tick(time.Second * 60000):
				GetToken()
			}
		}
	}()
	log.Println(GetAccount())
	//Buy(0.02, 1, 1)
	//CancelAllOrders()

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
					//log.Println("Request")
					resp, err := http.Get("https://data.fxbtc.com/api?op=query_ticker&symbol=" + T[i].symbol)
					if err != nil {
						log.Println(err)
						continue
					}
					defer resp.Body.Close()

					body, _ := ioutil.ReadAll(resp.Body)
					js, err := NewJson(body)
					if err != nil {
						log.Println(err)
						continue
					}

					ask, _ := js.Get("ticker").Get("ask").Float64()
					bid, _ := js.Get("ticker").Get("bid").Float64()

					//log.Println(T[i].symbol, ask, bid)

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
	abitrage1 := strconv.FormatFloat(M[CNY][BTC]*RATE*M[BTC][LTC]*RATE*M[LTC][CNY]*RATE, 'f', 6, 64)
	abitrage2 := strconv.FormatFloat(M[CNY][LTC]*RATE*M[LTC][BTC]*RATE*M[BTC][CNY]*RATE, 'f', 6, 64)
	return "hello: A1: " + abitrage1 + " A2: " + abitrage2 + " LOG:" + Log + val
}
