# 📱 Mobile Delivery + RabbitMQ - Guía de Implementación

## 🎯 Objetivo
Sistema completo para gestionar entregas con app móvil y procesamiento asíncrono de órdenes de trabajo mediante RabbitMQ.

## 🏗️ Arquitectura Implementada

```
Chatbot Infobip → API (Crea Delivery) → Empresa asigna Dispensers
                                               ↓
                                         App Móvil Repartidor
                                         ├─ Valida Token
                                         ├─ Escanea Dispensers
                                         └─ Completa Entrega
                                               ↓
                                         API publica a RabbitMQ
                                               ↓
                                         Worker Consumer procesa
                                         ├─ Crea Work Order
                                         ├─ Genera PDF
                                         ├─ Envía Email
                                         └─ (Futuro) Guarda en Storage
```

## 📂 Archivos Creados/Modificados

### Configuración
- ✅ `config/rabbitmq.go` - Configuración y conexión RabbitMQ
- ✅ `.env.example` - Variables de entorno agregadas

### DTOs
- ✅ `internal/dto/mobile_delivery_dto.go` - DTOs para app móvil
- ✅ `internal/dto/work_order_message_dto.go` - Mensaje RabbitMQ

### Servicios
- ✅ `internal/service/rabbitmq_publisher.go` - Publisher de mensajes
- ✅ `internal/service/mobile_delivery_service.go` - Lógica de validaciones
- ✅ `internal/service/work_order_consumer.go` - Worker consumer
- ✅ `internal/service/email_service.go` - Servicio de email (mock)

### Transport/Handlers
- ✅ `internal/transport/mobile_delivery_handler.go` - Endpoints móviles

### Rutas
- ✅ `internal/routes/mobile_routes.go` - Rutas para app móvil
- ✅ `internal/routes/router.go` - Integración de rutas móviles

### Main
- ✅ `api/cmd/main.go` - Inicialización de RabbitMQ y worker

### Documentación y Tests
- ✅ `docs/MOBILE_DELIVERY_FLOW.md` - Documentación completa del flujo
- ✅ `tests/test_mobile_flow.ps1` - Script de prueba end-to-end

## 🚀 Instalación y Configuración

### 1. Instalar Dependencias
```bash
go get github.com/rabbitmq/amqp091-go
go mod tidy
```

### 2. Configurar Variables de Entorno
Copia `.env.example` a `.env` y configura:

```env
RABBITMQ_HOST=192.168.0.250
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_QUEUE=q.workorder.generate
```

### 3. Verificar RabbitMQ
Asegúrate de que RabbitMQ esté corriendo:
- URL: http://192.168.0.250:15672
- Usuario: guest / Password: guest
- Cola debe existir: `q.workorder.generate`

### 4. Ejecutar la Aplicación
```bash
cd api/cmd
go run main.go
```

## 📡 Endpoints de la App Móvil

Base URL: `http://localhost:8080/api/v1/mobile`

### 1️⃣ Validar Token del Cliente
```http
POST /api/v1/mobile/validate-token
Content-Type: application/json

{
  "token": "ABC123"
}
```

### 2️⃣ Validar Dispenser Escaneado
```http
POST /api/v1/mobile/validate-dispenser
Content-Type: application/json

{
  "delivery_id": 1,
  "nro_serie": "LM123456789"
}
```

### 3️⃣ Completar Entrega
```http
POST /api/v1/mobile/complete-delivery
Content-Type: application/json

{
  "delivery_id": 1,
  "token": "ABC123",
  "validated_dispensers": ["LM123456789", "LM987654321"]
}
```

## 🧪 Pruebas

### Ejecutar Script de Prueba Completo
```powershell
cd tests
.\test_mobile_flow.ps1
```

Este script:
1. Crea un delivery de prueba
2. Agrega dispensers
3. Valida el token
4. Escanea todos los dispensers
5. Completa la entrega
6. Muestra el resultado del procesamiento

### Verificar en RabbitMQ
1. Acceder a http://192.168.0.250:15672
2. Ir a Queues → `q.workorder.generate`
3. Ver mensajes publicados y consumidos

## 📊 Monitoreo y Logs

Los logs estructurados muestran:
- 📨 Publicación de mensajes a RabbitMQ
- 📥 Consumo de mensajes por el worker
- ✅ Creación de órdenes de trabajo
- 📄 Generación de PDFs (mock)
- 📧 Envío de emails (mock)

Ejemplo de log:
```
INFO RabbitMQ Publisher connected successfully queue=q.workorder.generate host=192.168.0.250
INFO Work Order Consumer started. Waiting for messages...
INFO Processing work order message delivery_id=1
INFO Work order created order_number=OT-000001
INFO Work order workflow completed
```

## 🔧 Configuración de Producción

### RabbitMQ
- [ ] Habilitar autenticación robusta
- [ ] Configurar SSL/TLS
- [ ] Implementar Dead Letter Queue
- [ ] Configurar monitoring con Prometheus

### Servicios
- [ ] Implementar envío real de emails (SMTP/SendGrid)
- [ ] Implementar generación real de PDFs
- [ ] Implementar storage en S3/Azure Blob
- [ ] Agregar retry logic con backoff exponencial
- [ ] Implementar circuit breaker

### Base de Datos
- [ ] Agregar campos al modelo Delivery:
  - `client_name`
  - `client_address`
  - `client_locality`
  - `client_email`
  - `work_order_id` (FK a WorkOrder)

## 🐛 Troubleshooting

### RabbitMQ no conecta
- Verificar que el servidor esté corriendo
- Verificar credenciales en `.env`
- Verificar que la cola exista

### Worker no procesa mensajes
- Verificar logs del servidor
- Verificar que el consumer esté iniciado
- Verificar QoS de RabbitMQ

### Mensajes se quedan en la cola
- Verificar errores en logs del worker
- Implementar DLQ para mensajes fallidos
- Verificar que el ACK se esté enviando

## 📈 Métricas Sugeridas

- Tiempo promedio de procesamiento de mensajes
- Tasa de éxito/fallo de work orders
- Cantidad de mensajes en cola
- Tiempo de respuesta de endpoints móviles

## 🔐 Seguridad

### Mejoras Recomendadas
- [ ] Agregar autenticación JWT para endpoints móviles
- [ ] Validar permisos de repartidor
- [ ] Encriptar tokens sensibles
- [ ] Rate limiting en endpoints públicos
- [ ] Audit log de operaciones críticas

## 📚 Referencias

- [Documentación RabbitMQ](https://www.rabbitmq.com/documentation.html)
- [amqp091-go GitHub](https://github.com/rabbitmq/amqp091-go)
- [Flujo Completo](docs/MOBILE_DELIVERY_FLOW.md)

---

**¿Preguntas?** Consulta la documentación completa en `docs/MOBILE_DELIVERY_FLOW.md`
