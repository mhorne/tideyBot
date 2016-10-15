package modules

// Main interface for each module to implement
type module interface {
	GetModuleName()
	Initialize()
}
