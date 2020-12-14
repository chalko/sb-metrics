package client

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const namespace = "cable_modem"

var (

	// Metrics
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Has the cable modem is up",
		nil, nil,
	)
	// Startup Status
	acquire = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "acquire"),
		"Has the cable modem acquired a lock",
		nil, nil,
	)
	downChannel = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "down_channel"),
		"The downstream channel frequency in Hz",
		nil, nil,
	)

	channelLabels = []string{"id"}
	// down status Metrics
	downFreq = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "down_freq"),
		"The downstream frequency",
		channelLabels, nil,
	)

	downPower = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "down_power"),
		"The downstream power level",
		channelLabels, nil,
	)

	downCorr = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "down_corr"),
		"The downstream correctable errors",
		channelLabels, nil,
	)
	downUncorr = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "down_uncorr"),
		"The downstream uncorrectable errors",
		channelLabels, nil,
	)

	// upStream metrids
	upFreq = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up_freq"),
		"The upstream frequency",
		channelLabels, nil,
	)

	upPower = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up_power"),
		"The upstream power level",
		channelLabels, nil,
	)
)

type CableModemClient struct {
	connectionString string
	timeout          time.Duration
}

func (c *CableModemClient) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- acquire
	ch <- downChannel
	ch <- downFreq
	ch <- downPower
	ch <- downCorr
	ch <- downUncorr
	ch <- upFreq
	ch <- upPower
}
func (c *CableModemClient) Collect(ch chan<- prometheus.Metric) {
	status, err := c.GetModemStatus()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0,
		)
		log.Println(err)
		return
	}
	ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1)
	status.Startup.Collect(ch)
	for _, ds := range status.Ds {
		ds.Collect(ch)
	}
	for _, us := range status.Us {
		us.Collect(ch)
	}
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

func (s *StartupStatus) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(acquire, prometheus.GaugeValue, s.isAcquired())
	ch <- prometheus.MustNewConstMetric(downChannel, prometheus.GaugeValue, float64(s.DownFreq))
}

func (s *StartupStatus) isAcquired() float64 {
	var locked float64 = 0
	if s.Acquire == "Locked" {
		locked = 1
	}
	return locked
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

func (ds *DownStatus) Collect(ch chan<- prometheus.Metric) {
	id := ds.Id
	ch <- prometheus.MustNewConstMetric(downFreq, prometheus.GaugeValue, float64(ds.Freq), id)
	ch <- prometheus.MustNewConstMetric(downPower, prometheus.GaugeValue, ds.Power, id)
	ch <- prometheus.MustNewConstMetric(downCorr, prometheus.CounterValue, float64(ds.Corr), id)
	ch <- prometheus.MustNewConstMetric(downUncorr, prometheus.CounterValue, float64(ds.Uncorr), id)
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

func (us *UpStatus) Collect(ch chan<- prometheus.Metric) {
	id := us.Id
	ch <- prometheus.MustNewConstMetric(upFreq, prometheus.GaugeValue, float64(us.Freq), id)
	ch <- prometheus.MustNewConstMetric(upPower, prometheus.GaugeValue, us.Power, id)
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
