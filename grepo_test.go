package grepo

import (
	_ "github.com/Kingcxp/gonebot-plugin-test"
	"testing"

	"github.com/gonebot-dev/gonebot"
)

func TestMain(m *testing.M) {
	SetEntry("./grepo_test.go")

	Require("test", "v0.0.1")

	gonebot.StartBackend("onebot11")
}
