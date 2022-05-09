package models_test

import (
	"github.com/ditsuke/go-amizone/amizone/internal/models"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestClassSchedule_Sort(t *testing.T) {
	testCases := []struct {
		name     string
		schedule models.ClassSchedule
	}{
		{
			name: "2 classes - latter class in slice is earlier",
			schedule: models.ClassSchedule{
				{StartTime: time.Now()},
				{StartTime: time.Now().Add(-1 * time.Hour * 24)},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			tc.schedule.Sort()
			for i := 0; i < len(tc.schedule)-1; i++ {
				g.Expect(tc.schedule[i].StartTime.Before(tc.schedule[i+1].StartTime)).To(BeTrue())
			}
		})
	}
}

func TestClassSchedule_FilterByDate(t *testing.T) {
	testCases := []struct {
		name        string
		schedule    models.ClassSchedule
		filterDate  time.Time
		expectedLen int
	}{
		{
			name: "2 classes - one is on a past date",
			schedule: models.ClassSchedule{
				{StartTime: time.Now()},
				{StartTime: time.Now().Add(-1 * time.Hour * 24)},
			},
			filterDate:  time.Now(),
			expectedLen: 1,
		},
		{
			name: "2 classes - one is on a future date",
			schedule: models.ClassSchedule{
				{StartTime: time.Now()},
				{StartTime: time.Now().Add(1 * time.Hour * 24)},
			},
			filterDate:  time.Now(),
			expectedLen: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			filtered := tc.schedule.FilterByDate(tc.filterDate)
			g.Expect(len(filtered)).To(Equal(tc.expectedLen))
		})
	}
}
