
CREATE TABLE my_client (
  id SERIAL PRIMARY KEY,
  name VARCHAR(250) NOT NULL,
  slug VARCHAR(100) NOT NULL,
  is_project BOOLEAN NOT NULL DEFAULT false,
  self_capture CHAR(1) NOT NULL DEFAULT '1',
  client_logo VARCHAR(255) NOT NULL DEFAULT 'no-image.jpg',
  address TEXT,
  phone_number VARCHAR(50),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);
    