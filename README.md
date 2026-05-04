# 💸 Expense Tracker API

A REST API for tracking personal expenses, built with Go and PostgreSQL.

**Live API:** `https://expense-api-production-4f40.up.railway.app`

---

## Tech Stack

- **Go** — backend language
- **PostgreSQL** — database
- **GORM** — ORM for database operations
- **Railway** — cloud deployment

---

## Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/tambah` | Add a new expense |
| `GET` | `/pengeluaran` | Get all expenses |
| `GET` | `/total` | Get total amount spent |
| `GET` | `/filter?dari=YYYY-MM-DD&sampai=YYYY-MM-DD` | Filter expenses by date range |
| `PUT` | `/update?id=X` | Update an expense |
| `DELETE` | `/hapus?id=X` | Delete an expense |

---

## Request & Response Examples

### Add Expense
```http
POST /tambah
Content-Type: application/json

{
    "jumlah": 25000,
    "deskripsi": "Lunch at warteg",
    "kategori": "makan",
    "tanggal": "2026-05-04"
}
```

Response:
```json
{"pesan": "data berhasil dibuat"}
```

### Get All Expenses
```http
GET /pengeluaran
```

Response:
```json
[
    {
        "id": 1,
        "jumlah": 25000,
        "deskripsi": "Lunch at warteg",
        "kategori": "makan",
        "tanggal": "2026-05-04T00:00:00Z"
    }
]
```

### Get Total
```http
GET /total
```

Response:
```json
{"total": 25000}
```

### Filter by Date Range
```http
GET /filter?dari=2026-05-01&sampai=2026-05-31
```

---

## Validation Rules

- `jumlah` — required, must be greater than 0
- `deskripsi` — required
- `kategori` — required, one of: `makan`, `transport`, `hiburan`, `kesehatan`, `lainnya`
- `tanggal` — required, format: `YYYY-MM-DD`

---

## Local Development

**Prerequisites:** Go 1.22+, PostgreSQL

```bash
git clone https://github.com/Tarquished/expense-api.git
cd expense-api
go mod tidy
go run main.go
```

Server runs at `http://localhost:8080`.

**Environment Variables:**

| Variable | Description |
|----------|-------------|
| `DATABASE_URL` | PostgreSQL connection string |
| `PORT` | Server port (default: 8080) |
