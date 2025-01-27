package parse_test

import (
	"testing"

	"github.com/random2907/go-amizone/amizone/internal/mock"
	"github.com/random2907/go-amizone/amizone/internal/parse"
	"github.com/random2907/go-amizone/amizone/models"
	. "github.com/onsi/gomega"
)

func TestCourses(t *testing.T) {
	testCases := []struct {
		name           string
		bodyFile       mock.File
		coursesMatcher func(g *GomegaWithT, courses models.Courses)
		errMatcher     func(g *GomegaWithT, err error)
	}{
		{
			name:     "current courses page",
			bodyFile: mock.CoursesPage,
			coursesMatcher: func(g *GomegaWithT, courses models.Courses) {
				g.Expect(courses).ToNot(BeNil())
				g.Expect(len(courses)).To(Equal(8))
			},
			errMatcher: func(g *GomegaWithT, err error) {
				g.Expect(err).ToNot(HaveOccurred())
			},
		},
		{
			name:     "semester wise courses page",
			bodyFile: mock.CoursesPageSemWise,
			coursesMatcher: func(g *GomegaWithT, courses models.Courses) {
				g.Expect(courses).ToNot(BeNil())
				g.Expect(len(courses)).To(Equal(8))
			},
			errMatcher: func(g *GomegaWithT, err error) {
				g.Expect(err).ToNot(HaveOccurred())
			},
		},
		{
			name:     "invalid courses page (login page)",
			bodyFile: mock.LoginPage,
			coursesMatcher: func(g *GomegaWithT, courses models.Courses) {
				g.Expect(courses).To(BeNil())
			},
			errMatcher: func(g *GomegaWithT, err error) {
				g.Expect(err.Error()).To(ContainSubstring(parse.ErrFailedToParse))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			fileReader, err := testCase.bodyFile.Open()
			g.Expect(err).ToNot(HaveOccurred())
			courses, err := parse.Courses(fileReader)
			testCase.coursesMatcher(g, courses)
			testCase.errMatcher(g, err)
		})
	}
}
