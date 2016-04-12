package datahandling

/**
 * Interfaces.go describes the interfaces used in the data handling for the CodeCollaborate Server.
 */

// ProcessorInterface should be implemented by all request models.
// Provides standard interface for calling the processing.
type ProcessorInterface interface {
	Process()
}
