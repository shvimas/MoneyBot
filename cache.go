package MoneyBot

import (
	"github.com/shvimas/teleBot"
)

const CacheSize = 100

type cached struct {
	User   *user
	Weight int
}

type Cache struct {
	users map[int]cached
	prov  Provider
	chans LogChans
}

//func NewCache(prov Provider) *Cache { return &Cache{users: make(map[int]cached), prov: prov} }

func (cache Cache) Size() int {
	return len(cache.users)
}

func (cache Cache) Connect() error { return nil }
func (cache Cache) Close() error   { return nil }

func (cache Cache) GetUser(userId int) (*user, error) {
	cached, ok := cache.users[userId]
	if ok {
		cached.Weight++
		return cached.User, nil
	} else {
		return cache.prov.GetUser(userId)
	}
}

func (cache Cache) GetContainer(userId int) (*container, error) {
	cached, ok := cache.users[userId]
	if ok {
		cached.Weight++
		return cached.User.Container, nil
	} else {
		return cache.prov.GetContainer(userId)
	}
}

func (cache Cache) GetHistory(userId int) (*history, error) {
	cached, ok := cache.users[userId]
	if ok {
		cached.Weight++
		return cached.User.History, nil
	} else {
		return cache.prov.GetHistory(userId)
	}
}

func (cache Cache) PushContainer(user *user) error {
	cache.chans.LogChan <- "Pushing container " + user.Container.String() + " to cache"
	if cache.Size() > 0.7*CacheSize {
		err := cache.Clean()
		if err != nil {
			return err
		}
	}
	_, hasUser := cache.users[user.Id]
	if hasUser {
		err := cache.prov.PushContainer(user.Id, user.Container)
		if err != nil {
			return err
		}
		cache.users[user.Id].User.Container = user.Container
	}
	if !hasUser && cache.Size() < CacheSize {
		// FIXME: check and think
		cache.users[user.Id] = cached{User: user, Weight: 0}
		return nil
	}
	cache.chans.LogChan <- "Failed: cache full"
	cache.chans.LogChan <- "Pushing container " + user.Container.String() + " to provider"
	return cache.prov.PushContainer(user.Id, user.Container)
}

// ALL WRONG HERE
func (cache Cache) PushHistory(user *user) error {
	cache.chans.LogChan <- "Pushing history " + user.History.String() + " to cache"
	if cache.Size() > 0.7*CacheSize {
		err := cache.Clean()
		if err != nil {
			return err
		}
	}
	_, hasUser := cache.users[user.Id]
	if hasUser {
		err := cache.prov.PushHistory(user.Id, user.History)
		if err != nil {
			return err
		}
		cache.users[user.Id].User.History = user.History
	}
	if !hasUser && cache.Size() < CacheSize {
		// FIXME: check and think
		cache.users[user.Id] = cached{User: user, Weight: 0}
		return nil
	}
	cache.chans.LogChan <- "Failed: cache full"
	cache.chans.LogChan <- "Pushing history " + user.History.String() + " to provider"
	return cache.prov.PushHistory(user.Id, user.History)
}

func (cache Cache) RegisterUser(tUser teleBot.User) error {
	_, ok := cache.users[tUser.Id]
	if ok {
		return nil
	} else {
		// user := NewUser(tUser.Id, tUser.FullName(), tUser.Username)
		panic("not implemented")
		// return cache.PushUser(user)
	}
}

// FIXME: think about right strategy
func (cache Cache) Clean() error {
	for uid, cached := range cache.users {
		if cached.Weight == 0 {
			err := cache.prov.PushContainer(uid, cache.users[uid].User.Container)
			if err != nil {
				return err
			}
			delete(cache.users, uid)
		}
	}
	return nil
}
