package main

import (
	"flag"
	"fmt"
	"github.com/ziutek/syslog"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

var hostname string
var tagExcludeFilter string

type handler struct {
	*syslog.BaseHandler
}

func newHandler() *handler {
	h := handler{syslog.NewBaseHandler(5, func(m *syslog.Message) bool {
		match, _ := regexp.MatchString(tagExcludeFilter, m.Tag)
		if match {
			return false
		}
		return true
	}, false)}
	go h.mainLoop()
	return &h
}

func (h *handler) mainLoop() {
	for {
		m := h.Get()
		if m == nil {
			break
		}

		fmt.Println(m)
	}
	h.End()
}

func main() {
	hostname, _ = os.Hostname()
	tagExcludeFilter = os.Getenv("TAG_EXCLUDE_FILTER")
	fmt.Println(hostname)
	flag.Parse()
	s := syslog.NewServer()
	s.AddHandler(newHandler())
	s.Listen("0.0.0.0:1514")

	// Wait for terminating signal
	sc := make(chan os.Signal, 2)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT)
	<-sc

	fmt.Println("Shutdown the server...")
	s.Shutdown()
}
