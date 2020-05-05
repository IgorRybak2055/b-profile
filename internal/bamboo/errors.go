package bamboo

import "encoding/json"

// Error struct for return api error
type Error struct {
	Code int
	Err  error
}

func newError(code int, err error) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

// Error implements error interface.
func (e Error) Error() string {
	return e.Err.Error()
}

// MarshalJSON implements Marshaller interface.
func (e Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Code int    `json:"code"`
		Err  string `json:"error"`
	}{
		Code: e.Code,
		Err:  e.Err.Error(),
	})
}
