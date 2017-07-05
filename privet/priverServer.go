package privet

import (
	"log"

	"github.com/grandcat/zeroconf"
)

func PrivetPublish(printer Printer, port int) (*zeroconf.Server, error) {
	var (
		service = "_privet._tcp"
		//subType    = "_printer._sub._privet._tcp"
		domain     = "local"
		txtRecords = []string{
			"txtvers=1",
			"ty=" + printer.Name,
			"note=" + printer.Location,
			"url=https://www.google.com/cloudprint",
			"type=printer",
			"id=",
			"cs=offline"}
	)
	log.Print("Publishing printer ", printer.Name)
	return zeroconf.Register(printer.Name+" on iPrint", service, domain, port, txtRecords, nil)

}
