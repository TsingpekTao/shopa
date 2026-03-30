# user-profile-svc configuration

## Files

- `config.local.yaml`: local development
- `config.dev.yaml`: shared dev environment
- `config.prod.yaml`: production environment

## Run with specific config

Git Bash:

```bash
GF_GCFG_FILE=manifest/config/config.local.yaml gf run main.go
```

PowerShell:

```powershell
$env:GF_GCFG_FILE = "manifest/config/config.local.yaml"
gf run main.go
```

## Notes

- Keep secrets in config/secret manager, never hardcode in code.
- Local uses `127.0.0.1:3307` (MySQL) and `127.0.0.1:6379` (Redis).
- HTTP is `:8003`, gRPC is `:8002` (aligned with `edge-gateway` upstream).
