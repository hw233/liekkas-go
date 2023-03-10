package mysql

import (
	"log"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var handler *Handler

func TestMain(m *testing.M) {
	config := &Config{
		Addr:            "root:xuecm01,@tcp(localhost:3306)/overlord?charset=utf8mb4&parseTime=true&loc=Local",
		MaxIdleConn:     5,
		MaxOpenConn:     5,
		ConnMaxLifetime: 4 * time.Hour,
	}

	var err error
	handler, err = NewHandler(config)
	if err != nil {
		log.Printf("new manager error: %v", err)
	}

	m.Run()
}

type UserTest struct {
	ID           int64  `db:"id" major:"true"`
	Name         string `db:"name"`
	XXX          string `db:"-"`
	Info         string `db:"info"`
	*EmbedModule `db:"-"`
}

func TestLoad(t *testing.T) {
	user := &UserTest{}
	user.ID = 1
	user.EmbedModule = &EmbedModule{}

	err := handler.Load(user)
	if err != nil {
		t.Errorf("Load error: %v", err)
	}

	t.Logf("user: %+v", *user)
}

func TestSave(t *testing.T) {
	user := &UserTest{}
	user.ID = 1
	user.EmbedModule = &EmbedModule{}

	user.Name = "雪辙2222"

	err := handler.Save(user)
	if err != nil {
		t.Errorf("Load error: %v", err)
	}

	user.Name = "雪辙555555"
	user.Info = "wawawawawawa"

	err = handler.Save(user)
	if err != nil {
		t.Errorf("Save error: %v", err)
	}

	t.Logf("user: %+v", *user)
}

func TestCreate(t *testing.T) {
	user := &UserTest{}
	user.ID = 1
	user.EmbedModule = &EmbedModule{}

	user.Name = "雪辙555555"
	user.Info = "444444"

	err := handler.Create(user)
	if err != nil {
		t.Errorf("Load error: %v", err)
	}

	t.Logf("user: %+v", *user)
}

type UserTest2 struct {
	ID           int64  `db:"id" major:"true"`
	Name         string `db:"name"`
	XXX          string `db:"-"`
	Info         string `db:"info"`
	*EmbedModule `db:"-"`
}

func TestSetTable(t *testing.T) {
	user := &UserTest{}
	user.ID = 1
	user.EmbedModule = &EmbedModule{}
	user.SetTable("user_test")

	err := handler.Load(user)
	if err != nil {
		t.Errorf("Load error: %v", err)
	}

	t.Logf("user: %+v", *user)
}
