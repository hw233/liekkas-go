package entry

import "testing"

func TestGuideEntry_GetAllSortedConfig(t *testing.T) {
	t.Log(len(CSV.Guide.GetGuideBefore(10)))
}
