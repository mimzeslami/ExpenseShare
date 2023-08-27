-- File: 001_init_schema.down.sql
-- Description: Revert the changes made by the initial schema migration.

-- Drop foreign key constraints

-- Remove foreign key constraint on invitations table
ALTER TABLE invitations DROP CONSTRAINT IF EXISTS invitations_group_id_fkey;
ALTER TABLE invitations DROP CONSTRAINT IF EXISTS invitations_invitee_id_fkey;
ALTER TABLE invitations DROP CONSTRAINT IF EXISTS invitations_inviter_id_fkey;

-- Remove foreign key constraint on notifications table
ALTER TABLE notifications DROP CONSTRAINT IF EXISTS notifications_user_id_fkey;

-- Remove foreign key constraint on expense_shares table
ALTER TABLE expense_shares DROP CONSTRAINT IF EXISTS expense_shares_user_id_fkey;
ALTER TABLE expense_shares DROP CONSTRAINT IF EXISTS expense_shares_expense_id_fkey;

-- Remove foreign key constraint on expenses table
ALTER TABLE expenses DROP CONSTRAINT IF EXISTS expenses_paid_by_id_fkey;
ALTER TABLE expenses DROP CONSTRAINT IF EXISTS expenses_group_id_fkey;

-- Remove foreign key constraint on group_members table
ALTER TABLE group_members DROP CONSTRAINT IF EXISTS group_members_group_id_fkey;

-- Remove foreign key constraint on groups table
ALTER TABLE groups DROP CONSTRAINT IF EXISTS groups_category_id_fkey;
ALTER TABLE groups DROP CONSTRAINT IF EXISTS groups_created_by_id_fkey;

-- Drop indexes

-- Drop index on invitations table
DROP INDEX IF EXISTS invitations_group_id_idx;
DROP INDEX IF EXISTS invitations_invitee_id_idx;
DROP INDEX IF EXISTS invitations_inviter_id_idx;

-- Drop index on notifications table
DROP INDEX IF EXISTS notifications_user_id_idx;

-- Drop index on expense_shares table
DROP INDEX IF EXISTS expense_shares_user_id_idx;
DROP INDEX IF EXISTS expense_shares_expense_id_idx;

-- Drop index on expenses table
DROP INDEX IF EXISTS expenses_group_id_idx;
DROP INDEX IF EXISTS expenses_paid_by_id_idx;

-- Drop index on group_members table
DROP INDEX IF EXISTS group_members_group_id_idx;

-- Drop index on groups table
DROP INDEX IF EXISTS groups_category_id_idx;
DROP INDEX IF EXISTS groups_created_by_id_idx;

-- Drop tables

-- Drop table invitations
DROP TABLE IF EXISTS invitations;

-- Drop table notifications
DROP TABLE IF EXISTS notifications;

-- Drop table expense_shares
DROP TABLE IF EXISTS expense_shares;

-- Drop table expenses
DROP TABLE IF EXISTS expenses;

-- Drop table group_members
DROP TABLE IF EXISTS group_members;

-- Drop table groups
DROP TABLE IF EXISTS groups;

-- Drop table group_categories
DROP TABLE IF EXISTS group_categories;

-- Drop table currencies
DROP TABLE IF EXISTS currencies;

-- Drop table users
DROP TABLE IF EXISTS users;
