package observer

type event string

const (
	opGet         event = "op_get"
	opSet         event = "op_set"
	opWipe        event = "op_wipe"
	opDel         event = "op_del"
	verifiedEvent event = "verified_event"
)
