PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
  id INTEGER NOT NULL PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  auth_type TEXT CHECK( auth_type IN ('email','google') ) NOT NULL,
  password TEXT,
  socialID TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sessions (
  id TEXT UNIQUE NOT NULL,
  value JSON,
  time_accessed TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS stickey_connections (
  id TEXT UNIQUE NOT NULL,
  connection_api_key TEXT NOT NULL,
  topic_name TEXT NOT NULL,
  last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);