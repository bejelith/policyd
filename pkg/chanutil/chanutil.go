package chanutil

func IschannelClosed(c chan interface{}) bool {
	select {
	case _, ok := <-c:
		return ok
	default:
		return false
	}
}
