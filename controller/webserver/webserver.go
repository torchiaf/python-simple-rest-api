// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second
)

var (
	addr      = flag.String("addr", ":8080", "http service address")
	homeTempl = template.Must(template.New("").Parse(homeHTML))
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	wsVar *websocket.Conn
)

// func readFileIfModified(lastMod time.Time) ([]byte, time.Time, error) {
// 	fi, err := os.Stat(filename)
// 	if err != nil {
// 		return nil, lastMod, err
// 	}
// 	if !fi.ModTime().After(lastMod) {
// 		return nil, lastMod, nil
// 	}
// 	p, err := os.ReadFile(filename)
// 	if err != nil {
// 		return nil, fi.ModTime(), err
// 	}
// 	return p, fi.ModTime(), nil
// }

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

// func writer(ws *websocket.Conn) {
// 	lastError := ""
// 	pingTicker := time.NewTicker(pingPeriod)
// 	fileTicker := time.NewTicker(filePeriod)
// 	defer func() {
// 		pingTicker.Stop()
// 		fileTicker.Stop()
// 		ws.Close()
// 	}()
// 	for {
// 		select {
// 		case <-fileTicker.C:
// 			var p []byte
// 			var err error

// 			log.Println("Ticker")

// 			p, lastMod, err = readFileIfModified(lastMod)

// 			if err != nil {
// 				if s := err.Error(); s != lastError {
// 					lastError = s
// 					p = []byte(lastError)
// 				}
// 			} else {
// 				lastError = ""
// 			}

// 			if p != nil {
// 				ws.SetWriteDeadline(time.Now().Add(writeWait))
// 				if err := ws.WriteMessage(websocket.TextMessage, p); err != nil {
// 					return
// 				}
// 			}
// 		case <-pingTicker.C:
// 			ws.SetWriteDeadline(time.Now().Add(writeWait))
// 			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
// 				return
// 			}
// 		}
// 	}
// }

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	wsVar = ws
	// reader(ws)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var v = struct {
		Host string
		Data string
	}{
		r.Host,
		string(""),
	}
	homeTempl.Execute(w, &v)
}

func SendMessage(msg []byte) {
	if wsVar == nil {
		return
	}
	wsVar.SetWriteDeadline(time.Now().Add(writeWait))
	if err := wsVar.WriteMessage(websocket.TextMessage, msg); err != nil {
		return
	}
}

func InitWebServer() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

const homeHTML = `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>WebSocket Example</title>
    </head>
    <body>
        <pre id="fileData">{{.Data}}</pre>
        <script type="text/javascript">
            (function() {
                var data = document.getElementById("fileData");
                var conn = new WebSocket("ws://{{.Host}}/ws");
                conn.onclose = function(evt) {
                    data.textContent = 'Connection closed';
                }
                conn.onmessage = function(evt) {
                    console.log('file updated');
                    data.textContent = evt.data;
                }
            })();
        </script>
    </body>
</html>
`
