package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/uptrace/bunrouter"

	"github.com/KseniiaSalmina/Car-catalog/internal/models"
)

func (s *Server) DeleteCar(w http.ResponseWriter, r *http.Request) {
	regNum, ok := bunrouter.ParamsFromContext(r.Context()).Get("id")
	if !ok {
		http.Error(w, "registration number of the car is required", http.StatusBadRequest)
		return
	}

	if err := s.dbManager.DeleteCar(r.Context(), regNum); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) PostCars(w http.ResponseWriter, r *http.Request) {
	regNums := make([]string, 0, 10)
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&regNums); err != nil {
		http.Error(w, fmt.Sprintf("incorrect car registration numbers: %s", err.Error()), http.StatusBadRequest)
		return
	}

	cars := make([]models.Car, 0, len(regNums))
	for _, regNum := range regNums {
		car, err := s.infoReceiver.GetCarInfo(r.Context(), regNum)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cars = append(cars, *car)
	}

	if err := s.dbManager.AddCars(r.Context(), cars); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) PatchCar(w http.ResponseWriter, r *http.Request) {
	regNum, ok := bunrouter.ParamsFromContext(r.Context()).Get("id")
	if !ok {
		http.Error(w, "registration number of the car is required", http.StatusBadRequest)
		return
	}

	var car models.Car
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		http.Error(w, fmt.Sprintf("incorrect car data: %s", err.Error()), http.StatusBadRequest)
		return
	}
	car.RegNum = regNum

	if err := s.dbManager.ChangeCar(r.Context(), car); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetCars(w http.ResponseWriter, r *http.Request) {
	page, err := s.getPageInfo(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("incorrect page info: %s", err.Error()), http.StatusBadRequest)
		return
	}

	yearFilterMode := r.FormValue("yearFilterMode")
	if yearFilterMode == "" {
		yearFilterMode = "="
	}

	orderByMode := r.FormValue("orderByMode")
	if orderByMode == "" {
		orderByMode = "DESC"
	}
	if orderByMode != "DESC" && orderByMode != "ASC" {
		http.Error(w, "incorrect order by mode: should be ASC or DESC", http.StatusBadRequest)
		return
	}

	var car models.Car
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		http.Error(w, fmt.Sprintf("incorrect car data: %s", err.Error()), http.StatusBadRequest)
		return
	}

	carsPage, err := s.dbManager.GetCars(r.Context(), car, yearFilterMode, orderByMode, page.limit, page.offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	carsPage.PageNo = page.pageNumber

	_ = json.NewEncoder(w).Encode(carsPage)
}
