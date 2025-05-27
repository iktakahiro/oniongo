-- Create "project" table
CREATE TABLE `project` (`id` uuid NOT NULL, PRIMARY KEY (`id`));
-- Create "todo" table
CREATE TABLE `todo` (`id` uuid NOT NULL, `title` text NOT NULL, `body` text NULL, `status` text NOT NULL DEFAULT 'NOT_STARTED', `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, `deleted_at` datetime NULL, PRIMARY KEY (`id`));
-- Create index "todoschema_deleted_at_created_at" to table: "todo"
CREATE INDEX `todoschema_deleted_at_created_at` ON `todo` (`deleted_at`, `created_at`);
-- Create index "todoschema_deleted_at_updated_at" to table: "todo"
CREATE INDEX `todoschema_deleted_at_updated_at` ON `todo` (`deleted_at`, `updated_at`);
-- Create index "todoschema_deleted_at_status" to table: "todo"
CREATE INDEX `todoschema_deleted_at_status` ON `todo` (`deleted_at`, `status`);
