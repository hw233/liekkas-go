package param

import (
	"shared/utility/whitelist"
)

type WhiteListItem struct {
	Type  whitelist.WhiteListType `json:"type"`
	Value string                  `json:"value"`
}

type WhiteListItemSwitch struct {
	Switch bool `json:"switch"`
}
