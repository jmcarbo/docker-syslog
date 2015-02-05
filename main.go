package main

import (
	"flag"
	"fmt"
	"github.com/jmcarbo/syslog"
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

		match, _ := regexp.MatchString(tagExcludeFilter, m.Content)
		if match {
			//	fmt.Println(m)
		} else {
			fmt.Printf("%v\n", m)
			fmt.Printf("%#v\n", m)

		}
	}
	h.End()
}

func main() {
	hostname, _ = os.Hostname()
	tagExcludeFilter = os.Getenv("TAG_EXCLUDE_FILTER")
	fmt.Println(tagExcludeFilter)
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
