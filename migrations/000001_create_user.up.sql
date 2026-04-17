-- Buat tipe ENUM untuk role
CREATE TYPE user_role AS ENUM ('admin', 'petugas', 'owner');

CREATE TABLE tb_user (
    id SERIAL PRIMARY KEY,
    nama_lengkap VARCHAR(50) NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    role user_role DEFAULT 'petugas',
    status_aktif BOOLEAN DEFAULT TRUE, -- Di PostgreSQL, boolean lebih baik dari tinyint(1)
    deleted_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
