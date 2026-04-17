CREATE TABLE tb_kendaraan (
    id SERIAL PRIMARY KEY,
    plat_nomor VARCHAR(15) NOT NULL UNIQUE,
    jenis_kendaraan VARCHAR(20),
    warna VARCHAR(20),
    pemilik VARCHAR(100),
    id_user INT REFERENCES tb_user(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);