# Flujo Mobile de Entregas con RabbitMQ

## 📱 Resumen del Flujo Completo

### Fase 1: Chatbot Infobip (Ya implementado)
1. Cliente interactúa con chatbot de Infobip
2. Acepta términos y condiciones
3. Chatbot envía datos de entrega a la API
4. API crea la entrega y devuelve **token** al cliente

### Fase 2: Asignación de Dispensers (Interno)
1. En la empresa se asignan los dispensers a la entrega
2. Se vinculan mediante el endpoint `/api/v1/deliveries/{id}/dispensers`

### Fase 3: App Móvil del Repartidor (NUEVO)
1. **Validar Token del Cliente**
   - Endpoint: `POST /api/v1/mobile/validate-token`
   - Repartidor pide el token al cliente
   - App valida el token y muestra información de la entrega

2. **Escanear Dispensers**
   - Endpoint: `POST /api/v1/mobile/validate-dispenser`
   - Repartidor escanea cada código de dispenser
   - App valida que pertenezca a la entrega

3. **Completar Entrega**
   - Endpoint: `POST /api/v1/mobile/complete-delivery`
   - Una vez escaneados todos los dispensers
   - App completa la entrega

### Fase 4: Procesamiento Asíncrono con RabbitMQ (NUEVO)
1. Al completar la entrega, se publica mensaje a RabbitMQ
2. Worker consume el mensaje de la cola `q.workorder.generate`
3. Worker ejecuta:
   - Crea la orden de trabajo (WorkOrder)
   - Genera PDF con los detalles
   - Envía email al cliente
   - (Futuro) Guarda en storage

## 🔌 Endpoints Mobile

### 1. Validar Token
```http
POST /api/v1/mobile/validate-token
Content-Type: application/json

{
  "token": "ABC123"
}
```

**Respuesta exitosa:**
```json
{
  "valid": true,
  "message": "Token válido",
  "delivery": {
    "id": 1,
    "nro_cta": "12345",
    "nro_rto": "9",
    "cantidad": 2,
    "tipo_entrega": "Instalacion",
    "fecha_accion": "2025-11-12"
  },
  "dispensers": [
    {
      "id": 1,
      "marca": "LAMO",
      "nro_serie": "LM123456789",
      "tipo": "P",
      "validated": false
    },
    {
      "id": 2,
      "marca": "LAMO",
      "nro_serie": "LM987654321",
      "tipo": "M",
      "validated": false
    }
  ]
}
```

### 2. Validar Dispenser
```http
POST /api/v1/mobile/validate-dispenser
Content-Type: application/json

{
  "delivery_id": 1,
  "nro_serie": "LM123456789"
}
```

**Respuesta:**
```json
{
  "valid": true,
  "message": "Dispenser válido",
  "dispenser": {
    "id": 1,
    "marca": "LAMO",
    "nro_serie": "LM123456789",
    "tipo": "P",
    "validated": true
  }
}
```

### 3. Completar Entrega
```http
POST /api/v1/mobile/complete-delivery
Content-Type: application/json

{
  "delivery_id": 1,
  "token": "ABC123",
  "validated_dispensers": ["LM123456789", "LM987654321"]
}
```

**Respuesta:**
```json
{
  "success": true,
  "message": "Entrega completada exitosamente",
  "delivery_id": 1,
  "work_order_queued": true
}
```

## 🐰 Configuración RabbitMQ

### Variables de Entorno
Agregar al archivo `.env`:

```env
# RabbitMQ Configuration
RABBITMQ_HOST=192.168.0.250
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_QUEUE=q.workorder.generate
```

### Mensaje Publicado a RabbitMQ
Estructura del mensaje enviado a la cola:

```json
{
  "nroCta": "12345",
  "name": "Gonzalo Wiñazki",
  "address": "Santiago de Liniers 3118",
  "locality": "Ciudadela",
  "nroRto": "9",
  "createdAt": "2025-11-12",
  "tipoAccion": "Instalacion",
  "token": "ABC123",
  "dispensers": [
    {
      "marca": "LAMO",
      "nro_serie": "LM123456789"
    },
    {
      "marca": "LAMO",
      "nro_serie": "LM987654321"
    }
  ],
  "deliveryId": 1
}
```

## 🏗️ Arquitectura

```
┌─────────────┐
│ App Móvil   │
│ (Repartidor)│
└──────┬──────┘
       │ 1. validate-token
       │ 2. validate-dispenser (x N)
       │ 3. complete-delivery
       ▼
┌─────────────────┐
│   API REST      │
│   (GoFrioCalor) │
└────────┬────────┘
         │ Actualiza delivery
         │ Publica mensaje
         ▼
┌─────────────────┐
│   RabbitMQ      │
│ q.workorder.gen │
└────────┬────────┘
         │ Consume mensaje
         ▼
┌─────────────────┐
│  Worker         │
│  (Consumer)     │
└────────┬────────┘
         │
         ├─► 1. Crea WorkOrder
         ├─► 2. Genera PDF
         ├─► 3. Envía Email
         └─► 4. (Futuro) Guarda en Storage
```

## ✨ Ventajas de la Arquitectura

1. **Desacoplamiento**: La app móvil no espera que se procese todo
2. **Resiliencia**: Si algo falla, el mensaje persiste en RabbitMQ
3. **Escalabilidad**: Múltiples workers pueden procesar en paralelo
4. **Asincronía**: No afecta la experiencia del repartidor
5. **Observabilidad**: Logs estructurados en cada paso

## 🔄 Estados del Delivery

- `Pendiente`: Creado desde Infobip, esperando entrega
- `Completado`: Repartidor completó la entrega
- `Cancelado`: (Futuro) Entrega cancelada

## 📋 TODO: Mejoras Futuras

1. [ ] Agregar campos `name`, `address`, `locality` al modelo Delivery
2. [ ] Implementar servicio de email real (SMTP)
3. [ ] Implementar generación real de PDF
4. [ ] Implementar storage para guardar PDFs
5. [ ] Agregar Dead Letter Queue (DLQ) para mensajes fallidos
6. [ ] Implementar reintentos con backoff exponencial
7. [ ] Agregar métricas y monitoring (Prometheus, Grafana)
8. [ ] Implementar idempotencia con message ID
9. [ ] Agregar endpoint para consultar estado de work order
10. [ ] Implementar notificación push cuando se complete la work order

## 🧪 Testing

Ver archivo `tests/test_mobile_flow.ps1` para ejemplos de prueba del flujo completo.
