package MoneyBot

import (
	"testing"
	"fmt"
	"strings"
)

func TestBotParseUpdates(t *testing.T) {
	str := "add  food 100500     	rubles"
	fmt.Printf("%q", strings.Fields(str))
}
