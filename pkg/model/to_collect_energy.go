package model

type ToCollectEnergy struct {
	ID            int64 // 能量 ID，需要被某个用户采集
	UserId        string
	CollectEnergy int
	Status        string
}
