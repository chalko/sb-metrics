package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"testing"
)


func TestScrape(t *testing.T) {
	reader ,err:= os.Open(filepath.Join("testdata","index.html"))
	if err != nil {
		log.Fatal(err)
	}
	ans := scrape(reader)
	assert.Equal(t, "OK", ans.startup.bootState )

}

