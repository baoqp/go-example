package util

import (
	"time"
	"math/rand"
)

func Rand(up int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(up)
}


type TimeStat struct {
	m_llTime uint64
}

// TODO
func (timeStat *TimeStat) Point() int {
	return 0
}


