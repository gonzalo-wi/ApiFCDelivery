# API de Entregas (Deliveries)

**Base URL:** `http://<host>:8095/dispenser-operations/api/v1`

> Los endpoints de **escritura** (POST, PUT, DELETE) requieren autenticación (`x-api-key` o `Authorization: Bearer <token>`). Los endpoints de **lectura** (GET) son públicos.

---

## Índice

- [Listar entregas](#listar-entregas)
- [Obtener entrega por ID](#obtener-entrega-por-id)
- [Crear entrega](#crear-entrega)
- [Modificar entrega](#modificar-entrega)
- [Eliminar entrega](#eliminar-entrega)
- [Buscar por RTO](#buscar-por-rto)
- [Buscar por cuenta](#buscar-por-cuenta)
- [Pendientes por cuenta (Infobip)](#pendientes-por-cuenta-infobip)
- [Modelos de datos](#modelos-de-datos)

---

## Listar entregas

**`GET /deliveries`**

Devuelve un listado paginado de entregas, con filtros opcionales.

### Query params

| Param | Tipo | Default | Descripción |
|---|---|---|---|
| `page` | `integer` | `1` | Número de página |
| `page_size` | `integer` | `20` | Items por página (máximo: 100) |
| `estado` | `string` | — | Filtrar por estado: `Pendiente`, `Completado`, `Cancelado` |
| `nro_cta` | `string` | — | Filtrar por número de cuenta |
| `fecha_accion` | `string` | — | Filtrar por fecha (`YYYY-MM-DD`) |

### Ejemplos

```
GET /deliveries
GET /deliveries?estado=Pendiente
GET /deliveries?estado=Completado&page=2&page_size=50
GET /deliveries?nro_cta=12345&estado=Pendiente
GET /deliveries?nro_cta=12345&fecha_accion=2026-05-14
```

### Respuesta `200 OK`

```json
{
  "data": [
    {
      "id": 1,
      "nro_cta": "12345",
      "nro_rto": "RTO-001",
      "name": "Juan Pérez",
      "email": "juan@example.com",
      "address": "Av. Siempre Viva 742",
      "locality": "Buenos Aires",
      "cantidad": 2,
      "token": "abc123xyz",
      "estado": "Pendiente",
      "tipo_entrega": "Instalacion",
      "entregado_por": "Tecnico",
      "order_number": "ORD-0001",
      "fecha_accion": "2026-05-14T00:00:00Z",
      "item_dispensers": [
        { "tipo": "P", "cantidad": 1 },
        { "tipo": "M", "cantidad": 1 }
      ]
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

---

## Obtener entrega por ID

**`GET /deliveries/:id`**

### Ejemplo

```
GET /deliveries/42
```

### Respuesta `200 OK`

```json
{
  "id": 42,
  "nro_cta": "12345",
  "nro_rto": "RTO-001",
  "name": "Juan Pérez",
  "email": "juan@example.com",
  "address": "Av. Siempre Viva 742",
  "locality": "Buenos Aires",
  "cantidad": 1,
  "token": "abc123xyz",
  "estado": "Pendiente",
  "tipo_entrega": "Instalacion",
  "entregado_por": "Tecnico",
  "order_number": "ORD-0042",
  "fecha_accion": "2026-05-14T00:00:00Z",
  "item_dispensers": [
    { "tipo": "P", "cantidad": 1 }
  ]
}
```

---

## Crear entrega

**🔒 `POST /deliveries`**

### Body `application/json`

```json
{
  "nro_cta": "12345",
  "nro_rto": "RTO-001",
  "name": "Juan Pérez",
  "email": "juan@example.com",
  "address": "Av. Siempre Viva 742",
  "locality": "Buenos Aires",
  "cantidad": 2,
  "estado": "Pendiente",
  "tipo_entrega": "Instalacion",
  "entregado_por": "Tecnico",
  "order_number": "ORD-0001",
  "fecha_accion": "2026-05-14",
  "item_dispensers": [
    { "tipo": "P", "cantidad": 1 },
    { "tipo": "M", "cantidad": 1 }
  ]
}
```

### Campos obligatorios

| Campo | Tipo | Valores aceptados |
|---|---|---|
| `nro_cta` | `string` | — |
| `nro_rto` | `string` | — |
| `cantidad` | `integer` | `1` – `3` |
| `estado` | `string` | `Pendiente`, `Completado`, `Cancelado` |
| `tipo_entrega` | `string` | `Instalacion`, `Retiro`, `Recambio`, `Service`, `Mixto` |
| `entregado_por` | `string` | `Repartidor`, `Tecnico` |

### Respuesta `201 Created`

Devuelve el objeto delivery creado (mismo formato que GET por ID).

---

## Modificar entrega

**🔒 `PUT /deliveries/:id`**

Reemplaza todos los campos del delivery. Enviar el mismo body que POST.

### Ejemplo

```
PUT /deliveries/42
```

### Body `application/json`

```json
{
  "nro_cta": "12345",
  "nro_rto": "RTO-001",
  "cantidad": 1,
  "estado": "Completado",
  "tipo_entrega": "Instalacion",
  "entregado_por": "Tecnico",
  "fecha_accion": "2026-05-14"
}
```

### Respuesta `200 OK`

Devuelve el delivery actualizado.

---

## Eliminar entrega

**🔒 `DELETE /deliveries/:id`**

### Ejemplo

```
DELETE /deliveries/42
```

### Respuesta `200 OK`

```json
{
  "message": "Delivery eliminado correctamente"
}
```

---

## Buscar por RTO

**`GET /deliveries/by-rto`**

| Param | Tipo | Descripción |
|---|---|---|
| `nro_rto` | `string` | Número de RTO |
| `fecha_accion` | `string` | Fecha opcional (`YYYY-MM-DD`) |

### Ejemplo

```
GET /deliveries/by-rto?nro_rto=RTO-001&fecha_accion=2026-05-14
```

### Respuesta `200 OK`

Array de deliveries (mismo formato que `data` del listado paginado).

---

## Buscar por cuenta

**`GET /deliveries/by-cta`**

| Param | Tipo | Descripción |
|---|---|---|
| `nro_cta` | `string` | Número de cuenta |
| `fecha_accion` | `string` | Fecha opcional (`YYYY-MM-DD`) |

### Ejemplo

```
GET /deliveries/by-cta?nro_cta=12345&fecha_accion=2026-05-14
```

### Respuesta `200 OK`

Array de deliveries.

---

## Pendientes por cuenta (Infobip)

**`GET /deliveries/infobip/pending`**

| Param | Tipo | Descripción |
|---|---|---|
| `nro_cta` | `string` | **Obligatorio.** Número de cuenta |

### Ejemplo

```
GET /deliveries/infobip/pending?nro_cta=12345
```

### Respuesta `200 OK`

```json
{
  "nro_cta": "12345",
  "has_pending": true,
  "count": 2,
  "deliveries": [
    {
      "delivery_id": 1,
      "nro_cta": "12345",
      "nro_rto": "RTO-001",
      "cantidad": 2,
      "tipo_entrega": "Instalacion",
      "fecha_accion": "2026-05-14",
      "token": "abc123xyz"
    }
  ]
}
```

---

## Modelos de datos

### EstadoEntrega

| Valor | Descripción |
|---|---|
| `Pendiente` | Entrega creada, aún no completada |
| `Completado` | Entrega realizada |
| `Cancelado` | Entrega cancelada |

### TipoEntrega

| Valor |
|---|
| `Instalacion` |
| `Retiro` |
| `Recambio` |
| `Service` |
| `Mixto` |

### EntregadoPor

| Valor |
|---|
| `Repartidor` |
| `Tecnico` |

### TipoDispenser

| Valor | Descripción |
|---|---|
| `P` | Dispenser tipo P (frío) |
| `M` | Dispenser tipo M (calor) |
