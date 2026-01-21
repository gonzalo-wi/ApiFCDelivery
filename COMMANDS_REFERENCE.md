# 📝 Comandos Útiles - Términos y Condiciones

## 🚀 Comandos de Inicio

### Iniciar el servidor
```powershell
go run api/cmd/main.go
```

### Verificar instalación
```powershell
.\scripts\verify_installation.ps1
```

### Ejecutar pruebas automatizadas
```powershell
.\scripts\test_terms_flow.ps1
```

---

## 🧪 Comandos de Prueba Manual

### 1. Crear una sesión
```powershell
curl -X POST http://localhost:8080/api/v1/infobip/session `
  -H "Content-Type: application/json" `
  -d '{\"sessionId\": \"test-session-123\"}'
```

### 2. Consultar estado de un token
```powershell
$TOKEN = "tu-token-aqui"
curl http://localhost:8080/api/v1/terms/$TOKEN
```

### 3. Aceptar términos
```powershell
curl -X POST http://localhost:8080/api/v1/terms/$TOKEN/accept `
  -H "Content-Type: application/json"
```

### 4. Rechazar términos
```powershell
curl -X POST http://localhost:8080/api/v1/terms/$TOKEN/reject `
  -H "Content-Type: application/json"
```

---

## 🗄️ Comandos de Base de Datos

### Conectar a MySQL
```powershell
mysql -u root -p gofriocalor
```

### Ver todas las sesiones
```sql
SELECT * FROM terms_sessions ORDER BY created_at DESC;
```

### Ver sesiones por estado
```sql
SELECT status, COUNT(*) as total 
FROM terms_sessions 
GROUP BY status;
```

### Ver sesiones pendientes
```sql
SELECT token, session_id, created_at, expires_at 
FROM terms_sessions 
WHERE status = 'PENDING';
```

### Ver sesiones aceptadas hoy
```sql
SELECT token, session_id, accepted_at, ip, user_agent 
FROM terms_sessions 
WHERE status = 'ACCEPTED' 
AND DATE(accepted_at) = CURDATE();
```

### Ver fallos de notificación
```sql
SELECT token, session_id, notify_status, notify_attempts, last_error 
FROM terms_sessions 
WHERE notify_status = 'FAILED';
```

### Ver sesiones expiradas
```sql
SELECT token, session_id, created_at, expires_at 
FROM terms_sessions 
WHERE status = 'EXPIRED' 
OR (status = 'PENDING' AND expires_at < NOW());
```

### Limpiar sesiones antiguas (>30 días)
```sql
DELETE FROM terms_sessions 
WHERE created_at < DATE_SUB(NOW(), INTERVAL 30 DAY);
```

### Ver auditoría de una sesión específica
```sql
SELECT * FROM terms_sessions 
WHERE token = 'tu-token-aqui';
```

### Estadísticas generales
```sql
SELECT 
  status,
  COUNT(*) as total,
  COUNT(CASE WHEN notify_status = 'SENT' THEN 1 END) as notificaciones_exitosas,
  COUNT(CASE WHEN notify_status = 'FAILED' THEN 1 END) as notificaciones_fallidas
FROM terms_sessions 
GROUP BY status;
```

---

## 🔍 Comandos de Debugging

### Ver logs en tiempo real
```powershell
go run api/cmd/main.go | Select-String "términos"
```

### Ver logs de notificación
```powershell
go run api/cmd/main.go | Select-String "Notificación"
```

### Ver logs de errores
```powershell
go run api/cmd/main.go | Select-String "error|Error|ERROR"
```

### Verificar errores de compilación
```powershell
go build ./...
```

### Ejecutar tests (si existen)
```powershell
go test ./... -v
```

### Ver dependencias
```powershell
go list -m all
```

### Actualizar dependencias
```powershell
go mod tidy
```

---

## 📊 Comandos de Monitoreo

### Ver conexiones activas en el puerto 8080
```powershell
netstat -ano | Select-String ":8080"
```

### Ver procesos Go activos
```powershell
Get-Process | Where-Object {$_.ProcessName -eq "go"}
```

### Verificar conectividad con Infobip
```powershell
curl -I https://api2.infobip.com
```

### Test de endpoint de salud (si existe)
```powershell
curl http://localhost:8080/health
```

---

## 🔧 Comandos de Desarrollo

### Formatear código Go
```powershell
go fmt ./...
```

### Analizar código con go vet
```powershell
go vet ./...
```

### Instalar dependencia nueva
```powershell
go get github.com/nombre/paquete
```

### Ver documentación de un paquete
```powershell
go doc nombre/del/paquete
```

---

## 🌐 Comandos de Integración Frontend

### Probar CORS desde otro origen
```powershell
curl -X OPTIONS http://localhost:8080/api/v1/terms/test `
  -H "Origin: http://localhost:5173" `
  -H "Access-Control-Request-Method: POST"
```

### Simular request desde frontend
```powershell
curl -X POST http://localhost:8080/api/v1/terms/test-token/accept `
  -H "Content-Type: application/json" `
  -H "Origin: http://localhost:5173" `
  -H "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64)"
```

---

## 🔐 Comandos de Seguridad

### Verificar variables de entorno
```powershell
Get-Content .env
```

### Generar token de prueba (similar al backend)
```powershell
# En PowerShell
$bytes = New-Object byte[] 32
[Security.Cryptography.RandomNumberGenerator]::Fill($bytes)
[BitConverter]::ToString($bytes).Replace("-", "").ToLower()
```

---

## 📦 Comandos de Despliegue

### Compilar para producción
```powershell
go build -o bin/server.exe api/cmd/main.go
```

### Compilar para Linux
```powershell
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o bin/server api/cmd/main.go
```

### Ejecutar binario compilado
```powershell
.\bin\server.exe
```

---

## 🧹 Comandos de Limpieza

### Limpiar archivos temporales
```powershell
go clean
```

### Limpiar cache de módulos
```powershell
go clean -modcache
```

### Eliminar binarios compilados
```powershell
Remove-Item -Path "bin" -Recurse -Force
```

---

## 📚 Comandos de Documentación

### Ver documentación del proyecto
```powershell
Get-Content TERMS_README.md
Get-Content IMPLEMENTATION_SUMMARY.md
Get-Content docs\TERMS_INTEGRATION.md
```

### Abrir documentación en navegador
```powershell
# Si tienes markdown viewer
code TERMS_README.md
code IMPLEMENTATION_SUMMARY.md
```

---

## 🔄 Comandos de Git (si usas control de versiones)

### Ver archivos nuevos
```powershell
git status
```

### Agregar archivos de términos
```powershell
git add internal/models/terms_session.go
git add internal/dto/terms_dto.go
git add internal/store/terms_session_store.go
git add internal/service/infobip_client.go
git add internal/service/terms_session_service.go
git add internal/transport/terms_session_handler.go
git add internal/routes/terms_routes.go
git add migrations/001_create_terms_sessions.sql
git add docs/
git add scripts/
```

### Commit de la implementación
```powershell
git commit -m "feat: implementar flujo de términos y condiciones con Infobip

- Agregar modelo TermsSession con estados y auditoría
- Implementar cliente HTTP para Infobip con reintentos
- Agregar endpoints para crear, consultar, aceptar y rechazar términos
- Implementar notificación asíncrona a Infobip
- Agregar logging estructurado con zerolog
- Documentación completa del sistema
"
```

---

## 🎯 Comandos Rápidos por Escenario

### Escenario: Primera vez configurando el proyecto
```powershell
# 1. Copiar .env
cp .env.example .env

# 2. Editar .env (abrir en editor)
code .env

# 3. Verificar instalación
.\scripts\verify_installation.ps1

# 4. Iniciar servidor
go run api/cmd/main.go
```

### Escenario: Probar el flujo completo
```powershell
# En una terminal: iniciar servidor
go run api/cmd/main.go

# En otra terminal: ejecutar pruebas
.\scripts\test_terms_flow.ps1
```

### Escenario: Debugging de notificaciones fallidas
```powershell
# 1. Ver sesiones fallidas
mysql -u root -p gofriocalor -e "SELECT * FROM terms_sessions WHERE notify_status = 'FAILED';"

# 2. Ver logs de notificación
go run api/cmd/main.go | Select-String "Notificación"

# 3. Verificar conectividad con Infobip
curl -I https://api2.infobip.com

# 4. Verificar API Key en .env
Select-String "INFOBIP_API_KEY" .env
```

### Escenario: Limpiar y reiniciar desde cero
```powershell
# 1. Detener servidor (Ctrl+C)

# 2. Limpiar tabla
mysql -u root -p gofriocalor -e "TRUNCATE TABLE terms_sessions;"

# 3. Reiniciar servidor
go run api/cmd/main.go

# 4. Ejecutar pruebas
.\scripts\test_terms_flow.ps1
```

---

## 💡 Tips y Trucos

### Alias útiles (agregar a tu perfil de PowerShell)
```powershell
# Editar: $PROFILE
function Start-GoServer { go run api/cmd/main.go }
function Test-Terms { .\scripts\test_terms_flow.ps1 }
function Verify-Terms { .\scripts\verify_installation.ps1 }

# Usar como:
Start-GoServer
Test-Terms
Verify-Terms
```

### Variables de entorno temporales
```powershell
# Cambiar temporalmente el puerto
$env:PORT="9090"
go run api/cmd/main.go

# Cambiar nivel de log
$env:ENVIRONMENT="development"
go run api/cmd/main.go
```

### Watch mode (auto-reload al cambiar código)
```powershell
# Instalar air
go install github.com/cosmtrek/air@latest

# Ejecutar con auto-reload
air
```

---

**Comandos listos para usar! 🎉**

Para más información, consulta:
- TERMS_README.md
- IMPLEMENTATION_SUMMARY.md
- docs/TERMS_INTEGRATION.md
