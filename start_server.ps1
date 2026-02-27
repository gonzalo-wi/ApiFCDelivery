# Script para iniciar el servidor con las variables del archivo .env

# Cargar variables del archivo .env
Get-Content .env | ForEach-Object {
    if ($_ -match '^\s*([^#][^=]+)=(.+)$') {
        $name = $matches[1].Trim()
        $value = $matches[2].Trim()
        Set-Item -Path "env:$name" -Value $value
        Write-Host "OK $name configurado" -ForegroundColor Green
    }
}

Write-Host ""
Write-Host "Iniciando servidor..." -ForegroundColor Cyan
Write-Host ""

# Ejecutar el servidor
go run ./api/cmd/main.go
