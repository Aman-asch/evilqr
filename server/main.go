package main

import (
	"flag"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"evilqr-server/core"
	"evilqr-server/log"
)

var wwwdir = flag.String("d", "", "www content directory path")

func main() {
	flag.Parse()

	if *wwwdir == "" {
		log.Error("you need to set up www content directory path with '-d' parameter")
		return
	}

	// 1. Get the port Render assigned to us
	renderPort := os.Getenv("PORT")
	if renderPort == "" {
		renderPort = "10000"
	}

	// 2. Force the core server to run on a hidden internal port
	core.BIND_ADDRESS = "127.0.0.1:35000"

	// 3. Start the core server in the background
	server, err := core.NewHttpServer()
	if err != nil {
		log.Error("http: %s", err)
		return
	}
	go server.Run(*wwwdir)

	// Wait 1 second for the internal server to spin up
	time.Sleep(1 * time.Second)

	// 4. Set up a reverse proxy to the internal server
	target, _ := url.Parse("http://127.0.0.1:35000")
	proxy := httputil.NewSingleHostReverseProxy(target)

	// 🔥 THE FIX: Intercept the response and FORCE exactly one CORS header
	proxy.ModifyResponse = func(r *http.Response) error {
		r.Header.Set("Access-Control-Allow-Origin", "*")
		r.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		r.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		return nil
	}

	// 5. Create a custom mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Handle the preflight OPTIONS request directly (don't pass to backend)
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Forward the request to the core server
		proxy.ServeHTTP(w, r)
	})

	// 6. Start the proxy on the public Render port
	log.Info("Proxy listening on 0.0.0.0:%s", renderPort)
	if err := http.ListenAndServe("0.0.0.0:"+renderPort, mux); err != nil {
		log.Error("Server failed: %s", err)
	}
}
