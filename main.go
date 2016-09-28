/*
Author: Ben Polzin
Github Author: https://github.com/benpolzin
Github Repo: https://github.com/benpolzin/gobahnhof

Version: 0.1
*/

package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
)

type Recipients struct {
	Id   int64  `json:"id,string"` // 19063007
	Type string `json:"type"`      // RECIPIENT_GROUP
}

type RecipientMsg struct {
	Recipients []Recipients `json:"recipients"`
}

type codecJSONevent struct {
	Event struct {
		IDentification struct {
			IPAddress struct {
				Value net.IP `json:"Value"`
			} `json:"IPAddress"`
			MACAddress struct {
				Value string `json:"Value"`
			} `json:"MACAddress"`
			ProductID struct {
				Value string `json:"Value"`
			} `json:"ProductID"`
			ProductType struct {
				Value string `json:"Value"`
			} `json:"ProductType"`
			SWVersion struct {
				Value string `json:"Value"`
			} `json:"SWVersion"`
			SerialNumber struct {
				Value string `json:"Value"`
			} `json:"SerialNumber"`
			SystemName struct {
				Value string `json:"Value"`
			} `json:"SystemName"`
		} `json:"Identification"`
		UserInterface struct {
			ID         int64 `json:"id,string"`
			Extensions struct {
				ID     int64 `json:"id,string"`
				Widget struct {
					ID     int64 `json:"id,string"`
					Action struct {
						ID   int64 `json:"id,string"`
						Type struct {
							ID    int64  `json:"id,string"`
							Value string `json:"Value"`
						} `json:"Type"`
						Value struct {
							ID    int64  `json:"id,string"`
							Value string `json:"Value"`
						} `json:"Value"`
						WidgetID struct {
							ID    int64  `json:"id,string"`
							Value string `json:"Value"`
						} `json:"WidgetId"`
					} `json:"Action"`
				} `json:"Widget"`
			} `json:"Extensions"`
		} `json:"UserInterface"`
	} `json:"Event"`
}

func registerCodec() {
	// set up HTTP client with TLS transport
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	// Register codec via XML request for all Event feedback
	registerXML := `<Command><HttpFeedback><Register command="True" role="Admin" read="Admin"><FeedbackSlot>1</FeedbackSlot><Format>JSON</Format><ServerUrl>http://10.27.1.127:8080/codecFeedback</ServerUrl><Expression item="1">/Event/UserInterface/Extensions/Widget/Action</Expression></Register></HttpFeedback></Command>`
	req, err := http.NewRequest("POST", "https://10.27.2.151/putxml", bytes.NewBuffer([]byte(registerXML)))
	req.SetBasicAuth("****", "****")
	req.Header.Set("Content-Type", "text/xml")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Println("Register Req Status : ", resp.Status)
	fmt.Println("Register Req Status : ", resp.Body)
}

func sendIcastMsg() {
	// Include root CA of InformaCast server to allow HTTP client to trust server
	rootPEM := `-----BEGIN CERTIFICATE-----
MIIENTCCAx2gAwIBAgIJAK8nU+H0YZ1xMA0GCSqGSIb3DQEBBQUAMIGwMQswCQYD
VQQGEwJVUzESMBAGA1UECAwJTWlubmVzb3RhMQ0wCwYDVQQHDARIdWdvMRQwEgYD
VQQKDAtUaGUgUG9semluczEaMBgGA1UECwwRQ29sbGFib3JhdGlvbiBMYWIxITAf
BgNVBAMMGG5peC5wcmltZS50aGVwb2x6aW5zLm5ldDEpMCcGCSqGSIb3DQEJARYa
YWRtaW5AcHJpbWUudGhlcG9semlucy5uZXQwHhcNMTMxMTAzMDUyOTExWhcNMjMx
MTAxMDUyOTExWjCBsDELMAkGA1UEBhMCVVMxEjAQBgNVBAgMCU1pbm5lc290YTEN
MAsGA1UEBwwESHVnbzEUMBIGA1UECgwLVGhlIFBvbHppbnMxGjAYBgNVBAsMEUNv
bGxhYm9yYXRpb24gTGFiMSEwHwYDVQQDDBhuaXgucHJpbWUudGhlcG9semlucy5u
ZXQxKTAnBgkqhkiG9w0BCQEWGmFkbWluQHByaW1lLnRoZXBvbHppbnMubmV0MIIB
IjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAy91MCsSxwQExjzgnHxJmhiql
i20iUooJ71Vxt5Cttd/9MAW2aQu2B7JhkdmI//z5WA3jCs8yK/5ge6SI0Za/ZAdJ
uzI3XdEWcijUXlPJGOF4A4NBdmx4k5tyeA8ScRF8yomzLIa2Gk7mZdYsIsJujKYY
Ys1gNM/SblnCGXyiDQLNkoc4yeHhtyuSBZS54Na0j1JtHwA0wqEmJbO9kEGmbYA4
Vuk//2RrB1GFvBqs67ggkZtNatkR46aWcSoY9nqetyeHWeP3+SyiM2wONSjfbvcK
A4Hj+ryDppBBrL49yt6nZdCOsJdKgHDurMzNQk/dqd7vwPj8Rpit2OfXwFy4vQID
AQABo1AwTjAdBgNVHQ4EFgQUdocEhq/4SsAU/Z6myYxYVR43zbwwHwYDVR0jBBgw
FoAUdocEhq/4SsAU/Z6myYxYVR43zbwwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0B
AQUFAAOCAQEAAqVblaWVphvRrYvr2dvQs86AWBZ3tPSpxbeKBiOtjDDsly7C0Yus
sPgIYADFwxUSeOSb9GnLv+4dTegxJjgeGCL5fEgKlFQecp0rI/68e6Ok0X93MzbG
k8dS0rMeYgkh4HvczrZUW+WnxJEAyW12c6T9gYAxPxJlx84QfDA2gVbpu8NQyDFW
kQEeStydKfY0Mh5BwHNJnbmygrexjv3LcSJngL8yifZfwwYMyk4n2saRoMq781gE
LAiRXkzjvQqEIrSPqHkCUsJwPNB0YRU2C4TOI5bAqr23dhcYVv9FlRFtXp/dP8xi
GYy2sJ2of3XzBHdxOoPlCfdEg1sl1KVQ1Q==
-----END CERTIFICATE-----`

	// Add root cert to cert pool
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		panic("Failed to parse root certificate.")
	}

	// send HTTP API Call to InformaCast

	// set up HTTP client with TLS transport
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{RootCAs: roots},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	// Create JSON body for API call
	jsonSendMsg := RecipientMsg{Recipients: []Recipients{Recipients{Id: 19063007, Type: "RECIPIENT_GROUP"}}}
	fmt.Println(jsonSendMsg.Recipients)
	//	jsonMsg, _ := json.Marshal(jsonSendMsg)
	//	fmt.Println(string(jsonMsg))
	jsonb := new(bytes.Buffer)
	json.NewEncoder(jsonb).Encode(jsonSendMsg)

	// create new HTTP request for POST message with JSON body and auth header
	req, err := http.NewRequest("POST", "https://singlewire.prime.thepolzins.net:8444/InformaCast/RESTServices/V1/Messages/867", jsonb)
	req.SetBasicAuth("****", "****")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Println("Status : ", resp.Status)
}

func main() {
	// Register to receive feedback with Cisco codec
	registerCodec()

	// start web server to listen for feedback events
	http.HandleFunc("/codecFeedback", func(w http.ResponseWriter, r *http.Request) {
		u := new(codecJSONevent)
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		fmt.Println("Codec:", u.Event.IDentification.MACAddress, "Widget ID:", u.Event.UserInterface.Extensions.Widget.Action.WidgetID.Value)
		fmt.Println("Codec:", u.Event.IDentification.MACAddress, "Widget Event Type:", u.Event.UserInterface.Extensions.Widget.Action.Type.Value)
		if u.Event.UserInterface.Extensions.Widget.Action.WidgetID.Value == "humoctopus" && u.Event.UserInterface.Extensions.Widget.Action.Type.Value == "pressed" {
			fmt.Println("Send Notification!")
			sendIcastMsg()
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

/*
// web server
func main() {
  http.ListenAndServe(":8080", nil)
}
*/
