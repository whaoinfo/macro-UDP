package ioadapter

import (
	"github.com/whaoinfo/macro-UDP/internal/define"
	frame "github.com/whaoinfo/macro-UDP/pkg/gicframe"
)

func init() {
	frame.RegisterComponentInfo(define.IOAdapterComponentType, func() frame.IComponent {
		return &Component{}
	}, func() frame.IComponentKW {
		return &ComponentKW{}
	})
}
