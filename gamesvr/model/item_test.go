package model

import (
	"testing"
)

func TestItemPack_Add(t *testing.T) {
	itemPack := NewItemPack()
	itemPack.Add(1, 2)
	itemPack.Add(3, 4)

	js, err := itemPack.Marshal()
	if err != nil {
		t.Errorf("itemPack.Marshal error: %v", err)
	}

	t.Logf("itemPack json: %s", string(js))

	itemPackNew := NewItemPack()
	err = itemPackNew.Unmarshal(js)
	if err != nil {
		t.Errorf("itemPackNew.Unmarshal error: %v", err)
	}

	// for k, v := range itemPackNew.Items {
	// 	t.Logf("itemPackNew: k: %v, v: id: %d num: %d ex: %d", k, v.ID, v.Num, v.Expire)
	// }
}
