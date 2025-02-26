package main

import (
	"sync"

	"github.com/joaooliveirapro/xmlsyncgo/src/initializers"
	"github.com/joaooliveirapro/xmlsyncgo/src/parser"
	"github.com/joaooliveirapro/xmlsyncgo/src/server"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnecToDB()
}

func main() {
	// Run parser on a different go rountine
	// to avoid blocking the web server
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done() // Signal go routine is done after main is finished
		parser.Main()   // Parser main entry point
	}()
	wg.Wait() // Wait for go routine to finish

	// Run webserver
	webserver := server.NewTransport()
	webserver.ServerHTTP()
	select {}
}
