package datahandling

/**
 * requestmap.go provides the pseudo-factory map for looking up the associated request
 */

// pseudo
var requestMap = make(map[string]interface{})

func init() {

}

func GetRequestMap(name string) Request {
	return nil
}

type projectLookupRequest struct {
	ProjectId []uint64
	AbstractRequest
}

