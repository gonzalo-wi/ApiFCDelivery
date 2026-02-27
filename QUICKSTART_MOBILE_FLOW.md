# 🚀 Quick Start - Mobile Delivery Flow

Guía rápida para usar los endpoints móviles con Postman.

---

## 📋 Pre-requisitos
- ✅ Servidor corriendo: `.\start_server.ps1`
- ✅ RabbitMQ conectado en `192.168.0.250:5672`
- ✅ Base de datos MySQL en `192.168.0.227:3306`

---

## 📥 Paso 1: Importar Colección

1. Abre Postman
2. Click en **Import**
3. Selecciona: `postman/Mobile_Delivery_Flow.postman_collection.json`
4. Click **Import**

---

## 🎯 Paso 2: Ejecutar Flujo Completo

### A. Setup Inicial (ejecutar UNA VEZ)

```
1. Crear Delivery
   → Guarda automáticamente delivery_id y token

2. Agregar Dispenser 1
   → LM123456789 (Tipo: P)

3. Agregar Dispenser 2
   → LM987654321 (Tipo: M)
```

### B. Flujo Mobile (ejecutar EN ORDEN)

```
1. Validar Token
   Body: { 
     "token": "{{token}}", 
     "nro_cta": "12345",
     "fecha_accion": "2025-11-12"
   }
   → Retorna: delivery info + lista de dispensers

2. Validar Dispenser 1
   Body: { "delivery_id": {{delivery_id}}, "nro_serie": "LM123456789" }
   → Retorna: { "valid": true, "dispenser": {...} }

3. Validar Dispenser 2
   Body: { "delivery_id": {{delivery_id}}, "nro_serie": "LM987654321" }
   → Retorna: { "valid": true, "dispenser": {...} }

4. Completar Entrega
   Body: {
     "delivery_id": {{delivery_id}},
     "token": "{{token}}",
     "validated_dispensers": ["LM123456789", "LM987654321"]
   }
   → Retorna: { "work_order_queued": true, "message": "..." }
```

---

## ✅ Paso 3: Verificar Resultado

### En la respuesta de "Completar Entrega":
```json
{
  "delivery_id": 24,
  "status": "Completado",
  "work_order_queued": true,
  "message": "Entrega completada exitosamente. La orden de trabajo será procesada."
}
```

### En los logs del servidor:
```
✅ Message published to RabbitMQ
✅ Consumer received message
✅ WorkOrder created: OT-000020
✅ PDF generated: /tmp/work_order_OT-000020.pdf
✅ Email sent to: cliente@example.com
```

---

## 🔄 Variables Automáticas

Las siguientes variables se configuran automáticamente:

| Variable | Origen | Usado en |
|----------|--------|----------|
| `delivery_id` | Crear Delivery response | Todos los endpoints mobile |
| `token` | Crear Delivery response | Validar Token, Completar Entrega |
| `base_url` | Collection variable | Todos los endpoints |

**Nota:** `nro_cta` y `fecha_accion` deben coincidir con los valores usados al crear el Delivery.

**No necesitas copiar/pegar nada manualmente** - Postman lo hace por ti.

---

## 🛠️ Troubleshooting

### Error: "Token not found"
- Ejecuta "Crear Delivery" primero
- Verifica que la variable `{{token}}` tenga valor

### Error: "Dispenser no encontrado"
- Ejecuta "Agregar Dispenser 1" y "Agregar Dispenser 2" primero
- Verifica que uses los mismos números de serie

### Error: "Delivery not found"
- Ejecuta "Crear Delivery" primero
- Verifica que la variable `{{delivery_id}}` tenga valor

### Error: "Connection refused"
- Verifica que el servidor esté corriendo: `.\start_server.ps1`
- Verifica que esté en puerto 8080: `http://localhost:8080/health`

### RabbitMQ no conecta
- Verifica conexión: Abre `http://192.168.0.250:15672`
- Usuario: `admin-` / Password: `admin123`
- Verifica que el queue `q.workorder.generate` exista

---

## 📊 Flujo Completo Resumido

```
[Postman] → POST /deliveries
            ↓ (guarda delivery_id, token)
[Postman] → POST /dispensers (x2)
            ↓
[Mobile]  → POST /mobile/validate-token
            ↓
[Mobile]  → POST /mobile/validate-dispenser (x2)
            ↓
[Mobile]  → POST /mobile/complete-delivery
            ↓ (publica mensaje)
[RabbitMQ Queue: q.workorder.generate]
            ↓ (consume mensaje)
[Worker]  → Crea WorkOrder
            ↓
[Worker]  → Genera PDF
            ↓
[Worker]  → Envía Email
            ✅ DONE
```

---

## 📚 Más Documentación

- [Postman README](postman/README.md) - Instrucciones detalladas
- [Postman Mobile Endpoints](docs/POSTMAN_MOBILE_ENDPOINTS.md) - Ejemplos de requests/responses
- [Terms Integration](docs/TERMS_INTEGRATION.md) - Integración con términos y condiciones
- [Deployment Guide](DEPLOYMENT_GUIDE.md) - Guía de despliegue

---

## 💡 Tips

1. **Usa el Runner de Postman** para ejecutar toda la secuencia automáticamente
2. **Cambia los números de serie** en cada prueba para evitar duplicados
3. **Revisa los logs** del servidor para debugging en tiempo real
4. **Health Check**: Usa `GET /health` para verificar que el servidor funciona

---

## 🎉 ¡Listo!

Ahora tienes un flujo completo de trabajo con:
- ✅ API REST funcional
- ✅ Validación de tokens
- ✅ Escaneo de dispensers
- ✅ RabbitMQ async processing
- ✅ Generación automática de WorkOrders
- ✅ PDF y Email (mock por ahora)

**Próximos pasos sugeridos:**
1. Implementar generación real de PDF
2. Integrar servicio de email real (SMTP/SendGrid)
3. Agregar campos de cliente (nombre, dirección, localidad) al modelo Delivery
4. Implementar storage para PDFs (S3/Azure Blob)
