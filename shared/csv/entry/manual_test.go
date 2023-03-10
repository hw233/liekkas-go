package entry

import (
	"shared/csv/static"
	"testing"
)

func TestManual_GetManualDrops(t *testing.T) {
	manualDrops, _ := CSV.Manual.getManualDrop(101002)
	t.Log(manualDrops)
}

func TestManual_FindManualIdByRelatedId(t *testing.T) {
	manual, ok := CSV.Manual.FindManualIdByRelatedId(static.ManualTypeWorldItem, 10001)
	t.Log(manual)
	t.Log(ok)

}
