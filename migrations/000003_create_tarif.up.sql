CREATE TYPE jenis_kendaraan_enum AS ENUM ('motor', 'mobil', 'lainnya');

CREATE TABLE tb_tarif (
    id_tarif SERIAL PRIMARY KEY,
    jenis_kendaraan jenis_kendaraan_enum NOT NULL,
    tarif_per_jam DECIMAL(10,0) NOT NULL
);