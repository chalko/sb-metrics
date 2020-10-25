package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestScrape(t *testing.T) {
	reader, err := os.Open(filepath.Join("testdata", "index.html"))
	if err != nil {
		log.Fatal(err)
	}
	ans := scrape(reader)
	assert.Equal(t, ModemStatus{
		startup: StartupStatus{bootState: "OK"},
		ds: []DownStatus{
			{id: "16", lock: "Locked", mod: "QAM256", freq: 549000000, power: -4.6, snr: math.NaN(), corr: 204, uncorr: 443},
			{id: "1", lock: "Locked", mod: "QAM256", freq: 459000000, power: -0.2, snr: math.NaN(), corr: 984, uncorr: 3588},
			{id: "2", lock: "Locked", mod: "QAM256", freq: 465000000, power: -0.4, snr: math.NaN(), corr: 1050, uncorr: 3613},
			{id: "3", lock: "Locked", mod: "QAM256", freq: 471000000, power: -0.6, snr: math.NaN(), corr: 867, uncorr: 3231},
			{id: "4", lock: "Locked", mod: "QAM256", freq: 477000000, power: -0.7, snr: math.NaN(), corr: 885, uncorr: 3630},
			{id: "5", lock: "Locked", mod: "QAM256", freq: 483000000, power: -1, snr: math.NaN(), corr: 962, uncorr: 3633},
			{id: "6", lock: "Locked", mod: "QAM256", freq: 489000000, power: -1, snr: math.NaN(), corr: 1014, uncorr: 3502},
			{id: "7", lock: "Locked", mod: "QAM256", freq: 495000000, power: -1, snr: math.NaN(), corr: 1032, uncorr: 4042},
		},
	}, ans)

}
