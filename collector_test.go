package main

import (
	"github.com/stretchr/testify/assert"
	"log"
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
	assert.Equal(t,
		StartupStatus{bootState: "OK"},
		ans.startup)
	assert.Equal(t,
		[]DownStatus{
			{id: "16", lock: "Locked", mod: "QAM256", freq: 549000000, power: -4.6, snr: 38.8, corr: 204, uncorr: 443},
			{id: "1", lock: "Locked", mod: "QAM256", freq: 459000000, power: -0.2, snr: 42.1, corr: 984, uncorr: 3588},
			{id: "2", lock: "Locked", mod: "QAM256", freq: 465000000, power: -0.4, snr: 41.9, corr: 1050, uncorr: 3613},
			{id: "3", lock: "Locked", mod: "QAM256", freq: 471000000, power: -0.6, snr: 41.8, corr: 867, uncorr: 3231},
			{id: "4", lock: "Locked", mod: "QAM256", freq: 477000000, power: -0.7, snr: 41.8, corr: 885, uncorr: 3630},
			{id: "5", lock: "Locked", mod: "QAM256", freq: 483000000, power: -1, snr: 41.2, corr: 962, uncorr: 3633},
			{id: "6", lock: "Locked", mod: "QAM256", freq: 489000000, power: -1, snr: 41.4, corr: 1014, uncorr: 3502},
			{id: "7", lock: "Locked", mod: "QAM256", freq: 495000000, power: -1, snr: 41.4, corr: 1032, uncorr: 4042},
			{id: "8", lock: "Locked", mod: "QAM256", freq: 501000000, power: -1.3, snr: 41.2, corr: 1051, uncorr: 4068},
			{id: "9", lock: "Locked", mod: "QAM256", freq: 507000000, power: -1.9, snr: 40.9, corr: 1076, uncorr: 3985},
			{id: "10", lock: "Locked", mod: "QAM256", freq: 513000000, power: -2.2, snr: 40.6, corr: 1079, uncorr: 4247},
			{id: "11", lock: "Locked", mod: "QAM256", freq: 519000000, power: -2.4, snr: 40.5, corr: 991, uncorr: 3836},
			{id: "12", lock: "Locked", mod: "QAM256", freq: 525000000, power: -2.7, snr: 40.3, corr: 955, uncorr: 3859},
			{id: "13", lock: "Locked", mod: "QAM256", freq: 531000000, power: -3.1, snr: 40, corr: 980, uncorr: 3690},
			{id: "14", lock: "Locked", mod: "QAM256", freq: 537000000, power: -3.5, snr: 39.6, corr: 894, uncorr: 3257},
			{id: "15", lock: "Locked", mod: "QAM256", freq: 543000000, power: -4.1, snr: 39.3, corr: 1096, uncorr: 4190},
			{id: "17", lock: "Locked", mod: "QAM256", freq: 555000000, power: -5.2, snr: 38.2, corr: 1013, uncorr: 3879},
		}, ans.ds)

	assert.Equal(t,
		[]UpStatus{
			{num: "1", id: "57", lock: "Locked", chtype: "SC-QAM", width: 26000000, freq: 6400000, power: 50},
			{num: "2", id: "58", lock: "Locked", chtype: "SC-QAM", width: 19200000, freq: 6400000, power: 48},
		},
		ans.us)
}
