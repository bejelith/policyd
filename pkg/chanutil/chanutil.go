package chanutil

func IsChannelOpen(c chan interface{}) bool {
	select {
	case _, isOpen := <-c:
		return isOpen
	default:
		// Channel open with no elements
		return true
	}
}
