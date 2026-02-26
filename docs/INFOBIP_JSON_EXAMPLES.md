# Ejemplos JSON para Integración Infobip Chatbot

## Endpoint
```
POST https://tu-dominio.com/api/v1/deliveries/infobip
Content-Type: application/json
```

---

## Ejemplo 1: Instalación con 2 dispensers de pie

**Request:**
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

**Response (201 Created):**
```json
{
  "token": "1234",
  "message": "Entrega creada exitosamente"
}
```

---

## Ejemplo 2: Recambio con 1 de pie y 1 de mesada

**Request:**
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

**Response (201 Created):**
```json
{
  "token": "5678",
  "message": "Entrega creada exitosamente"
}
```

---

## Ejemplo 3: Retiro con 3 dispensers de mesada

**Request:**
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

**Response (201 Created):**
```json
{
  "token": "9012",
  "message": "Entrega creada exitosamente"
}
```

---

## Campos del Request

| Campo | Tipo | Obligatorio | Descripción | Valores Válidos |
|-------|------|-------------|-------------|-----------------|
| `nro_cta` | string | **SÍ** | Número de cuenta del cliente | Mínimo 1, máximo 50 caracteres |
| `nro_rto` | **SÍ** | string | Número de reparto | Mínimo 1, máximo 50 caracteres |
| `tipos` | object | **SÍ** | Cantidades por tipo de dispenser | Ver detalles abajo |
| `tipos.P` | number | NO | Cantidad de dispensers de **Pie** | 0 o más (default: 0) |
| `tipos.M` | number | NO | Cantidad de dispensers de **Mesada** | 0 o más (default: 0) |
| `tipo_entrega` | string | **SÍ** | Tipo de operación | `"Instalacion"`, `"Retiro"` o `"Recambio"` |
| `entregado_por` | string | **SÍ** | Responsable de la entrega | `"Repartidor"` o `"Tecnico"` |
| `session_id` | string | **SÍ** | ID de sesión del chatbot | Mínimo 1 carácter |
| `fecha_accion` | string | NO | Fecha programada de la entrega | Formato: `YYYY-MM-DD` o ISO 8601 |

### Importante sobre `tipos`:
- **Debe haber al menos 1 dispenser** (P + M >= 1)
- **Máximo 10 dispensers en total** (P + M <= 10)
- La cantidad total se calcula automáticamente

---

## Campos del Response

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `token` | string | Token de 4 dígitos que identifica la entrega |
| `message` | string | Mensaje de confirmación |

---

## Errores Posibles

### Error 400: Sin dispensers
**Request inválido:**
```json
{
  "nro_cta": "CTA12345",
  "nro_rto": "RTO001",
  "tipos": {
    "P": 0,
    "M": 0
  },
  "tipo_entrega": "Instalacion",
  "entregado_por": "Repartidor",
  "session_id": "INF-12345"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "validation_failed",
  "message": "Debe especificar al menos un dispenser (P o M)"
}
```

---

### Error 400: Campo faltante
**Request inválido (falta nro_cta):**
```json
{
  "nro_rto": "RTO001",
  "tipos": {
    "P": 1,
    "M": 0
  },
  "tipo_entrega": "Instalacion",
  "entregado_por": "Repartidor",
  "session_id": "INF-12345"
}
```

**Response (400 Bad Request):**
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

### Error 400: Tipo de entrega inválido
**Request inválido:**
```json
{
  "nro_cta": "CTA12345",
  "nro_rto": "RTO001",
  "tipos": {
    "P": 1,
    "M": 0
  },
  "tipo_entrega": "TipoInventado",
  "entregado_por": "Repartidor",
  "session_id": "INF-12345"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "invalid_input",
  "details": "tipo_entrega debe ser uno de: Instalacion, Retiro, Recambio"
}
```

---

## Ejemplo de Integración con el Chatbot

### Flujo de Conversación

1. **Chatbot pregunta al cliente:**
   ```
   ¿Cuál es su número de cuenta?
   → Cliente: CTA12345
   
   ¿Número de reparto?
   → Cliente: RTO001
   
   ¿Qué tipo de entrega necesita?
   1. Instalación
   2. Retiro  
   3. Recambio
   → Cliente: 1
   
   ¿Cuántos dispensers de pie necesita?
   → Cliente: 2
   
   ¿Cuántos dispensers de mesada necesita?
   → Cliente: 1
   
   ¿Fecha preferida? (YYYY-MM-DD o presione Enter para hoy)
   → Cliente: 2026-02-28
   ```

2. **Chatbot construye el JSON:**
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
     "session_id": "INF-SESSION-123456",
     "fecha_accion": "2026-02-28"
   }
   ```

3. **Chatbot recibe el token:**
   ```json
   {
     "token": "4567",
     "message": "Entrega creada exitosamente"
   }
   ```

4. **Chatbot informa al cliente:**
   ```
   ✅ ¡Entrega registrada exitosamente!
   
   📋 Su token de confirmación: 4567
   📅 Fecha programada: 28/02/2026
   📦 Total de dispensers: 3 (2 de pie, 1 de mesada)
   
   Este token identifica su entrega y será utilizado 
   por nuestro equipo de reparto.
   
   Recibirá una notificación cuando el repartidor 
   esté en camino.
   ```

---

## Códigos de Estado HTTP

| Código | Significado | Cuándo ocurre |
|--------|-------------|---------------|
| 201 | Created | Entrega creada exitosamente |
| 400 | Bad Request | Datos inválidos o incompletos |
| 500 | Internal Server Error | Error en el servidor |

---

## Testing con cURL

### Ejemplo básico:
```bash
curl -X POST http://localhost:8080/api/v1/deliveries/infobip \
  -H "Content-Type: application/json" \
  -d '{
    "nro_cta": "CTA12345",
    "nro_rto": "RTO001",
    "tipos": {
      "P": 2,
      "M": 1
    },
    "tipo_entrega": "Instalacion",
    "entregado_por": "Repartidor",
    "session_id": "TEST-123"
  }'
```

### Ejemplo con fecha:
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
    "session_id": "TEST-456",
    "fecha_accion": "2026-03-01"
  }'
```

---

## Notas Importantes

1. **Session ID único**: Cada llamada debe tener un `session_id` único del chatbot

2. **Tipos de dispenser**:
   - `P` = Pie (dispensers de pie)
   - `M` = Mesada (dispensers de mesada)

3. **Cantidad automática**: No es necesario enviar el campo `cantidad`, se calcula como P + M

4. **Token de 4 dígitos**: El sistema genera un token único de 4 dígitos para cada entrega

5. **Fecha opcional**: Si no se envía `fecha_accion`, se usa la fecha actual

6. **Case sensitive**: Los valores de `tipo_entrega` y `entregado_por` distinguen mayúsculas:
   - ✅ Correcto: `"Instalacion"`, `"Repartidor"`
   - ❌ Incorrecto: `"instalacion"`, `"repartidor"`

---

## Contacto Técnico

Para soporte o dudas técnicas sobre la integración:
- Revisar documentación completa en: `/docs/INFOBIP_DELIVERY_API.md`
- Ver ejemplos de código en: `/tests/test_infobip_delivery.ps1`
