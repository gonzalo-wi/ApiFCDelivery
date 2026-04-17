-- Migración 010: Agregar campo order_number a deliveries
-- El número de orden de trabajo ahora lo asigna la app mobile

ALTER TABLE deliveries ADD COLUMN IF NOT EXISTS order_number VARCHAR(50);
