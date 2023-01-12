package main

import (
	"flag"

	"github.com/code-to-go/safepool.lib/api"
	"github.com/code-to-go/safepool.lib/sql"
	"github.com/sirupsen/logrus"
)

func parseFlags() {
	var verbose int
	var dbname string

	flag.IntVar(&verbose, "v", 0, "verbose level - 0 to 2")
	flag.StringVar(&dbname, "d", "", "location of the SQLlite DB")
	flag.Parse()

	switch verbose {
	case 0:
		logrus.SetLevel(logrus.FatalLevel)
	case 1:
		logrus.SetLevel(logrus.InfoLevel)
	case 2:
		logrus.SetLevel(logrus.DebugLevel)
	}

	if dbname != "" {
		sql.DbName = dbname
	}
}

func main() {

	parseFlags()
	api.Start()
	SelectMain()
}
