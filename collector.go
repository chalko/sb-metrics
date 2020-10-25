package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"log"
	"math"
	"strconv"
	"strings"
)

type CableModemCollector struct {
}

type ModemStatus struct {
	startup StartupStatus
	ds      []DownStatus
}
type StartupStatus struct {
	bootState string
}

type DownStatus struct {
	id     string
	lock   string
	mod    string
	freq   int
	power  float64
	snr    float64
	corr   int
	uncorr int
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
		ds:      parseDs(ds),
	}
}

// Procedure Status Comment
type Psc struct {
	p string
	s string
	c string
}

func parseStartupStatus(table *goquery.Selection) StartupStatus {

	rows := nonHeaderRows(table)
	m := make(map[string]Psc, rows.Length())
	rows.Each(func(i int, s *goquery.Selection) {
		cells := s.Find("td")
		psc := Psc{
			p: cells.Get(0).FirstChild.Data,
			s: cells.Get(1).FirstChild.Data,
		}
		m[psc.p] = psc
	})

	return StartupStatus{
		bootState: m["Boot State"].s,
	}

}

func nonHeaderRows(table *goquery.Selection) *goquery.Selection {
	return table.Find("tr:has(td)").Filter("tr:not(:has(td strong))")
}

func parseDs(table *goquery.Selection) []DownStatus {
	rows := nonHeaderRows(table)
	ds := make([]DownStatus, rows.Length(), 32)
	rows.Each(func(i int, s *goquery.Selection) {
		cells := s.Find("td")
		ds[i] = DownStatus{
			id:     text(cells, 0),
			lock:   text(cells, 1),
			mod:    text(cells, 2),
			freq:   freq(text(cells, 3)),
			power:  float(text(cells, 4)," dBmV"),
			snr:    float(text(cells, 5)," dB"),
			corr:   corr(text(cells, 6)),
			uncorr: corr(text(cells, 7)),
		}
	})
	return ds
}

func corr(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return i

}

func float(s string, suffix string) float64 {
	f, err := strconv.ParseFloat(strings.TrimSuffix(s,suffix), 64)
	if err != nil {
		return math.NaN()
	}
	return f
}

func freq(s string) int {
	i, err := strconv.Atoi(strings.TrimSuffix(s," Hz"))
	if err != nil {
		return 0
	}
	return i
}

func text(cells *goquery.Selection, i int) string {
	node := cells.Get(i)
	return node.FirstChild.Data
}

func (cc CableModemCollector) Collect(ch chan<- prometheus.Metric) {

	//ms = scrape("http://192.168.100.1")

}
