package privet

import (
	"net/http"

	"github.com/grandcat/zeroconf"

	//"errors"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
)

type ApiServer struct {
	MdnsServer *zeroconf.Server
	Port       int
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

//"/privet/info"
func handleInfoRequest(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//w.WriteHeader(http.StatusOK)
	log.Print("Some body requested for info ")
	info := PrinterInfo{"1.0", "Test Xerox", "Near my Cube", "https://www.google.com/cloudprint",
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
func handleCapability(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	log.Print("Capabilities")
	w.Write([]byte(cddEx))
}

//"/privet/printer/submitdoc"
func handleSubmitDoc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Print("Content Type" + r.Header.Get("Content-Type"))
		out, err := os.Create("/tmp/test.pdf")
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

		submitDocResp := SubmitDocResp{"123", 500, "application/pdf", 123123, "My Pdf Document"}
		jdata, err := json.Marshal(submitDocResp)
		if err != nil {
			log.Panic(err)
		}
		w.Write(jdata)
	}
}
func PublishAsGCloudPrinter(printer Printer) (*ApiServer, error) {
	freePort, err := NewPort()
	if err != nil {
		return nil, errors.New("Failed to get a free port")
	}
	log.Print("Sharing on port ", freePort)
	mdnsServer, err := PrivetPublish(printer, freePort)
	if err != nil {
		return nil, err
	}
	http.HandleFunc("/privet/info", handleInfoRequest)
	http.HandleFunc("/privet/capabilities", handleCapability)
	http.HandleFunc("/privet/printer/submitdoc", handleSubmitDoc)

	//http.HandleFunc("/", any)

	http.ListenAndServe(":"+strconv.Itoa(freePort), nil)

	//httpServer := &http.Server{
	//	Addr:    strconv.Itoa(freePort),
	//	Handler: &privetApiServer{}}
	server := ApiServer{mdnsServer, freePort}
	return &server, nil
}
