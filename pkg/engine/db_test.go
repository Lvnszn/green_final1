package engine

import (
	"fmt"
	"testing"
)

func TestQuery(t *testing.T) {
	engine := NewMdb()
	r, e := engine.Query("select * from total_energy limit 100")
	t.Error(e)
	val, header, err := query(r)
	t.Error(err)
	fmt.Println(val)
	fmt.Println(header)

	tx := engine.Begin()
	_, err = tx.Query("select * from total_energy limit 100")
	t.Error(err)
	_, err = tx.Query("select * from to_collect_energy limit 10")
	t.Error(err)

}
