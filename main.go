package main

import (
	"context"
	"log"
	"os"

	"github.com/haung921209/nhn-cloud-mcp/config"
	"github.com/haung921209/nhn-cloud-mcp/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	serverName    = "nhn-cloud-mcp"
	serverVersion = "0.1.0"
)

func main() {
	// Load configuration from environment variables
	cfg := config.LoadFromEnv()

	// Validate configuration
	if !cfg.HasRDSCredentials() {
		log.Println("Warning: RDS credentials not configured. Set NHN_CLOUD_ACCESS_KEY_ID and NHN_CLOUD_SECRET_ACCESS_KEY.")
	}

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    serverName,
		Version: serverVersion,
	}, nil)

	// Register tools
	if cfg.HasRDSCredentials() {
		if cfg.MySQLAppKey != "" {
			tools.RegisterMySQLTools(server, cfg)
			log.Println("Registered MySQL tools")
		}
		// TODO: Add MariaDB tools
		// TODO: Add PostgreSQL tools
	}

	if cfg.HasComputeCredentials() {
		// TODO: Add Compute tools
		// TODO: Add Network tools
		log.Println("Compute credentials configured (tools not yet implemented)")
	}

	// Run server over stdio
	log.Printf("Starting %s v%s...\n", serverName, serverVersion)
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}
