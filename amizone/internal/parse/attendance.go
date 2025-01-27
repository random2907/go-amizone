package parse

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/random2907/go-amizone/amizone/models"
	"k8s.io/klog/v2"
)

// Attendance attempts to parse course attendance information from the Amizone home page
// into a models.AttendanceRecords instance.
func Attendance(body io.Reader) (models.AttendanceRecords, error) {
	const (
		AttendanceTableTitle = "My Attendance"
	)

	dom, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrFailedToParseDOM, err)
	}

	if !IsLoggedInDOM(dom) {
		return nil, errors.New(ErrNotLoggedIn)
	}

	// The attendance record is stored in a div-soup "widget". There are no semantic identifiers in the markup,
	// so we search this widget by title.
	attendanceWidgetHeader := dom.Find(".widget-header").
		Filter(fmt.Sprintf(":containsOwn('%s')", AttendanceTableTitle))
	attendanceList := attendanceWidgetHeader.Parent().Find("ul#tasks li")

	if attendanceWidgetHeader.Length() == 0 || attendanceList.Length() == 0 {
		klog.Warning("Failed to find the attendance widget header. Are we logged in and on the right page?")
		return nil, errors.New(ErrFailedToParse)
	}

	attendance := make(models.AttendanceRecords, attendanceList.Length())
	attendanceList.Each(func(i int, record *goquery.Selection) {
		attended, held := func() (int, int) {
			raw := record.Find("div.class-count span").Text()
			sanitized := strings.Trim(raw, " \"")
			divided := strings.Split(sanitized, "/")
			if len(divided) != 2 {
				klog.Warning("Attendance string has unexpected format!")
			}

			return parseToInt(divided[0]), parseToInt(divided[1])
		}()

		courseAttendance := models.AttendanceRecord{
			Course: models.CourseRef{
				Code: func() string {
					raw := record.Find("span.sub-code").Text()
					return strings.TrimSpace(raw)
				}(),
				Name: func() string {
					rawInner := record.Find("span.lbl").Text()
					spaceIndex := strings.IndexRune(rawInner, ' ')
					return strings.TrimSpace(rawInner[spaceIndex:])
				}(),
			},
			Attendance: models.Attendance{
				ClassesAttended: int32(attended),
				ClassesHeld:     int32(held),
			},
		}

		attendance[i] = courseAttendance
	})

	return attendance, nil
}

// parseToInt parses an integer to a string, logs on failure.
func parseToInt(raw string) int {
	i, err := strconv.Atoi(raw)
	if err != nil {
		klog.Errorf("Failed to parse string to int: %s", err.Error())
	}
	return i
}
