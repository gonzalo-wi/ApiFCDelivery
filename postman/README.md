# Colección Postman - Infobip Delivery API

## 📥 Importar a Postman

### Opción 1: Importar archivo JSON
1. Abre Postman
2. Click en **Import** (botón superior izquierdo)
3. Selecciona el archivo: `Infobip_Delivery_Collection.json`
4. Click en **Import**

### Opción 2: Arrastrar y soltar
1. Abre Postman
2. Arrastra el archivo `Infobip_Delivery_Collection.json` a la ventana de Postman
3. Se importará automáticamente

---

## 🚀 Uso Rápido

### 1. Asegúrate de que el servidor esté corriendo
```bash
go run api/cmd/main.go
```

El servidor debe estar en: `http://localhost:8080`

### 2. Abre la colección en Postman
- Busca "GoFrioCalor - Infobip Delivery API" en el panel izquierdo

### 3. Ejecuta los ejemplos
La colección incluye **9 ejemplos** organizados en 2 carpetas:

#### ✅ Casos Exitosos (5 ejemplos)
- Instalación - 2 Pie + 1 Mesada
- Recambio - 1 Pie + 1 Mesada
- Retiro - 3 Mesada
- Solo Pie - 3 Dispensers
- Solo Mesada - 2 Dispensers

#### ❌ Casos de Error (4 ejemplos)
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
