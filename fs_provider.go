package MoneyBot

import (
	"teleBot"
	"os"
	"errors"
	"strconv"
)

type FSProvider struct {
	root string
}

func (prov FSProvider) GetUserDir(uid int) string {
	return prov.root + strconv.Itoa(uid)
}

func (prov FSProvider) GetFile(uid int, file string) (*os.File, error) {
	name := prov.GetUserDir(uid) + "/" + file
	return os.OpenFile(name, os.O_RDWR, 0766)
}

func (prov FSProvider) HasUser(uid int) bool {
	file, err := os.Open(prov.GetUserDir(uid))
	if err != nil {
		return false
	}
	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	return stat.IsDir()
}

const (
	userFile      = "user" + EncodingExtension
	containerFile = "container" + EncodingExtension
	historyFile   = "history" + EncodingExtension
)

func (prov FSProvider) Connect() error {
	file, err := os.Open(prov.root)
	if err != nil {
		err = os.MkdirAll(prov.root, 0755)
		if err != nil {
			return err
		}
	}
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New(prov.root + " is not a directory")
	}
	return nil
}

func (prov FSProvider) Close() error {
	return nil
}

func (prov FSProvider) GetUser(uid int) (*user, error) {
	file, err := prov.GetFile(uid, userFile)
	if err != nil {
		return nil, err
	}
	target := &user{}
	if err = LoadFrom(file, target); err != nil {
		return nil, err
	}
	if !target.CheckFilled() {
		return target, errors.New("failed to fill target")
	}
	return target, nil
}

func (prov FSProvider) GetContainer(uid int) (*container, error) {
	file, err := prov.GetFile(uid, containerFile)
	if err != nil {
		return nil, err
	}
	target := &container{}
	if err = LoadFrom(file, target); err != nil {
		return nil, err
	}
	if !target.IsFilled() {
		return target, errors.New("failed to fill target")
	}
	return target, nil
}

func (prov FSProvider) GetHistory(uid int) (*history, error) {
	file, err := prov.GetFile(uid, historyFile)
	if err != nil {
		return nil, err
	}
	target := &history{}
	if err = LoadFrom(file, target); err != nil {
		return nil, err
	}
	if !target.IsFilled() {
		return target, errors.New("failed to fill target")
	}
	return target, nil
}

func (prov FSProvider) PushContainer(uid int, cont *container) error {
	file, err := prov.GetFile(uid, containerFile)
	if err != nil {
		return err
	}
	buf, err := cont.Serialize()
	if err != nil {
		return err
	}
	_, err = file.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (prov FSProvider) PushHistory(uid int, hist *history) error {
	file, err := prov.GetFile(uid, historyFile)
	if err != nil {
		return err
	}
	buf, err := hist.Serialize()
	if err != nil {
		return err
	}
	_, err = file.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (prov FSProvider) CreateFiles(uid int) error {
	err := os.MkdirAll(prov.GetUserDir(uid), 0755)
	if err != nil {
		return err
	}
	for _, file := range []string{userFile, containerFile, historyFile} {
		if _, err := os.Create(prov.GetUserDir(uid) + "/" + file); err != nil {
			return err
		}
	}
	return nil
}

func (prov FSProvider) RegisterUser(user teleBot.User) error {
	uid := user.Id
	if prov.HasUser(uid) {
		return nil
	}
	// create empty files
	if err := prov.CreateFiles(uid); err != nil {
		return err
	}
	// write zero user
	file, err := prov.GetFile(uid, userFile)
	if err != nil {
		return prov.HandleError(uid, err)
	}
	emptyUser, err := NewUser(uid, user.FullName(), user.Username).Serialize()
	if err != nil {
		return prov.HandleError(uid, err)
	}
	if _, err := file.Write(emptyUser); err != nil {
		return prov.HandleError(uid, err)
	}
	// write zero container
	file, err = prov.GetFile(uid, containerFile)
	if err != nil {
		return prov.HandleError(uid, err)
	}
	emptyCont, err := NewContainer(uid).Serialize()
	if err != nil {
		return prov.HandleError(uid, err)
	}
	if _, err := file.Write(emptyCont); err != nil {
		return prov.HandleError(uid, err)
	}
	// write zero history
	file, err = prov.GetFile(uid, historyFile)
	if err != nil {
		return prov.HandleError(uid, err)
	}
	emptyHist, err := NewHistory(uid).Serialize()
	if err != nil {
		return prov.HandleError(uid, err)
	}
	if _, err := file.Write(emptyHist); err != nil {
		return prov.HandleError(uid, err)
	}

	return nil
}

func (prov FSProvider) HandleError(uid int, err error) error {
	os.Remove(prov.GetUserDir(uid))
	return err
}
