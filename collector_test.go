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
	assert.Equal(t, ModemStatus{
		startup: StartupStatus{bootState: "OK"},
		ds: []DownStatus{
			{id: "16", lock: "", mod: "", freg: 0, power: 0, snr: 0, corr: 0, uncorr: 0},
			{id: "1", lock: "", mod: "", freg: 0, power: 0, snr: 0, corr: 0, uncorr: 0},
			{id: "2", lock: "", mod: "", freg: 0, power: 0, snr: 0, corr: 0, uncorr: 0},
			{id: "3", lock: "", mod: "", freg: 0, power: 0, snr: 0, corr: 0, uncorr: 0},
			{id: "4", lock: "", mod: "", freg: 0, power: 0, snr: 0, corr: 0, uncorr: 0},
			{id: "5", lock: "", mod: "", freg: 0, power: 0, snr: 0, corr: 0, uncorr: 0},
			{id: "6", lock: "", mod: "", freg: 0, power: 0, snr: 0, corr: 0, uncorr: 0},
			{id: "7", lock: "", mod: "", freg: 0, power: 0, snr: 0, corr: 0, uncorr: 0},
			{id: "8", lock: "", mod: "", freg: 0, power: 0, snr: 0, corr: 0, uncorr: 0},
			{id: "9", lock: "", mod: "", freg: 0, power: 0, snr: 0, corr: 0, uncorr: 0},
			{id: "10", lock: "", mod: "", freg: 0, power: 0, snr: 0, corr: 0, uncorr: 0},
			{id: "11", lock: "", mod: "", freg: 0, power: 0, snr: 0, corr: 0, uncorr: 0},
			{id: "12", lock: "", mod: "", freg: 0, power: 0, snr: 0, corr: 0, uncorr: 0},
			{id: "13", lock: "", mod: "", freg: 0, power: 0, snr: 0, corr: 0, uncorr: 0},
			{id: "14", lock: "", mod: "", freg: 0, power: 0, snr: 0, corr: 0, uncorr: 0},
			{id: "15", lock: "", mod: "", freg: 0, power: 0, snr: 0, corr: 0, uncorr: 0},
			{id: "17", lock: "", mod: "", freg: 0, power: 0, snr: 0, corr: 0, uncorr: 0},
		},
	}, ans)

}
