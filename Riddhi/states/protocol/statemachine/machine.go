package statemachine

type Role int

const (
	Requester Role = iota
	Provider
)

type State int

const (
	Idle State = iota

	WaitingForChunkResponse

	WaitingForLotteryTicket

	WaitingForKeyReveal

	Completed
)

//stores current protocol state
type Machine struct {
	Role Role
	State State
}

//Creates machine starting in:Idle state
func NewMachine(role Role) *Machine {
	return &Machine{
		Role: role,
		State: Idle,
	}
}

//This changes state
func (m *Machine) Transition(next State) {
	m.State = next
}