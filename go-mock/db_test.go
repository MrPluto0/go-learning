package mock

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestGetFromDB(t *testing.T) {
	// 1. initialize controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 2. create stubs, including params and returns
	m := NewMockDB(ctrl)
	m.EXPECT().Get(gomock.Eq("Tom")).Return(100, errors.New("not exist"))

	// 3. verify
	if v := GetFromDB(m, "Tom"); v != -1 {
		t.Fatal("expected -1, but got", v)
	}
}
