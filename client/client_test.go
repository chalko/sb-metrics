package client

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestScrape_good(t *testing.T) {
	reader, err := os.Open(filepath.Join("../testdata", "index.html"))
	if err != nil {
		log.Fatal(err)
	}
	ans, _ := scrape(reader)
	assert.Equal(t,
		StartupStatus{
			Acquire:       "Locked",
			DownFreq:      549000000,
			ConnState:     "OK",
			BootState:     "OK",
			ConfigFile:    "OK",
			SecurityState: "Enabled",
			SecurityType:  "BPI+",
			DocsisAccess:  "Allowed",
		},
		ans.Startup)
	assert.Equal(t,
		[]DownStatus{
			{Id: "16", Lock: "Locked", Mod: "QAM256", Freq: 549000000, Power: -4.6, Snr: 38.8, Corr: 204, Uncorr: 443},
			{Id: "1", Lock: "Locked", Mod: "QAM256", Freq: 459000000, Power: -0.2, Snr: 42.1, Corr: 984, Uncorr: 3588},
			{Id: "2", Lock: "Locked", Mod: "QAM256", Freq: 465000000, Power: -0.4, Snr: 41.9, Corr: 1050, Uncorr: 3613},
			{Id: "3", Lock: "Locked", Mod: "QAM256", Freq: 471000000, Power: -0.6, Snr: 41.8, Corr: 867, Uncorr: 3231},
			{Id: "4", Lock: "Locked", Mod: "QAM256", Freq: 477000000, Power: -0.7, Snr: 41.8, Corr: 885, Uncorr: 3630},
			{Id: "5", Lock: "Locked", Mod: "QAM256", Freq: 483000000, Power: -1, Snr: 41.2, Corr: 962, Uncorr: 3633},
			{Id: "6", Lock: "Locked", Mod: "QAM256", Freq: 489000000, Power: -1, Snr: 41.4, Corr: 1014, Uncorr: 3502},
			{Id: "7", Lock: "Locked", Mod: "QAM256", Freq: 495000000, Power: -1, Snr: 41.4, Corr: 1032, Uncorr: 4042},
			{Id: "8", Lock: "Locked", Mod: "QAM256", Freq: 501000000, Power: -1.3, Snr: 41.2, Corr: 1051, Uncorr: 4068},
			{Id: "9", Lock: "Locked", Mod: "QAM256", Freq: 507000000, Power: -1.9, Snr: 40.9, Corr: 1076, Uncorr: 3985},
			{Id: "10", Lock: "Locked", Mod: "QAM256", Freq: 513000000, Power: -2.2, Snr: 40.6, Corr: 1079, Uncorr: 4247},
			{Id: "11", Lock: "Locked", Mod: "QAM256", Freq: 519000000, Power: -2.4, Snr: 40.5, Corr: 991, Uncorr: 3836},
			{Id: "12", Lock: "Locked", Mod: "QAM256", Freq: 525000000, Power: -2.7, Snr: 40.3, Corr: 955, Uncorr: 3859},
			{Id: "13", Lock: "Locked", Mod: "QAM256", Freq: 531000000, Power: -3.1, Snr: 40, Corr: 980, Uncorr: 3690},
			{Id: "14", Lock: "Locked", Mod: "QAM256", Freq: 537000000, Power: -3.5, Snr: 39.6, Corr: 894, Uncorr: 3257},
			{Id: "15", Lock: "Locked", Mod: "QAM256", Freq: 543000000, Power: -4.1, Snr: 39.3, Corr: 1096, Uncorr: 4190},
			{Id: "17", Lock: "Locked", Mod: "QAM256", Freq: 555000000, Power: -5.2, Snr: 38.2, Corr: 1013, Uncorr: 3879},
		}, ans.Ds)

	assert.Equal(t,
		[]UpStatus{
			{Num: "1", Id: "57", Lock: "Locked", ChType: "SC-QAM", Width: 26000000, Freq: 6400000, Power: 50},
			{Num: "2", Id: "58", Lock: "Locked", ChType: "SC-QAM", Width: 19200000, Freq: 6400000, Power: 48},
		},
		ans.Us)
}

func TestScrape_down(t *testing.T) {
	reader, err := os.Open(filepath.Join("../testdata", "status-down.html"))
	if err != nil {
		log.Fatal(err)
	}
	ans, _ := scrape(reader)
	assert.Equal(t,
		StartupStatus{
			Acquire:       "In Progress",
			ConnState:     "In Progress",
			BootState:     "In Progress",
			ConfigFile:    "In Progress",
			SecurityState: "Disabled",
			SecurityType:  "Disabled",
			DocsisAccess:  "Denied",
		},

		ans.Startup)
}
