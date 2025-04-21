package tests

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestFullIntegrationFlow(t *testing.T) {
	connStr := "host=localhost port=5433 user=admin password=password dbname=TestDataBase sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		t.Fatalf("Не удалось подключиться к тестовой базе данных: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
		TRUNCATE TABLE products, receptions, pvz, users RESTART IDENTITY CASCADE
	`)
	if err != nil {
		t.Fatalf("Не удалось очистить таблицы: %v", err)
	}

	pvzID := uuid.New()
	registrationDate := time.Now()
	city := "Москва"

	_, err = db.Exec(`
		INSERT INTO pvz (id, registration_date, city)
		VALUES ($1, $2, $3)
	`, pvzID, registrationDate, city)
	if err != nil {
		t.Fatalf("Не удалось создать ПВЗ: %v", err)
	}

	receptionID := uuid.New()
	dateTime := time.Now()
	status := "open"

	_, err = db.Exec(`
		INSERT INTO receptions (id, date_time, pvz_id, reception_status)
		VALUES ($1, $2, $3, $4)
	`, receptionID, dateTime, pvzID, status)
	if err != nil {
		t.Fatalf("Не удалось создать приемку: %v", err)
	}

	for i := 0; i < 50; i++ {
		productID := uuid.New()
		productType := "Электроника"

		_, err = db.Exec(`
			INSERT INTO products (id, date_time, product_type, reception_id)
			VALUES ($1, $2, $3, $4)
		`, productID, dateTime, productType, receptionID)
		if err != nil {
			t.Fatalf("Не удалось добавить товар: %v", err)
		}
	}

	_, err = db.Exec(`
		UPDATE receptions SET reception_status = 'closed' WHERE id = $1
	`, receptionID)
	if err != nil {
		t.Fatalf("Не удалось закрыть приемку: %v", err)
	}

	var finalStatus string
	err = db.Get(&finalStatus, `SELECT reception_status FROM receptions WHERE id = $1`, receptionID)
	if err != nil {
		t.Fatalf("Не удалось получить статус приемки: %v", err)
	}
	assert.Equal(t, "closed", finalStatus, "Приемка не была закрыта")
}
