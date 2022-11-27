package service

import (
	"bytes"
	"fmt"
	"go.uber.org/atomic"
	"green_final1/pkg/engine"
	"green_final1/pkg/model"
	"math"
	"strconv"
	"sync"
	"time"
)

type memoryEnergy struct {
	e         engine.CommonEngine
	UsersSli  []*model.TotalEnergy
	Collected []*model.ToCollectEnergy
	Users     map[string]*model.TotalEnergy
	//CollectEnergy map[int64]*model.ToCollectEnergy
	lock          sync.RWMutex
	acnt          int
	cnt, aId, cId atomic.Int32
}

func (e *memoryEnergy) Collect(userId string, id int64) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	//e.cnt.CompareAndSwap(-1, e.cnt.Inc())
	//cnt := e.cnt.Load()
	e.acnt++
	if e.acnt%20000 == 0 {
		fmt.Printf(time.Now().Format("2006-01-02 15:04:05")+" "+"handle request count: %d, userId: %v, id %v\n", e.acnt, userId, id)
	}
	//if e.cnt > 100_0000-10_00 {
	//	time.Sleep(1 * time.Second)
	//}
	if id < 1 {
		return nil
	}
	user, totalFlag := e.Users[userId]
	cngine := e.Collected[id]
	if cngine.Status == "" && totalFlag {
		if cngine.Status == "all_collected" {
			return nil
		}
		if cngine.UserId == userId {
			//user.TotalEnergy += cngine.CollectEnergy
			//tmp := e.UsersSli[user.Idx]
			//tmp.TotalEnergyAtomic
			e.UsersSli[user.Idx].TotalEnergy += cngine.CollectEnergy
			//tmp.TotalEnergyAtomic.CompareAndSwap(0, tmp.TotalEnergyAtomic.Add(int32(cngine.CollectEnergy)))
			//tmp.TotalEnergyAtomic.CompareAndSwap(-1, tmp.TotalEnergyAtomic.Add(int32(cngine.CollectEnergy)))
			//user.TotalEnergy = 666
			//cngine.Status = "all_collected"
			//cngine.CollectEnergy = 0
			e.Collected[id].Status = "all_collected"
		} else {
			if "collected_by_other" == cngine.Status {
				return nil
			}
			collectedEnergy := int(math.Floor(0.3 * float64(cngine.CollectEnergy)))
			//tmp := e.UsersSli[user.Idx]
			//tmp.TotalEnergyAtomic.CompareAndSwap(-1, tmp.TotalEnergyAtomic.Add(collectedEnergy))
			e.UsersSli[user.Idx].TotalEnergy += collectedEnergy
			//cngine.CollectEnergy -= collectedEnergy
			//cngine.Status = "collected_by_other"
			e.Collected[id].Status = "collected_by_other"
		}
	}

	if e.acnt == 100_0000 {
		//go func() {
		//}()
		updateZeroBuffer := bytes.Buffer{}
		//updateSevenBuffer := bytes.Buffer{}
		//zIds := make([]int64, 0, 500000)
		//sIds := make([]int64, 0, 500000)
		updateZeroBuffer.WriteString("update to_collect_energy set to_collect_energy = 0 where id in (")
		//updateSevenBuffer.WriteString("update to_collect_energy set to_collect_energy = ceil(to_collect_energy*0.7) where id in (")
		for _, v := range e.Collected {
			if v == nil {
				continue
			}
			if v.Status == "collected_by_other" {
				//updateSevenBuffer.WriteString("" + strconv.FormatInt(v.ID, 10) + ",")
				//sIds = append(sIds, v.ID)
			} else if v.Status == "all_collected" {
				updateZeroBuffer.WriteString("" + strconv.FormatInt(v.ID, 10) + ",")
				//zIds = append(zIds, v.ID)
			} else {
				e.e.Update("update to_collect_energy set to_collect_energy = ? where id = ?", v.CollectEnergy, v.ID)
			}
		}

		updateZeroBuffer.Truncate(updateZeroBuffer.Len() - 1)
		updateZeroBuffer.WriteString(");")

		//updateSevenBuffer.Truncate(updateSevenBuffer.Len() - 1)
		//updateSevenBuffer.WriteString(");")
		go func() {
			fmt.Printf(time.Now().Format("2006-01-02 15:04:05")+" "+"updateZeroBuffer.String(): %v \n", updateZeroBuffer.String())
			e.e.Update(updateZeroBuffer.String())
		}()
		e.e.BulkUpdateTotalSlice(e.UsersSli)

		//go func() {
		//	defer wg.Done()
		//	fmt.Printf(time.Now().Format("2006-01-02 15:04:05")+" "+"updateSevenBuffer.String(): %d \n", updateSevenBuffer.String())
		//	e.e.Update(updateSevenBuffer.String())
		//}()
		//time.Sleep(5 * time.Second)
	}

	return nil
}

func NewMemoryService() Collector {
	db := engine.NewMdb()
	service := &memoryEnergy{
		e:        db,
		UsersSli: make([]*model.TotalEnergy, 100000),
		Users:    make(map[string]*model.TotalEnergy, 100000),
		//CollectEnergy: make(map[int64]*model.ToCollectEnergy, 1000000),
		Collected: make([]*model.ToCollectEnergy, 1000001),
	}

	toCollectEnergy, err := db.Query("select id, to_collect_energy, user_id, status from to_collect_energy")
	if err != nil {
		fmt.Printf("query toCollectEnergy err is %v", err)
	}
	defer toCollectEnergy.Close()

	totalEnergy, err := db.Query("select id, user_id, total_energy from total_energy")
	if err != nil {
		fmt.Printf("query totalEnergy err is %v", err)
	}
	defer totalEnergy.Close()

	var (
		collectEnergy, tEnergy    int
		tid, id                   int64
		toCollectUid, status, uid string
	)
	for toCollectEnergy.Next() {
		toCollectEnergy.Scan(&id, &collectEnergy, &toCollectUid, &status)
		//service.CollectEnergy[id] = &model.ToCollectEnergy{
		//	ID:            id,
		//	CollectEnergy: collectEnergy,
		//	UserId:        toCollectUid,
		//}
		service.Collected[id] = &model.ToCollectEnergy{
			ID:            id,
			CollectEnergy: collectEnergy,
			UserId:        toCollectUid,
		}
	}

	for totalEnergy.Next() {
		totalEnergy.Scan(&tid, &uid, &tEnergy)
		service.Users[uid] = &model.TotalEnergy{
			Idx:         tid,
			UserId:      uid,
			TotalEnergy: tEnergy,
		}
		service.UsersSli[tid] = &model.TotalEnergy{
			Idx:         tid,
			UserId:      uid,
			TotalEnergy: tEnergy,
			//TotalEnergyAtomic: atomic.NewInt32(int32(tEnergy)),
		}
	}

	service.e.Update("update to_collect_energy set to_collect_energy = ceil(to_collect_energy*0.7)")

	//db.DeleteAll("to_collect_energy")
	//db.DeleteAll("total_energy")
	return service
}
