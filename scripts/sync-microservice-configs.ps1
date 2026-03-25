param()

$ErrorActionPreference = 'Stop'

$services = @(
    @{ Name = 'edge-gateway';      HttpPort = 8000; HasGrpc = $false; GrpcPort = 0;    GrpcName = '' },
    @{ Name = 'seller-shop-svc';   HttpPort = 8003; HasGrpc = $false; GrpcPort = 0;    GrpcName = '' },
    @{ Name = 'catalog-svc';       HttpPort = 8004; HasGrpc = $false; GrpcPort = 0;    GrpcName = '' },
    @{ Name = 'review-svc';        HttpPort = 8005; HasGrpc = $false; GrpcPort = 0;    GrpcName = '' },
    @{ Name = 'media-svc';         HttpPort = 8006; HasGrpc = $false; GrpcPort = 0;    GrpcName = '' },
    @{ Name = 'cart-svc';          HttpPort = 8007; HasGrpc = $false; GrpcPort = 0;    GrpcName = '' },
    @{ Name = 'order-svc';         HttpPort = 8008; HasGrpc = $false; GrpcPort = 0;    GrpcName = '' },
    @{ Name = 'fulfillment-svc';   HttpPort = 8009; HasGrpc = $false; GrpcPort = 0;    GrpcName = '' },
    @{ Name = 'aftersale-svc';     HttpPort = 8010; HasGrpc = $false; GrpcPort = 0;    GrpcName = '' },
    @{ Name = 'user-profile-svc';  HttpPort = 8001; HasGrpc = $true;  GrpcPort = 8002; GrpcName = 'user-profile.svc' }
)

function Get-ConfigContent {
    param(
        [string]$Service,
        [int]$HttpPort,
        [bool]$HasGrpc,
        [int]$GrpcPort,
        [string]$GrpcName,
        [string]$EnvName
    )

    $dbName = ('shopa_' + $Service).Replace('-', '_')

    if ($EnvName -eq 'local') {
        $dbLink = "mysql:root:123456@tcp(127.0.0.1:3306)/$($dbName)?charset=utf8mb4&parseTime=true&loc=Local"
        $redisAddr = "127.0.0.1:6379"
        $redisPass = ""
        $mqEnabled = "false"
        $mqURL = "amqp://guest:guest@127.0.0.1:5672/"
        $loggerLevel = "all"
    }
    elseif ($EnvName -eq 'dev') {
        $dbLink = "mysql:dev_user:dev_password@tcp(dev-mysql:3306)/$($dbName)_dev?charset=utf8mb4&parseTime=true&loc=Local"
        $redisAddr = "dev-redis:6379"
        $redisPass = "replace-with-dev-redis-password"
        $mqEnabled = "false"
        $mqURL = "amqp://guest:guest@dev-rabbitmq:5672/"
        $loggerLevel = "all"
    }
    else {
        $dbLink = "mysql:prod_user:prod_password@tcp(prod-mysql:3306)/$($dbName)_prod?charset=utf8mb4&parseTime=true&loc=Local"
        $redisAddr = "prod-redis:6379"
        $redisPass = "replace-with-prod-redis-password"
        $mqEnabled = "true"
        $mqURL = "amqp://prod_user:prod_password@prod-rabbitmq:5672/"
        $loggerLevel = "warn"
    }

    $grpcBlock = ""
    if ($HasGrpc) {
        $grpcBlock = @"
# gRPC server configuration.
grpc:
  address: ":$GrpcPort"
  name: "$GrpcName"
  logPath: ""
  logStdout: true
  errorStack: true
  errorLogEnabled: true
  errorLogPattern: "error-{Ymd}.log"
  accessLogEnabled: true
  accessLogPattern: "access-{Ymd}.log"

"@
    }

    return @"
# ${Service} configuration (${EnvName} environment).
# Start command example:
#   gf run main.go --gf.gcfg.file=config.${EnvName}.yaml

$grpcBlock# HTTP server configuration (OpenAPI/Swagger + optional REST APIs).
server:
  address: ":$HttpPort"
  openapiPath: "/api.json"
  swaggerPath: "/swagger"

# Application identity, used for key prefix and tracing labels.
app:
  env: "$EnvName"
  project: "shopa"
  service: "${Service}"

# External dependencies.
database:
  default:
    link: "$dbLink"

redis:
  default:
    address: "$redisAddr"
    db: 0
    pass: "$redisPass"

mq:
  rabbitmq:
    enabled: $mqEnabled
    url: "$mqURL"
    exchange: "shopa.events"

# Service-to-service addresses (for RPC/http client usage, can be overridden per env).
upstream:
  edgeGateway: "edge-gateway:8000"
  iam: "iam-svc:8001"
  userProfile: "user-profile-svc:8001"

# Object storage / media placeholders.
storage:
  provider: "local"
  endpoint: ""
  accessKey: ""
  secretKey: ""
  bucket: "shopa"

logger:
  level: "$loggerLevel"
  stdout: true
"@
}

function Get-DockerfileContent {
    return @"
FROM loads/alpine:3.8

###############################################################################
#                                INSTALLATION
###############################################################################

ENV WORKDIR /app
# Default to production config file; can be overridden at runtime.
ENV GF_GCFG_FILE config.prod.yaml

ADD resource `$WORKDIR/
ADD manifest/config `$WORKDIR/manifest/config
ADD ./temp/linux_amd64/main `$WORKDIR/main
RUN chmod +x `$WORKDIR/main

###############################################################################
#                                   START
###############################################################################
WORKDIR `$WORKDIR
CMD ./main
"@
}

function Get-ConfigReadmeContent {
    param([string]$Service)

    return @"
# ${Service} configuration

## Files

- config.yaml: default entry (same as local template unless changed)
- config.local.yaml: local development
- config.dev.yaml: dev environment
- config.prod.yaml: production environment

## Switch config file

PowerShell:

```powershell
$env:GF_GCFG_FILE = "config.dev.yaml"
```

CLI args:

```bash
./main --gf.gcfg.file=config.prod.yaml
```

## Notes

- Put all third-party dependency configuration in this directory.
- Do not hardcode database/redis/mq secrets in code.
- For production, prefer injecting secrets via K8s Secret / CI variables.
"@
}

foreach ($svc in $services) {
    $serviceName = $svc.Name
    $basePath = Join-Path (Get-Location).Path $serviceName
    $configDir = Join-Path $basePath 'manifest\config'
    $dockerFile = Join-Path $basePath 'manifest\docker\Dockerfile'

    if (!(Test-Path $configDir)) {
        New-Item -ItemType Directory -Path $configDir -Force | Out-Null
    }

    $localContent = Get-ConfigContent -Service $serviceName -HttpPort $svc.HttpPort -HasGrpc $svc.HasGrpc -GrpcPort $svc.GrpcPort -GrpcName $svc.GrpcName -EnvName 'local'
    $devContent = Get-ConfigContent -Service $serviceName -HttpPort $svc.HttpPort -HasGrpc $svc.HasGrpc -GrpcPort $svc.GrpcPort -GrpcName $svc.GrpcName -EnvName 'dev'
    $prodContent = Get-ConfigContent -Service $serviceName -HttpPort $svc.HttpPort -HasGrpc $svc.HasGrpc -GrpcPort $svc.GrpcPort -GrpcName $svc.GrpcName -EnvName 'prod'

    Set-Content -Path (Join-Path $configDir 'config.local.yaml') -Value $localContent -Encoding UTF8
    Set-Content -Path (Join-Path $configDir 'config.dev.yaml') -Value $devContent -Encoding UTF8
    Set-Content -Path (Join-Path $configDir 'config.prod.yaml') -Value $prodContent -Encoding UTF8
    Set-Content -Path (Join-Path $configDir 'config.yaml') -Value $localContent -Encoding UTF8
    Set-Content -Path (Join-Path $configDir 'README.md') -Value (Get-ConfigReadmeContent -Service $serviceName) -Encoding UTF8

    Set-Content -Path $dockerFile -Value (Get-DockerfileContent) -Encoding UTF8

    Write-Output "synced: $serviceName"
}
