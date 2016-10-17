package redis

import (
	"log"

	"github.com/denbeigh2000/jfsi"
	"github.com/denbeigh2000/jfsi/metastore"

	"gopkg.in/redis.v5"
)

func NewStore(addr, password string, db int) metastore.MetaStore {
	store := &store{
		Addr:     addr,
		DB:       db,
		Password: password,
	}

	store.init()

	return store
}

type store struct {
	client *redis.Client

	Addr     string
	DB       int
	Password string
}

func (s *store) init() {
	s.client = redis.NewClient(&redis.Options{
		Addr:     s.Addr,
		DB:       s.DB,
		Password: s.Password,
	})
}

func (s *store) Create(key jfsi.ID, n int) (r metastore.Record, err error) {
	if n <= 0 {
		err = metastore.ZeroLengthCapacityRecordErr{}
		return
	}

	exists, err := s.client.Exists(string(key)).Result()
	if err != nil {
		return
	}
	if exists {
		err = metastore.KeyAlreadyExistsErr(key)
		return
	}

	items := make([]jfsi.ID, n)
	redisItems := make([]interface{}, n)
	for i := 0; i < n; i++ {
		id := jfsi.NewID()
		items[i] = id
		redisItems[i] = string(id)
	}

	err = s.client.RPush(string(key), redisItems...).Err()
	if err != nil {
		return
	}

	r = metastore.Record{
		Key:    key,
		Chunks: items,
	}
	return r, nil
}

func (s *store) Retrieve(key jfsi.ID) (r metastore.Record, err error) {
	exists, err := s.client.Exists(string(key)).Result()
	if err != nil {
		return
	}
	if !exists {
		err = metastore.KeyNotFoundErr(key)
		return
	}

	items, err := s.client.LRange(string(key), 0, -1).Result()
	if err != nil {
		return
	}

	ids := make([]jfsi.ID, len(items))
	for i, item := range items {
		ids[i] = jfsi.ID(item)
	}

	r = metastore.Record{
		Key:    key,
		Chunks: ids,
	}
	return
}

func (s *store) Update(key jfsi.ID, r metastore.Record) error {
	exists, err := s.client.Exists(string(key)).Result()
	if err != nil {
		return err
	}
	if !exists {
		err = metastore.KeyNotFoundErr(key)
		return err
	}

	result, err := s.client.Del(string(key)).Result()
	if err != nil {
		return err
	}

	if result != 1 {
		log.Printf("Expected 1 key to be removed, but for some reason %v keys were instead", result)
	}

	redisItems := make([]interface{}, len(r.Chunks))
	for i, chunkID := range r.Chunks {
		redisItems[i] = chunkID
	}

	err = s.client.RPush(string(key), redisItems...).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *store) Delete(key jfsi.ID) error {
	exists, err := s.client.Exists(string(key)).Result()
	if err != nil {
		return err
	}
	if !exists {
		err = metastore.KeyNotFoundErr(key)
		return err
	}

	err = s.client.Del(string(key)).Err()
	if err != nil {
		return err
	}

	return nil
}
