package gatt

import (
	"log"
	"sync"

	"github.com/krzysiekbielicki/ble"
	"github.com/krzysiekbielicki/ble/linux/att"
)

// NewServer ...
func NewServerWithName(name string) (*Server, error) {
	return &Server{
		name:name,
		svcs: defaultServices(name),
		db:   att.NewDB(defaultServices(name), uint16(1)),
	}, nil
}

// NewServer ...
func NewServer() (*Server, error) {
	return NewServerWithName("Gopher")
}

// Server ...
type Server struct {
	sync.Mutex
	name string

	svcs []*ble.Service
	db   *att.DB
}

// AddService ...
func (s *Server) AddService(svc *ble.Service) error {
	s.Lock()
	defer s.Unlock()
	s.svcs = append(s.svcs, svc)
	s.db = att.NewDB(s.svcs, uint16(1)) // ble attrs start at 1
	return nil
}

// RemoveAllServices ...
func (s *Server) RemoveAllServices() error {
	s.Lock()
	defer s.Unlock()
	s.svcs = defaultServices(s.name)
	s.db = att.NewDB(s.svcs, uint16(1)) // ble attrs start at 1
	return nil
}

// SetServices ...
func (s *Server) SetServices(svcs []*ble.Service) error {
	s.Lock()
	defer s.Unlock()
	s.svcs = append(defaultServices(s.name), svcs...)
	s.db = att.NewDB(s.svcs, uint16(1)) // ble attrs start at 1
	return nil
}

// DB ...
func (s *Server) DB() *att.DB {
	return s.db
}

func defaultServices(name string) []*ble.Service {
	// https://developer.bluetooth.org/gatt/characteristics/Pages/CharacteristicViewer.aspx?u=org.bluetooth.characteristic.ble.appearance.xml
	var gapCharAppearanceGenericComputer = []byte{0x00, 0x80}

	gapSvc := ble.NewService(ble.GAPUUID)
	gapSvc.NewCharacteristic(ble.DeviceNameUUID).SetValue([]byte(name))
	gapSvc.NewCharacteristic(ble.AppearanceUUID).SetValue(gapCharAppearanceGenericComputer)
	gapSvc.NewCharacteristic(ble.PeripheralPrivacyUUID).SetValue([]byte{0x00})
	gapSvc.NewCharacteristic(ble.ReconnectionAddrUUID).SetValue([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	gapSvc.NewCharacteristic(ble.PeferredParamsUUID).SetValue([]byte{0x06, 0x00, 0x06, 0x00, 0x00, 0x00, 0xd0, 0x07})

	gattSvc := ble.NewService(ble.GATTUUID)
	gattSvc.NewCharacteristic(ble.ServiceChangedUUID).HandleIndicate(
		ble.NotifyHandlerFunc(func(r ble.Request, n ble.Notifier) {
			log.Printf("TODO: indicate client when the services are changed")
			for {
				select {
				case <-n.Context().Done():
					log.Printf("count: Notification unsubscribed")
					return
				}
			}
		}))
	return []*ble.Service{gapSvc, gattSvc}
}
