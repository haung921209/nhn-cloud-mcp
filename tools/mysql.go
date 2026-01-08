package tools

import (
	"context"
	"fmt"

	"github.com/haung921209/nhn-cloud-mcp/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MySQL Tool Input/Output Types

// ListMySQLInstancesInput - input for listing MySQL instances
type ListMySQLInstancesInput struct{}

// MySQLInstance represents a MySQL instance
type MySQLInstance struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Version     string `json:"version"`
	StorageType string `json:"storage_type"`
	StorageSize int    `json:"storage_size_gb"`
}

// ListMySQLInstancesOutput - output for listing MySQL instances
type ListMySQLInstancesOutput struct {
	Instances []MySQLInstance `json:"instances"`
	Count     int             `json:"count"`
}

// GetMySQLInstanceInput - input for getting a single MySQL instance
type GetMySQLInstanceInput struct {
	InstanceID string `json:"instance_id" jsonschema:"required,description=The ID of the MySQL instance"`
}

// GetMySQLInstanceOutput - output for getting a single MySQL instance
type GetMySQLInstanceOutput struct {
	Instance MySQLInstance `json:"instance"`
}

// ListMySQLFlavorsInput - input for listing MySQL flavors
type ListMySQLFlavorsInput struct{}

// MySQLFlavor represents a MySQL flavor
type MySQLFlavor struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	VCPUs int    `json:"vcpus"`
	RAM   int    `json:"ram_mb"`
}

// ListMySQLFlavorsOutput - output for listing MySQL flavors
type ListMySQLFlavorsOutput struct {
	Flavors []MySQLFlavor `json:"flavors"`
	Count   int           `json:"count"`
}

// ListMySQLBackupsInput - input for listing MySQL backups
type ListMySQLBackupsInput struct {
	InstanceID string `json:"instance_id,omitempty" jsonschema:"description=Filter by instance ID (optional)"`
}

// MySQLBackup represents a MySQL backup
type MySQLBackup struct {
	ID         string `json:"id"`
	InstanceID string `json:"instance_id"`
	Status     string `json:"status"`
	Size       int64  `json:"size_gb"`
	CreatedAt  string `json:"created_at"`
}

// ListMySQLBackupsOutput - output for listing MySQL backups
type ListMySQLBackupsOutput struct {
	Backups []MySQLBackup `json:"backups"`
	Count   int           `json:"count"`
}

// RegisterMySQLTools registers all MySQL-related tools to the MCP server
func RegisterMySQLTools(server *mcp.Server, cfg *config.Config) {
	// List MySQL Instances
	mcp.AddTool(server, &mcp.Tool{
		Name:        "nhn_mysql_list_instances",
		Description: "List all NHN Cloud RDS MySQL instances. Returns instance ID, name, status, version, storage type, and storage size.",
	}, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[ListMySQLInstancesInput]) (*mcp.CallToolResultFor[ListMySQLInstancesOutput], error) {
		out, err := listMySQLInstances(ctx, cfg)
		if err != nil {
			return &mcp.CallToolResultFor[ListMySQLInstancesOutput]{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)}},
				IsError: true,
			}, nil
		}
		return &mcp.CallToolResultFor[ListMySQLInstancesOutput]{
			Content:           []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Found %d MySQL instances", out.Count)}},
			StructuredContent: out,
		}, nil
	})

	// Get MySQL Instance
	mcp.AddTool(server, &mcp.Tool{
		Name:        "nhn_mysql_get_instance",
		Description: "Get details of a specific NHN Cloud RDS MySQL instance by ID.",
	}, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[GetMySQLInstanceInput]) (*mcp.CallToolResultFor[GetMySQLInstanceOutput], error) {
		out, err := getMySQLInstance(ctx, cfg, params.Arguments.InstanceID)
		if err != nil {
			return &mcp.CallToolResultFor[GetMySQLInstanceOutput]{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)}},
				IsError: true,
			}, nil
		}
		return &mcp.CallToolResultFor[GetMySQLInstanceOutput]{
			Content:           []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Instance: %s (%s)", out.Instance.Name, out.Instance.Status)}},
			StructuredContent: out,
		}, nil
	})

	// List MySQL Flavors
	mcp.AddTool(server, &mcp.Tool{
		Name:        "nhn_mysql_list_flavors",
		Description: "List all available NHN Cloud RDS MySQL flavors (instance types). Returns flavor ID, name, vCPUs, and RAM.",
	}, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[ListMySQLFlavorsInput]) (*mcp.CallToolResultFor[ListMySQLFlavorsOutput], error) {
		out, err := listMySQLFlavors(ctx, cfg)
		if err != nil {
			return &mcp.CallToolResultFor[ListMySQLFlavorsOutput]{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)}},
				IsError: true,
			}, nil
		}
		return &mcp.CallToolResultFor[ListMySQLFlavorsOutput]{
			Content:           []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Found %d MySQL flavors", out.Count)}},
			StructuredContent: out,
		}, nil
	})

	// List MySQL Backups
	mcp.AddTool(server, &mcp.Tool{
		Name:        "nhn_mysql_list_backups",
		Description: "List NHN Cloud RDS MySQL backups. Optionally filter by instance ID.",
	}, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[ListMySQLBackupsInput]) (*mcp.CallToolResultFor[ListMySQLBackupsOutput], error) {
		out, err := listMySQLBackups(ctx, cfg, params.Arguments.InstanceID)
		if err != nil {
			return &mcp.CallToolResultFor[ListMySQLBackupsOutput]{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)}},
				IsError: true,
			}, nil
		}
		return &mcp.CallToolResultFor[ListMySQLBackupsOutput]{
			Content:           []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Found %d MySQL backups", out.Count)}},
			StructuredContent: out,
		}, nil
	})
}

// Tool implementations

func listMySQLInstances(ctx context.Context, cfg *config.Config) (ListMySQLInstancesOutput, error) {
	client, err := cfg.NewNHNCloudClient()
	if err != nil {
		return ListMySQLInstancesOutput{}, fmt.Errorf("failed to create client: %w", err)
	}

	result, err := client.MySQL().ListInstances(ctx)
	if err != nil {
		return ListMySQLInstancesOutput{}, fmt.Errorf("failed to list instances: %w", err)
	}

	instances := make([]MySQLInstance, 0, len(result.DBInstances))
	for _, inst := range result.DBInstances {
		instances = append(instances, MySQLInstance{
			ID:          inst.DBInstanceID,
			Name:        inst.DBInstanceName,
			Status:      inst.DBInstanceStatus,
			Version:     inst.DBVersion,
			StorageType: inst.StorageType,
			StorageSize: inst.StorageSize,
		})
	}

	return ListMySQLInstancesOutput{
		Instances: instances,
		Count:     len(instances),
	}, nil
}

func getMySQLInstance(ctx context.Context, cfg *config.Config, instanceID string) (GetMySQLInstanceOutput, error) {
	if instanceID == "" {
		return GetMySQLInstanceOutput{}, fmt.Errorf("instance_id is required")
	}

	client, err := cfg.NewNHNCloudClient()
	if err != nil {
		return GetMySQLInstanceOutput{}, fmt.Errorf("failed to create client: %w", err)
	}

	result, err := client.MySQL().GetInstance(ctx, instanceID)
	if err != nil {
		return GetMySQLInstanceOutput{}, fmt.Errorf("failed to get instance: %w", err)
	}

	// DatabaseInstanceResponse embeds DatabaseInstance, so fields are accessed directly
	return GetMySQLInstanceOutput{
		Instance: MySQLInstance{
			ID:          result.DBInstanceID,
			Name:        result.DBInstanceName,
			Status:      result.DBInstanceStatus,
			Version:     result.DBVersion,
			StorageType: result.StorageType,
			StorageSize: result.StorageSize,
		},
	}, nil
}

func listMySQLFlavors(ctx context.Context, cfg *config.Config) (ListMySQLFlavorsOutput, error) {
	client, err := cfg.NewNHNCloudClient()
	if err != nil {
		return ListMySQLFlavorsOutput{}, fmt.Errorf("failed to create client: %w", err)
	}

	result, err := client.MySQL().ListFlavors(ctx)
	if err != nil {
		return ListMySQLFlavorsOutput{}, fmt.Errorf("failed to list flavors: %w", err)
	}

	flavors := make([]MySQLFlavor, 0, len(result.DBFlavors))
	for _, f := range result.DBFlavors {
		flavors = append(flavors, MySQLFlavor{
			ID:    f.FlavorID,
			Name:  f.FlavorName,
			VCPUs: f.Vcpus,
			RAM:   f.Ram,
		})
	}

	return ListMySQLFlavorsOutput{
		Flavors: flavors,
		Count:   len(flavors),
	}, nil
}

func listMySQLBackups(ctx context.Context, cfg *config.Config, instanceID string) (ListMySQLBackupsOutput, error) {
	client, err := cfg.NewNHNCloudClient()
	if err != nil {
		return ListMySQLBackupsOutput{}, fmt.Errorf("failed to create client: %w", err)
	}

	result, err := client.MySQL().ListBackups(ctx, instanceID, "", 0, 100)
	if err != nil {
		return ListMySQLBackupsOutput{}, fmt.Errorf("failed to list backups: %w", err)
	}

	backups := make([]MySQLBackup, 0, len(result.Backups))
	for _, b := range result.Backups {
		backups = append(backups, MySQLBackup{
			ID:         b.BackupID,
			InstanceID: b.DBInstanceID,
			Status:     b.BackupStatus,
			Size:       b.BackupSize,
			CreatedAt:  b.CreatedYmdt,
		})
	}

	return ListMySQLBackupsOutput{
		Backups: backups,
		Count:   len(backups),
	}, nil
}
