package chanutil

func IsChannelClosed(c chan interface{}) bool {
	select {
	case _, open := <-c:
		return !open
	default:
		// Channel open with no elements
		return false
	}
}
