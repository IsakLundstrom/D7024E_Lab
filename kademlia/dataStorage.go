package kademlia

import (
	"sync"
	"time"
)

type DataStorage struct {
	store   map[string]StoreItem
	storeMu sync.Mutex
}

type StoreItem struct {
	Data     string
	ExpireAt time.Time
}

func NewDataStorage() DataStorage {
	return DataStorage{
		store: map[string]StoreItem{},
	}
}

// store data to the store, if key already exists, it will refresh the TTL for the data
func (storage *DataStorage) SetData(key string, data string) (item StoreItem, exist bool) {
	storage.storeMu.Lock()
	defer storage.storeMu.Unlock()

	if item, exist := storage.store[key]; exist {
		// Update the TTL for an existing item
		item.ExpireAt = time.Now().Add(DATA_TIME_TO_LIVE)
		storage.store[key] = item
		return item, true
	}

	expire := time.Now().Add(DATA_TIME_TO_LIVE)
	storeItem := StoreItem{Data: data, ExpireAt: expire}
	storage.store[key] = storeItem
	go storage.cleanupEvent(key)
	return storeItem, false
}

// get data stored, returns empty string if data expired or do not exist
func (storage *DataStorage) GetData(key string) (data string, exist bool) {
	storage.storeMu.Lock()
	defer storage.storeMu.Unlock()
	item, exists := storage.store[key]
	if !exists {
		return "", false
	}
	item.ExpireAt = time.Now().Add(DATA_TIME_TO_LIVE)
	storage.store[key] = item // refresh TTL
	return item.Data, true
}

func (storage *DataStorage) cleanupEvent(key string) {
	time.Sleep(DATA_TIME_TO_LIVE + 1*time.Millisecond) // 1ms to avoid race condition
	for {
		storage.storeMu.Lock()
		item := storage.store[key]
		if time.Now().After(item.ExpireAt) {
			delete(storage.store, key)
			storage.storeMu.Unlock()
			return
		}
		storage.storeMu.Unlock()
		time.Sleep(time.Until(item.ExpireAt) + 1*time.Millisecond) // 1ms to avoid race condition
	}
}
