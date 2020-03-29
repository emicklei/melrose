package melrose

var globalCurrentDevice AudioDevice

func CurrentDevice() AudioDevice     { return globalCurrentDevice }
func SetCurrentDevice(a AudioDevice) { globalCurrentDevice = a }
