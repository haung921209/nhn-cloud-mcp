# NHN Cloud MCP Server

Model Context Protocol (MCP) server for NHN Cloud services. Enables AI assistants (Claude, etc.) to interact with NHN Cloud resources.

## Features

### Implemented Tools

| Tool | Description |
|------|-------------|
| `nhn_set_credential` | Set credentials at runtime (interactive auth) |
| `nhn_get_credential_status` | Check which credentials are configured |
| `nhn_mysql_list_instances` | List all RDS MySQL instances |
| `nhn_mysql_get_instance` | Get details of a specific MySQL instance |
| `nhn_mysql_list_flavors` | List available MySQL flavors (instance types) |
| `nhn_mysql_list_backups` | List MySQL backups |

### Planned Tools

- MariaDB instance management
- PostgreSQL instance management
- Compute instance management
- Network/VPC management
- NKS (Kubernetes) cluster management

## Installation

### Build from Source

```bash
git clone https://github.com/haung921209/nhn-cloud-mcp.git
cd nhn-cloud-mcp
go build -o nhn-cloud-mcp .
```

### Go Install

```bash
go install github.com/haung921209/nhn-cloud-mcp@latest
```

## Configuration

Credentials are loaded with priority: **File > Environment > Interactive**

### Option 1: Credentials File (Recommended)

Share credentials with NHN Cloud CLI via `~/.nhncloud/credentials`:

```ini
[default]
access_key_id = your-access-key
secret_access_key = your-secret-key
region = kr1

# RDS App Keys
rds_app_key = your-mysql-appkey
rds_mariadb_app_key = your-mariadb-appkey
rds_postgresql_app_key = your-postgresql-appkey

# Compute/Network (optional)
username = your-email
api_password = your-api-password
tenant_id = your-tenant-id
```

### Option 2: Environment Variables

```bash
export NHN_CLOUD_ACCESS_KEY_ID="your-access-key"
export NHN_CLOUD_SECRET_ACCESS_KEY="your-secret-key"
export NHN_CLOUD_MYSQL_APPKEY="your-mysql-appkey"
export NHN_CLOUD_REGION="kr1"
```

### Option 3: Interactive (Runtime)

Use `nhn_set_credential` tool to set credentials at runtime:

```
nhn_set_credential(key="access_key_id", value="your-key")
nhn_set_credential(key="mysql_appkey", value="your-appkey")
```

Check status with `nhn_get_credential_status` tool.

## Usage with Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "nhn-cloud": {
      "command": "/path/to/nhn-cloud-mcp"
    }
  }
}
```

If `~/.nhncloud/credentials` exists, no additional configuration needed.

To override with environment variables:

```json
{
  "mcpServers": {
    "nhn-cloud": {
      "command": "/path/to/nhn-cloud-mcp",
      "env": {
        "NHN_CLOUD_ACCESS_KEY_ID": "override-key",
        "NHN_CLOUD_MYSQL_APPKEY": "override-appkey"
      }
    }
  }
}
```

## Usage with Cursor

Add to Cursor's MCP settings:

```json
{
  "mcpServers": {
    "nhn-cloud": {
      "command": "/path/to/nhn-cloud-mcp"
    }
  }
}
```

## Development

### Project Structure

```
nhn-cloud-mcp/
├── main.go           # MCP server entry point
├── config/
│   └── config.go     # Credential loading (file > env > interactive)
├── tools/
│   ├── auth.go       # Credential management tools
│   └── mysql.go      # MySQL tools
├── go.mod
└── README.md
```

### Adding New Tools

1. Create a new file in `tools/` directory (e.g., `tools/mariadb.go`)
2. Define input/output types for your tools
3. Implement tool handlers
4. Register tools in `Register*Tools()` function
5. Call the registration function in `main.go`

### Testing

```bash
# Build
go build -o nhn-cloud-mcp .

# Test with environment variables
NHN_CLOUD_ACCESS_KEY_ID=xxx \
NHN_CLOUD_SECRET_ACCESS_KEY=xxx \
NHN_CLOUD_MYSQL_APPKEY=xxx \
./nhn-cloud-mcp
```

## Dependencies

- [NHN Cloud SDK for Go](https://github.com/haung921209/nhn-cloud-sdk-go) - NHN Cloud API client
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) - Model Context Protocol implementation

## License

Apache License 2.0

## Related Projects

- [NHN Cloud CLI](https://github.com/haung921209/nhn-cloud-cli) - Command-line interface
- [NHN Cloud SDK for Go](https://github.com/haung921209/nhn-cloud-sdk-go) - Go SDK
