package history

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"sync"

	"github.com/ferama/pg/pkg/conf"
)

const (
	maxSize     = 100
	historyFile = "history.db"
)

var (
	once     sync.Once
	instance *History
)

type toFile struct {
	Items []string `yaml:"items"`
}

type History struct {
	cursor int
	list   []string

	lock sync.Mutex
}

func GetInstance() *History {
	once.Do(func() {
		instance = newHistory()
		instance.load()
	})
	return instance
}

func newHistory() *History {
	return &History{
		cursor: -1,
		list:   make([]string, 0),
	}
}

func (h *History) save() {
	items := make([]string, len(h.list))
	copy(items, h.list)

	j := toFile{Items: items}
	b, err := json.Marshal(j)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	usr, _ := user.Current()
	historyFile := filepath.Join(usr.HomeDir, conf.ConfDir, historyFile)
	fi, err := os.Create(historyFile)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	defer fi.Close()
	fi.Write(b)
}

func (h *History) load() {
	usr, _ := user.Current()
	historyFile := filepath.Join(usr.HomeDir, conf.ConfDir, historyFile)
	jsonFile, err := os.Open(historyFile)
	if err != nil {
		return
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var parsed toFile
	json.Unmarshal(byteValue, &parsed)
	h.list = make([]string, len(parsed.Items))
	copy(h.list, parsed.Items)

	h.cursor = 0
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
		h.cursor--
	}

	if len(h.list) > 0 {
		if h.list[h.cursor] != item {
			h.list = append(h.list, item)
			h.cursor++
		}
	} else {
		h.list = append(h.list, item)
		h.cursor++
	}
	h.save()
}

// GetAdIdx returns the value at index without moving
// the hisotry cursor
func (h *History) GetAtIdx(idx int) (string, error) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if len(h.list) > idx {
		return h.list[idx], nil
	}
	return "", errors.New("do not have element ad idx")
}

func (h *History) GoPrev() (string, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.cursor-1 >= 0 && len(h.list) > 0 {
		h.cursor--
		return h.list[h.cursor], nil
	}
	return "", errors.New("do not have prev element")
}

func (h *History) GoNext() (string, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.cursor+1 < len(h.list) {
		h.cursor++
		return h.list[h.cursor], nil
	}
	return "", errors.New("do not have next element")
}
