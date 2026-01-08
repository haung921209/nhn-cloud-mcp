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

### Credential Priority Chain

Credentials are loaded in order of priority. Higher priority sources override lower ones:

```
┌─────────────────────────────────────────────────────────────┐
│  1. Credentials File (~/.nhncloud/credentials)              │
│     - Shared with NHN Cloud CLI                             │
│     - Recommended for local development                     │
│     - Most secure (file permissions, not in process list)   │
├─────────────────────────────────────────────────────────────┤
│  2. Environment Variables                                   │
│     - Override credentials file values                      │
│     - Good for CI/CD, Docker, cloud deployments            │
│     - Set via shell or MCP client config                    │
├─────────────────────────────────────────────────────────────┤
│  3. Interactive (Runtime via MCP Tool)                      │
│     - Set credentials during conversation                   │
│     - Fallback when file/env not configured                │
│     - Not persisted (session only)                          │
└─────────────────────────────────────────────────────────────┘
```

### Required Credentials by Service

| Service | Required Credentials |
|---------|---------------------|
| **RDS MySQL** | `access_key_id` + `secret_access_key` + `rds_app_key` |
| **RDS MariaDB** | `access_key_id` + `secret_access_key` + `rds_mariadb_app_key` |
| **RDS PostgreSQL** | `access_key_id` + `secret_access_key` + `rds_postgresql_app_key` |
| **Compute/Network** | `username` + `api_password` + `tenant_id` |

### Option 1: Credentials File (Recommended)

Create `~/.nhncloud/credentials` (same format as NHN Cloud CLI):

```ini
[default]
# Required: OAuth credentials for RDS services
access_key_id = your-access-key
secret_access_key = your-secret-key
region = kr1

# RDS App Keys (get from NHN Cloud Console > each RDS service > App Key)
rds_app_key = your-mysql-appkey
rds_mariadb_app_key = your-mariadb-appkey
rds_postgresql_app_key = your-postgresql-appkey

# Optional: Identity credentials for Compute/Network services
username = your-email@example.com
api_password = your-api-password
tenant_id = your-tenant-id
```

Secure the file:
```bash
chmod 600 ~/.nhncloud/credentials
chmod 700 ~/.nhncloud
```

### Option 2: Environment Variables

Environment variables override credentials file values:

| File Key | Environment Variable | Description |
|----------|---------------------|-------------|
| `access_key_id` | `NHN_CLOUD_ACCESS_KEY_ID` | OAuth access key |
| `secret_access_key` | `NHN_CLOUD_SECRET_ACCESS_KEY` | OAuth secret key |
| `region` | `NHN_CLOUD_REGION` | Region (default: kr1) |
| `rds_app_key` | `NHN_CLOUD_MYSQL_APPKEY` | MySQL service app key |
| `rds_mariadb_app_key` | `NHN_CLOUD_MARIADB_APPKEY` | MariaDB service app key |
| `rds_postgresql_app_key` | `NHN_CLOUD_POSTGRESQL_APPKEY` | PostgreSQL service app key |
| `username` | `NHN_CLOUD_USERNAME` | API username (email) |
| `api_password` | `NHN_CLOUD_PASSWORD` | API password |
| `tenant_id` | `NHN_CLOUD_TENANT_ID` | Tenant ID |

Example:
```bash
export NHN_CLOUD_ACCESS_KEY_ID="your-access-key"
export NHN_CLOUD_SECRET_ACCESS_KEY="your-secret-key"
export NHN_CLOUD_MYSQL_APPKEY="your-mysql-appkey"
```

### Option 3: Interactive (Runtime)

When credentials are not configured via file or environment, use MCP tools:

**Set a credential:**
```
nhn_set_credential(key="access_key_id", value="your-key")
nhn_set_credential(key="secret_access_key", value="your-secret")
nhn_set_credential(key="mysql_appkey", value="your-appkey")
```

**Check credential status:**
```
nhn_get_credential_status()
```

Returns which credentials are configured and their source (file/env/interactive).

**Available keys for `nhn_set_credential`:**
- `access_key_id`, `secret_access_key`, `region`
- `mysql_appkey`, `mariadb_appkey`, `postgresql_appkey`
- `username`, `password`, `tenant_id`

### How to Get Credentials

1. **Access Key & Secret Key**: NHN Cloud Console > My Page > API Security Settings
2. **App Keys**: NHN Cloud Console > Select Service (e.g., RDS for MySQL) > URL & Appkey
3. **Tenant ID**: NHN Cloud Console > Compute > Instance > API Endpoint Information

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
