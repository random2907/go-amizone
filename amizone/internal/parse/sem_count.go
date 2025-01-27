package parse

import (
	"errors"
	"io"

	"github.com/PuerkitoBio/goquery"
	"github.com/random2907/go-amizone/amizone/models"
)

// Semesters returns the number of ongoing or passed semesters from the Amizone courses page.
func Semesters(body io.Reader) (models.SemesterList, error) {
	dom, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, errors.New(ErrFailedToParseDOM)
	}

	if !IsLoggedInDOM(dom) {
		return nil, errors.New(ErrNotLoggedIn)
	}

	if !isCoursesPage(dom) {
		return nil, errors.New(ErrFailedToParse)
	}

	var semesters models.SemesterList
	dom.Find("#CurrentSemesterInfo option").Each(func(_ int, opt *goquery.Selection) {
		if value := opt.AttrOr("value", ""); value != "" {
			sem := models.Semester{
				Name: CleanString(opt.Text()),
				Ref:  value,
			}
			semesters = append(semesters, sem)
		}
	})

	return semesters, nil
}
