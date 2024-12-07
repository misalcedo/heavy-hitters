package misra_gries

import (
	"container/list"
	"errors"
	"math"
)

type pair struct {
	key   string
	value int
}

type MisraGries struct {
	epsilon  float64
	capacity int
	pairs    list.List
	hits     float64
}

func NewMisraGries(epsilon float64) (*MisraGries, error) {
	if epsilon < 0 || epsilon > 1 {
		return nil, errors.New("epsilon should be between 0 and 1")
	}

	return &MisraGries{
		epsilon:  epsilon,
		capacity: int(math.Ceil(1.0 / epsilon)),
	}, nil
}

func (m *MisraGries) Hits() int {
	return int(m.hits)
}

func (m *MisraGries) Query(key string) (int, int) {
	element := m.find(key)
	err := int(math.Ceil(m.epsilon * m.hits))

	var value int

	if element != nil {
		value = element.Value.(pair).value
	}

	return value, value + err
}

func (m *MisraGries) Hit(key string) {
	m.hits += 1.0

	element := m.find(key)

	if element == nil {
		if m.pairs.Len() < m.capacity {
			m.pairs.PushBack(pair{key, 1})
		} else {
			m.decrement()
		}
	} else {
		p := element.Value.(pair)
		p.value += 1

		element.Value = p
	}
}

func (m *MisraGries) find(key string) *list.Element {
	for e := m.pairs.Front(); e != nil; e = e.Next() {
		if e.Value.(pair).key == key {
			return e
		}
	}

	return nil
}

func (m *MisraGries) decrement() {
	for e := m.pairs.Front(); e != nil; {
		p := e.Value.(pair)
		p.value -= 1

		next := e.Next()

		if p.value <= 0 {
			m.pairs.Remove(e)
		} else {
			e.Value = p
		}

		e = next
	}
}
