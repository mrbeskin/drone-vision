package control

import (
	"fmt"
	"sync"
	"time"

	"gobot.io/x/gobot/platforms/dji/tello"
)

const THROTTLE_INIT = 40

type Flyable interface {
	Forward(throttle int)
	Backward(throttle int)
	Left(throttle int)
	Right(throttle int)
	Up(throttle int)
	Down(throttle int)
	Clockwise(throttle int)
	CounterClockwise(throttle int)
}

type FlightController struct {
	ThrottleEvents         EventQueue
	ForwardEvents          EventQueue
	BackEvents             EventQueue
	LeftEvents             EventQueue
	RightEvents            EventQueue
	UpEvents               EventQueue
	DownEvents             EventQueue
	ClockwiseEvents        EventQueue
	CounterClockwiseEvents EventQueue
	vehicle                Flyable
	throttle               int
	mu                     *sync.Mutex
}

func NewFlightController(vehicle Flyable) *FlightController {
	return &FlightController{
		ThrottleEvents:         NewEventQueue(),
		ForwardEvents:          NewEventQueue(),
		BackEvents:             NewEventQueue(),
		LeftEvents:             NewEventQueue(),
		RightEvents:            NewEventQueue(),
		UpEvents:               NewEventQueue(),
		DownEvents:             NewEventQueue(),
		ClockwiseEvents:        NewEventQueue(),
		CounterClockwiseEvents: NewEventQueue(),
		vehicle:                vehicle,
		throttle:               THROTTLE_INIT,
		mu:                     &sync.Mutex{},
	}
}

func (fc *FlightController) StartControl() {
	go fc.processLoop()
}

func (fc *FlightController) processLoop() {
	for {
		fc.ProcessAll()
		time.Sleep(25 * time.Millisecond)
	}
}

func (fc *FlightController) ProcessAll() {
	fc.ProcessThrottleChange()
	fc.ProcessForwardEvents()
	fc.ProcessBackEvents()
	fc.ProcessLeftEvents()
	fc.ProcessRightEvents()
	fc.ProcessUpEvents()
	fc.ProcessDownEvents()
	fc.ProcessClockwiseEvents()
	fc.ProcessCounterClockwiseEvents()
}

func (fc *FlightController) ThrottleUp() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.ThrottleEvents = fc.ThrottleEvents.Push(1)
}

func (fc *FlightController) ThrottleDown() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.ThrottleEvents = fc.ThrottleEvents.Push(-1)
}

func (fc *FlightController) Forward() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.ForwardEvents = fc.ForwardEvents.Push(fc.throttle)
}

func (fc *FlightController) Backward() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.BackEvents = fc.BackEvents.Push(fc.throttle)
}

func (fc *FlightController) Left() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.LeftEvents = fc.LeftEvents.Push(fc.throttle)
}

func (fc *FlightController) Right() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.RightEvents = fc.RightEvents.Push(fc.throttle)
}

func (fc *FlightController) Up() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.UpEvents = fc.UpEvents.Push(fc.throttle)
}

func (fc *FlightController) Down() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.DownEvents = fc.DownEvents.Push(fc.throttle)
}

func (fc *FlightController) Clockwise() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.ClockwiseEvents = fc.ClockwiseEvents.Push(fc.throttle)
}

func (fc *FlightController) CounterClockwise() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.CounterClockwiseEvents = fc.CounterClockwiseEvents.Push(fc.throttle)
}

func (fc *FlightController) ProcessThrottleChange() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	if !fc.ThrottleEvents.isEmpty() {
		ev, q := fc.ThrottleEvents.Pop()
		fc.throttle += ev
		fc.ThrottleEvents = q
	}
}

func (fc *FlightController) ProcessForwardEvents() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	if !fc.ForwardEvents.isEmpty() {
		ev, q := fc.ForwardEvents.Pop()
		fc.vehicle.Forward(ev)
		fc.ForwardEvents = q
		fmt.Println("fwprocess")
	}
}

func (fc *FlightController) ProcessBackEvents() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	if !fc.BackEvents.isEmpty() {
		ev, q := fc.BackEvents.Pop()
		fc.vehicle.Backward(ev)
		fc.BackEvents = q
	}
}

func (fc *FlightController) ProcessLeftEvents() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	if !fc.LeftEvents.isEmpty() {
		ev, q := fc.LeftEvents.Pop()
		fc.vehicle.Left(ev)
		fc.LeftEvents = q
	}
}

func (fc *FlightController) ProcessRightEvents() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	if !fc.RightEvents.isEmpty() {
		ev, q := fc.RightEvents.Pop()
		fc.vehicle.Right(ev)
		fc.RightEvents = q
	} else {
		fc.vehicle.Right(0)
	}
}

func (fc *FlightController) ProcessUpEvents() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	if !fc.UpEvents.isEmpty() {
		ev, q := fc.UpEvents.Pop()
		fc.vehicle.Up(ev)
		fc.UpEvents = q
	}
}

func (fc *FlightController) ProcessDownEvents() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	if !fc.DownEvents.isEmpty() {
		ev, q := fc.DownEvents.Pop()
		fc.vehicle.Down(ev)
		fc.DownEvents = q
	}
}

func (fc *FlightController) ProcessClockwiseEvents() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	if !fc.ClockwiseEvents.isEmpty() {
		ev, q := fc.ClockwiseEvents.Pop()
		fc.vehicle.Clockwise(ev)
		fc.ClockwiseEvents = q
	}
}

func (fc *FlightController) ProcessCounterClockwiseEvents() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	if !fc.CounterClockwiseEvents.isEmpty() {
		ev, q := fc.CounterClockwiseEvents.Pop()
		fc.vehicle.CounterClockwise(ev)
		fc.CounterClockwiseEvents = q
	}
}

// TelloDrone implements Flyable and is in charge of actually calling
// the driver to make offboard communication to the drone
type TelloDrone struct {
	driver   *tello.Driver
	throttle int
}

func NewTelloDrone(driver *tello.Driver) *TelloDrone {
	return &TelloDrone{
		driver: driver,
	}
}

func (td *TelloDrone) Forward(throttle int) {
	td.driver.Forward(throttle)
}

func (td *TelloDrone) Backward(throttle int) {
	td.driver.Backward(throttle)
}

func (td *TelloDrone) Left(throttle int) {
	td.driver.Left(throttle)
}

func (td *TelloDrone) Right(throttle int) {
	td.driver.Right(throttle)
}

func (td *TelloDrone) Up(throttle int) {
	td.driver.Up(throttle)
}

func (td *TelloDrone) Down(throttle int) {
	td.driver.Down(throttle)
}

func (td *TelloDrone) Clockwise(throttle int) {
	td.driver.Clockwise(throttle)
}

func (td *TelloDrone) CounterClockwise(throttle int) {
	td.driver.CounterClockwise(throttle)
}

// EventQueue is used to process drone events
// it iss basically just a queue
type EventQueue []int

func NewEventQueue() EventQueue {
	return EventQueue(make([]int, 0))
}

func (eq EventQueue) Push(val int) EventQueue {
	eq = eq.PushZeroIfEmpty()
	return append(eq, val)
}

func (eq EventQueue) Peek() int {
	return eq[0]
}

func (eq EventQueue) Pop() (int, EventQueue) {

	next, eq := eq[0], eq[1:]
	if eq.isEmpty() && (next != 0) {
		eq = append(eq, 0)
	}
	return next, eq
}

func (eq EventQueue) isEmpty() bool {
	if len(eq) == 0 {
		return true
	}
	return false
}

func (eq EventQueue) PushZeroIfEmpty() EventQueue {
	if eq.isEmpty() {
		return append(eq, 0)
	}
	return eq
}
