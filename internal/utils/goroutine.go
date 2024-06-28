package utils

var goFuncSyncronous = false

func SetSyncGoroutine() {
	goFuncSyncronous = true
}

func ResetGoroutine() {
	goFuncSyncronous = false
}

func Go(fn func()) {
	if goFuncSyncronous {
		fn()
		return
	}
	go fn()
}
