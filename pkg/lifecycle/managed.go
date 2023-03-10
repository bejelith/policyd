package lifecycle

type Managed interface {
	Start()
	Stop()
}
