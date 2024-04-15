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

// @Summary Delete car
// @Tags cars
// @Description delete car by its registration number
// @Param regNum path string true "registration number of the car"
// @Success 200
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /cars/{regNum} [delete]
func (s *Server) DeleteCar(w http.ResponseWriter, r *http.Request) {
	regNum, ok := bunrouter.ParamsFromContext(r.Context()).Get("regNum")
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

// @Summary Post cars
// @Tags cars
// @Description get info about cars from outsider service and put it to the database
// @Accept json
// @Param regNums body []string true "array of regNums"
// @Success 200
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /cars [post]
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

// @Summary Patch car
// @Tags cars
// @Description update car info by accepted json
// @Accept json
// @Param regNum path string true "registration number of the car"
// @Param car body models.Car true "car with updated info"
// @Success 200
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /cars/{regNum} [patch]
func (s *Server) PatchCar(w http.ResponseWriter, r *http.Request) {
	regNum, ok := bunrouter.ParamsFromContext(r.Context()).Get("regNum")
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

// @Summary Get cars
// @Tags cars
// @Description get cars by filters. Filters are accepted in json format as a car structure
// @Produce json
// @Accept json
// @Param page query int false "page number"
// @Param limit query int false "limit records by page"
// @Param yearFilterMode query string false "can be =, >=, <=. By default =" Enums(=, >=, <=)
// @Param orderByMode query string false "relating to the car year. Can be ASC or DESC, by default DESC" Enums(ASC, DESC)
// @Param filters body models.Car true "empty fields do not affect the result. To search for car owners without a patronymic, the patronymic field must contain the string "-""
// @Success 200 {object} models.CarsPage
// @Success 204
// @Failure 400 {string} string
// @Failure 500	{string} string
// @Router /cars [get]
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
	if yearFilterMode != "=" && yearFilterMode != ">=" && yearFilterMode != "<=" {
		http.Error(w, fmt.Sprintf("incorrect year filter mode: should be =, >= or <= %s", err.Error()), http.StatusBadRequest)
		return
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
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	carsPage.PageNo = page.pageNumber

	_ = json.NewEncoder(w).Encode(carsPage)
}
