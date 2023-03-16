package tracer

import (
	"github.com/whaoinfo/macro-UDP/internal/define"
	frame "github.com/whaoinfo/macro-UDP/pkg/gicframe"
)

func init() {
	frame.RegisterComponentInfo(define.TracerComponentType, func() frame.IComponent {
		return &Component{}
	}, func() frame.IComponentKW {
		return &ComponentKW{}
	})
}
