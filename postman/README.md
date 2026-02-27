# Colecciones Postman - GoFrioCalor API

Este directorio contiene dos colecciones de Postman:

1. **Mobile_Delivery_Flow.postman_collection.json** - Flujo completo de entregas móviles con RabbitMQ
2. **Infobip_Delivery_Collection.json** - API de integración con Infobip (sistema externo)

---

## 📱 Mobile Delivery Flow (NUEVO)

### 📥 Importar a Postman
1. Abre Postman
2. Click en **Import** (botón superior izquierdo)
3. Selecciona el archivo: `Mobile_Delivery_Flow.postman_collection.json`
4. Click en **Import**

### 🚀 Uso Completo

#### 1. Inicia el servidor con RabbitMQ
```powershell
.\start_server.ps1
```

#### 2. Ejecuta los requests en orden:

**Paso 0: Setup (ejecutar una sola vez)**
- ✅ **Crear Delivery** - Crea un delivery y guarda automáticamente el ID y Token
- Agregar Dispenser 1
- Agregar Dispenser 2

**Pasos 1-3: Flujo Mobile (secuencial)**
1. **Validar Token** - El cliente proporciona el token
2. **Validar Dispenser 1** - Escanear primer dispenser
3. **Validar Dispenser 2** - Escanear segundo dispenser
4. **Completar Entrega** - Finaliza y envía mensaje a RabbitMQ

#### 3. Verifica el resultado
- El endpoint "Completar Entrega" responde con `work_order_queued: true`
- Revisa los logs del servidor para ver:
  - Message published to RabbitMQ
  - Consumer processing message
  - WorkOrder created: OT-XXXXXX
  - PDF generated
  - Email sent

### 🔄 Variables Automáticas
La colección usa variables que se configuran automáticamente:
- `delivery_id` - Se obtiene al crear el delivery
- `token` - Se obtiene al crear el delivery
- `base_url` - Por defecto: `http://localhost:8080/api/v1`

### ⚡ Flujo Rápido
1. Ejecuta "Crear Delivery" una vez
2. Ejecuta "Agregar Dispenser 1" y "Agregar Dispenser 2"
3. Ahora puedes ejecutar los 4 endpoints mobile en secuencia
4. ¡Listo! La WorkOrder se crea automáticamente en background

---

## 📨 Infobip Delivery API

### 📥 Importar a Postman
1. Abre Postman
2. Click en **Import** (botón superior izquierdo)
3. Selecciona el archivo: `Infobip_Delivery_Collection.json`
4. Click en **Import**

### 🚀 Uso Rápido

#### 1. Asegúrate de que el servidor esté corriendo
```bash
go run api/cmd/main.go
```

#### 2. Ejecuta los ejemplos
La colección incluye **9 ejemplos** organizados en 2 carpetas:

**✅ Casos Exitosos (5 ejemplos)**
- Instalación - 2 Pie + 1 Mesada
- Recambio - 1 Pie + 1 Mesada
- Retiro - 3 Mesada
- Solo Pie - 3 Dispensers
- Solo Mesada - 2 Dispensers

**❌ Casos de Error (4 ejemplos)**
- Error - Sin Dispensers (0 total)
- Error - Campo Faltante (nro_cta)
- Error - Tipo Entrega Inválido
- Error - Entregado Por Inválido

---

## ⚙️ Configuración

### Variable de Colección: baseUrl
Por defecto: `http://localhost:8080`

Para cambiarla:
1. Click derecho en la colección
2. **Edit**
3. Tab **Variables**
4. Cambiar `baseUrl` a tu servidor (ej: `https://api.gofricalor.com`)

---

## 📋 Respuestas Esperadas

### Ejemplo Exitoso (201 Created)
```json
{
  "token": "1234",
  "message": "Entrega creada exitosamente"
}
```

### Ejemplo de Error (400 Bad Request)
```json
{
  "error": "validation_failed",
  "message": "Debe especificar al menos un dispenser (P o M)"
}
```

---

## 🧪 Agregar Tests Automáticos

Puedes agregar estos scripts en la pestaña **Tests** de cada request:

```javascript
// Validar status code exitoso
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
    pm.expect(jsonData.token).to.match(/^\d{4}$/);
});

// Validar mensaje de éxito
pm.test("Response has success message", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.message).to.eql("Entrega creada exitosamente");
});
```

---

## 📚 Documentación Adicional

- Ver ejemplos detallados: `docs/POSTMAN_EXAMPLES.md`
- Ver especificación de API: `docs/INFOBIP_DELIVERY_API.md`
- Ver ejemplos JSON: `docs/INFOBIP_JSON_EXAMPLES.md`

---

## 🔍 Tips

1. **Duplicar requests**: Click derecho → Duplicate para crear variaciones
2. **Formatear JSON**: `Ctrl + B` en el body para formatear
3. **Vista Pretty**: Tab "Pretty" en la respuesta para mejor visualización
4. **Copiar como cURL**: Click en "Code" → "cURL" para compartir con otros
5. **Ejecutar todo**: Click en la colección → "Run" para ejecutar todos los tests

---

## 🌐 Para Producción

Cambia la variable `baseUrl` de:
```
http://localhost:8080
```

A tu servidor de producción:
```
https://api.gofricalor.com
```

No olvides agregar autenticación si es necesario.
