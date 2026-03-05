# Test Audit System - Prueba del sistema de auditoría
# Ejecutar: .\tests\test_audit_system.ps1

Write-Host "=== Test Sistema de Auditoría ===" -ForegroundColor Cyan
Write-Host ""

$BASE_URL = "http://localhost:8080/dispenser-operations"
$AUTH_API_KEY = "clave_compartidas"  # Cambiar por tu API key real

# Colores para output
function Write-Success { Write-Host $args -ForegroundColor Green }
function Write-Info { Write-Host $args -ForegroundColor Cyan }
function Write-Error { Write-Host $args -ForegroundColor Red }
function Write-Step { Write-Host "`n--- $args ---" -ForegroundColor Yellow }

# Paso 1: Obtener Token JWT
Write-Step "Paso 1: Obtener Token JWT"
Write-Info "POST $BASE_URL/auth/generar-token"

$authBody = @{
    api_key = $AUTH_API_KEY
} | ConvertTo-Json

try {
    $authResponse = Invoke-RestMethod -Uri "$BASE_URL/auth/generar-token" -Method Post -Body $authBody -ContentType "application/json"
    $TOKEN = $authResponse.token
    Write-Success "✓ Token obtenido exitosamente"
    Write-Info "Token: $($TOKEN.Substring(0, 50))..."
} catch {
    Write-Error "✗ Error al obtener token: $($_.Exception.Message)"
    exit 1
}

$headers = @{
    "Authorization" = "Bearer $TOKEN"
    "Content-Type" = "application/json"
}

# Paso 2: Consultar eventos recientes
Write-Step "Paso 2: Consultar Eventos Recientes"
Write-Info "GET $BASE_URL/api/v1/audit/recent?limit=10"

try {
    $recentEvents = Invoke-RestMethod -Uri "$BASE_URL/api/v1/audit/recent?limit=10" -Method Get -Headers $headers
    Write-Success "✓ Eventos recientes obtenidos"
    Write-Info "Total de eventos: $($recentEvents.Count)"
    
    if ($recentEvents.Count -gt 0) {
        Write-Info "Último evento:"
        $lastEvent = $recentEvents[0]
        Write-Info "  - ID: $($lastEvent.id)"
        Write-Info "  - Acción: $($lastEvent.action)"
        Write-Info "  - Entidad: $($lastEvent.entity_type) - $($lastEvent.entity_id)"
        Write-Info "  - Actor: $($lastEvent.actor_type) - $($lastEvent.actor_id)"
        Write-Info "  - Fecha: $($lastEvent.occurred_at)"
    }
} catch {
    Write-Error "✗ Error al consultar eventos recientes: $($_.Exception.Message)"
}

# Paso 3: Buscar por entidad específica
Write-Step "Paso 3: Buscar por Entidad"
Write-Info "GET $BASE_URL/api/v1/audit/entity/delivery/123"

try {
    $entityEvents = Invoke-RestMethod -Uri "$BASE_URL/api/v1/audit/entity/delivery/123" -Method Get -Headers $headers
    Write-Success "✓ Eventos de entidad obtenidos"
    Write-Info "Total de eventos para delivery 123: $($entityEvents.Count)"
} catch {
    # Puede que no exista, es normal
    Write-Info "ℹ Entidad 'delivery/123' no encontrada (normal si no hay datos de prueba)"
}

# Paso 4: Buscar por request ID
Write-Step "Paso 4: Buscar por Request ID"
if ($recentEvents.Count -gt 0 -and $recentEvents[0].request_id) {
    $requestId = $recentEvents[0].request_id
    Write-Info "GET $BASE_URL/api/v1/audit/request/$requestId"
    
    try {
        $requestEvents = Invoke-RestMethod -Uri "$BASE_URL/api/v1/audit/request/$requestId" -Method Get -Headers $headers
        Write-Success "✓ Eventos del request obtenidos"
        Write-Info "Total de eventos en request $requestId : $($requestEvents.Count)"
    } catch {
        Write-Error "✗ Error al buscar por request ID: $($_.Exception.Message)"
    }
} else {
    Write-Info "ℹ No hay request_id disponible en eventos recientes"
}

# Paso 5: Búsqueda avanzada con filtros
Write-Step "Paso 5: Búsqueda Avanzada"
Write-Info "POST $BASE_URL/api/v1/audit/search"

$searchBody = @{
    entity_type = "delivery"
    action = "CREATED"
    limit = 5
    offset = 0
} | ConvertTo-Json

try {
    $searchResults = Invoke-RestMethod -Uri "$BASE_URL/api/v1/audit/search" -Method Post -Body $searchBody -Headers $headers
    Write-Success "✓ Búsqueda completada"
    Write-Info "Eventos encontrados (CREATED deliveries): $($searchResults.Count)"
    
    if ($searchResults.Count -gt 0) {
        Write-Info "Primer resultado:"
        $firstResult = $searchResults[0]
        Write-Info "  - ID: $($firstResult.id)"
        Write-Info "  - Entity: $($firstResult.entity_type)/$($firstResult.entity_id)"
        Write-Info "  - Fecha: $($firstResult.occurred_at)"
    }
} catch {
    Write-Info "ℹ No se encontraron eventos con los filtros especificados"
}

# Paso 6: Estadísticas
Write-Step "Paso 6: Estadísticas de Auditoría"
Write-Info "GET $BASE_URL/api/v1/audit/stats?days=7"

try {
    $stats = Invoke-RestMethod -Uri "$BASE_URL/api/v1/audit/stats?days=7" -Method Get -Headers $headers
    Write-Success "✓ Estadísticas obtenidas"
    Write-Info "Estadísticas de los últimos 7 días:"
    Write-Info "  - Total de eventos: $($stats.total_events)"
    Write-Info "  - Por acción:"
    $stats.by_action | ForEach-Object {
        Write-Info "    $($_.action): $($_.count)"
    }
    Write-Info "  - Por entidad:"
    $stats.by_entity_type | ForEach-Object {
        Write-Info "    $($_.entity_type): $($_.count)"
    }
} catch {
    Write-Error "✗ Error al obtener estadísticas: $($_.Exception.Message)"
}

# Paso 7: Crear un evento de auditoría manualmente (para testing)
Write-Step "Paso 7: Generar Evento de Prueba"
Write-Info "Generando nueva entrega para crear evento de auditoría..."

$deliveryBody = @{
    codigo = "TEST-AUDIT-$(Get-Random -Maximum 99999)"
    tipo = "INSTALACION"
    dispensers = @(
        @{
            modelo_dispenser = "DFB_10"
            cantidad = 1
        }
    )
    cliente = @{
        dni_cuit = "12345678"
        nombre = "Test Audit"
        telefono = "1234567890"
        direccion = "Test Address"
        localidad = "Test City"
        provincia = "Test Province"
    }
    observaciones = "Test para sistema de auditoría"
} | ConvertTo-Json -Depth 10

try {
    $newDelivery = Invoke-RestMethod -Uri "$BASE_URL/api/v1/entregas" -Method Post -Body $deliveryBody -Headers $headers
    Write-Success "✓ Nueva entrega creada exitosamente"
    Write-Info "ID de entrega: $($newDelivery.id)"
    Write-Info "Código: $($newDelivery.codigo)"
    
    # Esperar un momento y buscar el evento de auditoría
    Start-Sleep -Seconds 1
    Write-Info "`nBuscando evento de auditoría generado..."
    
    $auditForNewDelivery = Invoke-RestMethod -Uri "$BASE_URL/api/v1/audit/entity/delivery/$($newDelivery.id)" -Method Get -Headers $headers
    if ($auditForNewDelivery.Count -gt 0) {
        Write-Success "✓ Evento de auditoría encontrado"
        Write-Info "  - Acción: $($auditForNewDelivery[0].action)"
        Write-Info "  - Fecha: $($auditForNewDelivery[0].occurred_at)"
        Write-Info "  - Request ID: $($auditForNewDelivery[0].request_id)"
    } else {
        Write-Info "ℹ Aún no se generó el evento de auditoría (puede que no esté integrado)"
    }
} catch {
    Write-Error "✗ Error al crear entrega de prueba: $($_.Exception.Message)"
}

Write-Step "Resumen"
Write-Success "✓ Tests del sistema de auditoría completados"
Write-Info "`nEndpoints disponibles:"
Write-Info "  GET  /audit/recent           - Eventos recientes"
Write-Info "  GET  /audit/entity/:type/:id - Historia de una entidad"
Write-Info "  GET  /audit/actor/:type/:id  - Acciones de un actor"
Write-Info "  GET  /audit/request/:id      - Traza de un request"
Write-Info "  POST /audit/search           - Búsqueda avanzada"
Write-Info "  GET  /audit/stats            - Estadísticas"
Write-Info "`nNota: Para que se generen eventos automáticamente, se debe integrar"
Write-Info "      auditService.LogEventAsync() en los handlers existentes"
Write-Host ""
