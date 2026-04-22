package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Response memetakan struktur utama dari API Open Data Jatim
type Response struct {
	Code       int              `json:"code"`
	Message    string           `json:"message"`
	Data       []DataPelatihan  `json:"data"`
	Error      interface{}      `json:"error"`
	Pagination PaginationDetail `json:"pagination"`
	Metadata   []string         `json:"metadata"`
}

// DataPelatihan memetakan isi array 'data'
type DataPelatihan struct {
	ID                 int    `json:"id"`
	NamaProvinsi       string `json:"nama_provinsi"`
	BalaiLatihanKerja  string `json:"balai_latihan_kerja"`
	Kategori           string `json:"kategori"`
	Jumlah             int    `json:"jumlah"`
	Satuan             string `json:"satuan"`
	Tahun              int    `json:"tahun"`
	PeriodeUpdate      string `json:"periode_update"`
}

// PaginationDetail memetakan informasi halaman
type PaginationDetail struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	TotalPage int `json:"total_page"`
}

func main() {
	// 1. Setup URL dan Query Parameters
	baseURL := "https://opendata.jatimprov.go.id/api/cleaned-bigdata/dinas_tenaga_kerja_dan_transmigrasi_provinsi_jawa_timur/jumlah_peserta_pelatihan_berbasis_kompetensi_pbk"
	
	params := url.Values{}
	params.Add("sort", "id:asc")
	params.Add("page", "1")
	params.Add("per_page", "10")
	params.Add("where", `{"periode_update":["2026-Q1"]}`)
	params.Add("where_or", "{}")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// 2. Buat HTTP Client dengan Timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 3. Execute Request
	resp, err := client.Get(fullURL)
	if err != nil {
		fmt.Printf("Error membuat request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 4. Read Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error membaca body: %v\n", err)
		return
	}

	// 5. Unmarshal JSON ke Struct
	var apiResponse Response
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	// 6. Tampilkan Hasil
	fmt.Printf("Status: %s\n", apiResponse.Message)
	fmt.Printf("Total Halaman: %d\n", apiResponse.Pagination.TotalPage)
	fmt.Println("---------------------------------------------------------")
	fmt.Printf("%-5s | %-30s | %-10s | %-5s\n", "ID", "Balai Latihan", "Kategori", "Jumlah")
	fmt.Println("---------------------------------------------------------")

	for _, item := range apiResponse.Data {
		fmt.Printf("%-5d | %-30.30s | %-10s | %-5d %s\n", 
			item.ID, 
			item.BalaiLatihanKerja, 
			item.Kategori, 
			item.Jumlah, 
			item.Satuan,
		)
	}
}