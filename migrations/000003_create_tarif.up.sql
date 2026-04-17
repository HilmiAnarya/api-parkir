CREATE TYPE jenis_kendaraan_enum AS ENUM ('motor', 'mobil', 'lainnya');

CREATE TABLE tb_tarif (
    id SERIAL PRIMARY KEY,
    jenis_kendaraan jenis_kendaraan_enum NOT NULL,
    tarif_per_jam DECIMAL(10,0) NOT NULL,
    deleted_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);