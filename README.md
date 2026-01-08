# NHN Cloud MCP Server

Model Context Protocol (MCP) server for NHN Cloud services. Enables AI assistants (Claude, etc.) to interact with NHN Cloud resources.

## Features

### Implemented Tools

| Tool | Description |
|------|-------------|
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

Set environment variables before running:

```bash
# Required for RDS services (MySQL, MariaDB, PostgreSQL)
export NHN_CLOUD_ACCESS_KEY_ID="your-access-key"
export NHN_CLOUD_SECRET_ACCESS_KEY="your-secret-key"
export NHN_CLOUD_MYSQL_APPKEY="your-mysql-appkey"

# Optional: Additional RDS services
export NHN_CLOUD_MARIADB_APPKEY="your-mariadb-appkey"
export NHN_CLOUD_POSTGRESQL_APPKEY="your-postgresql-appkey"

# Optional: For Compute/Network services
export NHN_CLOUD_USERNAME="your-email"
export NHN_CLOUD_PASSWORD="your-api-password"
export NHN_CLOUD_TENANT_ID="your-tenant-id"

# Optional: Region (default: kr1)
export NHN_CLOUD_REGION="kr1"
```

## Usage with Claude Desktop

Add to your Claude Desktop configuration (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "nhn-cloud": {
      "command": "/path/to/nhn-cloud-mcp",
      "env": {
        "NHN_CLOUD_ACCESS_KEY_ID": "your-access-key",
        "NHN_CLOUD_SECRET_ACCESS_KEY": "your-secret-key",
        "NHN_CLOUD_MYSQL_APPKEY": "your-mysql-appkey"
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

Make sure environment variables are set in your shell profile.

## Development

### Project Structure

```
nhn-cloud-mcp/
├── main.go           # MCP server entry point
├── config/
│   └── config.go     # Configuration management
├── tools/
│   └── mysql.go      # MySQL tools implementation
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
