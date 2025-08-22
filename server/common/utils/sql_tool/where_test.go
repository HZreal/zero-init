package sql_tool

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		id        int
		whereList []any
		query     string
		args      []any
		hasErr    bool
	}{
		{
			id:        1,
			whereList: []any{"id=1"},
			query:     "id=1",
			args:      []any{},
		},
		{
			id:        2,
			whereList: []any{"id=1", "status=2"},
			query:     "id=1 AND status=2",
			args:      []any{},
		},
		{
			id: 3,
			whereList: []any{
				[]any{"id=?", []any{1}},
			},
			query: "id=?",
			args:  []any{1},
		},
		{
			id: 4,
			whereList: []any{
				[]any{"id=?", []any{1}},
				[]any{"name=?", []any{"zhangsan"}},
			},
			query: "id=? AND name=?",
			args:  []any{1, "zhangsan"},
		},
		{
			id: 5,
			whereList: []any{
				[]any{"id=? AND name=?", []any{1, "zhangsan"}},
			},
			query: "id=? AND name=?",
			args:  []any{1, "zhangsan"},
		},
		{
			id: 6,
			whereList: []any{
				[]any{"id=? AND name=?", []any{1, "zhangsan"}, "xxx"},
			},
			query:  "id=? AND name=?",
			args:   []any{1, "zhangsan"},
			hasErr: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(fmt.Sprint(test.id), func(t *testing.T) {
			t.Parallel()

			//query, args, err := GetWhere(test.whereList)
			//query, args, err := parseWhereList(test.whereList, "AND")
			where, err := GetAndWhere(test.whereList, false)
			query := where.Where
			args := where.Args
			if test.hasErr {
				fmt.Println(err)
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, test.query, query)
				assert.Equal(t, test.args, args)
			}
		})
	}
}

// func TestWriteValue(t *testing.T) {
// 	var buf strings.Builder
// 	tm := time.Now()
// 	writeValue(&buf, &tm)
// 	assert.Equal(t, "'"+tm.String()+"'", buf.String())

// 	buf.Reset()
// 	writeValue(&buf, tm)
// 	assert.Equal(t, "'"+tm.String()+"'", buf.String())
// }

func TestWhere(t *testing.T) {
	whereList := []WhereItem{}

	whereList = append(whereList, WhereItem{
		Field:     "customer_id",
		Opterater: SQLOptE,
		Value:     1,
	})

	whereList = append(whereList, WhereItem{
		Field:     "id",
		Opterater: SQLOptIn,
		Value:     []any{1, 2, 3},
	})

	whereList = append(whereList, WhereItem{
		Field:     "name",
		Opterater: SQLOptLike,
		Value:     "aaa",
	})

	whereList = append(whereList, WhereItem{
		Field:     "",
		Opterater: SQLOptInvalid,
		SubWhere: []WhereItem{
			{
				Field:     "agent_id",
				Opterater: SQLOptE,
				Value:     2,
			},
			{
				Field:     "agent_id",
				Opterater: SQLOptE,
				Value:     5,
			},
		},
	})

	sql, args := GetWhereSet(&whereList, false)
	fmt.Printf("%v,%v", sql, args)
}

func TestDataset(t *testing.T) {
	dataset := map[string]any{}

	dataset["name"] = "2"
	dataset["remark"] = "55"
	dataset["status"] = 1

	sql, args := GetUpdateSet(&dataset)
	fmt.Printf("%v,%v", sql, args)
}

func TestInsert(t *testing.T) {
	type TblBillingPackage struct {
		Id         uint64 `db:"id"`
		CustomerId uint64 `db:"customer_id"`
		Name       string `db:"name"`
		Remark     string `db:"remark"`
		FatherId   uint64 `db:"father_id"`
		Ctime      uint64 `db:"ctime"`
		Mtime      uint64 `db:"mtime"`
		Dtime      uint64 `db:"dtime"`
		Cby        uint64 `db:"cby"`
		Mby        uint64 `db:"mby"`
		Dby        uint64 `db:"dby"`
		DelState   uint64 `db:"del_state"`
	}

	data := &TblBillingPackage{
		Id:         1,
		CustomerId: 2,
		Name:       "1",
		Remark:     "2",
		FatherId:   3,
		Ctime:      4,
		Mtime:      5,
		Dtime:      6,
		Cby:        7,
		Mby:        8,
		Dby:        9,
		DelState:   10,
	}

	s, v := GetInsertSet("test", false, false, data)

	fmt.Printf("\n%s\n%s\n", s, v)
}
