package locks

// armMutexKV is the instance of MutexKV for ARM resources
var armMutexKV = NewMutexKV()

func ByID(id string) {
	armMutexKV.Lock(id)
}

func UnlockByID(id string) {
	armMutexKV.Unlock(id)
}
