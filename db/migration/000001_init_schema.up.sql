-- File: 001_init_schema.up.sql
-- Description: Initial schema migration to create tables for users, groups, group members, expenses, and more.

-- Create users table
CREATE TABLE users (
  id bigserial PRIMARY KEY,
  first_name VARCHAR NOT NULL,
  last_name VARCHAR NOT NULL,
  email VARCHAR NOT NULL,
  password_hash VARCHAR NOT NULL,
  phone VARCHAR NOT NULL,
  image_path VARCHAR NOT NULL,
  time_zone VARCHAR NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now())
);

-- Create group_categories table
CREATE TABLE group_categories (
  id bigserial PRIMARY KEY,
  name VARCHAR NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now())
);

-- Create groups table
CREATE TABLE groups (
  id bigserial PRIMARY KEY,
  name VARCHAR NOT NULL,
  category_id bigint NOT NULL,
  image_path VARCHAR NOT NULL,
  created_by_id bigint NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now())
);

-- Create group_members table
CREATE TABLE group_members (
  id bigserial PRIMARY KEY,
  group_id bigint NOT NULL,
  user_id bigint NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now())
);

-- Create expenses table
CREATE TABLE expenses (
  id bigserial PRIMARY KEY,
  group_id bigint NOT NULL,
  paid_by_id bigint NOT NULL,
  amount VARCHAR NOT NULL,
  description VARCHAR NOT NULL,
  date timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now())
);

-- Create expense_shares table
CREATE TABLE expense_shares (
  id bigserial PRIMARY KEY,
  expense_id bigint NOT NULL,
  user_id bigint NOT NULL,
  share VARCHAR NOT NULL,
  paid_status boolean NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now())
);

-- Create notifications table
CREATE TABLE notifications (
  id bigserial PRIMARY KEY,
  user_id bigint NOT NULL,
  message VARCHAR NOT NULL,
  is_read boolean NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now())
);

-- Create invitations table
CREATE TABLE invitations (
  id bigserial PRIMARY KEY,
  inviter_id bigint NOT NULL,
  invitee_id bigint NOT NULL,
  group_id bigint NOT NULL,
  status VARCHAR NOT NULL,
  code VARCHAR NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now()),
  accepted_at timestamptz,
  rejected_at timestamptz
);

-- Create currency table
CREATE TABLE currencies (
  id bigserial PRIMARY KEY,
  code VARCHAR NOT NULL,
  name VARCHAR NOT NULL,
  symbol VARCHAR NOT NULL,
  exchange_rate double precision NOT NULL,
  updated_at timestamptz NOT NULL
);

-- ... Other tables ...

-- Add foreign key constraints

-- Add foreign key constraint on groups table
ALTER TABLE groups ADD FOREIGN KEY (category_id) REFERENCES group_categories (id);
ALTER TABLE groups ADD FOREIGN KEY (created_by_id) REFERENCES users (id);

-- Add foreign key constraint on group_members table
ALTER TABLE group_members ADD FOREIGN KEY (group_id) REFERENCES groups (id);

-- Add foreign key constraint on expenses table
ALTER TABLE expenses ADD FOREIGN KEY (group_id) REFERENCES groups (id);
ALTER TABLE expenses ADD FOREIGN KEY (paid_by_id) REFERENCES users (id);

-- Add foreign key constraint on expense_shares table
ALTER TABLE expense_shares ADD FOREIGN KEY (expense_id) REFERENCES expenses (id);
ALTER TABLE expense_shares ADD FOREIGN KEY (user_id) REFERENCES users (id);

-- Add foreign key constraint on notifications table
ALTER TABLE notifications ADD FOREIGN KEY (user_id) REFERENCES users (id);

-- Add foreign key constraint on invitations table
ALTER TABLE invitations ADD FOREIGN KEY (inviter_id) REFERENCES users (id);
ALTER TABLE invitations ADD FOREIGN KEY (invitee_id) REFERENCES users (id);
ALTER TABLE invitations ADD FOREIGN KEY (group_id) REFERENCES groups (id);




ALTER TABLE "users" ADD CONSTRAINT "email_key" UNIQUE ("email");


-- ... Other foreign key constraints ...

-- Create indexes

-- Create index on users table
CREATE INDEX ON users (email);

-- Create index on groups table
CREATE INDEX ON groups (category_id);
CREATE INDEX ON groups (created_by_id);

-- Create index on group_members table
CREATE INDEX ON group_members (group_id);

-- Create index on expenses table
CREATE INDEX ON expenses (group_id);
CREATE INDEX ON expenses (paid_by_id);

-- Create index on expense_shares table
CREATE INDEX ON expense_shares (expense_id);
CREATE INDEX ON expense_shares (user_id);

-- Create index on notifications table
CREATE INDEX ON notifications (user_id);

-- Create index on invitations table
CREATE INDEX ON invitations (inviter_id);
CREATE INDEX ON invitations (invitee_id);
CREATE INDEX ON invitations (group_id);

-- ... Other indexes ...
