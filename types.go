package MoneyBot

import (
	"encoding/json"
	"strings"
	"github.com/shvimas/teleBot"
	"os"
	"strconv"
	"bytes"
)

func LoadFrom(file *os.File, target interface{}) error {
	return json.NewDecoder(file).Decode(target)
}

func serialize(obj interface{}) ([]byte, error) {
	return json.MarshalIndent(obj, "", "\t")
}

const EncodingExtension = ".json"

// ------------------------------------------------------------------------------------

type user struct {
	Name      string
	Username  string
	Id        int
	Container *container
	History   *history
}

func NewUser(uid int, name, username string) *user {
	return &user{Name: name, Username: username, Id: uid, Container: NewContainer(uid), History: NewHistory(uid)}
}

func (user user) CheckFilled() bool {
	return user.Id != 0 && user.Username != "" && user.Name != "" && user.Container != nil && user.History != nil
}

func (user user) String() string {
	return teleBot.StructToString(user)
}

func (user user) Serialize() ([]byte, error) {
	return serialize(user)
}

// ------------------------------------------------------------------------------------

type container struct {
	Id      int
	Amounts map[string]float64
}

func NewContainer(uid int) *container {
	return &container{uid, make(map[string]float64)}
}

func (cont container) IsFilled() bool {
	return cont.Id != 0 && cont.Amounts != nil
}

func (cont container) Size() int {
	return len(cont.Amounts)
}

func (cont *container) Add(category string, amount float64) error {
	// zero values!
	cont.Amounts[strings.ToLower(category)] += amount
	return nil
}

func (cont *container) Delete(category string) error {
	delete(cont.Amounts, strings.ToLower(category))
	return nil
}

func (cont *container) Erase() error {
	cont.Amounts = make(map[string]float64)
	return nil
}

func (cont container) ToString() string {
	buf := bytes.Buffer{}
	// NOT MONOSPACE FONT
	//maxLen := 0
	//for category := range cont.Amounts {
	//	if len(category) > maxLen {
	//		maxLen = len(category)
	//	}
	//}
	for category, amount := range cont.Amounts {
		buf.WriteString(category)
		//buf.WriteString(strings.Repeat(" ", maxLen + 2 - len(category)))
		buf.WriteString(": ")
		buf.WriteString(strconv.FormatFloat(amount, 'f', 0, 64))
		buf.WriteString("\n")
	}
	if buf.Len() == 0 {
		buf.WriteString("<empty>")
	}
	return buf.String()
}

func (cont container) String() string {
	return teleBot.StructToString(cont)
}

func (cont container) Serialize() ([]byte, error) {
	return serialize(cont)
}

// ------------------------------------------------------------------------------------

type history struct {
	Id         int
	Containers map[string]container
}

func NewHistory(uid int) *history {
	return &history{uid, make(map[string]container)}
}

func (hist history) IsFilled() bool {
	return hist.Id != 0 && hist.Containers != nil
}

func (hist history) Size() int {
	return len(hist.Containers)
}

func (hist history) Get(name string) (container, bool) {
	cont, ok := hist.Containers[name]
	return cont, ok
}

func (hist *history) Add(container container, name string) string {
	if name == "" {
		counter := 1
		name = strconv.Itoa(hist.Size() + counter)
		for _, has := hist.Containers[name]; has; {
			counter++
			name = strconv.Itoa(hist.Size() + counter)
		}
	}
	hist.Containers[name] = container
	return name
}

func (hist *history) Delete(name string) bool {
	_, ok := hist.Containers[name]
	if ok {
		delete(hist.Containers, name)
	}
	return ok
}

func (hist history) ToString() string {
	buf := bytes.Buffer{}
	for name, container := range hist.Containers {
		buf.WriteString("###")
		buf.WriteString(name)
		buf.WriteString("###\n")
		buf.WriteString(container.ToString())
		buf.WriteString("\n\n")
	}
	if buf.Len() == 0 {
		buf.WriteString("<empty>")
	}
	return buf.String()
}

func (hist history) String() string {
	return teleBot.StructToString(hist)
}

func (hist history) Serialize() ([]byte, error) {
	return serialize(hist)
}
