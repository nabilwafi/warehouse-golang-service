BEGIN;

CREATE TYPE user_role AS ENUM ('admin', 'staff');

CREATE TABLE users (
  id UUID PRIMARY KEY,
  email VARCHAR(100) NOT NULL UNIQUE,
  password VARCHAR(100) NOT NULL,
  name VARCHAR(100) NOT NULL,
  role user_role NOT NULL DEFAULT 'staff'
);

COMMIT;