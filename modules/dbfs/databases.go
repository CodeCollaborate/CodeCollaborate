package dbfs

// DBConnection is the general interface to ensure all connections can be easily closed
type DBConnection interface {
	close() error
}
