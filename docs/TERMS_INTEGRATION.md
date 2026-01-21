# Integración de Términos y Condiciones con Infobip

## 📋 Descripción General

Este módulo implementa un flujo seguro de aceptación de términos y condiciones integrado con Infobip Bots. El flujo permite que Infobip solicite la aceptación de términos al cliente mediante un link único y seguro, y reciba notificación cuando el cliente acepta o rechaza.

## 🔄 Flujo del Proceso

```
1. Infobip → POST /api/v1/infobip/session { "sessionId": "..." }
   Backend genera token único y devuelve URL

2. Cliente → Accede a /terms/{token} (Frontend Vue)
   Frontend consulta GET /api/v1/terms/{token} para mostrar estado

3. Cliente → Acepta términos
   Frontend → POST /api/v1/terms/{token}/accept

4. Backend → Notifica a Infobip
   POST https://api2.infobip.com/bots/webhook/{sessionId}
   Con reintentos automáticos (3 intentos: 1s, 3s, 7s)
```

## 🔐 Características de Seguridad

- ✅ Token público de 64 caracteres (256 bits) generado con `crypto/rand`
- ✅ SessionID de Infobip NO se expone en URLs públicas
- ✅ Expiración configurable (default: 48 horas)
- ✅ Estados: PENDING, ACCEPTED, REJECTED, EXPIRED
- ✅ Idempotencia: aceptar/rechazar múltiples veces devuelve 200 sin renotificar
- ✅ Auditoría: guarda IP, User-Agent, timestamps
- ✅ Reintentos automáticos con backoff exponencial
- ✅ Tracking de intentos de notificación y errores

## 📦 Estructura de Archivos Creados

```
internal/
  ├── models/
  │   └── terms_session.go          # Modelo de datos
  ├── dto/
  │   └── terms_dto.go               # DTOs de request/response
  ├── store/
  │   └── terms_session_store.go     # Capa de persistencia
  ├── service/
  │   ├── infobip_client.go          # Cliente HTTP para Infobip
  │   └── terms_session_service.go   # Lógica de negocio
  ├── transport/
  │   └── terms_session_handler.go   # Handlers HTTP
  └── routes/
      └── terms_routes.go            # Definición de rutas

config/
  ├── env.go                         # Config actualizada
  └── database.go                    # Migration actualizada

migrations/
  └── 001_create_terms_sessions.sql  # Script SQL

.env.example                         # Variables de entorno
```

## 🗄️ Tabla de Base de Datos

```sql
CREATE TABLE terms_sessions (
    id INT PRIMARY KEY AUTO_INCREMENT,
    token VARCHAR(64) UNIQUE NOT NULL,
    session_id VARCHAR(255) NOT NULL,
    status ENUM('PENDING', 'ACCEPTED', 'REJECTED', 'EXPIRED'),
    created_at TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    accepted_at TIMESTAMP NULL,
    rejected_at TIMESTAMP NULL,
    ip VARCHAR(45) NULL,
    user_agent TEXT NULL,
    notify_status ENUM('PENDING', 'SENT', 'FAILED'),
    notify_attempts INT DEFAULT 0,
    last_error TEXT NULL
);
```

## 🔧 Configuración

### Variables de Entorno

Agregar al archivo `.env`:

```env
# Infobip
INFOBIP_BASE_URL=https://api2.infobip.com
INFOBIP_API_KEY=tu-api-key-de-infobip

# Aplicación
APP_BASE_URL=https://mi-dominio.com
TERMS_TTL_HOURS=48
```

### Valores por Defecto

- `INFOBIP_BASE_URL`: `https://api2.infobip.com`
- `APP_BASE_URL`: `http://localhost:5173`
- `TERMS_TTL_HOURS`: `48`

## 🌐 API Endpoints

### 1. Crear Sesión (desde Infobip)

**Endpoint:** `POST /api/v1/infobip/session`

**Descripción:** Infobip llama este endpoint para generar un link de términos.

**Request:**
```json
{
  "sessionId": "unique-infobip-session-id"
}
```

**Response:** `200 OK`
```json
{
  "token": "abc123...def456",
  "url": "https://mi-dominio.com/terms/abc123...def456",
  "expiresAt": "2025-12-26T10:30:00Z"
}
```

**Errores:**
- `400 Bad Request`: sessionId faltante o inválido
- `500 Internal Server Error`: Error en BD o generación de token

**Ejemplo curl:**
```bash
curl -X POST http://localhost:8080/api/v1/infobip/session \
  -H "Content-Type: application/json" \
  -d '{"sessionId": "test-session-123"}'
```

---

### 2. Consultar Estado del Token

**Endpoint:** `GET /api/v1/terms/:token`

**Descripción:** El frontend consulta el estado actual del token.

**Response:** `200 OK`
```json
{
  "status": "PENDING",
  "expiresAt": "2025-12-26T10:30:00Z",
  "acceptedAt": null,
  "rejectedAt": null
}
```

Posibles valores de `status`:
- `PENDING`: Esperando acción del cliente
- `ACCEPTED`: Términos aceptados
- `REJECTED`: Términos rechazados
- `EXPIRED`: Token expirado

**Errores:**
- `400 Bad Request`: Token no proporcionado
- `404 Not Found`: Token no existe

**Ejemplo curl:**
```bash
curl http://localhost:8080/api/v1/terms/abc123...def456
```

---

### 3. Aceptar Términos

**Endpoint:** `POST /api/v1/terms/:token/accept`

**Descripción:** El cliente acepta los términos. Backend notifica a Infobip automáticamente.

**Request:** (vacío o con datos opcionales)
```json
{}
```

**Response:** `200 OK`
```json
{
  "status": "ACCEPTED",
  "message": "Términos aceptados exitosamente",
  "acceptedAt": "2025-12-24T15:30:00Z"
}
```

**Idempotencia:** Si ya fue aceptado previamente:
```json
{
  "status": "ACCEPTED",
  "message": "Términos ya fueron aceptados previamente",
  "acceptedAt": "2025-12-24T15:30:00Z"
}
```

**Errores:**
- `400 Bad Request`: Token no proporcionado o ya está en estado no modificable
- `404 Not Found`: Token no existe
- `410 Gone`: Token expirado
- `500 Internal Server Error`: Error en BD o notificación

**Headers capturados automáticamente:**
- IP del cliente: `c.ClientIP()`
- User-Agent: `User-Agent` header

**Ejemplo curl:**
```bash
curl -X POST http://localhost:8080/api/v1/terms/abc123...def456/accept \
  -H "Content-Type: application/json" \
  -H "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64)"
```

**Notificación a Infobip:**

El backend automáticamente envía:
```http
POST https://api2.infobip.com/bots/webhook/{sessionId}
Authorization: App {INFOBIP_API_KEY}
Content-Type: application/json

{
  "event": "TERMS_ACCEPTED",
  "token": "abc123...def456",
  "acceptedAt": "2025-12-24T15:30:00Z"
}
```

Con reintentos:
- Intento 1: inmediato
- Intento 2: después de 1 segundo
- Intento 3: después de 3 segundos
- Intento 4: después de 7 segundos

---

### 4. Rechazar Términos (Opcional)

**Endpoint:** `POST /api/v1/terms/:token/reject`

**Descripción:** El cliente rechaza los términos.

**Request:** (vacío)
```json
{}
```

**Response:** `200 OK`
```json
{
  "status": "REJECTED",
  "message": "Términos rechazados",
  "rejectedAt": "2025-12-24T15:35:00Z"
}
```

**Ejemplo curl:**
```bash
curl -X POST http://localhost:8080/api/v1/terms/abc123...def456/reject \
  -H "Content-Type: application/json"
```

**Notificación a Infobip:**
```json
{
  "event": "TERMS_REJECTED",
  "token": "abc123...def456",
  "rejectedAt": "2025-12-24T15:35:00Z"
}
```

---

## 🔄 Flujo de Notificación a Infobip

### Payload Enviado

```json
{
  "event": "TERMS_ACCEPTED",  // o "TERMS_REJECTED"
  "token": "...",
  "acceptedAt": "2025-12-24T15:30:00Z",  // o rejectedAt
  "rejectedAt": null
}
```

### Headers

```
POST /bots/webhook/{sessionId}
Host: api2.infobip.com
Authorization: App {INFOBIP_API_KEY}
Content-Type: application/json
```

### Reintentos

- **Estrategia:** Backoff exponencial
- **Intentos:** 3 reintentos
- **Delays:** 1s → 3s → 7s
- **Timeout:** 10 segundos por request

### Estados de Notificación

En la tabla `terms_sessions`:
- `notify_status = 'PENDING'`: No enviado aún
- `notify_status = 'SENT'`: Enviado exitosamente
- `notify_status = 'FAILED'`: Falló después de todos los reintentos
- `notify_attempts`: Número de intentos realizados
- `last_error`: Último error registrado (si falló)

---

## 🎨 Integración con Frontend (Vue)

### Ruta del Frontend

El frontend debe tener una ruta: `/terms/:token`

### Flujo Sugerido

```vue
<template>
  <div v-if="loading">Cargando...</div>
  
  <div v-else-if="status === 'EXPIRED'">
    <h1>Link Expirado</h1>
    <p>Este link ha expirado. Por favor, solicita uno nuevo.</p>
  </div>
  
  <div v-else-if="status === 'ACCEPTED'">
    <h1>Términos Aceptados</h1>
    <p>Ya aceptaste los términos el {{ acceptedAt }}</p>
  </div>
  
  <div v-else-if="status === 'PENDING'">
    <h1>Términos y Condiciones</h1>
    <div>{{ termsContent }}</div>
    <button @click="acceptTerms">Aceptar</button>
    <button @click="rejectTerms">Rechazar</button>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const token = route.params.token
const status = ref('PENDING')
const loading = ref(true)
const acceptedAt = ref(null)

const API_BASE = 'http://localhost:8080/api/v1'

onMounted(async () => {
  try {
    const { data } = await axios.get(`${API_BASE}/terms/${token}`)
    status.value = data.status
    acceptedAt.value = data.acceptedAt
  } catch (error) {
    console.error('Error loading terms:', error)
  } finally {
    loading.value = false
  }
})

const acceptTerms = async () => {
  try {
    const { data } = await axios.post(`${API_BASE}/terms/${token}/accept`)
    status.value = data.status
    acceptedAt.value = data.acceptedAt
    alert(data.message)
  } catch (error) {
    alert('Error aceptando términos: ' + error.response?.data?.error)
  }
}

const rejectTerms = async () => {
  try {
    const { data } = await axios.post(`${API_BASE}/terms/${token}/reject`)
    status.value = data.status
    alert(data.message)
  } catch (error) {
    alert('Error rechazando términos: ' + error.response?.data?.error)
  }
}
</script>
```

---

## 🧪 Pruebas con cURL

### Flujo Completo de Prueba

#### 1. Crear una sesión (simular Infobip)

```bash
curl -X POST http://localhost:8080/api/v1/infobip/session \
  -H "Content-Type: application/json" \
  -d '{"sessionId": "test-session-abc123"}'
```

**Respuesta esperada:**
```json
{
  "token": "f4d3c2b1a0987654321fedcba0123456789abcdef0123456789abcdef01234567",
  "url": "http://localhost:5173/terms/f4d3c2b1a0987654321fedcba0123456789abcdef0123456789abcdef01234567",
  "expiresAt": "2025-12-26T15:30:00Z"
}
```

#### 2. Consultar estado del token

```bash
TOKEN="f4d3c2b1a0987654321fedcba0123456789abcdef0123456789abcdef01234567"
curl http://localhost:8080/api/v1/terms/$TOKEN
```

**Respuesta:**
```json
{
  "status": "PENDING",
  "expiresAt": "2025-12-26T15:30:00Z",
  "acceptedAt": null,
  "rejectedAt": null
}
```

#### 3. Aceptar términos

```bash
curl -X POST http://localhost:8080/api/v1/terms/$TOKEN/accept \
  -H "Content-Type: application/json" \
  -H "User-Agent: curl/test"
```

**Respuesta:**
```json
{
  "status": "ACCEPTED",
  "message": "Términos aceptados exitosamente",
  "acceptedAt": "2025-12-24T15:30:00Z"
}
```

#### 4. Verificar idempotencia (aceptar nuevamente)

```bash
curl -X POST http://localhost:8080/api/v1/terms/$TOKEN/accept \
  -H "Content-Type: application/json"
```

**Respuesta:**
```json
{
  "status": "ACCEPTED",
  "message": "Términos ya fueron aceptados previamente",
  "acceptedAt": "2025-12-24T15:30:00Z"
}
```

#### 5. Consultar estado final

```bash
curl http://localhost:8080/api/v1/terms/$TOKEN
```

**Respuesta:**
```json
{
  "status": "ACCEPTED",
  "expiresAt": "2025-12-26T15:30:00Z",
  "acceptedAt": "2025-12-24T15:30:00Z",
  "rejectedAt": null
}
```

---

## 📊 Monitoreo y Logs

### Logs Importantes

El sistema usa `zerolog` para logging estructurado:

```go
// Al crear sesión
log.Info().
  Str("session_id", sessionID).
  Str("token", token).
  Time("expires_at", expiresAt).
  Msg("Sesión de términos creada exitosamente")

// Al aceptar términos
log.Info().
  Str("token", token).
  Str("session_id", session.SessionID).
  Str("ip", ip).
  Msg("Términos aceptados, iniciando notificación a Infobip")

// Notificación exitosa
log.Info().
  Str("session_id", sessionID).
  Int("attempts", attempt+1).
  Msg("Notificación a Infobip exitosa")

// Notificación fallida
log.Error().
  Err(lastError).
  Str("session_id", sessionID).
  Int("attempts", maxRetries).
  Msg("Notificación a Infobip falló después de todos los reintentos")
```

### Consultas SQL Útiles

**Ver todas las sesiones:**
```sql
SELECT * FROM terms_sessions ORDER BY created_at DESC;
```

**Ver sesiones pendientes expiradas:**
```sql
SELECT * FROM terms_sessions 
WHERE status = 'PENDING' 
  AND expires_at < NOW();
```

**Ver fallos de notificación:**
```sql
SELECT id, token, session_id, notify_status, notify_attempts, last_error
FROM terms_sessions 
WHERE notify_status = 'FAILED';
```

**Estadísticas por estado:**
```sql
SELECT status, COUNT(*) as count
FROM terms_sessions
GROUP BY status;
```

---

## ⚠️ Consideraciones Importantes

### Seguridad

1. **NUNCA exponer sessionId**: Solo el token debe estar en URLs públicas
2. **Validar CORS**: Configurar origins permitidos en `.env`
3. **HTTPS en producción**: APP_BASE_URL debe usar https://
4. **API Key segura**: Guardar INFOBIP_API_KEY en secretos (no en repo)

### Rendimiento

1. **Notificaciones asíncronas**: Se ejecutan en goroutine separada
2. **Timeout HTTP**: 10 segundos por request a Infobip
3. **Context**: Todos los métodos usan `context.Context`
4. **Índices DB**: Token, session_id, status, expires_at

### Mantenimiento

1. **Limpiar sesiones expiradas**: Crear job periódico para eliminar registros antiguos
2. **Reintentar fallos**: Considerar un worker para reintentar notificaciones fallidas
3. **Logs**: Rotar logs en producción

### Personalización

**Cambiar payload de Infobip:**

Editar [infobip_client.go](internal/service/infobip_client.go):
```go
type InfobipWebhookPayload struct {
    Event      string     `json:"event"`
    Token      string     `json:"token"`
    // Agregar más campos según necesites
    CustomField string    `json:"customField"`
}
```

**Cambiar reintentos:**

Editar [terms_session_service.go](internal/service/terms_session_service.go):
```go
return &termsSessionService{
    store:         store,
    infobipClient: infobipClient,
    maxRetries:    5,  // Cambiar número de reintentos
    retryDelays:   []time.Duration{1*time.Second, 2*time.Second, 5*time.Second, 10*time.Second},
}
```

---

## 🚀 Despliegue

### Checklist de Producción

- [ ] Actualizar `.env` con valores de producción
- [ ] Configurar `INFOBIP_API_KEY` en secretos del servidor
- [ ] Cambiar `APP_BASE_URL` a dominio real con HTTPS
- [ ] Habilitar SSL/TLS en el servidor
- [ ] Configurar CORS apropiadamente
- [ ] Configurar rotación de logs
- [ ] Monitorear tabla `terms_sessions` para fallos
- [ ] Crear job para limpiar registros antiguos (>30 días)
- [ ] Probar el flujo completo end-to-end

### Variables de Entorno Producción

```env
ENVIRONMENT=production
INFOBIP_BASE_URL=https://api2.infobip.com
INFOBIP_API_KEY=prod-api-key-secreto
APP_BASE_URL=https://miapp.com
TERMS_TTL_HOURS=48
```

---

## 🤝 Soporte

Para cualquier duda o problema:

1. Revisar logs en `zerolog` con nivel DEBUG
2. Verificar tabla `terms_sessions` en BD
3. Comprobar conectividad con Infobip API
4. Validar configuración en `.env`

---

## 📝 Changelog

### v1.0.0 (2025-12-24)
- ✨ Implementación inicial del flujo de términos y condiciones
- ✨ Integración con Infobip Bots
- ✨ Sistema de reintentos con backoff exponencial
- ✨ Idempotencia en aceptación/rechazo
- ✨ Auditoría completa (IP, User-Agent, timestamps)
- ✨ Estados y expiración de tokens
- ✨ Logging estructurado con zerolog

---

## 📄 Licencia

[Especificar licencia del proyecto]
