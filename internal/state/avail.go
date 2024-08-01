package state

var avail bool

func init() {
	SetUnavailable()
}

func GetAvail() bool {
	return avail
}

func SetUnavailable() {
	avail = false
}

func SetAvailable() {
	avail = true
}
