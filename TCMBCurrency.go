package main

import (
	"encoding/xml"
	"log"
	"net/http"
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
	Unit            string `xml:"Unit,attr"`
	Isim            string `xml:"Isim,attr"`
	CurrencyName    string `xml:"CurrencyName,attr"`
	ForexBuying     string `xml:"ForexBuying,attr"`
	ForexSelling    string `xml:"ForexSelling,attr"`
	BanknoteBuying  string `xml:"BanknoteBuying,attr"`
	BanknoteSelling string `xml:"BanknoteSelling,attr"`
	CrossRateUSD    string `xml:"CrossRateUSD,attr"`
	CrossRateOther  string `xml:"CrossRateOther,attr"`
}

func (c *CurrencyDay) GetData(CurrencyDate time.Time) {

	xDate := CurrencyDate
	t := new(tarih_Date)
	currDay := t.getDate(CurrencyDate, xDate)

	for {
		if currDay == nil {
			CurrencyDate = CurrencyDate.AddDate(0, 0, -1)
			currDay = t.getDate(CurrencyDate, xDate)
			if currDay != nil {
				break
			}
		} else {
			break
		}

	}
}

func (c *tarih_Date) getDate(CurrencyDate time.Time, xDate time.Time) *CurrencyDay {

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
	}
}
