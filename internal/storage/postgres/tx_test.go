package postgres

import (
	"testing"

	"github.com/KseniiaSalmina/Car-catalog/internal/models"
)

func TestTransaction_filterToSQL(t1 *testing.T) {
	type args struct {
		isCount        bool
		filter         models.Car
		yearFilterMode string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "filter by regNum, not count", args: args{isCount: false, filter: models.Car{RegNum: "testREGnum"}, yearFilterMode: ""},
			want: `SELECT cars.reg_num,
       cars.mark, 
       cars.model, 
       cars.year, 
       persons.name, 
       persons.surname, 
       persons.patronymic 
FROM cars JOIN persons ON cars.owner_id = persons.id WHERE cars.reg_num =testREGnum `},
		{name: "filter by mark, count", args: args{isCount: true, filter: models.Car{Mark: "testMark"}, yearFilterMode: ""},
			want: `SELECT (*) FROM cars JOIN persons ON cars.owner_id = persons.id WHERE cars.mark =testMark `},
		{name: "filter by no patronymic, count", args: args{isCount: true, filter: models.Car{Owner: models.Person{Patronymic: "-"}}, yearFilterMode: "="},
			want: `SELECT (*) FROM cars JOIN persons ON cars.owner_id = persons.id WHERE persons.patronymic IS NULL`},
		{name: "filter by year, year filter mode more or equal, count", args: args{isCount: true, filter: models.Car{Year: 2015}, yearFilterMode: ">="},
			want: `SELECT (*) FROM cars JOIN persons ON cars.owner_id = persons.id WHERE cars.year >=2015 `},
		{name: "a lot of filters, count", args: args{isCount: true, filter: models.Car{RegNum: "testNUM", Model: "testModel", Mark: "testMark", Year: 2022, Owner: models.Person{Name: "Ivan"}}, yearFilterMode: "="},
			want: `SELECT (*) FROM cars JOIN persons ON cars.owner_id = persons.id WHERE cars.reg_num =testNUM AND cars.mark =testMark AND cars.model =testModel AND cars.year =2022 AND persons.name =Ivan `},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Transaction{}
			if got := t.filterToSQL(tt.args.isCount, tt.args.filter, tt.args.yearFilterMode); got != tt.want {
				t1.Errorf("filterToSQL(filter %v, yearFilterMode %s)\n got %v, \nwant %v", tt.args.filter, tt.args.yearFilterMode, got, tt.want)
			}
		})
	}
}
