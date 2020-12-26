package client

import (
	"errors"
	"github.com/vmihailenco/msgpack"
)

type Params struct {
	FunName string
	Data    *[]byte
}
type Consumer struct {
	Funs  map[string]func(*[]byte) ([]byte, error)
	index []string
}
type Index struct {
	Index []string
}

func (c *Consumer) Deal(p Params, ret *[]byte) error {
	var err error
	if c.checkSafe(p.FunName) {
		*ret, err = c.Funs[p.FunName](p.Data)
	} else {
		err = errors.New("Not Find: " + p.FunName)
	}
	return err
}

func (c *Consumer) GetIndex(_ Params, ret *[]byte) error {
	if c.index == nil {
		i := 0
		for range c.Funs {
			i++
		}
		index := make([]string, i)
		for k := range c.Funs {
			index = append(index, k)
		}
		c.index = index
		var err error
		*ret, err = msgpack.Marshal(&Index{Index: index})
		return err
	} else {
		var err error
		*ret, err = msgpack.Marshal(&Index{Index: c.index})
		return err
	}
}

func (c *Consumer) checkSafe(FunName string) bool {
	for k := range c.Funs {
		if k == FunName {
			return true
		}
	}
	return false
}
