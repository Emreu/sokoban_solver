package world

// TODO: replace map with augmented map in solver
type AugmentedMap struct {
	Map
	deadZones Bitmap
	hLock     Bitmap
	vLock     Bitmap
}

func Augment(m Map) AugmentedMap {
	deadZones := FindDeadZones(m)
	hLock, vLock := FindMoveLocks(m, deadZones)
	return AugmentedMap{
		Map:       m,
		deadZones: deadZones,
		hLock:     hLock,
		vLock:     vLock,
	}
}
