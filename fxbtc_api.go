package main

import (
	"fmt"
	. "github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	//"regexp"
	"strconv"
	"strings"
	"time"
)

/*
  API 访问密匙 (Access Key)
  API 秘密密匙 (Secret Key)
*/
var (
	ACCESS_KEY = "API 访问密匙 (Access Key)"
	SECURT_KEY = "API 秘密密匙 (Secret Key)"
	USERNAME   = "USERNAME"
	PASSWORD   = "PASSWORD"
	TOKEN      = ""
)

var gCurCookies []*http.Cookie

func GetMarket() (float64, float64, float64, error) {
	resp, err := http.Get("https://www.okcoin.com/api/ticker.do?symbol=ltc_cny")
	if err != nil {
		log.Println(err)
		return 0, 0, 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return 0, 0, 0, err
	}

	js, err := NewJson(body)
	if err != nil {
		log.Println(err)
		return 0, 0, 0, err
	}

	last_str, err := js.Get("ticker").Get("last").String()
	if err != nil {
		log.Println(err)
		return 0, 0, 0, err
	}
	last, err := strconv.ParseFloat(last_str, 64)
	/* Buy/Sell should be from customer view */
	buy_str, err := js.Get("ticker").Get("sell").String()
	if err != nil {
		log.Println(err)
		return 0, 0, 0, err
	}
	buy, err := strconv.ParseFloat(buy_str, 64)

	sell_str, err := js.Get("ticker").Get("buy").String()
	if err != nil {
		log.Println(err)
		return 0, 0, 0, err
	}
	sell, err := strconv.ParseFloat(sell_str, 64)

	return last, buy, sell, err
}

func GetAccount() (cny float64, btc float64, ltc float64, err error) {

	client := &http.Client{}

	form := url.Values{
		"token": {TOKEN},
		"op":    {"get_info"},
	}

	//data := "partner=" + ACCESS_KEY + "&sign=" + suffix
	//fmt.Println(data)
	req, err := http.NewRequest("POST", "https://trade.fxbtc.com/api", strings.NewReader(form.Encode()))
	if err != nil {
		log.Println(err)
		return 0, 0, 0, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return 0, 0, 0, err
	}
	defer resp.Body.Close()
	//......

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body), err)
	if err != nil {
		log.Println(err)
		return 0, 0, 0, err
	}

	js, err := NewJson(body)
	if err != nil {
		log.Println(err)
		return
	}
	btc_str, err := js.Get("info").Get("funds").Get("free").Get("btc").String()
	if err != nil {
		log.Println(err)
		return
	}
	ltc_str, err := js.Get("info").Get("funds").Get("free").Get("ltc").String()
	if err != nil {
		log.Println(err)
		return
	}
	cny_str, err := js.Get("info").Get("funds").Get("free").Get("cny").String()
	if err != nil {
		log.Println(err)
		return
	}
	cny, err = strconv.ParseFloat(cny_str, 64)
	btc, err = strconv.ParseFloat(btc_str, 64)
	ltc, err = strconv.ParseFloat(ltc_str, 64)

	fmt.Println(cny, btc, ltc)

	return cny, btc, ltc, err
}

func Buy(buy float64, btc float64, symbol uint) string {
	fmt.Printf("%s: $$$ BUY %f at %f -> %s", time.Now(), btc, buy, T[symbol].symbol)
	return MakeOrder(buy, btc, symbol, true)
}

func Sell(sell float64, btc float64, symbol uint) string {
	fmt.Printf("%s: $$$ SELL %f at %f -> %s", time.Now(), btc, sell, T[symbol].symbol)
	return MakeOrder(sell, btc, symbol, false)
}

func MakeOrder(price float64, amount float64, symbol uint, buying bool) (id string) {
	var buying_str string
	if buying {
		buying_str = "buy"
	} else {
		buying_str = "sell"
	}

	price_str := strconv.FormatFloat(price, 'f', 2, 64)
	amount_str := strconv.FormatFloat(amount, 'f', 3, 64)

	form := url.Values{
		"token":  {TOKEN},
		"op":     {"trade"},
		"symbol": {T[symbol].symbol},
		"type":   {buying_str},
		"rate":   {price_str},
		"vol":    {amount_str},
	}

	req, err := http.NewRequest("POST", "https://trade.fxbtc.com/api", strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(resp)
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	defer resp.Body.Close()

	return "ok"

}

func GetOrders(symbol uint) (src []string) {

	form := url.Values{
		"token":  {TOKEN},
		"op":     {"get_orders"},
		"symbol": {T[symbol].symbol},
	}

	//fmt.Println(data)
	req, err := http.NewRequest("POST", "https://trade.fxbtc.com/api", strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(resp)
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	defer resp.Body.Close()

	//fmt.Println(src)
	js, err := NewJson(body)
	if err != nil {
		log.Println(err)
		return
	}
	id, err := js.Get("orders").GetIndex(0).Get("id").Int64()
	if err != nil {
		log.Println(err)
		return
	}
	id_str := strconv.FormatInt(id, 10)
	b := make([]string, 1)
	b[0] = id_str
	fmt.Println(b)

	return b
}

func CancelAllOrders() {
	Orderlist := make([]string, 10)

	for i := 0; i < 3; i++ {
		Orderlist = GetOrders(uint(i))

		for _, v := range Orderlist {
			CancelOrder(v, uint(i))
			time.Sleep(time.Second)
		}
	}
}

func CancelOrder(cancelID string, symbol uint) {

	form := url.Values{
		"token":  {TOKEN},
		"op":     {"cancel_order"},
		"symbol": {T[symbol].symbol},
		"id":     {cancelID},
	}

	req, err := http.NewRequest("POST", "https://trade.fxbtc.com/api", strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(resp)
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	defer resp.Body.Close()

	return
}

func GetToken() {
	form := url.Values{
		"op":       {"get_token"},
		"username": {USERNAME},
		"password": {PASSWORD},
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://trade.fxbtc.com/api", strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	js, err := NewJson(body)
	if err != nil {
		log.Println(err)
		return
	}
	valid, err := js.Get("result").Bool()
	if err != nil {
		log.Println(err)
		return
	}

	if valid != true {
		log.Println(err)
		return
	}

	TOKEN, _ = js.Get("token").String()

	return
}
