package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gakubhat/privet-service/privet"
)

const waitTime = 100

func main() {
	printer := privet.Printer{"Test Xerox", "Near my Cube"}
	server, err := privet.PublishAsGCloudPrinter(printer)

	if err != nil {
		log.Panic("Something went wrong ,", err)
	}
	defer server.MdnsServer.Shutdown()
	//defer server.HttpServer.Close()

	// Clean exit.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	// Timeout timer.
	var tc <-chan time.Time
	if waitTime > 0 {
		tc = time.After(time.Second * time.Duration(waitTime))
	}

	select {
	case <-sig:
		// Exit by user
	case <-tc:
		// Exit by timeout
	}

	log.Println("Shutting down.")
}
