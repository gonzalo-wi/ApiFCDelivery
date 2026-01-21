# Flujo Integrado: Entregas con Términos y Condiciones

## 📋 Descripción General

Este documento describe el flujo completo de creación de entregas que **requiere la aceptación de términos y condiciones ANTES** de que la entrega sea creada en el sistema.

## 🔄 Flujo de Operación

### Flujo Anterior (Sin Términos)
```
1. POST /api/v1/deliveries → Crea entrega inmediatamente
```

### Flujo Nuevo (Con Términos - REQUISITO PREVIO)
```
1. POST /api/v1/deliveries/initiate
   ↓
2. Sistema crea sesión de términos y guarda datos de entrega temporalmente
   ↓
3. Cliente recibe URL de términos → Debe aceptar
   ↓
4. POST /api/v1/deliveries/complete/{token}
   ↓
5. Sistema valida aceptación → Crea entrega definitiva
```

**Importante:** Si el cliente rechaza o no acepta los términos, la entrega **NUNCA se crea**.

---

## 🎯 Endpoints del Flujo Integrado

### 1. Iniciar Creación de Entrega (Requiere Términos)

**Endpoint:** `POST /api/v1/deliveries/initiate`

**Descripción:** Inicia el proceso de creación de entrega. Los datos se guardan temporalmente y se genera una sesión de términos que el cliente debe aceptar.

**Request Body:**
```json
{
  "nro_cta": "CTA12345",
  "nro_rto": "RTO67890",
  "dispensers": [
    {
      "marca": "CocaCola",
      "nro_serie": "CC-001",
      "tipo": "Enfriador"
    },
    {
      "marca": "Pepsi",
      "nro_serie": "PP-002",
      "tipo": "Calentador"
    }
  ],
  "cantidad": 2,
  "tipo_entrega": "Instalacion",
  "fecha_accion": "2024-01-15"
}
```

**Response (200 OK):**
```json
{
  "token": "a1b2c3d4e5f6...",
  "terms_url": "https://app.com/terms/a1b2c3d4e5f6",
  "expires_at": "2024-01-17T14:30:00Z",
  "message": "Por favor, acepte los términos y condiciones para completar la entrega"
}
```

**Validaciones:**
- `nro_cta`: Requerido, 1-50 caracteres
- `nro_rto`: Requerido, 1-50 caracteres (se usa como `sessionId` en términos)
- `dispensers`: Requerido, debe tener al menos 1
- `cantidad`: Requerido, debe coincidir con el número de dispensers (1-3)
- `tipo_entrega`: Requerido, valores: `Instalacion`, `Retiro`, `Recambio`
- `fecha_accion`: Opcional, formato ISO 8601 o YYYY-MM-DD

**Errores Posibles:**
- `400 Bad Request`: Validación fallida
- `500 Internal Server Error`: Error al crear sesión de términos

---

### 2. Cliente Acepta Términos

El cliente debe acceder a la URL de términos (`terms_url`) y aceptar:

**Endpoint:** `POST /api/v1/terms/accept/{token}`

Ver documentación completa en `docs/TERMS_ENDPOINTS.md`

---

### 3. Completar Creación de Entrega

**Endpoint:** `POST /api/v1/deliveries/complete/{token}`

**Descripción:** Valida que los términos fueron aceptados y crea la entrega definitiva en la base de datos.

**Path Parameter:**
- `token`: Token de la sesión de términos (64 caracteres)

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Entrega creada exitosamente después de aceptar términos",
  "delivery": {
    "id": 123,
    "nro_cta": "CTA12345",
    "nro_rto": "RTO67890",
    "dispensers": [
      {
        "id": 1,
        "marca": "CocaCola",
        "nro_serie": "CC-001",
        "tipo": "Enfriador",
        "delivery_id": 123,
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
      }
    ],
    "cantidad": 2,
    "estado": "Completado",
    "tipo_entrega": "Instalacion",
    "token": "1234",
    "terms_session_id": 456,
    "fecha_accion": "2024-01-15T00:00:00Z",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Errores Posibles:**
- `400 Bad Request`: 
  - Términos no aceptados (estado: PENDING, REJECTED, EXPIRED)
  - Sesión expirada
- `404 Not Found`: Token de términos no encontrado
- `500 Internal Server Error`: Error al crear entrega

**Validaciones del Backend:**
1. Sesión de términos existe y es válida
2. Estado de términos es `ACCEPTED`
3. Sesión no está expirada (`expires_at > now()`)
4. Hay datos de entrega guardados en la sesión
5. Datos de entrega son deserializables

---

### 4. Verificar Estado (Opcional)

**Endpoint:** `GET /api/v1/deliveries/status/{token}`

**Descripción:** Redirige al endpoint de términos para verificar el estado.

**Response:**
```json
{
  "message": "Use el endpoint /api/v1/terms/status/:token para verificar el estado de los términos",
  "token": "a1b2c3d4e5f6..."
}
```

---

## 🗂️ Arquitectura de Datos

### Relación entre Modelos

```
TermsSession (1) ←→ (0..1) Delivery
```

**TermsSession:**
- `id`: ID único de la sesión
- `session_id`: Identificador de sesión (usa `nro_rto`)
- `token`: Token de 64 caracteres (nunca expuesto al cliente)
- `status`: PENDING → ACCEPTED → (Crear Delivery)
- `delivery_data`: JSON con datos temporales de la entrega
- `expires_at`: Fecha de expiración (default 48h)

**Delivery:**
- `id`: ID único de la entrega
- `terms_session_id`: FK a `terms_sessions` (nullable)
- `nro_cta`, `nro_rto`, etc.
- `estado`: Completado (automático si aceptó términos)

---

## 🔐 Seguridad

1. **Token Seguro (64 chars):** Generado con `crypto/rand`
2. **Expiration:** Sesiones expiran después de TTL configurado (default 48h)
3. **Estado Inmutable:** Una vez ACCEPTED/REJECTED, no se puede cambiar
4. **Datos Temporales:** `delivery_data` se borra después de crear entrega (opcional)
5. **Validación de Estado:** No se puede completar entrega si términos no están aceptados

---

## 📊 Flujo de Estados

```
[Initiate Delivery]
       ↓
  TermsSession
   (PENDING)
       ↓
   ┌───┴────┐
   ↓        ↓
ACCEPTED  REJECTED
   ↓        ↓
[Create   [No
Delivery]  Delivery]
```

---

## 🧪 Ejemplo de Uso Completo

### Paso 1: Iniciar Entrega
```bash
curl -X POST http://localhost:8080/api/v1/deliveries/initiate \
  -H "Content-Type: application/json" \
  -d '{
    "nro_cta": "CTA12345",
    "nro_rto": "RTO67890",
    "dispensers": [
      {
        "marca": "CocaCola",
        "nro_serie": "CC-001",
        "tipo": "Enfriador"
      }
    ],
    "cantidad": 1,
    "tipo_entrega": "Instalacion"
  }'
```

**Respuesta:**
```json
{
  "token": "a1b2c3d4e5f6...",
  "terms_url": "https://app.com/terms/a1b2c3d4e5f6",
  "expires_at": "2024-01-17T14:30:00Z",
  "message": "Por favor, acepte los términos y condiciones para completar la entrega"
}
```

### Paso 2: Cliente Acepta Términos
```bash
curl -X POST http://localhost:8080/api/v1/terms/accept/a1b2c3d4e5f6 \
  -H "Content-Type: application/json" \
  -d '{
    "webhook_url": "https://infobip.com/webhook/terms-accepted"
  }'
```

### Paso 3: Completar Entrega
```bash
curl -X POST http://localhost:8080/api/v1/deliveries/complete/a1b2c3d4e5f6
```

**Respuesta:**
```json
{
  "success": true,
  "message": "Entrega creada exitosamente después de aceptar términos",
  "delivery": {
    "id": 123,
    "nro_rto": "RTO67890",
    "estado": "Completado",
    ...
  }
}
```

---

## ⚙️ Configuración Requerida

En `.env`:
```env
# URLs
APP_BASE_URL=https://tu-app.com
INFOBIP_BASE_URL=https://api.infobip.com
INFOBIP_API_KEY=tu_api_key

# Configuración de términos
TERMS_TTL_HOURS=48
```

---

## 🚀 Archivos Modificados/Creados

### Nuevos Archivos:
1. `internal/service/delivery_with_terms_service.go` - Servicio integrado
2. `internal/transport/delivery_with_terms_handler.go` - Handlers HTTP
3. `internal/dto/delivery_integration_dto.go` - DTOs de integración

### Archivos Modificados:
1. `internal/models/delivery.go` - Agregado `TermsSessionID`
2. `internal/models/terms_session.go` - Agregado `DeliveryData`
3. `internal/routes/delivery_routes.go` - Nuevas rutas
4. `internal/routes/router.go` - Registro de rutas
5. `api/cmd/main.go` - Inicialización de componentes

---

## 📝 Notas Importantes

1. **Flujo Obligatorio:** Este es el nuevo flujo recomendado para todas las entregas que requieren términos
2. **Flujo Legacy:** El endpoint `POST /api/v1/deliveries` sigue existiendo para retrocompatibilidad
3. **Infobip:** El webhook a Infobip se envía cuando el cliente acepta/rechaza términos
4. **TTL:** Configurar `TERMS_TTL_HOURS` según necesidades del negocio
5. **NroRto:** Se usa como `sessionId` en términos para trazabilidad

---

## 🐛 Troubleshooting

### Error: "los términos no han sido aceptados"
- **Causa:** El cliente aún no aceptó términos o los rechazó
- **Solución:** Verificar estado con `GET /api/v1/terms/status/{token}`

### Error: "la sesión de términos ha expirado"
- **Causa:** Pasaron más de 48h (o TTL configurado)
- **Solución:** Reiniciar el flujo con `/initiate`

### Error: "sesión de términos no encontrada"
- **Causa:** Token inválido o sesión no existe
- **Solución:** Verificar que el token sea correcto

---

## 📚 Referencias

- [TERMS_README.md](./TERMS_README.md) - Documentación completa de términos
- [TERMS_ENDPOINTS.md](./TERMS_ENDPOINTS.md) - Todos los endpoints de términos
- [IMPLEMENTATION_SUMMARY.md](./IMPLEMENTATION_SUMMARY.md) - Resumen de implementación

---

**Última actualización:** 2024-01-15
**Versión:** 2.0 (Flujo integrado)