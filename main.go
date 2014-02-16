package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

//
//   flk
//
//   a simple snowflake server
//
//   serves guaranteed unique ids
//
//

var (
	id_chan chan []byte

	bufflen  = flag.Int("bufflen", 1000, "how many ids to buffer internally")
	maxprocs = flag.Int("maxprocs", runtime.NumCPU(), "GOMAXPROCS setting")
	host     = flag.String("host", ":8080", "host for http requests")
)

func main() {

	flag.Parse()

	runtime.GOMAXPROCS(*maxprocs)

	id_chan = make(chan []byte, *bufflen)

	go func() {

		inc := time.Now().UnixNano()

		hostname, err := os.Hostname()
		if err != nil {
			log.Fatalf("Unable to get hostname: %s\n", err)
		}

		for {
			id_chan <- []byte(fmt.Sprintf("%s:%d", hostname, inc))
			inc++
		}

	}()

	http.HandleFunc("/id", func(w http.ResponseWriter, _ *http.Request) {
		w.Write(<-id_chan)
	})

	log.Fatal(http.ListenAndServe(*host, nil))

}
