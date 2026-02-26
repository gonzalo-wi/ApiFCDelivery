# API de Entrega desde Infobip

## Descripción
Este endpoint permite al chatbot de Infobip crear entregas directamente, especificando la cantidad de dispensers por tipo (Pie o Mesada). El sistema genera automáticamente un token de 4 dígitos que se devuelve al chatbot para confirmar la creación.

## Endpoint

**POST** `/api/v1/deliveries/infobip`

### Request Body

```json
{
  "nro_cta": "CTA12345",
  "nro_rto": "RTO001",
  "tipos": {
    "P": 2,
    "M": 1
  },
  "tipo_entrega": "Instalacion",
  "entregado_por": "Repartidor",
  "session_id": "54656",
  "fecha_accion": "2026-02-25"
}
```

### Campos del Request

| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `nro_cta` | string | Sí | Número de cuenta del cliente |
| `nro_rto` | string | Sí | Número de reparto |
| `tipos` | object | Sí | Objeto con cantidades por tipo de dispenser |
| `tipos.P` | uint | No | Cantidad de dispensers de Pie (0 si no se especifica) |
| `tipos.M` | uint | No | Cantidad de dispensers de Mesada (0 si no se especifica) |
| `tipo_entrega` | string | Sí | Tipo de entrega: `Instalacion`, `Retiro` o `Recambio` |
| `entregado_por` | string | Sí | Responsable: `Repartidor` o `Tecnico` |
| `session_id` | string | Sí | ID de sesión del chatbot de Infobip |
| `fecha_accion` | string | No | Fecha de la acción en formato YYYY-MM-DD o ISO 8601 |

### Response (201 Created)

```json
{
  "token": "1234",
  "message": "Entrega creada exitosamente"
}
```

### Campos del Response

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `token` | string | Token de 4 dígitos generado para la entrega |
| `message` | string | Mensaje de confirmación |

## Características Importantes

### 1. Cálculo Automático de Cantidad
La cantidad total de dispensers se calcula automáticamente sumando `P + M`. No es necesario enviar el campo `cantidad` en el request.

**Ejemplo:**
- `P: 2` + `M: 1` = 3 dispensers totales

### 2. Tipos de Dispensers
- **P**: Dispenser de Pie
- **M**: Dispenser de Mesada

Debe especificarse al menos un dispenser (P o M mayor a 0).

### 3. Dispensers Placeholder
Los dispensers se crean con información placeholder que será actualizada posteriormente:
- **Marca**: `"PENDIENTE"`
- **NroSerie**: Generado automáticamente (ej: `P-RTO001-1`, `M-RTO001-1`)
- **Tipo**: `P` (Pie) o `M` (Mesada)

### 4. Concurrencia
El endpoint está diseñado para manejar múltiples solicitudes simultáneas de forma segura:
- Generación de tokens thread-safe usando crypto/rand
- Transacciones de base de datos para integridad
- Sin bloqueos explícitos necesarios

### 5. Estado de la Entrega
Las entregas creadas desde Infobip tienen estado inicial `Pendiente` y pueden ser completadas posteriormente.

## Códigos de Error

| Código | Descripción |
|--------|-------------|
| 400 | Bad Request - Datos inválidos o falta información requerida |
| 500 | Internal Server Error - Error al crear la entrega en la base de datos |

### Ejemplos de Errores

**Ningún dispenser especificado:**
```json
{
  "error": "validation_failed",
  "message": "Debe especificar al menos un dispenser (P o M)"
}
```

**Datos inválidos:**
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

## Ejemplos de Uso

### Ejemplo 1: Instalación con 2 dispensers de pie

```bash
curl -X POST http://localhost:8080/api/v1/deliveries/infobip \
  -H "Content-Type: application/json" \
  -d '{
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
  }'
```

**Response:**
```json
{
  "token": "7234",
  "message": "Entrega creada exitosamente"
}
```

### Ejemplo 2: Recambio con 1 de cada tipo

```bash
curl -X POST http://localhost:8080/api/v1/deliveries/infobip \
  -H "Content-Type: application/json" \
  -d '{
    "nro_cta": "CTA99999",
    "nro_rto": "RTO999",
    "tipos": {
      "P": 1,
      "M": 1
    },
    "tipo_entrega": "Recambio",
    "entregado_por": "Tecnico",
    "session_id": "INF-78910"
  }'
```

**Response:**
```json
{
  "token": "4891",
  "message": "Entrega creada exitosamente"
}
```

### Ejemplo 3: Retiro con 3 dispensers de mesada

```bash
curl -X POST http://localhost:8080/api/v1/deliveries/infobip \
  -H "Content-Type: application/json" \
  -d '{
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
  }'
```

**Response:**
```json
{
  "token": "2156",
  "message": "Entrega creada exitosamente"
}
```

## Integración con el Chatbot

### Flujo de Conversación Recomendado

1. **Chatbot recoge información del cliente:**
   - Número de cuenta
   - Número de reparto
   - Tipo de entrega (Instalación, Retiro, Recambio)
   - Cantidad de dispensers de pie
   - Cantidad de dispensers de mesada
   - Fecha preferida (opcional)

2. **Chatbot valida que al menos hay un dispenser**

3. **Chatbot envía request al endpoint**

4. **Sistema responde con el token de 4 dígitos**

5. **Chatbot informa al cliente el token de confirmación**

### Ejemplo de Mensaje del Chatbot

```
✅ ¡Entrega registrada exitosamente!

Su token de confirmación es: 1234

Este token identifica su entrega y será utilizado por nuestro equipo.
Fecha programada: 25/02/2026

Recibirá una notificación cuando el repartidor esté en camino.
```

## Notas Técnicas

### Thread Safety
El generador de tokens utiliza `crypto/rand` que es thread-safe y proporciona números criptográficamente seguros.

### Límites
- Cantidad máxima total de dispensers: 10
- Cantidad mínima: 1
- Longitud del token: exactamente 4 dígitos

### Formato de Fecha
Acepta dos formatos:
- **Simple**: `YYYY-MM-DD` (ej: `2026-02-25`)
- **ISO 8601**: Con hora y zona horaria (ej: `2026-02-25T10:30:00Z`)

Si no se especifica fecha, se usa la fecha y hora actual del sistema.

## Ver también
- [Documentación de Términos](TERMS_INTEGRATION.md)
- [Flow Diagram](FLOW_DIAGRAM.md)
- [Frontend Integration](FRONTEND_INTEGRATION.md)
