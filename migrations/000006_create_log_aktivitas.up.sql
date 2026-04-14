CREATE TABLE tb_log_aktivitas (
    id_log SERIAL PRIMARY KEY,
    id_user INT NOT NULL REFERENCES tb_user(id_user) ON DELETE CASCADE,
    aktivitas VARCHAR(255) NOT NULL,
    waktu_aktivitas TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);