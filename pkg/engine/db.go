package engine

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_const "green_final1/pkg/const"
	"green_final1/pkg/model"
	"strconv"
	"strings"
)

type CommonEngine interface {
	Begin() *sql.Tx
	Query(sql string, args ...any) (*sql.Rows, error)
	BulkInsertTotal(unsavedRows map[string]*model.TotalEnergy) error
	BulkUpdateTotal(unsavedRows map[string]*model.TotalEnergy) error
	BulkUpdateTotalSlice(unsavedRows []*model.TotalEnergy) error
	BulkInsertCollect(unsavedRows map[int64]*model.ToCollectEnergy) error
	DeleteAll(name string) error
	Update(sql string, args ...any) error
}

func (m *mysqlEngine) DeleteAll(name string) error {
	r, err := m.Db.Exec(fmt.Sprintf("delete from %s", name))
	if err != nil {
		return err
	}
	affected, err := r.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Printf("delete %d records", affected)
	return nil
}

type mysqlEngine struct {
	Db *sql.DB
}

func (m *mysqlEngine) BulkUpdateTotal(unsavedRows map[string]*model.TotalEnergy) error {
	for _, post := range unsavedRows {
		m.Db.Exec("update total_energy set total_energy = ? where user_id = ?", post.TotalEnergyAtomic.Load(), post.UserId)
	}

	return nil
}

func (m *mysqlEngine) BulkUpdateTotalSlice(unsavedRows []*model.TotalEnergy) error {
	energies := map[int][]string{}
	for _, te := range unsavedRows {
		if te != nil {
			total := te.TotalEnergy
			energies[total] = append(energies[total], te.UserId)
		}
	}
	sqlBuffer := bytes.Buffer{}
	for total, userIds := range energies {
		sqlBuffer.WriteString("update `total_energy` set `total_energy` = " + strconv.Itoa(total) + " where user_id in (")
		for _, userId := range userIds {
			sqlBuffer.WriteString("'" + userId + "',")
		}
		sqlBuffer.Truncate(sqlBuffer.Len() - 1)
		sqlBuffer.WriteString(");")
	}
	m.Db.Exec(sqlBuffer.String())
	//for _, post := range unsavedRows {
	//	m.Db.Exec("update total_energy set total_energy = ? where user_id = ?", post.TotalEnergyAtomic.Load(), post.UserId)
	//}
	return nil
}

func (m *mysqlEngine) Begin() *sql.Tx {
	b, _ := m.Db.Begin()
	return b
}

func (m *mysqlEngine) BulkInsertTotal(unsavedRows map[string]*model.TotalEnergy) error {
	//valueStrings := make([]string, 0, len(unsavedRows))
	//valueArgs := make([]interface{}, 0, len(unsavedRows)*3)
	//for _, post := range unsavedRows {
	//	valueStrings = append(valueStrings, "(?, ?)")
	//	valueArgs = append(valueArgs, post.UserId)
	//	//valueArgs = append(valueArgs, post.TotalEnergy)
	//}
	//stmt := fmt.Sprintf("INSERT INTO total_energy (user_id, total_energy) VALUES %s",
	//	strings.Join(valueStrings, ","))
	//_, err := m.Db.Exec(stmt, valueArgs...)
	//return err
	return nil
}

func (m *mysqlEngine) BulkInsertCollect(unsavedRows map[int64]*model.ToCollectEnergy) error {
	valueStrings := make([]string, 0, len(unsavedRows))
	valueArgs := make([]interface{}, 0, len(unsavedRows)*3)
	for _, post := range unsavedRows {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, post.ID)
		valueArgs = append(valueArgs, post.UserId)
		valueArgs = append(valueArgs, post.CollectEnergy)
		valueArgs = append(valueArgs, post.Status)
	}
	stmt := fmt.Sprintf("INSERT INTO to_collect_energy (id, user_id, collect_energy, status) VALUES %s",
		strings.Join(valueStrings, ","))
	_, err := m.Db.Exec(stmt, valueArgs...)
	return err
}

func (m *mysqlEngine) Update(sql string, args ...any) error {
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//tx, _ := m.Db.BeginTx(ctx, nil)
	_, err := m.Db.Exec(sql, args...)
	if err != nil {
		fmt.Printf("exec err %v\n", err)
	}
	//cancel()
	return err
}

func (m *mysqlEngine) Batch(sql string, rows []interface{}) error {
	_, err := m.Db.Exec(sql, rows...)
	if err != nil {
		return err
	}

	return nil
}

func (m *mysqlEngine) Query(sql string, args ...any) (*sql.Rows, error) {
	return m.Db.Query(sql, args...)
}

func NewMdb() CommonEngine {
	db, err := sql.Open("mysql", _const.DBPATH)
	if err != nil {
		panic(err)
	}
	fmt.Printf("connect db success \n")
	db.SetMaxIdleConns(1200)
	db.SetMaxOpenConns(1200)
	return &mysqlEngine{
		Db: db,
	}
}

func query(res *sql.Rows) ([][]interface{}, []string, error) {
	result := make([][]interface{}, 0, 100)
	cols, err := res.Columns()
	if err != nil {
		return nil, nil, err
	}
	colTypes, err := res.ColumnTypes()
	if err != nil {
		return nil, nil, err
	}
	colLen := len(cols)
	rawRow := make([]interface{}, colLen)
	dest := make([]interface{}, colLen)
	for i := range rawRow {
		dest[i] = new(sql.NullString)
	}

	for res.Next() {
		if err = res.Scan(dest...); err != nil {
			continue
		}
		row := make([]interface{}, colLen)
		for i := 0; i < colLen; i++ {
			v := dest[i].(*sql.NullString)
			if !v.Valid {
				row[i] = nil
				continue
			}
			switch colTypes[i].DatabaseTypeName() {
			case "NUMERIC":
				idx := strings.LastIndex(v.String, ".")
				if idx > 0 && len(v.String) > idx+5 {
					row[i] = v.String[:idx+5]
				} else {
					row[i] = v.String
				}
			default:
				row[i] = v.String
			}
		}
		result = append(result, row)
	}
	if err = res.Close(); err != nil {
		fmt.Errorf("result err %v", err)
	}

	return result, cols, nil
}
