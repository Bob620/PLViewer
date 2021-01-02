package backend

import (
	"encoding/json"
	"fmt"
	"github.com/bob620/baka-rpc-go/parameters"
	"github.com/bob620/baka-rpc-go/rpc"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os/exec"
	"strconv"
)

type Backend struct {
	Client   *rpc.BakaRpc
	upgrader *websocket.Upgrader
	port     int
	await    chan bool
}

func startBackendConnection() *Backend {
	backend := &Backend{
		upgrader: &websocket.Upgrader{},
		Client:   rpc.CreateBakaRpc(nil, nil),
		port:     7789,
		await:    make(chan bool),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c, err := backend.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()
		go func() {
			backend.await <- true
		}()
		backend.Client.UseChannels(rpc.MakeSocketReaderChan(c), rpc.MakeSocketWriterChan(c))
		backend.await <- false
	})

	go func() {
		http.ListenAndServe(fmt.Sprintf("localhost:%d", backend.port), nil)
	}()
	initializeBackend(backend)

	return backend
}

func startNodeClient(port int) error {
	backend := exec.Command("./node.exe", "./index.mjs", strconv.Itoa(port))
	backend.Dir = "./node"

	return backend.Run()
}

func StartBackend() *Backend {
	backend := startBackendConnection()

	go func() {
		err := startNodeClient(backend.port)
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-backend.await
	return backend
}

func initializeBackend(backend *Backend) {
	backend.Client.RegisterMethod("log", []parameters.Param{
		&parameters.StringParam{
			Name:     "text",
			Default:  "",
			Required: false,
		},
	}, func(params map[string]parameters.Param) (returnMessage json.RawMessage, err error) {
		text, _ := params["text"].(*parameters.StringParam).GetString()
		log.Println(text)
		return
	})
}
