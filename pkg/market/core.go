package market

import (
	"math"

	"github.com/umichan0621/steam/pkg/common"
)

type Core struct {
	language string
	currency string
	country  string
}

func (core *Core) Init() {
	core.language = "english"
	core.currency = common.CurrencyUSD
	core.country = "CN"
}

func (core *Core) SetLanguage(language string) { core.language = language }

func (core *Core) SetCurrency(currency string) { core.currency = currency }

func (core *Core) SetCountry(country string) { core.country = country }

func ReceivedPrice(payment float64) float64 {
	received := payment / 1.15
	received = math.Round(received*100.0) / 100.0
	return received
}
