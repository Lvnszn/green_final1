package service

import (
	"fmt"
	"go.uber.org/atomic"
	"green_final1/pkg/engine"
	"green_final1/pkg/model"
	"math"
	"sync"
)

type memoryEnergy struct {
	e             engine.CommonEngine
	UsersSli      []*model.TotalEnergy
	Users         map[string]*model.TotalEnergy
	CollectEnergy map[int64]*model.ToCollectEnergy
	lock          sync.RWMutex
	//cnt           int
	cnt atomic.Int32
}

func (e *memoryEnergy) Collect(userId string, id int64) error {
	//e.lock.Lock()
	//defer e.lock.Unlock()
	e.cnt.Store(e.cnt.Inc())
	if e.cnt.Load()%20000 == 0 {
		fmt.Printf("handle request count: %d \n", e.cnt.Load())
	}

	//if e.cnt > 100_0000-10_00 {
	//	time.Sleep(1 * time.Second)
	//}
	user, totalFlag := e.Users[userId]
	cngine, toCollectFlag := e.CollectEnergy[id]
	if toCollectFlag && totalFlag {
		//if cngine.Status == "all_collected" {
		//	return nil
		//}
		if cngine.UserId == userId {
			//user.TotalEnergy += cngine.CollectEnergy
			tmp := e.UsersSli[user.Idx]
			tmp.TotalEnergyAtomic.Store(tmp.TotalEnergyAtomic.Add(int32(cngine.CollectEnergy)))
			//user.TotalEnergy = 666
			//cngine.Status = "all_collected"
			//cngine.CollectEnergy = 0
		} else {
			//if "collected_by_other" == cngine.Status {
			//	return nil
			//}
			collectedEnergy := int32(math.Floor(0.3 * float64(cngine.CollectEnergy)))
			tmp := e.UsersSli[user.Idx]
			tmp.TotalEnergyAtomic.Store(tmp.TotalEnergyAtomic.Add(collectedEnergy))
			//user.TotalEnergy += collectedEnergy
			//cngine.CollectEnergy -= collectedEnergy
			//cngine.Status = "collected_by_other"
		}
	}

	if e.cnt.Load() == 100_0000 {
		e.e.BulkUpdateTotal(e.Users)
		//time.Sleep(5 * time.Second)
	}

	return nil
}

func NewMemoryService() Collector {
	db := engine.NewMdb()
	service := &memoryEnergy{
		e:             db,
		UsersSli:      make([]*model.TotalEnergy, 100000),
		Users:         make(map[string]*model.TotalEnergy, 100000),
		CollectEnergy: make(map[int64]*model.ToCollectEnergy, 1000000),
	}

	tx1 := db.Begin()
	toCollectEnergy, err := tx1.Query("select id, to_collect_energy, user_id, status from to_collect_energy")
	if err != nil {
		fmt.Printf("query toCollectEnergy err is %v", err)
	}
	defer toCollectEnergy.Close()
	tx1.Commit()

	tx2 := db.Begin()
	totalEnergy, err := tx2.Query("select id, user_id, total_energy from total_energy")
	if err != nil {
		fmt.Printf("query totalEnergy err is %v", err)
	}
	tx1.Commit()
	defer totalEnergy.Close()

	var (
		collectEnergy, tEnergy    int
		tid, id                   int64
		toCollectUid, status, uid string
	)
	for toCollectEnergy.Next() {
		toCollectEnergy.Scan(&id, &collectEnergy, &toCollectUid, &status)
		service.CollectEnergy[id] = &model.ToCollectEnergy{
			ID:            id,
			CollectEnergy: collectEnergy,
			UserId:        toCollectUid,
		}
	}

	for totalEnergy.Next() {
		totalEnergy.Scan(&tid, &uid, &tEnergy)
		service.Users[uid] = &model.TotalEnergy{
			Idx:    tid % 100000,
			UserId: uid,
			//TotalEnergy: tEnergy,
		}
		service.UsersSli[tid] = &model.TotalEnergy{
			Idx:               tid % 100000,
			UserId:            uid,
			TotalEnergyAtomic: atomic.NewInt32(int32(tEnergy)),
		}
	}

	//db.DeleteAll("to_collect_energy")
	//db.DeleteAll("total_energy")
	return service
}
