package service

import (
	"database/sql"
	"fmt"
	"green_final1/pkg/engine"
	"green_final1/pkg/model"
	"math"
	"strings"
	"sync"
)

type energy struct {
	e             engine.CommonEngine
	Users         map[int64]*model.TotalEnergy
	CollectEnergy map[int64]*model.ToCollectEnergy
	UserLock      sync.RWMutex
	CollectLock   sync.RWMutex
}

func (e *energy) Collect(userId string, id int64) error {
	tx1 := e.e.Begin()
	toCollectEnergy, err := tx1.Query("select to_collect_energy, user_id, status from to_collect_energy where id = ?", id)
	if err != nil {
		fmt.Printf("query toCollectEnergy err is %v", err)
		return err
	}
	tx1.Commit()
	defer toCollectEnergy.Close()
	tx2 := e.e.Begin()
	totalEnergy, err := tx2.Query("select total_energy from total_energy where user_id = ?", userId)
	if err != nil {
		fmt.Printf("query totalEnergy err is %v", err)
		return err
	}
	tx1.Commit()
	defer totalEnergy.Close()

	tx := e.e.Begin()
	defer tx.Commit()
	if tx != nil {

		var (
			collectEnergy, tEnergy   int
			toCollectUid, status     string
			toCollectFlag, totalFlag bool
		)
		for toCollectEnergy.Next() {
			toCollectFlag = true
			toCollectEnergy.Scan(&collectEnergy, &toCollectUid, &status)
		}

		for totalEnergy.Next() {
			totalFlag = true
			totalEnergy.Scan(&tEnergy)
		}

		if toCollectFlag && totalFlag {
			if status == "all_collected" {
				return nil
			}
			if toCollectUid == userId {
				tEnergy += collectEnergy
				status = "all_collected"
				collectEnergy = 0
			} else {
				if "collected_by_other" == status {
					return nil
				}
				collectedEnergy := int(math.Floor(0.3 * float64(collectEnergy)))
				tEnergy += collectedEnergy
				collectEnergy -= collectedEnergy
				status = "collected_by_other"
			}
		}

		tx.Exec("update total_energy set total_energy = ? where user_id = ?", tEnergy, userId)
		tx.Exec("update to_collect_energy set to_collect_energy = ?, status = ? where id = ?", collectEnergy, status, id)

		//} else {
		//	e.e.Update("update total_energy set total_energy = ? where user_id = ?", tEnergy, userId)
		//	e.e.Update("update to_collect_energy set to_collect_energy = ?, status = ? where id = ?", collectEnergy, status, id)
	}
	return nil
}

type Collector interface {
	Collect(userId string, id int64) error
}

func NewService() Collector {
	db := engine.NewMdb()
	service := &energy{
		e:             db,
		Users:         make(map[int64]*model.TotalEnergy, 100000),
		CollectEnergy: make(map[int64]*model.ToCollectEnergy, 1000000),
	}

	return service
}

//
//func getTotalEnergy(res *sql.Rows) ([][]interface{}, []string, error) {
//	result := make(map[int64]*model.TotalEnergy, 100000)
//	cols, err := res.Columns()
//	if err != nil {
//		return nil, nil, err
//	}
//	colTypes, err := res.ColumnTypes()
//	if err != nil {
//		return nil, nil, err
//	}
//	colLen := len(cols)
//	rawRow := make([]interface{}, colLen)
//	dest := make([]interface{}, colLen)
//	for i := range rawRow {
//		dest[i] = new(sql.NullString)
//	}
//
//	for res.Next() {
//		if err = res.Scan(dest...); err != nil {
//			continue
//		}
//		row := make([]interface{}, colLen)
//		for i := 0; i < colLen; i++ {
//			v := dest[i].(*sql.NullString)
//			if !v.Valid {
//				row[i] = nil
//				continue
//			}
//			switch colTypes[i].DatabaseTypeName() {
//			case "NUMERIC":
//				idx := strings.LastIndex(v.String, ".")
//				if idx > 0 && len(v.String) > idx+5 {
//					row[i] = v.String[:idx+5]
//				} else {
//					row[i] = v.String
//				}
//			default:
//				row[i] = v.String
//			}
//		}
//		result = append(result, row)
//	}
//	if err = res.Close(); err != nil {
//		fmt.Errorf("result err %v", err)
//	}
//
//	return result, cols, nil
//}

func getToCollectEnergy(res *sql.Rows) ([][]interface{}, []string, error) {
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
