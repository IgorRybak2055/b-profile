package bamboo

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/IgorRybak2055/bamboo/internal/repository"
	"github.com/IgorRybak2055/bamboo/internal/services"
)

// HTTPConfig stores configs for bamboo.
type HTTPConfig struct {
	Addr     string `config:"HTTP_ADDR,required"`
	LogLevel string `config:"LOG_LEVEL,required"`
}

// App is a structure for bamboo server.
type App struct {
	cfg    *HTTPConfig
	Logger *logrus.Logger
	Srv    *http.Server
	DBC    *sqlx.DB

	profileService services.Profile
}

// New returns App with HTTPConfig and logger for start bamboo.
func New(cfg *HTTPConfig) *App {
	return &App{
		cfg:    cfg,
		Logger: logrus.New(),
	}
}

// Start configures all needs App fields and start bamboo.
func (a *App) Start() error {
	a.profileService = services.NewProfileService(repository.NewAccountRepository(a.DBC), a.Logger)

	router := mux.NewRouter()
	router.Use(recovery)

	router.HandleFunc("/health", a.handleHealth).Methods(http.MethodGet)

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/order", a.handleOrder).Methods(http.MethodPut)

	a.Srv = &http.Server{
		Handler: router,
		Addr:    a.cfg.Addr,
	}

	a.Logger.Info("starting ...")

	return a.Srv.ListenAndServe()
}

// Respond does answer with data on API request.
func Respond(w http.ResponseWriter, statusCode int, data [][]string) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=TheCSVFileName.csv")

	w.WriteHeader(statusCode)

	wr := csv.NewWriter(w)

	for _, v := range data {
		if err := wr.Write(v); err != nil {
			log.Printf("failed to write data {%s} response: %s", v, err)
		}
	}

	wr.Flush()

	log.Println("success respond")
}

// RespondError does answer with error on API request.
func RespondError(w http.ResponseWriter, err Error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(err.Code)

	if err := json.NewEncoder(w).Encode(err); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
