package id

import (
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"sync/atomic"
	"time"
)

var rval, rcounter uint32

func init() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rval = r.Uint32()
	rcounter = r.Uint32()
}

func New() string {
	var id [12]byte

	v := uint32(time.Now().Unix())
	binary.BigEndian.PutUint32(id[0:], v)

	binary.BigEndian.PutUint32(id[4:], rval)

	c := atomic.AddUint32(&rcounter, 1)
	binary.BigEndian.PutUint32(id[8:], c)

	return hex.EncodeToString(id[:])
}
