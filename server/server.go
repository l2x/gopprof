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

// Init at first
func Init(cfg string) error {
	if err := initConfig(cfg); err != nil {
		return err
	}
	if err := initStoreSaver(); err != nil {
		return err
	}
	return nil
}

// Close at last
func Close() {
	if storeSaver != nil {
		storeSaver.Close()
	}
}
