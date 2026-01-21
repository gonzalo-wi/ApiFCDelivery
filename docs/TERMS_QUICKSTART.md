# 🚀 Quick Start - Términos y Condiciones con Infobip

## Configuración Rápida

### 1. Actualizar `.env`

```bash
cp .env.example .env
```

Editar `.env` y agregar:
```env
INFOBIP_BASE_URL=https://api2.infobip.com
INFOBIP_API_KEY=tu-api-key-aqui
APP_BASE_URL=http://localhost:5173
TERMS_TTL_HOURS=48
```

### 2. Ejecutar la aplicación

La tabla `terms_sessions` se creará automáticamente gracias a GORM AutoMigrate.

```bash
go run api/cmd/main.go
```

### 3. Probar con cURL

```bash
# Crear sesión (simular Infobip)
curl -X POST http://localhost:8080/api/v1/infobip/session \
  -H "Content-Type: application/json" \
  -d '{"sessionId": "test-123"}'

# Guardar el token de la respuesta
TOKEN="<token-recibido>"

# Consultar estado
curl http://localhost:8080/api/v1/terms/$TOKEN

# Aceptar términos
curl -X POST http://localhost:8080/api/v1/terms/$TOKEN/accept
```

## 🎯 Endpoints Disponibles

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/api/v1/infobip/session` | Crear sesión (Infobip) |
| GET | `/api/v1/terms/:token` | Consultar estado |
| POST | `/api/v1/terms/:token/accept` | Aceptar términos |
| POST | `/api/v1/terms/:token/reject` | Rechazar términos |

## 📖 Documentación Completa

Ver [TERMS_INTEGRATION.md](TERMS_INTEGRATION.md) para:
- Flujo completo del proceso
- Integración con frontend Vue
- Detalles de seguridad
- Ejemplos avanzados
- Troubleshooting

## 🔍 Verificar en la BD

```sql
SELECT * FROM terms_sessions ORDER BY created_at DESC LIMIT 10;
```

## ⚠️ Recordatorios

- El `sessionId` de Infobip NO debe aparecer en URLs públicas ✅
- Solo el `token` generado por el backend se usa en URLs públicas ✅
- Las notificaciones a Infobip son asíncronas con reintentos automáticos ✅
- El sistema es idempotente: aceptar múltiples veces es seguro ✅
