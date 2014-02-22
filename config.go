package main

import (
	//"bufio"
	"encoding/json"
	"fmt"
	//"io"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Id        uint
	Type      uint
	Email     string
	Password  string
	AccessKey string
	SecurtKey string
	Quick     uint
	Slow      uint
	QuickInit float64
	SlowInit  float64
	Delta     float64
	Pulse     uint
}

const CONFIG_FILE = "config.json"

func SaveConfig(config *Config) (err error) {
	fout, err := os.Create(CONFIG_FILE)
	defer fout.Close()
	if err != nil {
		fmt.Println(fout, err)
		return
	}
	/* pretty print */
	b, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		log.Println(err)
		return
	}

	//fout.WriteString("Just a config!\r\n")
	fout.Write(b)
	log.Println("CONFIG SAVED.")
	return
}

func LoadConfig(config *Config) (err error) {

	file, err := os.Open(CONFIG_FILE) // For read access.
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	//r := bufio.NewReader(file)
	//meta_json := r.ReadBytes(io.EOF)
	meta_json, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(meta_json, config)
	if err != nil {
		log.Panic(err)
	}
	/*ra, _ := ioutil.ReadFile("C:\\Windows\\win.ini")*/
	//fmt.Println(config)
	log.Println("CONFIG LOADED.")

	return
}
