# Guía de Uso desde Postman - Endpoints Mobile

## 🔧 Configuración Base

**Base URL:** `http://localhost:8080/api/v1`

**Headers para todas las peticiones:**
```
Content-Type: application/json
```

---

## 📱 Flujo Completo desde Postman

### Pre-requisito: Crear un Delivery con Token

Primero necesitas crear un delivery que tenga dispensers asignados.

#### 1️⃣ POST `/api/v1/deliveries` - Crear Delivery
```http
POST http://localhost:8080/api/v1/deliveries
Content-Type: application/json

{
  "nro_cta": "12345",
  "nro_rto": "9",
  "cantidad": 2,
  "estado": "Pendiente",
  "tipo_entrega": "Instalacion",
  "entregado_por": "Tecnico",
  "fecha_accion": "2025-11-12"
}
```

**Respuesta:**
```json
{
  "id": 24,
  "nro_cta": "12345",
  "nro_rto": "9",
  "dispensers": [],
  "cantidad": 2,
  "token": "2780",    👈 GUARDA ESTE TOKEN
  "estado": "Pendiente",
  "tipo_entrega": "Instalacion",
  "entregado_por": "Tecnico",
  "created_at": "2025-11-12T15:46:53Z"
}
```

#### 2️⃣ POST `/api/v1/dispensers` - Agregar Dispensers
Debes agregar los dispensers uno por uno:

**Dispenser 1:**
```http
POST http://localhost:8080/api/v1/dispensers
Content-Type: application/json

{
  "marca": "LAMO",
  "nro_serie": "LM123456789",
  "tipo": "P",
  "delivery_id": 24    👈 USA EL ID DEL DELIVERY
}
```

**Dispenser 2:**
```http
POST http://localhost:8080/api/v1/dispensers
Content-Type: application/json

{
  "marca": "LAMO",
  "nro_serie": "LM987654321",
  "tipo": "M",
  "delivery_id": 24
}
```

---

## 🚀 Endpoints Mobile (Simulan la App del Repartidor)

### 3️⃣ GET `/api/v1/mobile/deliveries/search` - Buscar Deliveries (OPCIONAL)

**💡 Útil para proveedores externos:** Este endpoint permite buscar deliveries por fecha y reparto antes de iniciar el flujo. Devuelve el **ID del delivery** necesario para completar la entrega posteriormente.

**Request:**
```http
GET http://localhost:8080/api/v1/mobile/deliveries/search?fecha_accion=2025-11-12&nro_rto=9
```

**Parámetros:**
- `fecha_accion` (requerido): Fecha en formato YYYY-MM-DD
- `nro_rto` (opcional): Número de reparto para filtrar

**Respuesta Exitosa:**
```json
[
  {
    "id": 24,                     👈 ID necesario para complete-delivery
    "fecha_accion": "2025-11-12",
    "nro_cta": "12345",
    "token": "2780"
  },
  {
    "id": 25,
    "fecha_accion": "2025-11-12",
    "nro_cta": "67890",
    "token": "3891"
  }
]
```

**Respuesta si no hay resultados:**
```json
[]
```

---

### 4️⃣ POST `/api/v1/mobile/validate-token` - Validar Token del Cliente

El repartidor le pide al cliente: **token + número de cuenta + fecha de entrega** para validar la entrega.

**Request:**
```http
POST http://localhost:8080/api/v1/mobile/validate-token
Content-Type: application/json

{
  "token": "2780",           👈 Token proporcionado por el cliente
  "nro_cta": "12345",        👈 Número de cuenta del cliente
  "fecha_accion": "2025-11-12"  👈 Fecha de la entrega (YYYY-MM-DD)
}
```

**Respuesta Exitosa:**
```json
{
  "valid": true,
  "message": "Token válido",
  "delivery": {
    "id": 24,
    "nro_cta": "12345",
    "nro_rto": "9",
    "cantidad": 2,
    "tipo_entrega": "Instalacion",
    "fecha_accion": "2025-11-12"
  },
  "dispensers": [
    {
      "id": 45,
      "marca": "LAMO",
      "nro_serie": "LM123456789",
      "tipo": "P",
      "validated": false
    },
    {
      "id": 46,
      "marca": "LAMO",
      "nro_serie": "LM987654321",
      "tipo": "M",
      "validated": false
    }
  ]
}
```

**Respuesta si el token es inválido:**
```json
{
  "valid": false,
  "message": "Token inválido o entrega ya completada"
}
```

---

### 5️⃣ POST `/api/v1/mobile/validate-dispenser` - Validar Dispenser Escaneado

El repartidor escanea cada dispenser (código QR/barras) y valida que pertenezca al delivery.

**Request - Dispenser 1:**
```http
POST http://localhost:8080/api/v1/mobile/validate-dispenser
Content-Type: application/json

{
  "delivery_id": 24,
  "nro_serie": "LM123456789"
}
```

**Respuesta Exitosa:**
```json
{
  "valid": true,
  "message": "Dispenser válido",
  "dispenser": {
    "id": 45,
    "marca": "LAMO",
    "nro_serie": "LM123456789",
    "tipo": "P",
    "validated": true
  }
}
```

**Respuesta si NO pertenece:**
```json
{
  "valid": false,
  "message": "El dispenser no pertenece a esta entrega"
}
```

**Repite para el Dispenser 2:**
```http
POST http://localhost:8080/api/v1/mobile/validate-dispenser
Content-Type: application/json

{
  "delivery_id": 24,
  "nro_serie": "LM987654321"
}
```

---

### 6️⃣ POST `/api/v1/mobile/complete-delivery` - Completar Entrega

Una vez validados todos los dispensers, el repartidor completa la entrega.

**Request:**
```http
POST http://localhost:8080/api/v1/mobile/complete-delivery
Content-Type: application/json

{
  "delivery_id": 24,
  "token": "2780",
  "validated_dispensers": [
    "LM123456789",
    "LM987654321"
  ]
}
```

**Respuesta Exitosa:**
```json
{
  "success": true,
  "message": "Entrega completada exitosamente",
  "delivery_id": 24,
  "work_order_queued": true    👈 Mensaje enviado a RabbitMQ
}
```

**Respuestas de Error:**

❌ Token inválido:
```json
{
  "success": false,
  "message": "Token inválido"
}
```

❌ Entrega ya procesada:
```json
{
  "success": false,
  "message": "La entrega ya fue procesada (estado: Completado)"
}
```

❌ Faltan dispensers por escanear:
```json
{
  "success": false,
  "message": "Faltan dispensers por escanear (esperados: 2, validados: 1)"
}
```

---

## 📊 ¿Qué pasa después de completar la entrega?

1. ✅ El delivery cambia de estado "Pendiente" → "Completado"
2. 📨 Se publica un mensaje a RabbitMQ en la cola `q.workorder.generate`
3. 🔄 El Worker consume el mensaje automáticamente
4. 📝 Se crea una WorkOrder (OT-XXXXXX)
5. 📄 Se genera el PDF (actualmente mock)
6. 📧 Se envía el email al cliente (actualmente mock)

**Puedes ver el proceso en los logs del servidor:**
```
INF Work order message published successfully
INF Processing work order message
INF Work order created order_number=OT-000020
INF PDF generated successfully
INF Email sent successfully
```

---

## 🎯 Colección de Postman (JSON)

Si quieres importar directamente a Postman, copia este JSON:

```json
{
  "info": {
    "name": "GoFrioCalor - Mobile Endpoints",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "1. Validar Token",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"token\": \"2780\"\n}"
        },
        "url": {
          "raw": "http://localhost:8080/api/v1/mobile/validate-token",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["api", "v1", "mobile", "validate-token"]
        }
      }
    },
    {
      "name": "2. Validar Dispenser",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"delivery_id\": 24,\n  \"nro_serie\": \"LM123456789\"\n}"
        },
        "url": {
          "raw": "http://localhost:8080/api/v1/mobile/validate-dispenser",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["api", "v1", "mobile", "validate-dispenser"]
        }
      }
    },
    {
      "name": "3. Completar Entrega",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"delivery_id\": 24,\n  \"token\": \"2780\",\n  \"validated_dispensers\": [\n    \"LM123456789\",\n    \"LM987654321\"\n  ]\n}"
        },
        "url": {
          "raw": "http://localhost:8080/api/v1/mobile/complete-delivery",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["api", "v1", "mobile", "complete-delivery"]
        }
      }
    }
  ]
}
```

---

## 🔍 Verificar en RabbitMQ

Después de completar la entrega, puedes verificar en:
- **URL:** http://192.168.0.250:15672
- **Usuario:** admin-
- **Password:** admin123
- **Cola:** q.workorder.generate

Verás el mensaje publicado y consumido.

---

## ⚠️ Notas Importantes

1. **Token único:** Cada delivery tiene un token único
2. **Estado:** Solo puedes completar deliveries en estado "Pendiente"
3. **Validación completa:** Debes validar TODOS los dispensers antes de completar
4. **Idempotencia:** Una vez completado, no se puede volver a completar
5. **RabbitMQ:** El servidor debe estar conectado a RabbitMQ para encolar mensajes
