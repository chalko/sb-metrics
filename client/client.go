package client

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CableModemClient struct {
	connectionString string
	timeout          time.Duration
}

type ModemStatus struct {
	Startup StartupStatus
	Ds      []DownStatus
	Us      []UpStatus
}
type StartupStatus struct {
	Acquire       string
	ConnState     string
	BootState     string
	DownFreq      int
	ConfigFile    string
	SecurityState string
	SecurityType  string
	DocsisAccess  string
}

type DownStatus struct {
	Id     string
	Lock   string
	Mod    string
	Freq   int
	Power  float64
	Snr    float64
	Corr   int
	Uncorr int
}

type UpStatus struct {
	Num    string
	Id     string
	Lock   string
	ChType string
	Width  int
	Freq   int
	Power  float64
}

func NewCableModemClient(connectionString string, timeout time.Duration) *CableModemClient {
	return &CableModemClient{
		connectionString: connectionString,
		timeout:          timeout,
	}
}

func (c *CableModemClient) GetModemStatus() (ModemStatus, error) {
	//ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	client := http.Client{
		Timeout: c.timeout,
	}
	resp, err := client.Get(c.connectionString)
	if err != nil {
		return ModemStatus{}, err
	}
	defer resp.Body.Close()
	return scrape(resp.Body)

}

func scrape(reader io.Reader) (ModemStatus, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return ModemStatus{}, err
	}
	tables := doc.Find("table")
	startup := tables.Filter("table:has(th:contains('Startup Procedure'))")
	ds := tables.Filter("table:has(th:contains('Downstream Bonded Channels'))")
	us := tables.Filter("table:has(th:contains('Upstream Bonded Channels'))")

	return ModemStatus{
		Startup: parseStartupStatus(startup),
		Ds:      parseDs(ds),
		Us:      parseUs(us),
	}, nil
}

// Status Comment
type StatComm struct {
	s string
	c string
}

func parseStartupStatus(table *goquery.Selection) StartupStatus {
	rows := nonHeaderRows(table)
	m := make(map[string]StatComm, rows.Length())
	if rows.Length() > 0 {
		rows.Each(func(i int, s *goquery.Selection) {
			cells := s.Find("td")
			td := cells.First()
			m[td.Text()] = StatComm{
				s: td.Next().Text(),
				c: td.Next().Next().Text(),
			}
		})
	}
	return StartupStatus{
		Acquire:       m["Acquire Downstream Channel"].c,
		DownFreq:      hz(m["Acquire Downstream Channel"].s),
		ConnState:     m["Connectivity State"].s,
		BootState:     m["Boot State"].s,
		ConfigFile:    m["Configuration File"].s,
		SecurityState: m["Security"].s,
		SecurityType:  m["Security"].c,
		DocsisAccess:  m["DOCSIS Network Access Enabled"].s,
	}

}

func nonHeaderRows(table *goquery.Selection) *goquery.Selection {
	return table.Find("tr:has(td)").Filter("tr:not(:has(td strong))")
}

func parseDs(table *goquery.Selection) []DownStatus {
	rows := nonHeaderRows(table)
	ds := make([]DownStatus, rows.Length())
	if rows.Length() > 0 {
		rows.Each(func(i int, s *goquery.Selection) {
			cells := s.Find("td")
			ds[i] = DownStatus{
				Id:     text(cells, 0),
				Lock:   text(cells, 1),
				Mod:    text(cells, 2),
				Freq:   hz(text(cells, 3)),
				Power:  float(text(cells, 4), " dBmV"),
				Snr:    float(text(cells, 5), " dB"),
				Corr:   corr(text(cells, 6)),
				Uncorr: corr(text(cells, 7)),
			}
		})
	}
	return ds
}

func parseUs(table *goquery.Selection) []UpStatus {
	rows := nonHeaderRows(table)
	us := make([]UpStatus, rows.Length())
	if rows.Length() > 0 {
		rows.Each(func(i int, s *goquery.Selection) {
			cells := s.Find("td")
			us[i] = UpStatus{
				Num:    text(cells, 0),
				Id:     text(cells, 1),
				Lock:   text(cells, 2),
				ChType: text(cells, 3),
				Width:  hz(text(cells, 4)),
				Freq:   hz(text(cells, 5)),
				Power:  float(text(cells, 6), " dBmV"),
			}
		})
	}
	return us
}

func corr(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return i

}

func float(s string, suffix string) float64 {
	f, err := strconv.ParseFloat(strings.TrimSuffix(s, suffix), 64)
	if err != nil {
		return math.NaN()
	}
	return f
}

func hz(s string) int {
	i, err := strconv.Atoi(strings.TrimSuffix(s, " Hz"))
	if err != nil {
		return 0
	}
	return i
}

func text(cells *goquery.Selection, i int) string {
	node := cells.Get(i)
	return node.FirstChild.Data
}
