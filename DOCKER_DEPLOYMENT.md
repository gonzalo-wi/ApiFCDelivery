# Despliegue con Docker en servidor 192.168.0.250

## Requisitos previos
- Docker instalado
- Docker Compose instalado
- Acceso al servidor MySQL en 192.168.0.227

## Configuración

### 1. Editar archivo `.env.docker`
```bash
# Agregar tu API Key de Infobip
INFOBIP_API_KEY=tu_api_key_real_aqui
```

### 2. Construir y levantar el contenedor

```bash
# Construir la imagen
docker-compose build

# Levantar el contenedor
docker-compose up -d
```

### 3. Verificar que está corriendo

```bash
# Ver logs
docker-compose logs -f

# Ver estado
docker-compose ps

# Probar la API
curl http://192.168.0.250:8095/health
```

## Comandos útiles

### Ver logs en tiempo real
```bash
docker-compose logs -f app
```

### Reiniciar el servicio
```bash
docker-compose restart
```

### Detener el servicio
```bash
docker-compose down
```

### Reconstruir después de cambios en el código
```bash
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

### Entrar al contenedor
```bash
docker exec -it gofrocalor-api sh
```

## Acceso

La API estará disponible en:
- **URL interna:** http://localhost:8080 (dentro del contenedor)
- **URL externa:** http://192.168.0.250:8095

## Endpoints principales

- Health check: `http://192.168.0.250:8095/health`
- API v1: `http://192.168.0.250:8095/api/v1/`

## Variables de entorno

Las siguientes variables están configuradas en `docker-compose.yml`:

- `DB_HOST`: 192.168.0.227
- `DB_PORT`: 3306
- `DB_NAME`: friocalor
- `PORT`: 8080 (interno)
- `ENVIRONMENT`: production
- `APP_BASE_URL`: https://www.somoselagua.com.ar/api/dist
- `CORS_ORIGINS`: Múltiples orígenes permitidos

## Troubleshooting

### El contenedor no inicia
```bash
# Ver logs completos
docker-compose logs app

# Verificar que no haya otro servicio en el puerto 8095
netstat -an | grep 8095
```

### No puede conectar a la base de datos
```bash
# Verificar conectividad desde el contenedor
docker exec -it gofrocalor-api ping 192.168.0.227
```

### Actualizar solo variables de entorno
```bash
# Editar docker-compose.yml o .env.docker
# Luego:
docker-compose up -d
```

## Monitoreo

### Health check automático
Docker verificará automáticamente la salud del servicio cada 30 segundos mediante el endpoint `/health`.

### Ver estado del health check
```bash
docker inspect gofrocalor-api | grep Health -A 10
```
