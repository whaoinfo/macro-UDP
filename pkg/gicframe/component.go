package gicframe

import "fmt"

type ComponentType string
type ComponentID string
type ComponentStatus uint16

const (
	ComponentBaseInitStatus = iota
	ComponentInitStatus
	ComponentStartStatus
	ComponentStopStatus
)

type NewComponent func() IComponent

type IComponentKW interface{}

type NewComponentKW func() IComponentKW

type RegComponentInfo struct {
	Tpy            ComponentType
	NewComponent   NewComponent
	NewComponentKW NewComponentKW
}

type IComponent interface {
	baseInitialize(num int, tpy ComponentType) error
	Initialize(kw IComponentKW) error
	GetNum() int
	GetID() ComponentID
	GetType() ComponentType
	Start() error
	Stop() error
	//AcceptAssociation(args ...interface{})
	//Associate()
}

var (
	regComponentInfoMap = make(map[ComponentType]*RegComponentInfo)
)

func RegisterComponentInfo(tpy ComponentType, newComponent NewComponent, newComponentKW NewComponentKW) {
	regComponentInfoMap[tpy] = &RegComponentInfo{
		Tpy:            tpy,
		NewComponent:   newComponent,
		NewComponentKW: newComponentKW,
	}
}

type BaseComponent struct {
	num      int
	tpy      ComponentType
	id       ComponentID
	status   ComponentStatus
	appProxy IApplication
}

func (t *BaseComponent) baseInitialize(num int, tpy ComponentType) error {
	t.num = num
	t.tpy = tpy
	t.id = ComponentID(fmt.Sprintf("%v_%d", tpy, num))

	return nil
}

func (t *BaseComponent) Initialize(args ...interface{}) error {
	return nil
}

func (t *BaseComponent) GetType() ComponentType {
	return t.tpy
}

func (t *BaseComponent) GetNum() int {
	return t.num
}

func (t *BaseComponent) GetID() ComponentID {
	return t.id
}

func (t *BaseComponent) GetStatus() ComponentStatus {
	return t.status
}

func (t *BaseComponent) setPartStatus(status ComponentStatus) {
	t.status = status
}

func (t *BaseComponent) Start() error {
	return nil
}

func (t *BaseComponent) Stop() error {
	return nil
}
