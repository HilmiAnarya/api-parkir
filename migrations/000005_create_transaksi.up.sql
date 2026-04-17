CREATE TYPE status_transaksi AS ENUM ('masuk', 'keluar');

CREATE TABLE tb_transaksi (
    id_parkir SERIAL PRIMARY KEY,
    id_kendaraan INT NOT NULL REFERENCES tb_kendaraan(id) ON DELETE CASCADE,
    waktu_masuk TIMESTAMP NOT NULL,
    waktu_keluar TIMESTAMP,
    id_tarif INT REFERENCES tb_tarif(id) ON DELETE SET NULL,
    durasi_jam INT,
    biaya_total DECIMAL(10,0),
    status status_transaksi DEFAULT 'masuk',
    id_user INT NOT NULL REFERENCES tb_user(id),
    id_area INT NOT NULL REFERENCES tb_area_parkir(id),
    foto_masuk VARCHAR(255) NOT NULL,
    foto_keluar VARCHAR(255),
    deleted_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);