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
	type res struct {
		sql  string
		args []interface{}
	}
	tests := []struct {
		name string
		args args
		want res
	}{
		{name: "filter by regNum, not count", args: args{isCount: false, filter: models.Car{RegNum: "testREGnum"}, yearFilterMode: ""},
			want: res{sql: `SELECT cars.reg_num,
       cars.mark, 
       cars.model, 
       cars.year, 
       persons.name, 
       persons.surname, 
       persons.patronymic 
FROM cars JOIN persons ON cars.owner_id = persons.id WHERE cars.reg_num = $1 `, args: make([]interface{}, 1)}},
		{name: "filter by mark, count", args: args{isCount: true, filter: models.Car{Mark: "testMark"}, yearFilterMode: ""},
			want: res{sql: `SELECT COUNT(*) FROM cars JOIN persons ON cars.owner_id = persons.id WHERE cars.mark = $1 `, args: make([]interface{}, 1)}},
		{name: "filter by no patronymic, count", args: args{isCount: true, filter: models.Car{Owner: models.Person{Patronymic: "-"}}, yearFilterMode: "="},
			want: res{sql: `SELECT COUNT(*) FROM cars JOIN persons ON cars.owner_id = persons.id WHERE persons.patronymic IS NULL`, args: make([]interface{}, 0)}},
		{name: "filter by year, year filter mode more or equal, count", args: args{isCount: true, filter: models.Car{Year: 2015}, yearFilterMode: ">="},
			want: res{sql: `SELECT COUNT(*) FROM cars JOIN persons ON cars.owner_id = persons.id WHERE cars.year >= $1 `, args: make([]interface{}, 1)}},
		{name: "a lot of filters, count", args: args{isCount: true, filter: models.Car{RegNum: "testNUM", Model: "testModel", Mark: "testMark", Year: 2022, Owner: models.Person{Name: "Ivan"}}, yearFilterMode: "="},
			want: res{sql: `SELECT COUNT(*) FROM cars JOIN persons ON cars.owner_id = persons.id WHERE cars.reg_num = $1 AND cars.mark = $2 AND cars.model = $3 AND cars.year = $4 AND persons.name = $5 `, args: make([]interface{}, 5)}},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Transaction{}
			gotSql, gotArgs := t.filterToSQL(tt.args.isCount, tt.args.filter, tt.args.yearFilterMode)
			if gotSql != tt.want.sql {
				t1.Errorf("filterToSQL(filter %v, yearFilterMode %s) SQL\n got %v, \nwant %v", tt.args.filter, tt.args.yearFilterMode, gotSql, tt.want)
			}
			if len(gotArgs) != len(tt.want.args) {
				t1.Errorf("filterToSQL(filter %v, yearFilterMode %s) ARGS \n got %v, \nwant %v", tt.args.filter, tt.args.yearFilterMode, gotSql, tt.want)
			}
		})
	}
}
