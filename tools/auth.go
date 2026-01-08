package tools

import (
	"context"
	"fmt"

	"github.com/haung921209/nhn-cloud-mcp/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SetCredentialInput struct {
	Key   string `json:"key" jsonschema_description:"Credential key: access_key_id, secret_access_key, mysql_appkey, mariadb_appkey, postgresql_appkey, username, password, tenant_id, region"`
	Value string `json:"value" jsonschema_description:"Credential value"`
}

type SetCredentialOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type GetCredentialStatusInput struct{}

type CredentialStatusItem struct {
	Name       string `json:"name"`
	Configured bool   `json:"configured"`
	Source     string `json:"source"`
}

type GetCredentialStatusOutput struct {
	Credentials  []CredentialStatusItem `json:"credentials"`
	RDSReady     bool                   `json:"rds_ready"`
	ComputeReady bool                   `json:"compute_ready"`
}

func RegisterAuthTools(server *mcp.Server, cfg *config.Config) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "nhn_set_credential",
		Description: "Set NHN Cloud credential at runtime. Use when credentials are not configured via file or environment variables. Keys: access_key_id, secret_access_key, mysql_appkey, mariadb_appkey, postgresql_appkey, username, password, tenant_id, region",
	}, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[SetCredentialInput]) (*mcp.CallToolResultFor[SetCredentialOutput], error) {
		validKeys := map[string]bool{
			"access_key_id":     true,
			"secret_access_key": true,
			"mysql_appkey":      true,
			"mariadb_appkey":    true,
			"postgresql_appkey": true,
			"username":          true,
			"password":          true,
			"tenant_id":         true,
			"region":            true,
		}

		if !validKeys[params.Arguments.Key] {
			return &mcp.CallToolResultFor[SetCredentialOutput]{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Invalid key: %s", params.Arguments.Key)}},
				StructuredContent: SetCredentialOutput{
					Success: false,
					Message: fmt.Sprintf("Invalid key: %s. Valid keys: access_key_id, secret_access_key, mysql_appkey, mariadb_appkey, postgresql_appkey, username, password, tenant_id, region", params.Arguments.Key),
				},
				IsError: true,
			}, nil
		}

		cfg.SetInteractive(params.Arguments.Key, params.Arguments.Value)

		return &mcp.CallToolResultFor[SetCredentialOutput]{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Credential '%s' set successfully", params.Arguments.Key)}},
			StructuredContent: SetCredentialOutput{
				Success: true,
				Message: fmt.Sprintf("Credential '%s' set successfully (source: interactive)", params.Arguments.Key),
			},
		}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "nhn_get_credential_status",
		Description: "Check which NHN Cloud credentials are configured and their source (file, env, interactive, or none). Does not expose credential values.",
	}, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[GetCredentialStatusInput]) (*mcp.CallToolResultFor[GetCredentialStatusOutput], error) {
		status := cfg.GetStatus()

		creds := make([]CredentialStatusItem, 0, len(status))
		for name, info := range status {
			creds = append(creds, CredentialStatusItem{
				Name:       name,
				Configured: info["configured"] == "yes",
				Source:     info["source"],
			})
		}

		out := GetCredentialStatusOutput{
			Credentials:  creds,
			RDSReady:     cfg.HasRDSCredentials(),
			ComputeReady: cfg.HasComputeCredentials(),
		}

		summary := fmt.Sprintf("RDS Ready: %v, Compute Ready: %v", out.RDSReady, out.ComputeReady)
		return &mcp.CallToolResultFor[GetCredentialStatusOutput]{
			Content:           []mcp.Content{&mcp.TextContent{Text: summary}},
			StructuredContent: out,
		}, nil
	})
}
