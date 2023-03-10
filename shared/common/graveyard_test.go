package common

import (
	"testing"
	"time"
)

func TestGraveyardHelpRequests_GetRequestsExcept(t *testing.T) {

	requests := NewGraveyardRequests()
	requests.Add(NewGraveyardHelpRequest(1, 1, 1, 20, 1, 2, 2, time.Now().Unix()+1000))
	except := requests.GetRequestsExcept(2, nil)
	t.Log(except)
}
