package domain

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewAccountDomain, NewUsernameDomain, NewGoogleDomain, NewAppleDomain, NewFacebookDomain)
