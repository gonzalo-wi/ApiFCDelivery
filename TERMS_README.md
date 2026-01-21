# 🚀 Flujo de Términos y Condiciones con Infobip - Quick Reference

## 📚 Documentación Disponible

- **[IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md)** - Resumen completo de la implementación
- **[docs/TERMS_INTEGRATION.md](docs/TERMS_INTEGRATION.md)** - Documentación técnica detallada
- **[docs/TERMS_QUICKSTART.md](docs/TERMS_QUICKSTART.md)** - Guía rápida de inicio
- **[docs/FRONTEND_INTEGRATION.md](docs/FRONTEND_INTEGRATION.md)** - Integración con Vue.js

## ⚡ Quick Start

### 1. Configurar Variables de Entorno

```bash
# Copiar ejemplo
cp .env.example .env

# Editar .env y agregar:
INFOBIP_BASE_URL=https://api2.infobip.com
INFOBIP_API_KEY=tu-api-key-aqui
APP_BASE_URL=http://localhost:5173
TERMS_TTL_HOURS=48
```

### 2. Ejecutar Aplicación

```bash
go run api/cmd/main.go
```

La tabla `terms_sessions` se crea automáticamente.

### 3. Probar con PowerShell

```powershell
.\scripts\test_terms_flow.ps1
```

## 🌐 Endpoints

| Método | Ruta | Descripción |
|--------|------|-------------|
| `POST` | `/api/v1/infobip/session` | Crear sesión (Infobip) |
| `GET` | `/api/v1/terms/:token` | Consultar estado |
| `POST` | `/api/v1/terms/:token/accept` | Aceptar términos |
| `POST` | `/api/v1/terms/:token/reject` | Rechazar términos |

## 🧪 Test Rápido

```bash
# Crear sesión
curl -X POST http://localhost:8080/api/v1/infobip/session \
  -H "Content-Type: application/json" \
  -d '{"sessionId": "test-123"}'

# Guardar el token de la respuesta
$TOKEN = "token-aqui"

# Aceptar términos
curl -X POST http://localhost:8080/api/v1/terms/$TOKEN/accept
```

## 📁 Archivos Nuevos Creados

```
internal/
  ├── models/terms_session.go
  ├── dto/terms_dto.go
  ├── store/terms_session_store.go
  ├── service/
  │   ├── infobip_client.go
  │   └── terms_session_service.go
  ├── transport/terms_session_handler.go
  └── routes/terms_routes.go

migrations/001_create_terms_sessions.sql
scripts/test_terms_flow.ps1
docs/
  ├── TERMS_INTEGRATION.md
  ├── TERMS_QUICKSTART.md
  └── FRONTEND_INTEGRATION.md
```

## 🔐 Características

- ✅ Token seguro de 64 caracteres (crypto/rand)
- ✅ SessionID nunca expuesto en URLs públicas
- ✅ Expiración configurable (48h default)
- ✅ Estados: PENDING, ACCEPTED, REJECTED, EXPIRED
- ✅ Reintentos automáticos a Infobip (1s, 3s, 7s)
- ✅ Idempotencia garantizada
- ✅ Auditoría completa (IP, User-Agent, timestamps)
- ✅ Logging estructurado con zerolog

## 🔍 Troubleshooting

**Token no encontrado:**
```sql
SELECT * FROM terms_sessions WHERE token = 'tu-token';
```

**Notificaciones fallidas:**
```sql
SELECT * FROM terms_sessions WHERE notify_status = 'FAILED';
```

**Ver logs:**
```bash
# Logs en tiempo real
go run api/cmd/main.go | grep "términos"
```

## 📖 Más Información

Ver [IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md) para detalles completos.

---

**Implementado y listo para usar! 🎉**
