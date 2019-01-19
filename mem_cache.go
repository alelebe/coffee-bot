package main

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/memcachier/mc"
)

const (
	ordersKey = "ORDERS"
	expTime   = 12 * 3600 // 12 hours
	expTime2  = 2 * 3600  // 2 hours to collect after 1st collection attemp

	watchersKey = "WATCHERS"
)

var mcClient *mc.Client

func initMemCache() *mc.Client {
	username := getVar("MEMCACHIER_USERNAME")
	password := getVar("MEMCACHIER_PASSWORD")
	servers := getVar("MEMCACHIER_SERVERS")

	mcClient = mc.NewMC(servers, username, password)
	return mcClient
}

// CoffeeOrder :
type CoffeeOrder struct {
	UserID    int
	UserName  string
	ChatID    int64
	Beverage  string
	Price     float64
	OrderTime time.Time
}

func placeOrder(order CoffeeOrder) bool {
	var err error

	jsonBytes, err := json.Marshal(order)
	if err != nil {
		log.Printf("MemCache: json.Marshal failed on CoffeeOrder %+v, error: %v", order, err)
		return false
	}

	orderStr := string(jsonBytes)

	cas, err := mcClient.Add(ordersKey, orderStr, 0, expTime)
	if err == mc.ErrKeyExists {
		log.Printf("MemCache: Key '%s' already exists in the cache. Appending...", ordersKey)
		_, err = mcClient.Append(ordersKey, "\n"+orderStr, cas)
	}

	if err == nil {
		log.Printf("MemCache: Order successfully added: %s", orderStr)
		return true

	} else {
		log.Printf("MemCache: Order wasn't stored: %v", err)
	}
	return false
}

func ordersReadyForCollection() ([]CoffeeOrder, uint64) {
	var err error

	allData, _, cas, err := mcClient.GAT(ordersKey, expTime2)
	if err == mc.ErrNotFound {
		return make([]CoffeeOrder, 0), 0
	} else if err != nil {
		log.Printf("Mem Cache: Failed to GAT('%s') : %v", ordersKey, err)
	}

	lines := strings.Split(string(allData), "\n")
	orders := make([]CoffeeOrder, 0, len(lines))

	for _, strVal := range lines {
		if strVal == "" {
			continue
		}
		var obj CoffeeOrder
		err = json.Unmarshal([]byte(strVal), &obj)
		if err != nil {
			log.Printf("CoffeeOrder: json.Unmarshal failed: %v", err)
		}
		//skip broken orders...
		if obj.UserID == 0 || obj.UserName == "" || obj.Beverage == "" || obj.Price == 0 {
			continue
		}
		orders = append(orders, obj)
	}
	return orders, cas
}

func collectOrdes(cas uint64) bool {
	var err error

	err = mcClient.DelCAS(ordersKey, cas)
	if err != nil {
		log.Printf("Mem Cache: Failed to DelCAS('%s') : %v", ordersKey, err)
		return false
	}
	return true
}

// CoffeeWatcher :
type CoffeeWatcher struct {
	UserID   int
	UserName string
	ChatID   int64
}

func (p CoffeeWatcher) isEmpty() bool {
	if p.UserID == 0 || p.UserName == "" || p.ChatID == 0 {
		return true
	}
	return false
}

func allCoffeeWatchers() []CoffeeWatcher {
	var err error

	allData, _, _, err := mcClient.Get(watchersKey)
	if err == mc.ErrNotFound {
		return make([]CoffeeWatcher, 0)
	} else if err != nil {
		log.Printf("Mem Cache: Failed to Get('%s') : %v", watchersKey, err)
		return make([]CoffeeWatcher, 0)
	}

	lines := strings.Split(string(allData), "\n")
	watchers := make([]CoffeeWatcher, 0, len(lines))

	for _, strVal := range lines {
		if strVal == "" {
			continue
		}
		var obj CoffeeWatcher
		err = json.Unmarshal([]byte(strVal), &obj)
		if err != nil {
			log.Printf("CoffeeWatcher: json.Unmarshal failed: %v", err)
		}

		//skip broken watchers...
		if obj.isEmpty() {
			continue
		}

		watchers = append(watchers, obj)
	}
	return watchers
}

func amIcoffeeWatcher(userID int) (bool, []CoffeeWatcher) {

	allWatchers := allCoffeeWatchers()
	found := false

	for _, obj := range allWatchers {

		if userID == obj.UserID {
			// I am Coffee Watcher!
			found = true
		}
	}
	return found, allWatchers
}

func addCoffeeWatcher(watcher CoffeeWatcher) error {
	var err error

	jsonBytes, err := json.Marshal(watcher)
	if err != nil {
		log.Printf("MemCache: json.Marshal failed on CoffeeWatcher: %+v, error: %v", watcher, err)
		return err
	}

	watcherStr := string(jsonBytes)

	cas, err := mcClient.Add(watchersKey, watcherStr, 0, 0)
	if err == mc.ErrKeyExists {
		log.Printf("MemCache: Key '%s' already exists in the cache. Appending...", watchersKey)
		_, err = mcClient.Append(watchersKey, "\n"+watcherStr, cas)
	}

	if err == nil {
		log.Printf("MemCache: Watcher successfully added: %s", watcherStr)
	} else {
		log.Printf("MemCache: Watcher addition failed: %v", err)
	}
	return err
}

func removeCoffeeWatcher(watcher CoffeeWatcher) error {
	var err error

	allData, _, cas, err := mcClient.Get(watchersKey)
	if err == mc.ErrNotFound {
		return nil
	} else if err != nil {
		log.Printf("Mem Cache: Failed to Get('%s') : %v", watchersKey, err)
		return err
	}

	lines := strings.Split(string(allData), "\n")
	var newData string
	found := false

	for _, strVal := range lines {
		if strVal == "" {
			continue
		}
		var obj CoffeeWatcher
		err = json.Unmarshal([]byte(strVal), &obj)
		if err != nil {
			log.Printf("CoffeeWatcher: json.Unmarshal failed: %v", err)
		}

		//skip broken watchers...
		if obj.isEmpty() {
			continue
		}
		//skip the watcher
		if watcher.UserID == obj.UserID {
			found = true
			continue
		}

		if len(newData) > 0 {
			newData += "\n" + strVal
		} else {
			newData = strVal
		}
	}

	if found {
		if len(newData) > 0 {
			_, err = mcClient.Replace(watchersKey, newData, 0, 0, cas)
		} else {
			err = mcClient.DelCAS(watchersKey, cas)
		}
	}
	return err
}
