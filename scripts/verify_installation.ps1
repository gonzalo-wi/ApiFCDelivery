# Script de verificación del sistema de Términos y Condiciones
# Verifica que todos los componentes estén correctamente instalados

Write-Host "🔍 Verificación del Sistema de Términos y Condiciones" -ForegroundColor Cyan
Write-Host "=====================================================" -ForegroundColor Cyan
Write-Host ""

$errors = 0
$warnings = 0

# Verificar estructura de directorios
Write-Host "📁 Verificando estructura de archivos..." -ForegroundColor Yellow

$requiredFiles = @(
    "internal\models\terms_session.go",
    "internal\dto\terms_dto.go",
    "internal\store\terms_session_store.go",
    "internal\service\infobip_client.go",
    "internal\service\terms_session_service.go",
    "internal\transport\terms_session_handler.go",
    "internal\routes\terms_routes.go",
    "migrations\001_create_terms_sessions.sql",
    "docs\TERMS_INTEGRATION.md",
    "docs\TERMS_QUICKSTART.md",
    "docs\FRONTEND_INTEGRATION.md",
    "docs\FLOW_DIAGRAM.md",
    "scripts\test_terms_flow.ps1",
    ".env.example"
)

foreach ($file in $requiredFiles) {
    if (Test-Path $file) {
        Write-Host "  ✓ $file" -ForegroundColor Green
    } else {
        Write-Host "  ✗ $file (FALTANTE)" -ForegroundColor Red
        $errors++
    }
}

Write-Host ""

# Verificar archivo .env
Write-Host "⚙️  Verificando configuración..." -ForegroundColor Yellow

if (Test-Path ".env") {
    Write-Host "  ✓ Archivo .env existe" -ForegroundColor Green
    
    $envContent = Get-Content ".env" -Raw
    
    $requiredVars = @(
        "INFOBIP_BASE_URL",
        "INFOBIP_API_KEY",
        "APP_BASE_URL",
        "TERMS_TTL_HOURS"
    )
    
    foreach ($var in $requiredVars) {
        if ($envContent -match $var) {
            Write-Host "  ✓ Variable $var definida" -ForegroundColor Green
        } else {
            Write-Host "  ⚠ Variable $var no encontrada" -ForegroundColor Yellow
            $warnings++
        }
    }
} else {
    Write-Host "  ⚠ Archivo .env no existe (usar .env.example como base)" -ForegroundColor Yellow
    $warnings++
}

Write-Host ""

# Verificar que Go esté instalado
Write-Host "🔧 Verificando herramientas..." -ForegroundColor Yellow

try {
    $goVersion = go version 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  ✓ Go instalado: $goVersion" -ForegroundColor Green
    } else {
        Write-Host "  ✗ Go no está instalado o no está en PATH" -ForegroundColor Red
        $errors++
    }
} catch {
    Write-Host "  ✗ Go no está instalado o no está en PATH" -ForegroundColor Red
    $errors++
}

Write-Host ""

# Verificar módulos Go
Write-Host "📦 Verificando dependencias Go..." -ForegroundColor Yellow

if (Test-Path "go.mod") {
    Write-Host "  ✓ go.mod existe" -ForegroundColor Green
    
    $requiredPackages = @(
        "github.com/gin-gonic/gin",
        "gorm.io/gorm",
        "github.com/rs/zerolog"
    )
    
    $goModContent = Get-Content "go.mod" -Raw
    
    foreach ($package in $requiredPackages) {
        if ($goModContent -match [regex]::Escape($package)) {
            Write-Host "  ✓ $package" -ForegroundColor Green
        } else {
            Write-Host "  ⚠ $package no encontrado" -ForegroundColor Yellow
            $warnings++
        }
    }
} else {
    Write-Host "  ✗ go.mod no encontrado" -ForegroundColor Red
    $errors++
}

Write-Host ""

# Verificar sintaxis de archivos Go
Write-Host "🔍 Verificando sintaxis Go..." -ForegroundColor Yellow

$goFiles = Get-ChildItem -Path "internal" -Filter "*.go" -Recurse | Where-Object { $_.Name -like "*terms*" }

$syntaxErrors = 0
foreach ($file in $goFiles) {
    go vet $file.FullName 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) {
        Write-Host "  ⚠ Advertencias en $($file.Name)" -ForegroundColor Yellow
        $syntaxErrors++
        $warnings++
    }
}

if ($syntaxErrors -eq 0) {
    Write-Host "  ✓ Sin errores de sintaxis" -ForegroundColor Green
}

Write-Host ""

# Verificar que el servidor no esté corriendo
Write-Host "🌐 Verificando estado del servidor..." -ForegroundColor Yellow

try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/infobip/session" -Method POST -ErrorAction SilentlyContinue -TimeoutSec 2
    if ($response.StatusCode -eq 400 -or $response.StatusCode -eq 404) {
        Write-Host "  ✓ Servidor respondiendo en puerto 8080" -ForegroundColor Green
        Write-Host "    (Se puede ejecutar pruebas)" -ForegroundColor Gray
    }
} catch {
    if ($_.Exception.Message -like "*Connection refused*" -or $_.Exception.Message -like "*No connection*") {
        Write-Host "  ⚠ Servidor no está corriendo" -ForegroundColor Yellow
        Write-Host "    Ejecutar: go run api/cmd/main.go" -ForegroundColor Gray
    } else {
        Write-Host "  ⚠ No se pudo verificar el servidor" -ForegroundColor Yellow
    }
}

Write-Host ""

# Resumen
Write-Host "=====================================================" -ForegroundColor Cyan
Write-Host "📊 Resumen de Verificación" -ForegroundColor Cyan
Write-Host "=====================================================" -ForegroundColor Cyan

if ($errors -eq 0 -and $warnings -eq 0) {
    Write-Host "✅ Sistema completamente verificado!" -ForegroundColor Green
    Write-Host ""
    Write-Host "🚀 Próximos pasos:" -ForegroundColor Yellow
    Write-Host "  1. Configurar .env con tus credenciales" -ForegroundColor Gray
    Write-Host "  2. Ejecutar: go run api/cmd/main.go" -ForegroundColor Gray
    Write-Host "  3. Ejecutar pruebas: .\scripts\test_terms_flow.ps1" -ForegroundColor Gray
} elseif ($errors -eq 0) {
    Write-Host "⚠️  Sistema verificado con $warnings advertencia(s)" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "El sistema debería funcionar, pero revisa las advertencias." -ForegroundColor Gray
} else {
    Write-Host "❌ Se encontraron $errors error(es) y $warnings advertencia(s)" -ForegroundColor Red
    Write-Host ""
    Write-Host "Por favor, corrige los errores antes de continuar." -ForegroundColor Gray
}

Write-Host ""
Write-Host "📚 Documentación disponible:" -ForegroundColor Cyan
Write-Host "  - IMPLEMENTATION_SUMMARY.md" -ForegroundColor Gray
Write-Host "  - docs\TERMS_INTEGRATION.md" -ForegroundColor Gray
Write-Host "  - docs\TERMS_QUICKSTART.md" -ForegroundColor Gray
Write-Host "  - docs\FRONTEND_INTEGRATION.md" -ForegroundColor Gray
Write-Host "  - docs\FLOW_DIAGRAM.md" -ForegroundColor Gray
Write-Host ""
