package lifecycle

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type managed struct {
	startCount int
	stopCount  int
}

func (m *managed) Start() {
	m.startCount = m.startCount + 1

}
func (m *managed) Stop() {
	m.stopCount += 1
}

func TestStart(t *testing.T) {
	l := New()
	m := &managed{}
	l.Manage(m)
	l.Start()
	assert.Equal(t, 1, m.startCount)
	l.Start()
	assert.Equal(t, 1, m.startCount)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		l.Wait()
		wg.Done()
	}()
	l.Stop()
	assert.Equal(t, 1, m.stopCount)
	wg.Wait()
}

func TestManageAfterStar(t *testing.T) {
	l := New()
	m := &managed{}
	l.Start()
	l.Manage(m)
	assert.Zero(t, len(l.objects), "No supervised objects should be present")
	assert.Equal(t, 0, m.startCount)
}
