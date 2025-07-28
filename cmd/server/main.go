package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/server"
	. "github.com/Sister20/if3230-tubes-dark-syster/lib/util"
)

func main() {
	var (
		host        = flag.String("host", "127.0.0.1", "Server host address")
		port        = flag.String("port", "8080", "Server port")
		contactHost = flag.String("contact-host", "", "Contact node host address (optional)")
		contactPort = flag.String("contact-port", "", "Contact node port (optional)")
		help        = flag.Bool("help", false, "Show help message")
	)
	
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "nanorediz - A distributed Redis-like key-value store\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  NANOREDIZ_HOST              Server host (default: 0.0.0.0)\n")
		fmt.Fprintf(os.Stderr, "  NANOREDIZ_PORT              Server port (default: 8080)\n")
		fmt.Fprintf(os.Stderr, "  NANOREDIZ_LOG_LEVEL         Log level (debug, info, warn, error)\n")
		fmt.Fprintf(os.Stderr, "  NANOREDIZ_GRPC_TIMEOUT      gRPC timeout (default: 30s)\n")
		fmt.Fprintf(os.Stderr, "  NANOREDIZ_SHUTDOWN_TIMEOUT  Graceful shutdown timeout (default: 10s)\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -host 127.0.0.1 -port 8080\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -host 127.0.0.1 -port 8081 -contact-host 127.0.0.1 -contact-port 8080\n", os.Args[0])
	}
	
	flag.Parse()
	
	if *help {
		flag.Usage()
		return
	}

	// Validate required arguments
	if *host == "" || *port == "" {
		fmt.Fprintf(os.Stderr, "Error: host and port are required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	var isContact bool = false
	var address, contactAddress Address

	address = *NewAddress(*host, *port)

	// Check if contact node is specified
	if *contactHost != "" && *contactPort != "" {
		contactAddress = *NewAddress(*contactHost, *contactPort)
		isContact = true
		fmt.Printf("Starting node with contact address: %s:%s\n", *contactHost, *contactPort)
	} else if *contactHost != "" || *contactPort != "" {
		fmt.Fprintf(os.Stderr, "Error: both contact-host and contact-port must be specified together\n")
		os.Exit(1)
	}

	fmt.Printf("Starting nanorediz server on %s:%s\n", *host, *port)
	
	serverInstance := server.NewServer(&address, isContact, &contactAddress)
	serverInstance.Serve()
}