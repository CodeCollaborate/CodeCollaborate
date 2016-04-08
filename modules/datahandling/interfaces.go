package datahandling

// ProcessorInterface should be implemented by all request models.
// Provides standard interface for calling the processing.
type ProcessorInterface interface {
	Process()
}
