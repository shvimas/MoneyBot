package MoneyBot

import (
	"testing"
	"fmt"
)

func TestContainer_Add(t *testing.T) {

}

func TestHistory_ToString(t *testing.T) {
	history := NewHistory(111)
	container := NewContainer(111)
	container.Add("test", 110)
	container.Add("test_2", 2)
	history.Add(*container, "first")
	fmt.Println(history.ToString())

}
