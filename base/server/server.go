// File: `Server Implementation File` source/server/server.go
package PhoeniciaDigitalServer

import (
	PhoeniciaDigitalUtils "Phoenicia-Digital-Base-API/base/utils"
	PhoeniciaDigitalConfig "Phoenicia-Digital-Base-API/config"
	"Phoenicia-Digital-Base-API/source"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Initialize Server Ecosystem Variables
var multiplexer *http.ServeMux = http.NewServeMux()

var PhoeniciaDigitalServer *http.Server = &http.Server{
	Addr:    PhoeniciaDigitalConfig.Config.Port,
	Handler: multiplexer,
}

func StartServer() {
	if PhoeniciaDigitalServer.Addr != ":" {
		if portNumber, err := strconv.Atoi(PhoeniciaDigitalServer.Addr[1:]); err != nil {
			log.Printf("Given PORT is Invalid: %s != int | Change in ./config/.env", PhoeniciaDigitalServer.Addr[1:])
			PhoeniciaDigitalUtils.Log(fmt.Sprintf("Given PORT is Invalid: %s != int | Change in ./config/.env", PhoeniciaDigitalServer.Addr[1:]))
		} else {
			if portNumber >= 0 && portNumber <= 65535 {
				log.Printf("Server Running on http://localhost%s", PhoeniciaDigitalServer.Addr)
				PhoeniciaDigitalUtils.Log(fmt.Sprintf("Server started on PORT --> %s", PhoeniciaDigitalServer.Addr))
				log.Fatal(PhoeniciaDigitalServer.ListenAndServe())
			} else {
				log.Printf("Given PORT: %s is OUT OF RANGE 0 --> 65535 | Change in ./config/.env", PhoeniciaDigitalServer.Addr[1:])
				PhoeniciaDigitalUtils.Log(fmt.Sprintf("Given PORT: %s is OUT OF RANGE 0 --> 65535 | Change in ./config/.env", PhoeniciaDigitalServer.Addr[1:]))
			}
		}
	} else {
		log.Printf("Given PORT is empty | Change in ./config/.env")
		PhoeniciaDigitalUtils.Log("Given PORT is empty | Change in ~/config/.env")
	}
}

// Initialize Server Logic
func init() {
	multiplexer.HandleFunc("/sensor", source.HandleMeasureDistance)

	multiplexer.HandleFunc("OPTIONS /loiter", func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for all requests (can be more specific if needed)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3001") // Allow requests from any origin (http://localhost:3000 in your case)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	})
	multiplexer.Handle("GET /loiter", PhoeniciaDigitalUtils.PhoeniciaDigitalHandler(source.HandleLoiter))

	multiplexer.HandleFunc("OPTIONS /rotate-right", func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for all requests (can be more specific if needed)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3001") // Allow requests from any origin (http://localhost:3000 in your case)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	})
	multiplexer.Handle("GET /rotate-right", PhoeniciaDigitalUtils.PhoeniciaDigitalHandler(source.HandleRotateRight))

	multiplexer.HandleFunc("OPTIONS /rotate-left", func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for all requests (can be more specific if needed)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3001") // Allow requests from any origin (http://localhost:3000 in your case)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	})
	multiplexer.Handle("GET /rotate-left", PhoeniciaDigitalUtils.PhoeniciaDigitalHandler(source.HandleRotateLeft))

	// multiplexer.HandleFunc("OPTIONS /sensor", func(w http.ResponseWriter, r *http.Request) {
	// 	// Set CORS headers for all requests (can be more specific if needed)
	// 	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow requests from any origin (http://localhost:3000 in your case)
	// 	w.Header().Set("Access-Control-Allow-Credentials", "true")
	// 	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	// })

	// multiplexer.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintln(w, "Hello, world!")
	// })
}
