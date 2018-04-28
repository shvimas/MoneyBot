package MoneyBot

import (
	"testing"
	"github.com/shvimas/teleBot"
	"os"
)

var prov = FSProvider{"test_data/fs_provider/"}
var id = 101

func TestFSProvider_RegisterUser(t *testing.T) {
	user := teleBot.User{Id: id}
	err := prov.RegisterUser(user)
	if err != nil {
		t.Error(err)
	}
}

func TestFSProvider_PushContainer(t *testing.T) {
	id := 101
	user := NewUser(id, "Test", "test")
	user.Container.Add("test", 100500)
	err := prov.PushContainer(user.Id, user.Container)
	if err != nil {
		t.Error(err)
	}
}

func TestFSProvider_Close(t *testing.T) {
	user := teleBot.User{Id: id}
	err := os.RemoveAll(prov.GetUserDir(user.Id))
	if err != nil {
		t.Error(err)
	}
}
