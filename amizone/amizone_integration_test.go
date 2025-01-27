//go:build integration

package amizone_test

import (
	"os"
	"testing"

	"github.com/random2907/go-amizone/amizone"
	. "github.com/onsi/gomega"
)

func TestIntegrateNewClient(t *testing.T) {
	g := NewGomegaWithT(t)

	validUser := os.Getenv("AMIZONE_USERNAME")
	validPassword := os.Getenv("AMIZONE_PASSWORD")

	g.Expect(validUser).ToNot(BeEmpty(), "AMIZONE_USERNAME environment variable is not set")
	g.Expect(validPassword).ToNot(BeEmpty(), "AMIZONE_PASSWORD environment variable is not set")

	testCases := []struct {
		name          string
		credentials   amizone.Credentials
		errorMatcher  func(g *GomegaWithT, err error)
		clientMatcher func(g *GomegaWithT, client amizone.ClientInterface)
	}{
		{
			name:        "valid credentials",
			credentials: amizone.Credentials{Username: validUser, Password: validPassword},
			errorMatcher: func(g *GomegaWithT, err error) {
				g.Expect(err).To(BeNil())
			},
			clientMatcher: func(g *GomegaWithT, client amizone.ClientInterface) {
				g.Expect(client).ToNot(BeNil())
				g.Expect(client.DidLogin()).To(BeTrue())
			},
		},
		{
			name:        "invalid credentials",
			credentials: amizone.Credentials{Username: "this-user-does-not-exist", Password: "neither-does-this-password"},
			errorMatcher: func(g *GomegaWithT, err error) {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(ContainSubstring(amizone.ErrFailedLogin))
				g.Expect(err.Error()).To(ContainSubstring(amizone.ErrInvalidCredentials))
			},
			clientMatcher: func(g *GomegaWithT, client amizone.ClientInterface) {
				g.Expect(client).ToNot(BeNil())
				g.Expect(client.DidLogin()).To(BeFalse())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			client, err := amizone.NewClient(tc.credentials, nil)
			tc.errorMatcher(g, err)
			tc.clientMatcher(g, client)
		})
	}
}

func TestIntegrateAmizone_GetProfile(t *testing.T) {
	// g := NewWithT(t)

	// goal: test that we can get a profile matching
	// the information in the environment (need to add name, UUID, etc as environment variables)
}
