-- Migration 003: Agregar tipos de dispensers P (Pie) y M (Mesada)
-- Fecha: 2026-02-26
-- Descripción: Amplía el enum de tipos de dispensers para incluir Pie (P) y Mesada (M)
-- Estos tipos son utilizados por el endpoint de Infobip para especificar dispensers

-- Nota: GORM AutoMigrate debería manejar esto automáticamente al arrancar
-- Esta migración es para documentación y ejecución manual si es necesario

-- MySQL: Modificar constraint de CHECK si existe
-- (En la mayoría de casos, GORM maneja las validaciones a nivel de aplicación)

-- Verificar tipos actuales
SELECT DISTINCT tipo FROM dispensers;

-- Los nuevos tipos 'P' y 'M' serán aceptados automáticamente por GORM
-- cuando se actualice el modelo en el código

-- Si tu base de datos usa enum explícito (no es el caso de GORM por defecto), ejecuta:
-- ALTER TABLE dispensers MODIFY COLUMN tipo VARCHAR(20);

-- Validar que el cambio fue aplicado
-- SELECT COLUMN_TYPE FROM INFORMATION_SCHEMA.COLUMNS 
-- WHERE TABLE_NAME = 'dispensers' AND COLUMN_NAME = 'tipo';

-- Esta migración es principalmente informativa ya que GORM AutoMigrate
-- maneja los cambios de validación automáticamente
