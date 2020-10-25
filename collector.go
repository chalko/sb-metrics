package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"log"
)

type CableModemCollector struct {
}

type ModemStatus struct {
	startup StartupStatus
	ds [] DownStatus
}
type StartupStatus struct {
	bootState string
}

type DownStatus struct {

	id string
    lock string
	mod string
	freg int64
	power float32
	snr float32
	corr int64
	uncorr int64

}

func scrape(reader io.Reader) ModemStatus {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}
	tables := doc.Find("table")
    startup := tables.Filter("table:has(th:contains('Startup Procedure'))")
    ds := tables.Filter("table:has(th:contains('Downstream Bonded Channels'))")

	return ModemStatus{
		startup: parseStartupStatus(startup),
		ds: parseDs(ds),
	}
}


// Procedure Status Comment
type Psc struct {
	p string
	s string
	c string

}
func parseStartupStatus(table *goquery.Selection) StartupStatus {
	m := make(map[string]Psc)
	table.Find("tr:has(td)") .Each(func(i int, s *goquery.Selection) {
		log.Print(s.Html())
		cells := s.Find("td")
		psc := Psc{
			p: cells.Get(0).FirstChild.Data,
			s: cells.Get(1).FirstChild.Data,

		}
		m[psc.p]=psc
	})

	return StartupStatus{
		bootState: m["Boot State"].s,
	}

}

func parseDs(table *goquery.Selection) []DownStatus {
	rows := table.Find("tr:has(td)").Filter("tr:not(:has(td strong))")
	ds :=  make([]DownStatus, rows.Length(), 32)
	rows.Each(func(i int, s *goquery.Selection) {
		cells := s.Find("td")
		ds[i] = DownStatus{
			id: cells.Get(0).FirstChild.Data,
		}
	})
    return ds
}


func (cc CableModemCollector) Collect(ch chan<- prometheus.Metric) {

   //ms = scrape("http://192.168.100.1")

}

