package grepo

import (
	"testing"

	_ "github.com/Kingcxp/gonebot-plugin-test"

	"github.com/gonebot-dev/gonebot"
)

func TestMain(m *testing.M) {
	SetEntry("./grepo_test.go")

	Require("tester", "v0.0.1")

	gonebot.StartBackend("onebot11")
}
