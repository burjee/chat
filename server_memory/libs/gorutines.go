package libs

import (
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/viper"
)

func NewGoPool() *ants.Pool {
	pool_size := viper.GetInt("pool.size")
	pool, err := ants.NewPool(pool_size)
	if err != nil {
		panic(err)
	}
	return pool
}
