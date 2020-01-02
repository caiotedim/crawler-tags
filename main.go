package main

import (
	"flag"
	"os"

	"github.com/caiotedim/crawler-tags/webapp"
	"github.com/golang/glog"
)

var (
	port *int
	bind *string
)

const version = "1.0.0-alpha"

func main() {
	glog.Infof("Starting CRAWLER_TAGS on version %s: bind:[%s] port:[%d]", version, *bind, *port)
	webapp.Server(bind, port, version)
}

func usage() {
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	bind = flag.String("bind", "0.0.0.0", "bind address")
	port = flag.Int("port", 8080, "port")
	flag.Usage = usage
	flag.Set("logtostderr", "true")
	flag.Parse()
}