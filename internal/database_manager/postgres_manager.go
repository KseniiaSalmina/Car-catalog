package database_manager

import (
	"context"
	"fmt"

	"github.com/KseniiaSalmina/Car-catalog/internal/models"
	"github.com/KseniiaSalmina/Car-catalog/internal/storage/postgres"
)

type PostgresManager struct {
	db *postgres.DB
}

func NewPostgresManager(db *postgres.DB) *PostgresManager {
	return &PostgresManager{db: db}
}

func (pm *PostgresManager) DeleteCar(ctx context.Context, regNum string) error {
	tx, err := pm.db.NewTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete car: transaction error: %w", err)
	}

	if err := tx.DeleteCar(ctx, regNum); err != nil {
		defer tx.Rollback(ctx)
		return fmt.Errorf("failed to delete car: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to delete car: %w", err)
	}

	return nil
}

func (pm *PostgresManager) ChangeCar(ctx context.Context, car models.Car) error {
	tx, err := pm.db.NewTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to change car: transaction error: %w", err)
	}

	if err := tx.ChangeCar(ctx, car); err != nil {
		defer tx.Rollback(ctx)
		return fmt.Errorf("failed to change car: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to change car: %w", err)
	}

	return nil
}

func (pm *PostgresManager) AddCars(ctx context.Context, cars []models.Car) error {
	tx, err := pm.db.NewTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to add cars: transaction error: %w", err)
	}

	for _, car := range cars {
		if err := tx.NewCar(ctx, car); err != nil {
			defer tx.Rollback(ctx)
			return fmt.Errorf("failed to add cars: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to add cars: %w", err)
	}

	return nil
}

func (pm *PostgresManager) GetCars(ctx context.Context, filters models.Car, yearFilterMode string, orderByMode string, limit, offset int) (*models.CarsPage, error) {
	tx, err := pm.db.NewTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get cars: transaction error: %w", err)
	}
	defer tx.Rollback(ctx)

	cars, err := tx.FindCars(ctx, filters, yearFilterMode, orderByMode, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get cars: %w", err)
	}

	rowsAmount, err := tx.CountCars(ctx, filters, yearFilterMode)
	if err != nil {
		return nil, fmt.Errorf("failed to count cars: %w", err)
	}

	_ = tx.Commit(ctx)

	var pagesAmount int
	if rowsAmount%limit != 0 {
		pagesAmount = rowsAmount/limit + 1
	} else {
		pagesAmount = rowsAmount / limit
	}

	resultPage := models.CarsPage{
		Cars:        cars,
		Limit:       limit,
		PagesAmount: pagesAmount,
	}

	return &resultPage, nil
}
