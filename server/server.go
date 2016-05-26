package server

// Main func
func Main() {
	go ListenHTTP(":8080")
	go ListenRPC(":8081")
	select {}
}

// Exit func
func Exit() {
}
