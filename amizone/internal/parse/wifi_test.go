package parse_test

import (
	"net"
	"testing"

	"github.com/random2907/go-amizone/amizone/internal/mock"
	"github.com/random2907/go-amizone/amizone/internal/parse"
	"github.com/random2907/go-amizone/amizone/models"
	. "github.com/onsi/gomega"
)

func TestWifi(t *testing.T) {
	testCases := []struct {
		name             string
		bodyFile         mock.File
		AddressesMatcher func(g *GomegaWithT, info *models.WifiMacInfo)
		errMatcher       func(g *GomegaWithT, err error)
	}{
		{
			name:     "both mac addresses populated",
			bodyFile: mock.WifiPage,
			AddressesMatcher: func(g *GomegaWithT, info *models.WifiMacInfo) {
				g.Expect(info).ToNot(BeNil())
				g.Expect(info.Slots).To(Equal(2))
				g.Expect(info.FreeSlots).To(Equal(0))
				g.Expect(info.RegisteredAddresses).To(HaveLen(2))
				g.Expect(info.RegisteredAddresses).To(ConsistOf(net.HardwareAddr{85, 4, 45, 231, 190, 164}, net.HardwareAddr{253, 213, 20, 24, 12, 139}))
				g.Expect(info.Slots).To(Equal(2))
			},
			errMatcher: func(g *GomegaWithT, err error) {
				g.Expect(err).ToNot(HaveOccurred())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			body, err := tc.bodyFile.Open()
			g.Expect(err).ToNot(HaveOccurred())
			addresses, err := parse.WifiMacInfo(body)
			tc.AddressesMatcher(g, addresses)
			tc.errMatcher(g, err)
		})
	}
}
