# Endpoint Contact Center - Obtener Token

## Descripción
Endpoint público (sin autenticación) para que el contact center pueda obtener el token de un delivery especificando la fecha de acción y el número de cuenta.

## Información del Endpoint

### URL
```
GET /dispenser-operations/api/v1/deliveries/contact-center/token
```

### Autenticación
**NO requiere autenticación** - Endpoint público

### Parámetros Query

| Parámetro | Tipo | Obligatorio | Descripción | Ejemplo |
|-----------|------|-------------|-------------|---------|
| `fecha_accion` | string | Sí | Fecha de acción en formato YYYY-MM-DD | `2026-03-21` |
| `nro_cta` | string | Sí | Número de cuenta del cliente | `43534` |

## Respuestas

### Éxito (200 OK)
```json
{
    "id": 18,
    "fecha_accion": "2026-03-21",
    "nro_cta": "43534",
    "token": "2181"
}
```

### Error (400 Bad Request) - Parámetros faltantes
```json
{
    "error": "Parámetros 'fecha_accion' y 'nro_cta' son requeridos"
}
```

### Error (400 Bad Request) - Fecha inválida
```json
{
    "error": "Fecha inválida. Formato esperado: YYYY-MM-DD"
}
```

### Error (404 Not Found) - Delivery no encontrado
```json
{
    "error": "No se encontró delivery con los parámetros especificados"
}
```

### Error (500 Internal Server Error)
```json
{
    "error": "mensaje de error del servidor"
}
```

## Ejemplos de Uso

### cURL
```bash
curl -X GET "http://localhost:9090/dispenser-operations/api/v1/deliveries/contact-center/token?fecha_accion=2026-03-21&nro_cta=43534"
```

### PowerShell
```powershell
$fechaAccion = "2026-03-21"
$nroCta = "43534"
$url = "http://localhost:9090/dispenser-operations/api/v1/deliveries/contact-center/token?fecha_accion=$fechaAccion&nro_cta=$nroCta"

$response = Invoke-RestMethod -Uri $url -Method Get -ContentType "application/json"
Write-Host "Token: $($response.token)"
```

### JavaScript/Fetch
```javascript
const fechaAccion = '2026-03-21';
const nroCta = '43534';
const url = `http://localhost:9090/dispenser-operations/api/v1/deliveries/contact-center/token?fecha_accion=${fechaAccion}&nro_cta=${nroCta}`;

fetch(url)
    .then(response => response.json())
    .then(data => {
        console.log('Token:', data.token);
        console.log('ID:', data.id);
    })
    .catch(error => console.error('Error:', error));
```

### Python
```python
import requests

fecha_accion = "2026-03-21"
nro_cta = "43534"
url = f"http://localhost:9090/dispenser-operations/api/v1/deliveries/contact-center/token"

params = {
    "fecha_accion": fecha_accion,
    "nro_cta": nro_cta
}

response = requests.get(url, params=params)
data = response.json()

print(f"Token: {data['token']}")
print(f"ID: {data['id']}")
```

## Prueba con el Script de PowerShell

Ejecuta el script de prueba incluido:
```powershell
.\tests\test_contact_center_token.ps1
```

## Notas Importantes

1. **Sin Autenticación**: Este endpoint NO requiere el header `x-api-key` ni ningún tipo de autenticación
2. **Formato de Fecha**: Debe ser estrictamente YYYY-MM-DD (ejemplo: 2026-03-21)
3. **Parámetros Obligatorios**: Ambos parámetros (`fecha_accion` y `nro_cta`) son obligatorios
4. **Búsqueda Exacta**: Busca un delivery que coincida exactamente con la fecha y el número de cuenta
5. **Zona Horaria**: La búsqueda se realiza en UTC para evitar problemas de zona horaria

## Casos de Uso

### Contact Center
El personal del contact center puede usar este endpoint para:
- Consultar el token de un cliente cuando llama
- Verificar que el cliente tenga un delivery programado
- Proporcionar el token al cliente para que lo use en la app móvil

### Integración con Panel Web
Este endpoint puede integrarse en un panel web del contact center donde:
1. El operador ingresa la fecha de acción y el número de cuenta
2. El sistema consulta automáticamente el token
3. Se muestra el token en pantalla para comunicárselo al cliente

## Seguridad

Aunque es un endpoint público, solo devuelve información básica:
- ID del delivery
- Fecha de acción
- Número de cuenta
- Token

No expone información sensible como:
- Nombre del cliente
- Dirección
- Email
- Detalles de los dispensers

## Diagrama de Flujo

```
Cliente llama al Contact Center
         ↓
Operador solicita fecha y nro_cta
         ↓
Sistema consulta: GET /contact-center/token?fecha_accion=X&nro_cta=Y
         ↓
    ¿Encontrado?
    ↙         ↘
  Sí          No
  ↓           ↓
Devolver    Devolver
token       404
  ↓
Operador comunica
token al cliente
```
