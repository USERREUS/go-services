package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"order-service/internal/app/model"
	"order-service/internal/app/store"
	"time"

	"github.com/IBM/sarama"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.logRequest)
	s.router.HandleFunc("/order", s.handleOrderCreate()).Methods("POST")
	s.router.HandleFunc("/order", s.handleOrderFindAll()).Methods("GET")
	s.router.HandleFunc("/order/{id}", s.handleOrderFindOne()).Methods("GET")
}

func prod(message string) error {
	config := sarama.NewConfig()
	producer, err := sarama.NewAsyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		return err
	}
	defer producer.Close()

	producerMessage := &sarama.ProducerMessage{
		Topic: "order-topic",
		Value: sarama.StringEncoder(message),
	}
	producer.Input() <- producerMessage
	if err != nil {
		return err
	}

	return nil
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		var level logrus.Level
		switch {
		case rw.code >= 500:
			level = logrus.ErrorLevel
		case rw.code >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}
		logger.Logf(
			level,
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (s *server) handleOrderCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := &model.Model{}
		if err := json.NewDecoder(r.Body).Decode(m); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		url := fmt.Sprintf("http://localhost:8082/inventory/%s", m.ProductCode)
		resp, err := http.Get(url)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		record := &model.Model{}
		if err := json.NewDecoder(resp.Body).Decode(record); err != nil {
			s.error(w, r, http.StatusInternalServerError, errors.New("2"))
			return
		}

		if record.Count < m.Count {
			s.error(w, r, http.StatusInternalServerError, errors.New("3"))
			return
		}

		if err := s.store.Repository().Create(m); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		//PUT INVENTORY

		if err := prod(fmt.Sprintf("productCode: %s", m.ProductCode)); err != nil {
			s.error(w, r, http.StatusInternalServerError, errors.New("4"))
			return
		}

		s.respond(w, r, http.StatusCreated, m)
	}
}

func (s *server) handleOrderFindOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr, _ := mux.Vars(r)["id"]
		record, err := s.store.Repository().FindOne(idStr)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusFound, record)
	}
}

func (s *server) handleOrderFindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		records, err := s.store.Repository().FindAll()
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusFound, records)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
