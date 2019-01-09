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
	Bewerage  string
	Price     float64
	OrderTime time.Time
}

func placeOrder(order CoffeeOrder) error {
	var err error

	jsonBytes, err := json.Marshal(order)
	if err != nil {
		return err
	}
	orderStr := string(jsonBytes)

	cas, err := mcClient.Add(ordersKey, orderStr, 0, expTime)
	if err == mc.ErrKeyExists {
		log.Printf("MemCache: Key '%s' already exists in the cache. Appending...", ordersKey)
		_, err = mcClient.Append(ordersKey, "\n"+orderStr, cas)
	}

	if err == nil {
		log.Printf("MemCache: Order successfully added: %s", orderStr)
	} else {
		log.Printf("MemCache: Order wasn't stored: %v", err)
	}
	return err
}

func ordersReadyForCollection() ([]CoffeeOrder, uint64) {
	var err error

	allValues, _, cas, err := mcClient.GAT(ordersKey, expTime2)
	if err == mc.ErrNotFound {
		return make([]CoffeeOrder, 0), 0
	} else if err != nil {
		log.Printf("Mem Cache: Failed to GAT('%s') : %v", ordersKey, err)
	}

	lines := strings.Split(string(allValues), "\n")
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
		if obj.UserID == 0 || obj.UserName == "" || obj.Bewerage == "" || obj.Price == 0 {
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
