package MoneyBot

import "teleBot"

type Provider interface {
	Connect() error
	Close() error
	GetUser(int) (*user, error)
	GetContainer(int) (*container, error)
	GetHistory(int) (*history, error)
	PushContainer(int, *container) error
	PushHistory(int, *history) error
	RegisterUser(teleBot.User) error // rewrites user and all his data if already had
}

type DummyProvider struct {
	Error error
}

func (dp DummyProvider) Connect() error                       { return dp.Error }
func (dp DummyProvider) Close() error                         { return dp.Error }
func (dp DummyProvider) GetUser(int) (*user, error)           { return nil, dp.Error }
func (dp DummyProvider) GetContainer(int) (*container, error) { return nil, dp.Error }
func (dp DummyProvider) GetHistory(int) (*history, error)     { return nil, dp.Error }
func (dp DummyProvider) PushContainer(*user) error            { return dp.Error }
func (dp DummyProvider) PushHistory(*user) error              { return dp.Error }
func (dp DummyProvider) RegisterUser(teleBot.User) error      { return dp.Error }
