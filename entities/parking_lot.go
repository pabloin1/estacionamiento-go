package entities

import "sync"

type ParkingLot struct {
	spots      chan int
	entryGate  chan struct{}
	lock       *sync.Mutex
	observers  []Observer
}

func NewParkingLot(spots chan int, entryGate chan struct{}, lock *sync.Mutex) *ParkingLot {
	return &ParkingLot{
		spots:     spots,
		entryGate: entryGate,
		lock:      lock,
		observers: []Observer{},
	}
}

func (lot *ParkingLot) Register(observer Observer) {
	lot.lock.Lock()
	defer lot.lock.Unlock()
	lot.observers = append(lot.observers, observer)
}

func (lot *ParkingLot) NotifyObservers(vehicleID int, status string) {
	lot.lock.Lock()
	defer lot.lock.Unlock()
	for _, observer := range lot.observers {
		if observer.GetID() == vehicleID {
			observer.Update(status)
		}
	}
}
