package parse

import (
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
	"k8s.io/klog/v2"
)

const verificationTokenName = "__RequestVerificationToken"

func LoginFormTokens(body io.Reader) (map[string]string) {
	dom, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		klog.Errorf("failed to parse login page: %s. Was the right page passed?", err.Error())
	}
	formValues := make(map[string]string)
	formValues["__RequestVerificationToken"] = dom.Find("input[name='__RequestVerificationToken']").AttrOr("value", "")
	formValues["Salt"] = dom.Find("input[name='Salt']").AttrOr("value", "")
	formValues["SecretNumber"] = dom.Find("input[name='SecretNumber']").AttrOr("value", "")
	formValues["Signature"] = dom.Find("input[name='Signature']").AttrOr("value", "")
	formValues["Challenge"] = dom.Find("input[name='Challenge']").AttrOr("value", "")
	for name, value := range formValues {
		if value == "" {
			klog.Warningf("Form value for '%s' is missing.", name)
		}
	}
	return formValues
}


func VerificationTokenFromDom(dom *goquery.Document) string {
	return dom.Find(fmt.Sprintf("input[name='%s']", verificationTokenName)).AttrOr("value", "")
}
