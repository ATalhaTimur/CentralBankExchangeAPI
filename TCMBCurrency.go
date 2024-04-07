package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

type CurrencyDay struct {
	ID         string
	Date       time.Time
	DayNo      string
	Currencies []Currency
}

type Currency struct {
	Code            string
	CrossOrder      int
	Unit            int
	CurrencyNameTR  string
	CurrencyNameEN  string
	ForexBuying     float64
	ForexSelling    float64
	BanknoteBuying  float64
	BanknoteSelling float64
	CrossRateUSD    float64
	CrossRateOther  float64
}

type tarih_Date struct {
	XMlName   xml.Name `xml:"Tarih_Date"`
	Tarih     string   `xml:"Tarih,attr"`
	Date      string   `xml:"Date,attr"`
	Bulten_No string   `xml:"Bulten_No,attr"`
	Currency  []currency
}

type currency struct {
	Kod             string `xml:"Kod,attr"`
	CrossOrder      string `xml:"CrossOrder,attr"`
	CurrencyCode    string `xml:"CurrencyCode,attr"`
	Unit            string `xml:"Unit"`
	Isim            string `xml:"Isim"`
	CurrencyName    string `xml:"CurrencyName"`
	ForexBuying     string `xml:"ForexBuying"`
	ForexSelling    string `xml:"ForexSelling"`
	BanknoteBuying  string `xml:"BanknoteBuying"`
	BanknoteSelling string `xml:"BanknoteSelling"`
	CrossRateUSD    string `xml:"CrossRateUSD"`
	CrossRateOther  string `xml:"CrossRateOther"`
}

func (c *CurrencyDay) GetData(CurrencyDate time.Time) {

	xDate := CurrencyDate
	t := new(tarih_Date)
	currDay := t.GetData(CurrencyDate, xDate)

	for {
		if currDay == nil {
			CurrencyDate = CurrencyDate.AddDate(0, 0, -1)
			currDay := t.GetData(CurrencyDate, xDate)
			if currDay != nil {
				break
			}
		} else {
			break
		}

	}
}

func (c *tarih_Date) GetData(CurrencyDate time.Time, xDate time.Time) *CurrencyDay {

	currDay := new(CurrencyDay)
	var resp *http.Response
	var err error
	var url string
	currDay = new(CurrencyDay)
	url = "http://www.tcmb.gov.tr/kurlar/" + CurrencyDate.Format("200601") + "/" + CurrencyDate.Format("02012006") + ".xml"

	resp, err = http.Get(url)
	if err != nil {
		return nil
	} else {

		defer resp.Body.Close()
	}
	if resp.StatusCode != http.StatusNotFound {
		tarih := new(tarih_Date)
		d := xml.NewDecoder(resp.Body)
		marshalErr := d.Decode(&tarih)
		if marshalErr != nil {
			log.Printf("Error : %v", marshalErr)
		}

		c = &tarih_Date{}
		currDay.ID = xDate.Format("20060102")
		currDay.Date = xDate
		currDay.DayNo = tarih.Bulten_No
		currDay.Currencies = make([]Currency, len(tarih.Currency))
		for i, curr := range tarih.Currency {
			currDay.Currencies[i].Code = curr.CurrencyCode
			currDay.Currencies[i].CurrencyNameEN = curr.CurrencyName
			currDay.Currencies[i].CurrencyNameTR = curr.Isim
			currDay.Currencies[i].BanknoteBuying, err = strconv.ParseFloat(curr.BanknoteBuying, 64)
			currDay.Currencies[i].BanknoteSelling, err = strconv.ParseFloat(curr.BanknoteSelling, 64)
			currDay.Currencies[i].ForexBuying, err = strconv.ParseFloat(curr.ForexBuying, 64)
			currDay.Currencies[i].ForexSelling, err = strconv.ParseFloat(curr.ForexSelling, 64)
			currDay.Currencies[i].CrossOrder, err = strconv.Atoi(curr.CrossOrder)
			currDay.Currencies[i].CrossRateUSD, err = strconv.ParseFloat(curr.CrossRateUSD, 64)
			currDay.Currencies[i].CrossRateOther, err = strconv.ParseFloat(curr.CrossRateOther, 64)
			currDay.Currencies[i].Unit, err = strconv.Atoi(curr.Unit)
		}
		fmt.Println(currDay)
		SaveJson("CurrencyDay.json", currDay)

	} else {
		currDay = nil
	}
	return currDay
}

func SaveJson(fileName string, key interface{}) {
	outFile, err := os.Create(fileName)
	CheckError(err)
	defer outFile.Close()
	encoder := json.NewEncoder(outFile)
	err = encoder.Encode(key)
	CheckError(err)
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func main() {
	runtime.GOMAXPROCS(2)
	startTime := time.Now()
	CurrencyDay := new(CurrencyDay)
	CurrencyDate := time.Now()
	CurrencyDay.GetData(CurrencyDate)

	elepsedTime := time.Since(startTime)
	fmt.Println("Time : ", elepsedTime)

}
