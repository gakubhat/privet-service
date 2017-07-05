package privet

import (
	"net/http"

	"github.com/grandcat/zeroconf"

	//"errors"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

type ApiServer struct {
	MdnsServer *zeroconf.Server
	HttpServer *http.Server
	Port       int
	logger     *log.Logger
	mux        *http.ServeMux
	printer    Printer
	jobId      int
}

func New(printer Printer) (*ApiServer, error) {
	freePort, err := NewPort()
	if err != nil {
		return nil, errors.New("Failed to get a free port")
	}
	log.Print("Sharing "+printer.Name+" on port ", freePort)
	apiserver := &ApiServer{printer: printer, mux: http.NewServeMux(), Port: freePort}

	mdnsServer, err := PrivetPublish(printer, freePort)
	if err != nil {
		return nil, err
	}
	apiserver.MdnsServer = mdnsServer
	apiserver.mux.HandleFunc("/privet/info", apiserver.handleInfoRequest)
	apiserver.mux.HandleFunc("/privet/capabilities", apiserver.handleCapability)
	apiserver.mux.HandleFunc("/privet/printer/submitdoc", apiserver.handleSubmitDoc)

	//startHttpServer(freePort)
	return apiserver, nil
}

func Publish(printer Printer) (*ApiServer, error) {
	apiServer, err := New(printer)
	if err != nil {
		return nil, err
	}
	apiServer.HttpServer = &http.Server{
		Addr:    ":" + strconv.Itoa(apiServer.Port),
		Handler: apiServer,
	}
	log.Print("Starting Publish ", printer.Name)

	go func() {
		if err := apiServer.HttpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	log.Print("Completed Publish ", printer.Name)
	return apiServer, nil
}

func (apiServer *ApiServer) Shutdown() error {
	if apiServer.MdnsServer != nil {
		apiServer.MdnsServer.Shutdown()
	}
	if apiServer.HttpServer != nil {
		if err := apiServer.HttpServer.Shutdown(context.Background()); err != nil {
			log.Printf("Error: %v\n", err)
			return err
		} else {
			log.Println("Unpublising " + apiServer.printer.Name)
		}
	}
	return nil
}

var (
	cddEx = `{
  "version": "1.0",
  "printer": {
    "supported_content_type": [
      {"content_type": "application/pdf", "min_version": "1.5"},
      {"content_type": "image/jpeg"},
      {"content_type": "text/plain"}
    ],
    "input_tray_unit": [
      {
        "vendor_id": "tray",
        "type": "INPUT_TRAY"
      }
    ],
    "marker": [
      {
        "vendor_id": "black",
        "type": "INK",
        "color": {"type": "BLACK"}
      },
      {
        "vendor_id": "color",
        "type": "INK",
        "color": {"type": "COLOR"}
      }
    ],
    "cover": [
      {
        "vendor_id": "front",
        "type": "CUSTOM",
        "custom_display_name": "front cover"
      }
    ],
    "vendor_capability": [],
    "color": {
      "option": [
        {"type": "STANDARD_MONOCHROME"},
        {"type": "STANDARD_COLOR", "is_default": true},
        {
          "vendor_id": "ultra-color",
          "type": "CUSTOM_COLOR",
          "custom_display_name": "Best Color"
        }
      ]
    },
    "copies": {
      "default": 1,
      "max": 100
    },
    "media_size": {
      "option": [
        {
          "name": "ISO_A4",
          "width_microns": 210000,
          "height_microns": 297000,
          "is_default": true
        },
        {
          "name": "NA_LEGAL",
          "width_microns": 215900,
          "height_microns": 355600
        },
        {
          "name": "NA_LETTER",
          "width_microns": 215900,
          "height_microns": 279400
        }
      ]
    }
  }
}`
)

func (apiServer *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	apiServer.mux.ServeHTTP(w, r)
}

//"/privet/info"
func (apiServer *ApiServer) handleInfoRequest(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//w.WriteHeader(http.StatusOK)
	log.Print("Sending info of ", apiServer.printer.Name)
	info := PrinterInfo{"1.0", apiServer.printer.Name, apiServer.printer.Location, "https://www.google.com/cloudprint",
						[]string{"printer"}, "", "idle", "offline", "Google", "Google Chrome", "1111-22222-33333-4444",
						"24.0.1312.52", 600, "http://support.google.com/cloudprint/answer/1686197/?hl=en",
						"http://support.google.com/cloudprint/?hl=en", "http://support.google.com/cloudprint/?hl=en",

		"AIp06DjQd80yMoGYuGmT_VDAApuBZbInsQ:1358377509659",
		[]string{
			"/privet/accesstoken",
			"/privet/capabilities",
			"/privet/printer/submitdoc",
		},
	}
	jdata, err := json.Marshal(info)
	if err != nil {
		log.Panic(err)
	}
	w.Write(jdata)
}

//"/privet/capabilities"
func (apiServer *ApiServer) handleCapability(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	log.Print("Sending Capabilities ", apiServer.printer.Name)
	w.Write([]byte(cddEx))
}
func (apiServer *ApiServer) getNextJobId() int {
	apiServer.jobId++
	return apiServer.jobId
}

//"/privet/printer/submitdoc"
func (apiServer *ApiServer) handleSubmitDoc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//log.Print("Content Type" + r.Header.Get("Content-Type"))
		requestedIp := net.ParseIP(r.RemoteAddr)
		log.Print("Request from ", requestedIp)
		log.Print("Printing Document for ", apiServer.printer.Name)
		newpath := filepath.Join("/tmp", apiServer.printer.Name+"-jobs/")
		os.MkdirAll(newpath, os.ModePerm)
		jobId := apiServer.getNextJobId()
		jobName := "job" + strconv.Itoa(jobId) + ".pdf"
		jobFilePath := filepath.Join(newpath, jobName)
		out, err := os.Create(jobFilePath)
		if err != nil {
			log.Panic(err)
		}
		defer out.Close()
		_, err = io.Copy(out, r.Body)
		if err != nil {
			log.Panic(err)
		}

		//{
		//	"job_id": "123",
		//	"expires_in": 500,
		//	"job_type": "application/pdf",
		//	"job_size": 123456,
		//	"job_name": "My PDF document"
		//}

		submitDocResp := SubmitDocResp{strconv.Itoa(jobId), 500, "application/pdf", 123123, jobName}
		jdata, err := json.Marshal(submitDocResp)
		if err != nil {
			log.Panic(err)
		}
		log.Print("Printing done.  JobID ", jobId, " Job Name:", jobName)
		w.Write(jdata)
	}
}

//func startHttpServer(freePort int) *http.Server {
//	srv := &http.Server{Addr: ":" + strconv.Itoa(freePort)}
//
//	go func() {
//		if err := srv.ListenAndServe(); err != nil {
//			// cannot panic, because this probably is an intentional close
//			log.Printf("Httpserver: ListenAndServe() error: %s", err)
//		}
//	}()
//
//	// returning reference so caller can call Shutdown()
//	return srv
//}
