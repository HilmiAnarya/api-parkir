package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// AIResponse adalah kerangka untuk membaca balasan dari Python
type AIResponse struct {
	Success   bool    `json:"success"`
	PlatNomor string  `json:"plat_nomor"`
	Error     string  `json:"error"`
}

// DetectPlate mengirim foto ke API Python dan mengembalikan string plat nomor
func DetectPlate(fileBytes []byte, filename string) (string, error) {
	// 1. Siapkan body request bertipe "multipart/form-data"
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	// Cocokkan "file" dengan nama parameter yang diminta oleh API Python/FastAPI milikmu
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	part.Write(fileBytes)
	writer.Close()

	// 2. Buat HTTP Request ke Python (Asumsi Python jalan di port 8000)
	// Sesuaikan URL ini jika port Python-mu berbeda
	req, err := http.NewRequest("POST", "http://localhost:8000/detect-plate", body)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	// 3. Eksekusi Request dengan Timeout 10 detik agar Go tidak hang jika Python mati
	client := &http.Client{Timeout: 600 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("gagal menghubungi AI Python. Pastikan server Python menyala")
	}
	defer resp.Body.Close()

	// 4. Baca balasan JSON dari Python
	respBytes, _ := io.ReadAll(resp.Body)
	var aiResult AIResponse
	if err := json.Unmarshal(respBytes, &aiResult); err != nil {
		return "", errors.New("balasan dari AI tidak valid")
	}

	// 5. Cek apakah Python berhasil membaca
	if !aiResult.Success {
		return "", errors.New("AI Error: " + aiResult.Error)
	}

	return aiResult.PlatNomor, nil
}