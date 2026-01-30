# ecotachos-automation

Automatización de despliegue, monitoreo, backups y serverless en Go, con CI/CD en GitHub Actions.

## Qué hace
- Despliegue de backend y frontend vía SSH y Docker.
- Health check en CLI y función serverless.
- Backups tipo snapshot locales (extensible a cloud SDKs).
- Provisionamiento cloud simulado, con interfaz para integrar SDKs.

## Estructura
- cmd/deploy/main.go: Automatiza actualización y reinicio de servicios.
- cmd/healthcheck/main.go: Verifica endpoints y guarda estado.
- cmd/backup/main.go: Genera snapshots locales.
- serverless/health_lambda.go: Handler serverless en Go (AWS/GCP).
- .github/workflows/deploy.yml: Pipeline CI/CD para backend y frontend.
- scripts/provision.go: Flags para gestionar recursos cloud.
- diagrams/arquitectura.png: Placeholder del diagrama.

## Uso rápido
### Despliegue (local)
```bash
HOST=1.2.3.4 USER=root KEY_PATH=~/.ssh/id_rsa go run cmd/deploy/main.go
```

### Health check
```bash
HEALTH_URL=http://localhost:8000/api/ia/health/ go run cmd/healthcheck/main.go
```

### Backup
```bash
go run cmd/backup/main.go
```

### Serverless (prueba local)
```bash
go run serverless/health_lambda.go
```

### Provisionamiento simulado
```bash
go run scripts/provision.go -provider aws -action list -region us-east-1
```

## CI/CD
Configura secretos en el repositorio:
- DROPLET_IP, DROPLET_USER, DROPLET_SSH_KEY

El workflow se ejecuta en cada push a `main`/`master` o manualmente.

## Tecnologías
- Go 1.22
- GitHub Actions (appleboy/ssh-action, scp-action)
- Docker/Nginx en servidor

## Diagrama
Ver `diagrams/arquitectura.png` (placeholder). Añade tu diagrama para el PDF.
