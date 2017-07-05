package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"errors"

	"github.com/gakubhat/privet-service/privet"
)

const waitTime = 30

var data = make(map[string]*privet.ApiServer)

func main() {
	printer1 := privet.Printer{"Test Xerox", "Near my Cube"}
	printer2 := privet.Printer{"Ricoh1", "Far away from my Cube"}

	AddPrinter(printer1)
	AddPrinter(printer2)
	defer func() {
		for k, v := range data {
			if err := v.Shutdown(); err != nil {
				log.Print("Failed to shutdown " + k)
			} else {
				log.Print("Shut down of server for " + k + " completed")
			}
		}
	}()
	//defer server.HttpServer.Close()

	// Clean exit.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	// Timeout timer.
	var tc <-chan time.Time
	if waitTime > 0 {
		tc = time.After(time.Second * time.Duration(waitTime))
	}

	<-tc
	if err := RemovePrinter("Test Xerox"); err != nil {
		log.Panic("Failed to shutdown printer Test Xerox")
	}
	select {
	case <-sig:
		// Exit by user
		//case <-tc:
		// Exit by timeout
	}

	log.Println("Shutting down privet print service.")
}

func RemovePrinter(name string) error {
	if server, ok := data[name]; ok {
		if err := server.Shutdown(); err != nil {
			return err
		}
		delete(data, name)
	} else {
		return errors.New("Printer not found")
	}
	return nil
}
func AddPrinter(printer privet.Printer) error {
	server, err := privet.Publish(printer)
	data[printer.Name] = server
	if err != nil {
		log.Panic("Something went wrong ,"+printer.Name, err)
	}
	return nil
}
