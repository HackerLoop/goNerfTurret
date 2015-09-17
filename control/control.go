package control

import (
	"fmt"
	"encoding/json"
	"net/http"

	"github.com/vitaminwater/turret/orient"
	"github.com/vitaminwater/turret/shoot"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	log "github.com/Sirupsen/logrus"
)

type Packet struct {
	Type string `json:"type"`
	Payload interface{} `json:"payload"`
}

func Start(orientChan chan orient.Event, shootChan chan shoot.Event) {
	port := 4242

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	http.HandleFunc("/uav", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Connection received")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}

		defer conn.Close()

		for {
			messageType, reader, err := conn.NextReader()
			if err != nil {
				log.Println(err)
				return
			}
			if messageType == websocket.TextMessage {
				packet := Packet{}
				decoder := json.NewDecoder(reader)
				if err := decoder.Decode(&packet); err != nil {
					log.Warning(err)
					continue
				}

				switch packet.Type {
				case "orient":
					event := orient.Event{}
					mapstructure.Decode(packet.Payload, &event)
					orientChan <- event
				case "shoot":
					event := shoot.Event{}
					mapstructure.Decode(packet.Payload, &event)
					shootChan <- event
			default:
				log.Info("Oops unknown packet: ", packet)
				}
			}
		}
	})

	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	log.Infof("Websocket server started on port %d", port)
}
