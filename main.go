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
	cfg := config.Load()

	server := mcp.NewServer(&mcp.Implementation{
		Name:    serverName,
		Version: serverVersion,
	}, nil)

	tools.RegisterAuthTools(server, cfg)
	log.Println("Registered auth tools (nhn_set_credential, nhn_get_credential_status)")

	if cfg.MySQLAppKey != "" {
		tools.RegisterMySQLTools(server, cfg)
		log.Println("Registered MySQL tools")
	}

	logCredentialStatus(cfg)

	log.Printf("Starting %s v%s...\n", serverName, serverVersion)
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}

func logCredentialStatus(cfg *config.Config) {
	if cfg.HasRDSCredentials() {
		log.Printf("RDS credentials: configured (source: %s)", cfg.GetSource("AccessKeyID"))
	} else {
		log.Println("RDS credentials: not configured - use nhn_set_credential tool or configure ~/.nhncloud/credentials")
	}

	if cfg.HasComputeCredentials() {
		log.Printf("Compute credentials: configured (source: %s)", cfg.GetSource("Username"))
	}
}
