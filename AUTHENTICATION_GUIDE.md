# Autenticación JWT con Servicio Externo

## Descripción General

La API de GoFrioCalor ahora está protegida con autenticación JWT utilizando un servicio externo de autenticación en `http://192.168.0.55:8087`.

## Flujo de Autenticación

### 1. Obtener Token (Proveedor)

Los proveedores deben obtener un token JWT válido antes de consumir los endpoints de la API.

**Endpoint:** `http://localhost:8095/dispenser-operations/auth/generar-token` (o tu dominio público)

**Method:** `GET`

**Headers:**
```
x-api-key: MOBEUS_kG7pX2sV9nQ4aJ1cL8rT0yZ5wH3eU6mF2dC9bA1sR4xP
```

**Respuesta Exitosa:**
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "proveedor": "mobeus",
    "detail": "El token es válido por 30 minutos.",
    "expires_at_local": "2026-03-05T12:00:00-03:00"
}
```

⚠️ **Importante:** 
- Este endpoint actúa como proxy hacia el servicio interno de autenticación
- El proveedor NO necesita acceso directo a la red interna
- El token expira en 30 minutos. Después de ese tiempo, debe generarse uno nuevo.

### 2. Consumir Endpoints de la API

Una vez obtenido el token, debe incluirse en cada request a los endpoints de GoFrioCalor.

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Ejemplo de Request:**
```bash
curl -X POST http://localhost:8095/dispenser-operations/api/v1/deliveries \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "work_order_id": 12345,
    "delivery_status": "en_curso"
  }'
```

## Endpoints Protegidos

Todos los endpoints bajo `/dispenser-operations/api/v1/*` requieren autenticación:

- ✅ `/dispenser-operations/api/v1/deliveries`
- ✅ `/dispenser-operations/api/v1/work-orders`
- ✅ `/dispenser-operations/api/v1/terms`
- ✅ `/dispenser-operations/api/v1/deliveries-with-terms`
- ✅ `/dispenser-operations/api/v1/mobile/*`

## Endpoints Públicos (Sin Autenticación)

- ✅ `/health` - Health check para monitoreo
- ✅ `/dispenser-operations/auth/generar-token` - Obtener token JWT (requiere x-api-key)

## Respuestas de Error

### Token No Proporcionado
**Status:** `401 Unauthorized`
```json
{
    "error": "Token no proporcionado",
    "detail": "Se requiere el header Authorization con un Bearer token"
}
```

### Formato de Token Inválido
**Status:** `401 Unauthorized`
```json
{
    "error": "Formato de token inválido",
    "detail": "El header Authorization debe tener el formato: Bearer <token>"
}
```

### Token Inválido o Expirado
**Status:** `401 Unauthorized`
```json
{
    "error": "Token inválido",
    "detail": "Token inválido o expirado"
}
```

### Servicio de Autenticación No Disponible
**Status:** `401 Unauthorized`
```json
{
    "error": "Token inválido",
    "detail": "Servicio de autenticación no disponible"
}
```

## Configuración

La URL del servicio de autenticación se configura en el archivo `.env`:

```env
AUTH_SERVICE_URL=http://192.168.0.55:8087
```

## Pruebas con PowerShell

### 1. Obtener un Token

```powershell
$headers = @{
    "x-api-key" = "MOBEUS_kG7pX2sV9nQ4aJ1cL8rT0yZ5wH3eU6mF2dC9bA1sR4xP"
}

$response = Invoke-RestMethod -Uri "http://localhost:8095/dispenser-operations/auth/generar-token" -Method GET -Headers $headers
$token = $response.token
Write-Host "Token obtenido: $token"
```

### 2. Usar el Token en una Request

```powershell
$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

$body = @{
    work_order_id = 12345
    delivery_status = "en_curso"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8095/dispenser-operations/api/v1/deliveries" -Method POST -Headers $headers -Body $body
```

### 3. Probar sin Token (Debe Fallar)

```powershell
try {
    Invoke-RestMethod -Uri "http://localhost:8095/dispenser-operations/api/v1/deliveries" -Method GET
} catch {
    Write-Host "Error esperado - 401 Unauthorized"
    $_.Exception.Response.StatusCode
}
```

## Logs

El middleware de autenticación genera logs para monitoreo:

- **Token Válido:** `Token validado exitosamente`
- **Token Inválido:** `Token inválido o expirado`
- **Sin Token:** `Request sin header Authorization`
- **Error de Servicio:** `Error llamando al servicio de validación`

## Implementación Técnica

### Middleware de Autenticación
- **Archivo:** `internal/middleware/auth.go`
- **Función:** `AuthMiddleware(authServiceURL string)`
- **Flujo:**
  1. Extrae el token del header `Authorization`
  2. Valida el formato `Bearer <token>`
  3. Hace una llamada POST a `{authServiceURL}/validar-token`
  4. Permite o rechaza la request según la respuesta

### Validación de Token
- **Timeout:** 5 segundos
- **Endpoint:** `POST {AUTH_SERVICE_URL}/validar-token`
- **Header:** `Authorization: Bearer <token>`

### Configuración
- **Archivo:** `config/env.go`
- **Variable:** `AuthServiceURL`
- **Valor por defecto:** `http://192.168.0.55:8087`

## Notas de Seguridad

1. ✅ Los tokens expiran automáticamente (30 minutos)
2. ✅ Cada proveedor tiene su propia API key
3. ✅ Los endpoints sensibles están protegidos
4. ✅ El health check permanece público para monitoreo
5. ✅ Se implementa timeout para evitar requests colgadas
6. ✅ Los logs permiten auditoría de accesos

## Cambios Realizados

1. ✅ Agregada configuración `AUTH_SERVICE_URL` en `.env`
2. ✅ Creado middleware de autenticación en `internal/middleware/auth.go`
3. ✅ Aplicado middleware a todas las rutas de la API
4. ✅ Mantenido `/health` sin protección para monitoreo
5. ✅ Agregado soporte para header `Authorization` en CORS
