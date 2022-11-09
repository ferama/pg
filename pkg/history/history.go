package history

import (
	"errors"
	"sync"
)

const maxSize = 100

var (
	once     sync.Once
	instance *History
)

type History struct {
	currentIndex int
	list         []string

	lock sync.Mutex
}

func GetInstance() *History {
	once.Do(func() {
		instance = newHistory()
	})
	return instance
}

func newHistory() *History {
	return &History{
		currentIndex: -1,
		list:         make([]string, 0),
	}
}

func (h *History) GetList() []string {
	h.lock.Lock()
	defer h.lock.Unlock()

	return h.list
}

func (h *History) Append(item string) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if len(h.list)+1 > maxSize {
		h.list = h.list[1:]
		h.currentIndex--
	}

	if len(h.list) > 0 {
		if h.list[h.currentIndex] != item {
			h.list = append(h.list, item)
			h.currentIndex++
		}
	} else {
		h.list = append(h.list, item)
		h.currentIndex++
	}
}

func (h *History) GetPrev() (string, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.currentIndex-1 >= 0 && len(h.list) > 0 {
		h.currentIndex--
		return h.list[h.currentIndex], nil
	}
	return "", errors.New("do not have prev element")
}

func (h *History) GetNext() (string, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.currentIndex+1 < len(h.list) {
		h.currentIndex++
		return h.list[h.currentIndex], nil
	}
	return "", errors.New("do not have next element")
}
