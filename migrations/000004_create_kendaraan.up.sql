CREATE TABLE tb_kendaraan (
    id_kendaraan SERIAL PRIMARY KEY,
    plat_nomor VARCHAR(15) NOT NULL UNIQUE,
    jenis_kendaraan VARCHAR(20),
    warna VARCHAR(20),
    pemilik VARCHAR(100),
    id_user INT REFERENCES tb_user(id_user) ON DELETE SET NULL
);