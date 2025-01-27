package amizone

import (
	"net/http"
	"time"

	"github.com/random2907/go-amizone/amizone/models"
)

// ClientInterface is an exported interface for client to make mocking and testing more convenient.
type ClientInterface interface {
	DidLogin() bool
	GetAttendance() (models.AttendanceRecords, error)
	GetClassSchedule(year int, month time.Month, date int) (models.ClassSchedule, error)
	GetExamSchedule() (*models.ExaminationSchedule, error)
}

// Interface compliance constraint for Client
var _ ClientInterface = &Client{}

// ClientFactoryInterface is a type for functions that return ClientInterface
// instances. Functions returning concrete types need to be wrapped by functions
// that return the interface; apparently a limitation of the Go compiler's type
// inference.
type ClientFactoryInterface func(cred Credentials, httpClient *http.Client) (ClientInterface, error)

// Interface compliance constraint for NewClient [requires wrapping]
var _ ClientFactoryInterface = func(cred Credentials, httpClient *http.Client) (ClientInterface, error) {
	return NewClient(cred, httpClient)
}
