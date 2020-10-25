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
}
type StartupStatus struct {
	bootState string
}

func scrape(reader io.Reader) ModemStatus {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}
	tables := doc.Find("table")
    startup := tables.Filter("table:has(th:contains('Startup Procedure'))")
    startupStatus := parseStartupStatus(startup)
    log.Print(startupStatus)

	return ModemStatus{
		startup: startupStatus,
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
func (cc CableModemCollector) Collect(ch chan<- prometheus.Metric) {

   //ms = scrape("http://192.168.100.1")

}

