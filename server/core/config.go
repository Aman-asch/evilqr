package core

import "os"

const API_TOKEN = "00000000-0000-0000-0000-000000000000"
const QRCODE_ID = "11111111-1111-1111-1111-111111111111"

var BIND_ADDRESS = getBindAddress()

func getBindAddress() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "35000" // Fallback port for local testing
	}
	return "0.0.0.0:" + port
}
