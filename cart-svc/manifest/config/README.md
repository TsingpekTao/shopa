# cart-svc configuration

## Files

- config.yaml: default entry (same as local template unless changed)
- config.local.yaml: local development
- config.dev.yaml: dev environment
- config.prod.yaml: production environment

## Switch config file

PowerShell:

`powershell
 = "config.dev.yaml"
`

CLI args:

`ash
./main --gf.gcfg.file=config.prod.yaml
`

## Notes

- Put all third-party dependency configuration in this directory.
- Do not hardcode database/redis/mq secrets in code.
- For production, prefer injecting secrets via K8s Secret / CI variables.
