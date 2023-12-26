package _tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ArtemFed/mts-final-taxi/pkg/mylogger"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/location/internal/config"
	http2 "gitlab.com/ArtemFed/mts-final-taxi/projects/location/internal/handler"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/location/internal/handler/generated"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/location/internal/repository"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/location/internal/service"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Выполнять внутри папки .tests

type ReqAns struct {
	req *http.Request
	ans int
}

func TestLocationService(t *testing.T) {
	cfg, err := config.NewConfig("../../config/config.test.yml", "LOCATION")
	if err != nil {
		t.Fatal("Error reading config", zap.Error(err))
	}

	logger, err := mylogger.InitLogger(cfg.Logger, cfg.App.Name)
	if err != nil {
		t.Fatal("Error init logger", zap.Error(err))
	}

	router := gin.Default()

	db, err := repository.NewPostgresDB(cfg.Postgres)
	if err != nil {
		t.Fatal("Error creating postgre", zap.Error(err))
	}

	locationRepo := repository.NewLocalRepository(db)

	locationService := service.NewLocationService(locationRepo)

	http2.InitHandler(router, logger, nil, locationService, cfg.App)

	// Подготовка базы данных с использованием миграций
	err = runMigrations(db)
	if err != nil {
		t.Fatal("Error running migrations:", err)
	}

	// Создаём запросы
	driverId1 := "e4141f80-2d8c-4864-9b45-8f9eaf60dc1f"
	reqPostOk1 := createRequest(
		t, "POST",
		fmt.Sprintf("http://localhost:%v/api/v1/drivers/%s/location", cfg.Http.Port, driverId1),
		&map[string]interface{}{"Lat": 2, "Lng": 2},
	)

	driverId2 := "e5141f80-2d8c-4864-9b45-8f9eaf60dc1f"
	reqPostOk2 := createRequest(
		t, "POST",
		fmt.Sprintf("http://localhost:%v/api/v1/drivers/%s/location", cfg.Http.Port, driverId2),
		&map[string]interface{}{"Lat": 90.0, "Lng": 180.0},
	)

	driverId3 := "e6141f80-2d8c-4864-9b45-8f9eaf60dc1f"
	reqPostOk3 := createRequest(
		t, "POST",
		fmt.Sprintf("http://localhost:%v/api/v1/drivers/%s/location", cfg.Http.Port, driverId3),
		&map[string]interface{}{"Lat": -1, "Lng": 0},
	)

	driverId4 := "e7141f80-2d8c-4864-9b45-8f9eaf60dc1f"
	reqPostBad1 := createRequest(
		t, "POST",
		fmt.Sprintf("http://localhost:%v/api/v1/drivers/%s/location", cfg.Http.Port, driverId4),
		&map[string]interface{}{"Lat": 91, "Lng": 2},
	)

	reqPostBad2 := createRequest(
		t, "POST",
		fmt.Sprintf("http://localhost:%v/api/v1/drivers/%s/location", cfg.Http.Port, driverId4),
		&map[string]interface{}{"Lat": -91, "Lng": 2},
	)

	reqPostBad3 := createRequest(
		t, "POST",
		fmt.Sprintf("http://localhost:%v/api/v1/drivers/%s/location", cfg.Http.Port, driverId4),
		&map[string]interface{}{"Lat": 2, "Lng": 181},
	)

	reqPostBad4 := createRequest(
		t, "POST",
		fmt.Sprintf("http://localhost:%v/api/v1/drivers/%s/location", cfg.Http.Port, driverId4),
		&map[string]interface{}{"Lat": 2, "Lng": -181},
	)

	reqGet := createRequest(
		t, "GET",
		//fmt.Sprintf("http://localhost:%v/api/v1/drivers?lat=%f&lng=%f&radius=%f", cfg.Http.Port, 1.0, 1.0, 111282.0),
		fmt.Sprintf("http://localhost:%v/api/v1/drivers?lat=%f&lng=%f&radius=%f", cfg.Http.Port, 1.0, 1.0, 2111282.0),
		nil,
	)

	reqs := []ReqAns{
		{reqPostOk1, http.StatusOK},
		{reqPostOk2, http.StatusOK},
		{reqPostOk3, http.StatusOK},
		{reqPostBad1, http.StatusBadRequest},
		{reqPostBad2, http.StatusBadRequest},
		{reqPostBad3, http.StatusBadRequest},
		{reqPostBad4, http.StatusBadRequest},
	}

	for i, r := range reqs {
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r.req)
		if status := w.Code; status != r.ans {
			t.Errorf("Wrong status code returned for UpdateDriver req#%d: got %v want %v", i, status, r.ans)
		}
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqGet)
	if status := w.Code; status != http.StatusOK {
		t.Errorf("Wrong status code returned for GetDrivers: got %v want %v", status, http.StatusOK)
	}
	var response []generated.Driver
	t.Log("Body" + w.Body.String())
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal("Error decoding response:", zap.Error(err))
	}

	fmt.Printf("Drivers: %v\n", response)
	if *response[0].Id != driverId1 || response[0].Lat != 2 || response[0].Lng != 2 {
		t.Error(fmt.Sprintf("Answer!=Correct: %s!=%s or !=2 or !=2", *response[0].Id, driverId1), zap.Error(err))
	}
	if *response[1].Id != driverId3 || response[1].Lat != -1 || response[1].Lng != 0 {
		t.Error(fmt.Sprintf("Answer!=Correct: %s!=%s or !=-1 or !=0", *response[1].Id, driverId3), zap.Error(err))
	}
}

func createRequest(t *testing.T, method string, path string, dict *map[string]interface{}) *http.Request {
	var err error
	var jsonBody []byte
	if dict != nil {
		jsonBody, err = json.Marshal(dict)
		if err != nil {
			t.Fatal("Error encoding JSON:", err)
		}
	}

	req, err := http.NewRequest(method, path, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal("Error creating request:", err)
	}

	req.Header.Set("Content-Type", "application/json")

	return req
}

func runMigrations(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations",
		"postgres",
		driver,
	)

	// Накатываем миграции
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
