package session

import (
	"log"
	"testing"
	"time"
)

type sessBuilder struct{}

func (b *sessBuilder) NewSession() Session {
	return &sessTest{EmbedManagedSession: &EmbedManagedSession{}}
}

type sessTest struct {
	*EmbedManagedSession

	ID int64
}

func (s *sessTest) OnCreated(opts OnCreatedOpts) error {
	log.Println("on created")
	s.ID = opts.ID
	return nil
}

func (s *sessTest) OnTriggered() {
	log.Println("on triggered")
}

func (s *sessTest) OnClosed() {
	log.Println("on closed")
}

func TestSession(t *testing.T) {
	manager := NewManager(&sessBuilder{}, &ManagerConfig{Expire: time.Second * 3, Capacity: 1000})
	sess, err := manager.GetSession(1)
	if err != nil {
		t.Errorf("get session error: %v", err)
		return
	}

	st := sess.(*sessTest)
	t.Logf("session: %v", st.ID)
	_, _ = manager.GetSession(1)

	time.Sleep(5 * time.Second)

	manager.Close()
}
