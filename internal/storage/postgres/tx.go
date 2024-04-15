package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/KseniiaSalmina/Car-catalog/internal/models"
)

type Transaction struct {
	tx pgx.Tx
}

func (t *Transaction) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *Transaction) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

func (t *Transaction) DeleteCar(ctx context.Context, num string) error {
	var personID pgtype.Int8
	if err := t.tx.QueryRow(ctx, `DELETE FROM cars WHERE reg_num = $1 RETURNING owner_id`, num).Scan(&personID); err != nil {
		return fmt.Errorf("failed to get owner id: %w", err)
	}

	if err := t.deletePersonWithoutCars(ctx, personID.Int64); err != nil {
		return fmt.Errorf("failed to delete car: %w", err)
	}

	return nil
}

func (t *Transaction) deletePersonWithoutCars(ctx context.Context, personID int64) error {
	var ownersCarsAmount pgtype.Int8
	if err := t.tx.QueryRow(ctx, `SELECT COUNT(*) FROM cars WHERE owner_id = $1`, personID).Scan(&ownersCarsAmount); err != nil {
		return fmt.Errorf("failed to count owners cars: %w", err)
	}

	if ownersCarsAmount.Int64 == 0 {
		if _, err := t.tx.Exec(ctx, `DELETE FROM persons WHERE id = $1`, personID); err != nil {
			return fmt.Errorf("failed to delete person: %w", err)
		}
	}

	return nil
}

func (t *Transaction) ChangeCar(ctx context.Context, car models.Car) error {
	var isCarExist pgtype.Int8
	if err := t.tx.QueryRow(ctx, `SELECT COUNT(*) FROM cars WHERE reg_num=$1`, car.RegNum).Scan(&isCarExist); err != nil {
		return fmt.Errorf("failed to change car: %w", err)
	}
	if isCarExist.Int64 == 0 {
		return errors.New("car with this number does not exist")
	}

	ownerID, err := t.findOrCreatePerson(ctx, car.Owner)
	if err != nil {
		return fmt.Errorf("failed to change car: %w", err)
	}

	var oldOwnerID pgtype.Int8
	if err := t.tx.QueryRow(ctx, `SELECT owner_id FROM cars WHERE reg_num=$1`, car.RegNum).Scan(&oldOwnerID); err != nil {
		return fmt.Errorf("failed to change car: %w", err)
	}

	if int(oldOwnerID.Int64) != ownerID {
		if err := t.deletePersonWithoutCars(ctx, oldOwnerID.Int64); err != nil {
			return fmt.Errorf("failed to change car: %w", err)
		}
	}

	if _, err := t.tx.Exec(ctx, `UPDATE cars SET mark=$1, model=$2, year=$3, owner_id=$4 WHERE reg_num=$5`, car.Mark, car.Model, car.Year, ownerID, car.RegNum); err != nil {
		return fmt.Errorf("failed to change car: %w", err)
	}

	return nil
}

func (t *Transaction) findOrCreatePerson(ctx context.Context, person models.Person) (int, error) {
	var ownerID int

	ownerID, err := t.personID(ctx, person)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("failed to find or create person: %w", err)
		}

		ownerID, err = t.newPerson(ctx, person)
		if err != nil {
			return 0, fmt.Errorf("failed to find or create person: %w", err)
		}
	}

	return ownerID, nil
}

func (t *Transaction) personID(ctx context.Context, person models.Person) (int, error) {
	var personID pgtype.Int8
	if person.Patronymic == "" {
		if err := t.tx.QueryRow(ctx, `SELECT id FROM persons WHERE name = $1 AND surname =$2 AND patronymic IS NULL`, person.Name, person.Surname).Scan(&personID); err != nil {
			return 0, fmt.Errorf("failed to get owner id: %w", err)
		}
	} else {
		if err := t.tx.QueryRow(ctx, `SELECT id FROM persons WHERE name = $1 AND surname =$2 AND patronymic =$3`, person.Name, person.Surname, person.Patronymic).Scan(&personID); err != nil {
			return 0, fmt.Errorf("failed to get owner id: %w", err)
		}
	}

	return int(personID.Int64), nil
}

func (t *Transaction) newPerson(ctx context.Context, person models.Person) (int, error) {
	var id pgtype.Int8

	if person.Patronymic == "" {
		if err := t.tx.QueryRow(ctx, `INSERT INTO persons (name, surname) VALUES ($1, $2) RETURNING id`, person.Name, person.Surname).Scan(&id); err != nil {
			return 0, fmt.Errorf("failed to create person: %w", err)
		}
	} else {
		if err := t.tx.QueryRow(ctx, `INSERT INTO persons (name, surname, patronymic) VALUES ($1, $2, $3) RETURNING id`, person.Name, person.Surname, person.Patronymic).Scan(&id); err != nil {
			return 0, fmt.Errorf("failed to create person: %w", err)
		}
	}

	return int(id.Int64), nil
}

func (t *Transaction) NewCar(ctx context.Context, car models.Car) error {
	ownerID, err := t.findOrCreatePerson(ctx, car.Owner)
	if err != nil {
		return fmt.Errorf("failed to create car: %w", err)
	}

	if _, err := t.tx.Exec(ctx, `INSERT INTO cars (reg_num, mark, model, year, owner_id) VALUES ($1, $2, $3, $4, $5)`, car.RegNum, car.Mark, car.Model, car.Year, ownerID); err != nil {
		return fmt.Errorf("failed to create car: %w", err)
	}

	return nil
}

func (t *Transaction) FindCars(ctx context.Context, filters models.Car, yearFilterMode string, orderByMode string, limit, offset int) ([]models.Car, error) {
	sqlFilters, args := t.filterToSQL(false, filters, yearFilterMode)
	sql := fmt.Sprintf("%s ORDER BY cars.year %s LIMIT %d OFFSET %d", sqlFilters, orderByMode, limit, offset)

	rows, err := t.tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get cars: %w", err)
	}

	cars := make([]models.Car, 0, limit)
	for rows.Next() {
		var regNum, mark, model, ownerName, ownerSurname, ownerPatronymic pgtype.Text
		var year pgtype.Int8

		if err := rows.Scan(&regNum, &mark, &model, &year, &ownerName, &ownerSurname, &ownerPatronymic); err != nil {
			return nil, fmt.Errorf("failed to get cars: scan error: %w", err)
		}

		car := models.Car{
			RegNum: regNum.String,
			Mark:   mark.String,
			Model:  model.String,
			Year:   int(year.Int64),
			Owner: models.Person{
				Name:       ownerName.String,
				Surname:    ownerSurname.String,
				Patronymic: ownerPatronymic.String,
			},
		}
		cars = append(cars, car)
	}

	return cars, nil
}

func (t *Transaction) CountCars(ctx context.Context, filters models.Car, yearFilterMode string) (int, error) {
	sql, args := t.filterToSQL(true, filters, yearFilterMode)

	var amount pgtype.Int8
	if err := t.tx.QueryRow(ctx, sql, args...).Scan(&amount); err != nil {
		return 0, fmt.Errorf("failed to count cars: %w", err)
	}

	return int(amount.Int64), nil
}

func (t *Transaction) filterToSQL(isCount bool, filter models.Car, yearFilterMode string) (string, []interface{}) {
	emptyFilter := models.Car{}
	var sql string

	if !isCount {
		sql = `SELECT cars.reg_num,
       cars.mark, 
       cars.model, 
       cars.year, 
       persons.name, 
       persons.surname, 
       persons.patronymic 
FROM cars JOIN persons ON cars.owner_id = persons.id`
	} else {
		sql = `SELECT COUNT(*) FROM cars JOIN persons ON cars.owner_id = persons.id`
	}

	if filter == emptyFilter {
		return sql, nil
	}

	filters := make([]string, 0, 7)
	nextArg := 1
	args := make([]interface{}, 0, 7)

	if filter.RegNum != "" {
		filters = append(filters, fmt.Sprintf(" cars.reg_num = $%d ", nextArg))
		nextArg++
		args = append(args, filter.RegNum)
	}

	if filter.Mark != "" {
		filters = append(filters, fmt.Sprintf(" cars.mark = $%d ", nextArg))
		nextArg++
		args = append(args, filter.Mark)
	}

	if filter.Model != "" {
		filters = append(filters, fmt.Sprintf(" cars.model = $%d ", nextArg))
		nextArg++
		args = append(args, filter.Model)
	}

	if filter.Year != 0 {
		filters = append(filters, fmt.Sprintf(" cars.year %s $%d ", yearFilterMode, nextArg))
		nextArg++
		args = append(args, filter.Year)
	}

	if filter.Owner.Name != "" {
		filters = append(filters, fmt.Sprintf(" persons.name = $%d ", nextArg))
		nextArg++
		args = append(args, filter.Owner.Name)
	}

	if filter.Owner.Surname != "" {
		filters = append(filters, fmt.Sprintf(" persons.surname = $%d ", nextArg))
		nextArg++
		args = append(args, filter.Owner.Surname)
	}

	if filter.Owner.Patronymic != "" {
		if filter.Owner.Patronymic == "-" {
			filters = append(filters, " persons.patronymic IS NULL")
		} else {
			filters = append(filters, fmt.Sprintf(" persons.patronymic = $%d ", nextArg))
			args = append(args, filter.Owner.Patronymic)
		}
	}

	if len(filters) != 0 {
		sql += ` WHERE` + strings.Join(filters, "AND")
	}

	return sql, args
}
