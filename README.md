# ldapSSHPubkeyReader

OpenSSH relies on an external application to provide SSH keys when AuthorizedKeysCommand directive is used. In order to use SSH with keys stored in LDAP, a suitable script or program needs to provide the keys. Traditionally, some shell script that calls ldapsearch is being used for this purpose.

This go program can be used to fullfill the same need, only faster and in a lightweight manner.

### Compile

The Makefile supports building for multiple platforms with optimized settings:

#### Quick Start
```bash
make                # Build for all platforms (Linux amd64, macOS amd64, macOS arm64)
make info           # Show available build targets and configuration
```

#### Platform-Specific Builds
```bash
make linux-amd64    # Linux x86_64
make darwin-amd64   # macOS Intel
make darwin-arm64   # macOS Apple Silicon
```

#### Development Builds (with race detection)
```bash
make dev-all        # Development builds for all platforms
make dev-linux-amd64 # Development build for Linux
```

#### Build Features
- **CGO disabled**: Creates fully static binaries with no external dependencies
- **Optimized**: Uses `-w -s` flags to strip debug info and reduce binary size
- **Secure**: Uses `-trimpath` to remove filesystem paths from binaries
- **Cross-platform**: Supports Linux amd64, macOS Intel, and macOS Apple Silicon

#### Clean Up
```bash
make clean          # Remove all built binaries
```

The build produces optimized, statically-linked binaries named with platform suffixes (e.g., `ldapPubKeyReader-linux-amd64`). If you need debug information, use the development builds instead.

### Configuration

The application supports multiple configuration methods with the following priority order:
1. Environment variables (highest priority)
2. Configuration file
3. Default values (lowest priority)

#### Configuration File Locations

The application searches for `ldapPubKeyReader.json` in these directories:
```
/etc/ssh/ldapPubKeyReader.json
/etc/ldapPubKeyReader.json
./ldapPubKeyReader.json (relative to binary location)
```

You can also specify additional config paths as command line arguments.

#### Configuration Format

```json
{
    "ldap_server": {
        "url": "ldaps://ldap.example.com:636",
        "bind_dn": "cn=readonly,dc=example,dc=com",
        "bind_password": "password",
        "connection_timeout": "10s",
        "search_timeout": "30s",
        "max_retries": 3,
        "retry_delay": "1s"
    },
    "base_dn": "ou=People,dc=example,dc=com",
    "public_key_attribute": "sshPublicKey",
    "user_attribute": "uid",
    "search_filter": "(%s=%s)",
    "tls": {
        "insecure_skip_verify": false,
        "cert_file": "/path/to/client.crt",
        "key_file": "/path/to/client.key",
        "ca_file": "/path/to/ca.crt"
    },
    "debug": false
}
```

#### Environment Variables

All configuration options can be set via environment variables with the `LDAP_SSH_` prefix:

```bash
export LDAP_SSH_LDAP_SERVER_URL="ldaps://ldap.example.com:636"
export LDAP_SSH_BASE_DN="ou=People,dc=example,dc=com"
export LDAP_SSH_LDAP_SERVER_BIND_DN="cn=readonly,dc=example,dc=com"
export LDAP_SSH_LDAP_SERVER_BIND_PASSWORD="password"
export LDAP_SSH_DEBUG="true"
```

#### Configuration Options

| Option | Description | Default | Required |
|--------|-------------|---------|----------|
| `ldap_server.url` | LDAP server URL | - | Yes |
| `ldap_server.bind_dn` | Bind DN for authentication | - | No |
| `ldap_server.bind_password` | Bind password | - | No |
| `ldap_server.connection_timeout` | Connection timeout | 10s | No |
| `ldap_server.search_timeout` | Search timeout | 30s | No |
| `ldap_server.max_retries` | Max connection retries | 3 | No |
| `ldap_server.retry_delay` | Delay between retries | 1s | No |
| `base_dn` | LDAP search base DN | - | Yes |
| `public_key_attribute` | SSH public key attribute | sshPublicKey | No |
| `user_attribute` | User identifier attribute | uid | No |
| `search_filter` | LDAP search filter template | (%s=%s) | No |
| `tls.insecure_skip_verify` | Skip TLS certificate verification | false | No |
| `tls.cert_file` | Client certificate file | - | No |
| `tls.key_file` | Client private key file | - | No |
| `tls.ca_file` | CA certificate file | - | No |
| `debug` | Enable debug output | false | No |

#### Timeout Configuration Examples

```bash
# Set a very short connection timeout (useful for testing)
export LDAP_SSH_LDAP_SERVER_CONNECTION_TIMEOUT="2s"

# Set a longer search timeout for slow LDAP servers
export LDAP_SSH_LDAP_SERVER_SEARCH_TIMEOUT="60s"

# Disable retries (set max_retries to 0)
export LDAP_SSH_LDAP_SERVER_MAX_RETRIES="0"

# Quick retry with minimal delay
export LDAP_SSH_LDAP_SERVER_RETRY_DELAY="100ms"
```

**Timeout Format**: Use Go duration format: `10s`, `1m30s`, `500ms`, `2h`, etc.

### Usage

```bash
ldapPubKeyReader <username> [config-path...]

# Examples:
ldapPubKeyReader john.doe
ldapPubKeyReader john.doe /custom/config/path
LDAP_SSH_DEBUG=true ldapPubKeyReader john.doe
```

#### Arguments

- `username`: LDAP username to search for SSH public keys
- `config-path`: Optional additional configuration search paths

#### Security Features

- **LDAP Injection Protection**: Input sanitization prevents LDAP injection attacks
- **Connection Retry Logic**: Automatic retry with configurable delays and max attempts
- **Comprehensive Timeout Controls**: 
  - Connection timeout for establishing LDAP connections
  - Search timeout for LDAP query operations
  - Network-level timeout via custom dialer
  - Per-operation timeout on LDAP connection
- **TLS Support**: Full TLS configuration including custom certificates
- **Bind Authentication**: Support for authenticated LDAP connections

#### Error Codes

- `0`: Success
- `1`: Configuration error or missing arguments
- `2`: LDAP search error 



