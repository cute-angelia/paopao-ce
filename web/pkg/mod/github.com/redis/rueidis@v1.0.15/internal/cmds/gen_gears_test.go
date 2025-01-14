// Code generated DO NOT EDIT

package cmds

import "testing"

func gears0(s Builder) {
	s.RgAbortexecution().Id("1").Build()
	s.RgConfigget().Key("1").Key("1").Build()
	s.RgConfigset().KeyValue().KeyValue("1", "1").KeyValue("1", "1").Build()
	s.RgDropexecution().Id("1").Build()
	s.RgDumpexecutions().Build()
	s.RgDumpregistrations().Build()
	s.RgGetexecution().Id("1").Shard().Build()
	s.RgGetexecution().Id("1").Cluster().Build()
	s.RgGetexecution().Id("1").Build()
	s.RgGetresults().Id("1").Build()
	s.RgGetresultsblocking().Id("1").Build()
	s.RgInfocluster().Build()
	s.RgPydumpreqs().Build()
	s.RgPyexecute().Function("1").Unblocking().Id("1").Description("1").Upgrade().ReplaceWith("1").Requirements("1").Requirements("1").Build()
	s.RgPyexecute().Function("1").Unblocking().Id("1").Description("1").Upgrade().ReplaceWith("1").Build()
	s.RgPyexecute().Function("1").Unblocking().Id("1").Description("1").Upgrade().Requirements("1").Requirements("1").Build()
	s.RgPyexecute().Function("1").Unblocking().Id("1").Description("1").Upgrade().Build()
	s.RgPyexecute().Function("1").Unblocking().Id("1").Description("1").ReplaceWith("1").Build()
	s.RgPyexecute().Function("1").Unblocking().Id("1").Description("1").Requirements("1").Requirements("1").Build()
	s.RgPyexecute().Function("1").Unblocking().Id("1").Description("1").Build()
	s.RgPyexecute().Function("1").Unblocking().Id("1").Upgrade().Build()
	s.RgPyexecute().Function("1").Unblocking().Id("1").ReplaceWith("1").Build()
	s.RgPyexecute().Function("1").Unblocking().Id("1").Requirements("1").Requirements("1").Build()
	s.RgPyexecute().Function("1").Unblocking().Id("1").Build()
	s.RgPyexecute().Function("1").Unblocking().Description("1").Build()
	s.RgPyexecute().Function("1").Unblocking().Upgrade().Build()
	s.RgPyexecute().Function("1").Unblocking().ReplaceWith("1").Build()
	s.RgPyexecute().Function("1").Unblocking().Requirements("1").Requirements("1").Build()
	s.RgPyexecute().Function("1").Unblocking().Build()
	s.RgPyexecute().Function("1").Id("1").Build()
	s.RgPyexecute().Function("1").Description("1").Build()
	s.RgPyexecute().Function("1").Upgrade().Build()
	s.RgPyexecute().Function("1").ReplaceWith("1").Build()
	s.RgPyexecute().Function("1").Requirements("1").Requirements("1").Build()
	s.RgPyexecute().Function("1").Build()
	s.RgPystats().Build()
	s.RgRefreshcluster().Build()
	s.RgTrigger().Trigger("1").Argument("1").Argument("1").Build()
	s.RgUnregister().Id("1").Build()
}

func TestCommand_InitSlot_gears(t *testing.T) {
	var s = NewBuilder(InitSlot)
	t.Run("0", func(t *testing.T) { gears0(s) })
}

func TestCommand_NoSlot_gears(t *testing.T) {
	var s = NewBuilder(NoSlot)
	t.Run("0", func(t *testing.T) { gears0(s) })
}
