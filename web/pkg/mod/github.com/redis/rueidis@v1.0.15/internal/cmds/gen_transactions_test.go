// Code generated DO NOT EDIT

package cmds

import "testing"

func transactions0(s Builder) {
	s.Discard().Build()
	s.Exec().Build()
	s.Multi().Build()
	s.Unwatch().Build()
	s.Watch().Key("1").Key("1").Build()
}

func TestCommand_InitSlot_transactions(t *testing.T) {
	var s = NewBuilder(InitSlot)
	t.Run("0", func(t *testing.T) { transactions0(s) })
}

func TestCommand_NoSlot_transactions(t *testing.T) {
	var s = NewBuilder(NoSlot)
	t.Run("0", func(t *testing.T) { transactions0(s) })
}
