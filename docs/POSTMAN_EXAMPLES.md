# Ejemplos para Postman - API Infobip Delivery

## Configuración Inicial de Postman

### 1. Crear Nueva Request
- Method: **POST**
- URL: `http://localhost:8080/api/v1/deliveries/infobip`

### 2. Headers
```
Content-Type: application/json
```

---

## 📋 Ejemplo 1: Instalación con 2 dispensers de pie

### Request Body (JSON)
```json
{
  "nro_cta": "CTA12345",
  "nro_rto": "RTO001",
  "tipos": {
    "P": 2,
    "M": 0
  },
  "tipo_entrega": "Instalacion",
  "entregado_por": "Repartidor",
  "session_id": "INF-54656",
  "fecha_accion": "2026-02-25"
}
```

### Expected Response (201 Created)
```json
{
  "token": "1234",
  "message": "Entrega creada exitosamente"
}
```

---

## 📋 Ejemplo 2: Recambio con 1 de pie y 1 de mesada

### Request Body (JSON)
```json
{
  "nro_cta": "CTA99999",
  "nro_rto": "RTO999",
  "tipos": {
    "P": 1,
    "M": 1
  },
  "tipo_entrega": "Recambio",
  "entregado_por": "Tecnico",
  "session_id": "INF-78910"
}
```

### Expected Response (201 Created)
```json
{
  "token": "5678",
  "message": "Entrega creada exitosamente"
}
```

---

## 📋 Ejemplo 3: Retiro con 3 dispensers de mesada

### Request Body (JSON)
```json
{
  "nro_cta": "CTA77777",
  "nro_rto": "RTO777",
  "tipos": {
    "P": 0,
    "M": 3
  },
  "tipo_entrega": "Retiro",
  "entregado_por": "Repartidor",
  "session_id": "INF-12345",
  "fecha_accion": "2026-02-26T10:30:00Z"
}
```

### Expected Response (201 Created)
```json
{
  "token": "9012",
  "message": "Entrega creada exitosamente"
}
```

---

## 📋 Ejemplo 4: Solo dispensers de pie

### Request Body (JSON)
```json
{
  "nro_cta": "CTA11111",
  "nro_rto": "RTO111",
  "tipos": {
    "P": 3,
    "M": 0
  },
  "tipo_entrega": "Instalacion",
  "entregado_por": "Repartidor",
  "session_id": "INF-11111"
}
```

---

## 📋 Ejemplo 5: Solo dispensers de mesada

### Request Body (JSON)
```json
{
  "nro_cta": "CTA22222",
  "nro_rto": "RTO222",
  "tipos": {
    "P": 0,
    "M": 2
  },
  "tipo_entrega": "Instalacion",
  "entregado_por": "Tecnico",
  "session_id": "INF-22222"
}
```

---

## 📋 Ejemplo 6: Con fecha ISO 8601

### Request Body (JSON)
```json
{
  "nro_cta": "CTA33333",
  "nro_rto": "RTO333",
  "tipos": {
    "P": 1,
    "M": 1
  },
  "tipo_entrega": "Recambio",
  "entregado_por": "Repartidor",
  "session_id": "INF-33333",
  "fecha_accion": "2026-03-01T14:30:00Z"
}
```

---

## 📋 Ejemplo 7: Sin fecha (usa fecha actual)

### Request Body (JSON)
```json
{
  "nro_cta": "CTA44444",
  "nro_rto": "RTO444",
  "tipos": {
    "P": 2,
    "M": 2
  },
  "tipo_entrega": "Instalacion",
  "entregado_por": "Repartidor",
  "session_id": "INF-44444"
}
```

---

## ❌ Ejemplos de Errores

### Error 1: Sin dispensers (cantidad 0)

#### Request Body (JSON)
```json
{
  "nro_cta": "CTA00000",
  "nro_rto": "RTO000",
  "tipos": {
    "P": 0,
    "M": 0
  },
  "tipo_entrega": "Instalacion",
  "entregado_por": "Repartidor",
  "session_id": "INF-00000"
}
```

#### Expected Response (400 Bad Request)
```json
{
  "error": "validation_failed",
  "message": "Debe especificar al menos un dispenser (P o M)"
}
```

---

### Error 2: Campo faltante (sin nro_cta)

#### Request Body (JSON)
```json
{
  "nro_rto": "RTO555",
  "tipos": {
    "P": 1,
    "M": 0
  },
  "tipo_entrega": "Instalacion",
  "entregado_por": "Repartidor",
  "session_id": "INF-555"
}
```

#### Expected Response (400 Bad Request)
```json
{
  "error": "invalid_input",
  "details": [
    {
      "field": "nro_cta",
      "error": "required"
    }
  ]
}
```

---

### Error 3: Tipo de entrega inválido

#### Request Body (JSON)
```json
{
  "nro_cta": "CTA66666",
  "nro_rto": "RTO666",
  "tipos": {
    "P": 1,
    "M": 1
  },
  "tipo_entrega": "TipoInvalido",
  "entregado_por": "Repartidor",
  "session_id": "INF-666"
}
```

#### Expected Response (400 Bad Request)
```json
{
  "error": "invalid_input",
  "details": "tipo_entrega: Key: 'InfobipDeliveryRequest.TipoEntrega' Error:Field validation for 'TipoEntrega' failed on the 'oneof' tag"
}
```

---

### Error 4: Entregado por inválido

#### Request Body (JSON)
```json
{
  "nro_cta": "CTA77777",
  "nro_rto": "RTO777",
  "tipos": {
    "P": 1,
    "M": 1
  },
  "tipo_entrega": "Instalacion",
  "entregado_por": "PersonaInvalida",
  "session_id": "INF-777"
}
```

#### Expected Response (400 Bad Request)
```json
{
  "error": "invalid_input",
  "details": "entregado_por: debe ser 'Repartidor' o 'Tecnico'"
}
```

---

## 🔧 Configuración de Postman Collection

### Variables de Colección

Puedes crear variables para reutilizar:

```
baseUrl = http://localhost:8080
apiVersion = v1
endpoint = /api/{{apiVersion}}/deliveries/infobip
```

URL completa: `{{baseUrl}}{{endpoint}}`

### Pre-request Script para generar datos dinámicos

```javascript
// Generar session_id único
pm.environment.set("session_id", "INF-" + Date.now());

// Generar número de reparto único
pm.environment.set("nro_rto", "RTO" + Math.floor(Math.random() * 10000));

// Fecha de hoy en formato YYYY-MM-DD
const today = new Date().toISOString().split('T')[0];
pm.environment.set("fecha_hoy", today);
```

### Request Body con variables

```json
{
  "nro_cta": "CTA12345",
  "nro_rto": "{{nro_rto}}",
  "tipos": {
    "P": 2,
    "M": 1
  },
  "tipo_entrega": "Instalacion",
  "entregado_por": "Repartidor",
  "session_id": "{{session_id}}",
  "fecha_accion": "{{fecha_hoy}}"
}
```

### Tests para validar respuesta

```javascript
// Validar status code
pm.test("Status code is 201", function () {
    pm.response.to.have.status(201);
});

// Validar que existe el token
pm.test("Response has token", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('token');
});

// Validar que el token tiene 4 dígitos
pm.test("Token has 4 digits", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.token).to.have.lengthOf(4);
    pm.expect(jsonData.token).to.match(/^\d{4}$/);
});

// Validar mensaje de éxito
pm.test("Response has success message", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.message).to.eql("Entrega creada exitosamente");
});

// Guardar token en variable de entorno
var jsonData = pm.response.json();
pm.environment.set("last_token", jsonData.token);
```

---

## 📦 Importar Colección a Postman

### Formato JSON para importar

```json
{
  "info": {
    "name": "Infobip Delivery API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Crear Entrega - Instalación (2P + 1M)",
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
          "raw": "{\n  \"nro_cta\": \"CTA12345\",\n  \"nro_rto\": \"RTO001\",\n  \"tipos\": {\n    \"P\": 2,\n    \"M\": 1\n  },\n  \"tipo_entrega\": \"Instalacion\",\n  \"entregado_por\": \"Repartidor\",\n  \"session_id\": \"INF-54656\",\n  \"fecha_accion\": \"2026-02-25\"\n}"
        },
        "url": {
          "raw": "http://localhost:8080/api/v1/deliveries/infobip",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["api", "v1", "deliveries", "infobip"]
        }
      }
    },
    {
      "name": "Crear Entrega - Recambio (1P + 1M)",
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
          "raw": "{\n  \"nro_cta\": \"CTA99999\",\n  \"nro_rto\": \"RTO999\",\n  \"tipos\": {\n    \"P\": 1,\n    \"M\": 1\n  },\n  \"tipo_entrega\": \"Recambio\",\n  \"entregado_por\": \"Tecnico\",\n  \"session_id\": \"INF-78910\"\n}"
        },
        "url": {
          "raw": "http://localhost:8080/api/v1/deliveries/infobip",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["api", "v1", "deliveries", "infobip"]
        }
      }
    },
    {
      "name": "Crear Entrega - Retiro (3M)",
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
          "raw": "{\n  \"nro_cta\": \"CTA77777\",\n  \"nro_rto\": \"RTO777\",\n  \"tipos\": {\n    \"P\": 0,\n    \"M\": 3\n  },\n  \"tipo_entrega\": \"Retiro\",\n  \"entregado_por\": \"Repartidor\",\n  \"session_id\": \"INF-12345\"\n}"
        },
        "url": {
          "raw": "http://localhost:8080/api/v1/deliveries/infobip",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["api", "v1", "deliveries", "infobip"]
        }
      }
    },
    {
      "name": "Error - Sin Dispensers",
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
          "raw": "{\n  \"nro_cta\": \"CTA00000\",\n  \"nro_rto\": \"RTO000\",\n  \"tipos\": {\n    \"P\": 0,\n    \"M\": 0\n  },\n  \"tipo_entrega\": \"Instalacion\",\n  \"entregado_por\": \"Repartidor\",\n  \"session_id\": \"INF-00000\"\n}"
        },
        "url": {
          "raw": "http://localhost:8080/api/v1/deliveries/infobip",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["api", "v1", "deliveries", "infobip"]
        }
      }
    }
  ]
}
```

**Para importar:**
1. Abre Postman
2. Click en "Import" (esquina superior izquierda)
3. Pega el JSON anterior
4. Click en "Import"

---

## ✅ Checklist de Pruebas en Postman

- [ ] Instalación con 2 de pie
- [ ] Instalación con 1 de mesada
- [ ] Instalación con mix (P y M)
- [ ] Recambio
- [ ] Retiro
- [ ] Con fecha específica
- [ ] Sin fecha (usa actual)
- [ ] Error: sin dispensers
- [ ] Error: campo faltante
- [ ] Error: tipo entrega inválido
- [ ] Validar token de 4 dígitos
- [ ] Validar formato JSON response

---

## 🔍 Tips para Postman

1. **Formatear JSON automáticamente**: Usa `Ctrl + B` en el body
2. **Ver response bonito**: Tab "Pretty" en la respuesta
3. **Copiar como cURL**: Click en "Code" → "cURL"
4. **Duplicar request**: Click derecho → "Duplicate"
5. **Organizar en carpetas**: Agrupa requests por tipo de entrega

---

## 🌐 Para Producción

Cambia la URL de:
```
http://localhost:8080/api/v1/deliveries/infobip
```

A:
```
https://tu-dominio.com/api/v1/deliveries/infobip
```

Y agrega headers de autenticación si es necesario:
```
Authorization: Bearer {tu-token}
```
