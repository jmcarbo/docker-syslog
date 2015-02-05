package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ziutek/syslog"
)

var hostname string

type handler struct {
	*syslog.BaseHandler
}

func newHandler() *handler {
	h := handler{syslog.NewBaseHandler(5, func(m *syslog.Message) bool {
		if m.Hostname != hostname {
			return true
		}
		return false
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
