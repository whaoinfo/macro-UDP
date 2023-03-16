package gicframe

import (
	"encoding/json"
	"fmt"
	"github.com/whaoinfo/go-box/ctime"
	"github.com/whaoinfo/go-box/logger"
	"github.com/whaoinfo/go-box/ossignal"
	"github.com/whaoinfo/go-box/pubsub"
	"os"
	"syscall"
)

const (
	APPDefaultSignChanSize = 1
)

type ApplicationID string
type NewApplication func() IApplication
type TopicType uint16
type TopicSubKey interface{}

type IApplication interface {
	baseInitialize(appID ApplicationID) error
	Initialize(args ...interface{}) error
	GetID() ApplicationID
	importComponents(infoList []componentConfigModel) error
	Pub(topic TopicType, args ...interface{}) error
	Sub(topic TopicType, subKey TopicSubKey, handle pubsub.TopicFunc, preArgs ...interface{}) error
	start() error
	AfterStart() error
	stop() error
	StopBefore() error
	forever()
}

type BaseApplication struct {
	id            ApplicationID
	signalHandler *ossignal.SignalHandler
	ob            *pubsub.ObServer
	componentMap  map[ComponentID]IComponent
}

func (t *BaseApplication) baseInitialize(id ApplicationID) error {
	t.id = id
	t.componentMap = make(map[ComponentID]IComponent)

	t.signalHandler = &ossignal.SignalHandler{}
	if err := t.signalHandler.InitSignalHandler(APPDefaultSignChanSize); err != nil {
		return err
	}

	ob, newObErr := pubsub.NewObServer(true)
	if newObErr != nil {
		return newObErr
	}
	t.ob = ob

	for _, sig := range []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT} {
		t.signalHandler.RegisterSignal(sig, func() {
			t.signalHandler.CloseSignalHandler()
		})
	}

	return nil
}

func (t *BaseApplication) Initialize(args ...interface{}) error {
	return nil
}

func (t *BaseApplication) GetID() ApplicationID {
	return t.id
}

func (t *BaseApplication) Pub(topic TopicType, args ...interface{}) error {
	return t.ob.Publish(topic, false, args...)
}

func (t *BaseApplication) Sub(topic TopicType, subKey TopicSubKey,
	handle pubsub.TopicFunc, preArgs ...interface{}) error {

	return t.ob.Subscribe(topic, subKey, handle, preArgs...)
}

func (t *BaseApplication) start() error {
	for componentID, component := range t.componentMap {
		logger.InfoFmt("Start %v component...", componentID)
		if err := component.Start(); err != nil {
			return fmt.Errorf("the %v component has failed to start, %v", componentID, err)
		}
		component.setStartTimestamp(ctime.CurrentTimestamp())
		logger.InfoFmt("The %v component has started", componentID)
	}
	return nil
}

func (t *BaseApplication) AfterStart() error {
	return nil
}

func (t *BaseApplication) stop() error {
	for componentID, component := range t.componentMap {
		logger.InfoFmt("Stop %v component...", componentID)
		if err := component.Stop(); err != nil {
			return fmt.Errorf("failed to stop component, component ID: %v, err: %v", componentID, err)
		}
		logger.InfoFmt("The component %v has Stopped", componentID)
	}
	return nil
}

func (t *BaseApplication) StopBefore() error {

	return nil
}

func (t *BaseApplication) forever() {
	t.signalHandler.ListenSignal()
}

func (t *BaseApplication) importComponents(cfgList []componentConfigModel) error {
	componentNum := 1
	for _, cfg := range cfgList {
		if cfg.Disable {
			continue
		}

		compTpy := ComponentType(cfg.ComponentType)
		regInfo, exist := regComponentInfoMap[compTpy]
		logger.InfoFmt("Import %v component...", compTpy)
		if !exist {
			return fmt.Errorf("component type %v dose not exist", compTpy)
		}

		component := regInfo.NewComponent()
		if err := component.baseInitialize(componentNum, compTpy); err != nil {
			return err
		}

		kwData, msErr := json.Marshal(cfg.Kw)
		if msErr != nil {
			return msErr
		}
		kw := regInfo.NewComponentKW()
		if kw != nil {
			if err := json.Unmarshal(kwData, kw); err != nil {
				return err
			}
		}

		if err := component.Initialize(kw); err != nil {
			return err
		}

		t.componentMap[component.GetID()] = component
		logger.InfoFmt("The %v component has imported, id: %v", compTpy, component.GetID())
	}

	return nil
}
