# Ejemplos de Autenticación para Postman/cURL

## Colección Postman

### 1. Generar Token

**Request:**
```
GET http://localhost:8095/dispenser-operations/auth/generar-token
```

**Headers:**
```
x-api-key: MOBEUS_kG7pX2sV9nQ4aJ1cL8rT0yZ5wH3eU6mF2dC9bA1sR4xP
```

**Respuesta Exitosa:**
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJtb2JldXMiLCJleHAiOjE3NzI3MjI4MDAsImlhdCI6MTc3MjcyMTAwMH0.nXI1B0aWtRm4HgLD0sEylCXt_CCdI6YaRbj3yvPLnW4",
    "proveedor": "mobeus",
    "detail": "El token es válido por 30 minutos.",
    "expires_at_local": "2026-03-05T12:00:00-03:00"
}
```

### 2. Usar Token en Requests a GoFrioCalor

**Request:**
```
POST http://localhost:8095/dispenser-operations/api/v1/deliveries
```

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

**Body:**
```json
{
    "work_order_id": 12345,
    "delivery_status": "en_curso",
    "truck_id": 101,
    "client_name": "Cliente Test",
    "client_address": "Dirección Test"
}
```

---

## cURL - Linux/Mac

### 1. Obtener Token
```bash
curl -X GET "http://localhost:8095/dispenser-operations/auth/generar-token" \
  -H "x-api-key: MOBEUS_kG7pX2sV9nQ4aJ1cL8rT0yZ5wH3eU6mF2dC9bA1sR4xP"
```

### 2. Guardar Token en Variable
```bash
TOKEN=$(curl -s -X GET "http://localhost:8095/dispenser-operations/auth/generar-token" \
  -H "x-api-key: MOBEUS_kG7pX2sV9nQ4aJ1cL8rT0yZ5wH3eU6mF2dC9bA1sR4xP" | jq -r '.token')

echo "Token: $TOKEN"
```

### 3. Usar Token en Request
```bash
curl -X POST "http://localhost:8095/dispenser-operations/api/v1/deliveries" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "work_order_id": 12345,
    "delivery_status": "en_curso",
    "truck_id": 101
  }'
```

### 4. GET Request con Token
```bash
curl -X GET "http://localhost:8095/dispenser-operations/api/v1/work-orders" \
  -H "Authorization: Bearer $TOKEN"
```

---

## PowerShell

### 1. Obtener y Guardar Token
```powershell
$headers = @{
    "x-api-key" = "MOBEUS_kG7pX2sV9nQ4aJ1cL8rT0yZ5wH3eU6mF2dC9bA1sR4xP"
}

$tokenResponse = Invoke-RestMethod -Uri "http://localhost:8095/dispenser-operations/auth/generar-token" -Method GET -Headers $headers
$token = $tokenResponse.token
Write-Host "Token obtenido: $token"
```

### 2. POST Request con Token
```powershell
$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

$body = @{
    work_order_id = 12345
    delivery_status = "en_curso"
    truck_id = 101
    client_name = "Cliente Test"
    client_address = "Dirección Test"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:8095/dispenser-operations/api/v1/deliveries" -Method POST -Headers $headers -Body $body
$response | ConvertTo-Json -Depth 10
```

### 3. GET Request con Token
```powershell
$headers = @{
    "Authorization" = "Bearer $token"
}

$response = Invoke-RestMethod -Uri "http://localhost:8095/dispenser-operations/api/v1/work-orders" -Method GET -Headers $headers
$response | ConvertTo-Json -Depth 10
```

---

## JavaScript (Fetch API)

### 1. Obtener Token
```javascript
async function getToken() {
    const response = await fetch('http://localhost:8095/dispenser-operations/auth/generar-token', {
        method: 'GET',
        headers: {
            'x-api-key': 'MOBEUS_kG7pX2sV9nQ4aJ1cL8rT0yZ5wH3eU6mF2dC9bA1sR4xP'
        }
    });
    
    const data = await response.json();
    return data.token;
}
```

### 2. Usar Token en Request
```javascript
async function createDelivery(token, deliveryData) {
    const response = await fetch('http://localhost:8095/dispenser-operations/api/v1/deliveries', {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(deliveryData)
    });
    
    return await response.json();
}
```

### 3. Flujo Completo
```javascript
async function main() {
    try {
        // 1. Obtener token
        const token = await getToken();
        console.log('Token obtenido:', token);
        
        // 2. Crear delivery
        const deliveryData = {
            work_order_id: 12345,
            delivery_status: 'en_curso',
            truck_id: 101,
            client_name: 'Cliente Test'
        };
        
        const result = await createDelivery(token, deliveryData);
        console.log('Delivery creado:', result);
    } catch (error) {
        console.error('Error:', error);
    }
}

main();
```

---

## Python (Requests)

### 1. Obtener Token
```python
import requests

def get_token():
    url = "http://localhost:8095/dispenser-operations/auth/generar-token"
    headers = {
        "x-api-key": "MOBEUS_kG7pX2sV9nQ4aJ1cL8rT0yZ5wH3eU6mF2dC9bA1sR4xP"
    }
    
    response = requests.get(url, headers=headers)
    return response.json()['token']
```

### 2. Crear Delivery con Token
```python
def create_delivery(token, delivery_data):
    url = "http://localhost:8095/dispenser-operations/api/v1/deliveries"
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json"
    }
    
    response = requests.post(url, json=delivery_data, headers=headers)
    return response.json()
```

### 3. Flujo Completo
```python
def main():
    # 1. Obtener token
    token = get_token()
    print(f"Token obtenido: {token[:50]}...")
    
    # 2. Crear delivery
    delivery_data = {
        "work_order_id": 12345,
        "delivery_status": "en_curso",
        "truck_id": 101,
        "client_name": "Cliente Test"
    }
    
    result = create_delivery(token, delivery_data)
    print("Delivery creado:", result)

if __name__ == "__main__":
    main()
```

---

## Configuración en Postman

### Variables de Entorno

Crea un Environment en Postman con estas variables:

```
auth_url = http://localhost:8095/dispenser-operations/auth/generar-token
api_base_url = http://localhost:8095/dispenser-operations/api/v1
api_key = MOBEUS_kG7pX2sV9nQ4aJ1cL8rT0yZ5wH3eU6mF2dC9bA1sR4xP
token = (se llenará automáticamente)
```

### Pre-request Script (Para Auto-generar Token)

Agrega este script en la pestaña "Pre-request Script" de tu colección o carpeta:

```javascript
// Verificar si el token existe y no ha expirado
const token = pm.environment.get("token");
const tokenExpiry = pm.environment.get("token_expiry");
const now = new Date().getTime();

if (!token || !tokenExpiry || now >= tokenExpiry) {
    // Obtener nuevo token
    pm.sendRequest({
        url: pm.environment.get("auth_url"),
        method: 'GET',
        header: {
            'x-api-key': pm.environment.get("api_key")
        }
    }, function (err, response) {
        if (err) {
            console.error('Error obteniendo token:', err);
        } else {
            const jsonData = response.json();
            pm.environment.set("token", jsonData.token);
            
            // Calcular tiempo de expiración (30 minutos - 1 minuto de margen)
            const expiryTime = now + (29 * 60 * 1000);
            pm.environment.set("token_expiry", expiryTime);
            
            console.log('Token actualizado:', jsonData.token);
        }
    });
}
```

### Configuración de Authorization

En cada request de la API de GoFrioCalor:

1. Ir a la pestaña "Authorization"
2. Seleccionar Type: "Bearer Token"
3. Token: `{{token}}`

O manualmente en Headers:
```
Authorization: Bearer {{token}}
```

---

## Manejo de Errores Comunes

### Error 401: "Token no proporcionado"
```bash
# Verifica que incluyes el header Authorization
-H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Error 401: "Formato de token inválido"
```bash
# Asegúrate de usar el formato correcto:
Authorization: Bearer <token>
# NO uses:
Authorization: <token>
```

### Error 401: "Token inválido o expirado"
```bash
# El token expira en 30 minutos. Genera uno nuevo:
curl -X GET "http://192.168.0.55:8087/generar-token" \
  -H "x-api-key: YOUR_API_KEY"
```

### Error 401: "Servicio de autenticación no disponible"
```bash
# Verifica que el servicio de auth esté corriendo:
curl http://192.168.0.55:8087/health
```

---

## Testing Automatizado

### Script Bash Completo
```bash
#!/bin/bash

# Configuración
AUTH_URL="http://localhost:8095/dispenser-operations/auth/generar-token"
API_BASE="http://localhost:8095/dispenser-operations/api/v1"
API_KEY="MOBEUS_kG7pX2sV9nQ4aJ1cL8rT0yZ5wH3eU6mF2dC9bA1sR4xP"

# 1. Obtener token
echo "Obteniendo token..."
TOKEN=$(curl -s -X GET "$AUTH_URL" \
  -H "x-api-key: $API_KEY" | jq -r '.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
    echo "❌ Error: No se pudo obtener el token"
    exit 1
fi

echo "✓ Token obtenido: ${TOKEN:0:50}..."

# 2. Probar GET
echo "Probando GET..."
curl -X GET "$API_BASE/work-orders" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -w "\nStatus: %{http_code}\n"

# 3. Probar POST
echo "Probando POST..."
curl -X POST "$API_BASE/deliveries" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "work_order_id": 12345,
    "delivery_status": "en_curso"
  }' \
  -w "\nStatus: %{http_code}\n"

echo "✓ Pruebas completadas"
```
