<?php

// 1. Setup URL dan Query Parameters
$baseURL = "https://opendata.jatimprov.go.id/api/cleaned-bigdata/dinas_tenaga_kerja_dan_transmigrasi_provinsi_jawa_timur/jumlah_peserta_pelatihan_berbasis_kompetensi_pbk";

$params = [
    "sort"     => "id:asc",
    "page"     => 1,
    "per_page" => 10,
    "where"    => json_encode(["periode_update" => ["2026-Q1"]]),
    "where_or" => "{}"
];

$fullURL = $baseURL . "?" . http_build_query($params);

// 2. Init cURL
$ch = curl_init();

curl_setopt_array($ch, [
    CURLOPT_URL            => $fullURL,
    CURLOPT_RETURNTRANSFER => true,       // Kembalikan hasil sebagai string
    CURLOPT_TIMEOUT        => 10,         // Timeout 10 detik
    CURLOPT_HTTPGET        => true,       // Gunakan GET method
    CURLOPT_USERAGENT      => 'Mozilla/5.0 (X11; Linux x86_64; rv:149.0) Gecko/20100101 Firefox/149.0' // Tanpa UA kena blokir Cloudflare
]);

// 3. Execute Request
$response = curl_exec($ch);
$err = curl_error($ch);
$info = curl_getinfo($ch);

curl_close($ch);

// 4. Cek Error
if ($err) {
    echo "Error membuat request: " . $err;
    exit;
}

// 5. Decode JSON ke Associative Array
$apiResponse = json_decode($response, true);

if (json_last_error() !== JSON_ERROR_NONE) {
    echo "Error parsing JSON: " . json_last_error_msg();
    exit;
}

// 6. Tampilkan Hasil
echo "Status: " . ($apiResponse['message'] ?? 'N/A') . "\n";
echo "Total Halaman: " . ($apiResponse['pagination']['total_page'] ?? 0) . "\n";
echo str_repeat("-", 70) . "\n";
printf("%-5s | %-35s | %-15s | %-10s\n", "ID", "Balai Latihan", "Kategori", "Jumlah");
echo str_repeat("-", 70) . "\n";

if (!empty($apiResponse['data'])) {
    foreach ($apiResponse['data'] as $item) {
        printf(
            "%-5d | %-35.35s | %-15s | %-5d %s\n",
            $item['id'],
            $item['balai_latihan_kerja'],
            $item['kategori'],
            $item['jumlah'],
            $item['satuan']
        );
    }
} else {
    echo "Tidak ada data ditemukan.\n";
}