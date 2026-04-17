CREATE TABLE tb_area_parkir (
    id_area SERIAL PRIMARY KEY,
    nama_area VARCHAR(50) NOT NULL,
    kapasitas INT NOT NULL DEFAULT 0,
    terisi INT NOT NULL DEFAULT 0,
    deleted_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);