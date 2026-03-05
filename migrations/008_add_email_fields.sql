-- Migration 008: Add email field to deliveries and work_orders tables
-- Description: Adds email column to store customer email addresses for notifications

-- Add email column to deliveries table
ALTER TABLE deliveries 
ADD COLUMN IF NOT EXISTS email VARCHAR(200);

-- Add email column to work_orders table  
ALTER TABLE work_orders
ADD COLUMN IF NOT EXISTS email VARCHAR(200);

-- Add index on email for potential email-based queries
CREATE INDEX IF NOT EXISTS idx_deliveries_email ON deliveries(email);
CREATE INDEX IF NOT EXISTS idx_work_orders_email ON work_orders(email);

-- Comment the columns for documentation
COMMENT ON COLUMN deliveries.email IS 'Customer email address for delivery notifications';
COMMENT ON COLUMN work_orders.email IS 'Customer email address for work order notifications';
