-- Migration 005: Add client information fields to deliveries table
-- Adds name, address, and locality columns for customer details

ALTER TABLE deliveries
ADD COLUMN name VARCHAR(200) DEFAULT '' AFTER nro_cta,
ADD COLUMN address VARCHAR(300) DEFAULT '' AFTER name,
ADD COLUMN locality VARCHAR(100) DEFAULT '' AFTER address;

-- Add indexes for searching
CREATE INDEX idx_deliveries_name ON deliveries(name);
CREATE INDEX idx_deliveries_locality ON deliveries(locality);
