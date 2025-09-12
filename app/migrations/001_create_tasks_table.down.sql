DROP TRIGGER IF EXISTS update_tasks_updated_at ON tasks;

DROP FUNCTION IF EXISTS update_updated_at_column();

DROP INDEX IF EXISTS idx_tasks_status_due_date;
DROP INDEX IF EXISTS idx_tasks_status_priority;
DROP INDEX IF EXISTS idx_tasks_updated_at;
DROP INDEX IF EXISTS idx_tasks_created_at;
DROP INDEX IF EXISTS idx_tasks_due_date;
DROP INDEX IF EXISTS idx_tasks_priority;
DROP INDEX IF EXISTS idx_tasks_status;

DROP TABLE IF EXISTS tasks;