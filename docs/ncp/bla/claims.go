package bla

import (
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nats.go"
)

func validateClaims(claims jwt.Claims) error {
	vr := jwt.ValidationResults{}
	claims.Validate(&vr)
	if vr.IsBlocking(true) {
		var vErr error
		for _, iss := range vr.Issues {
			vErr = errors.Join(vErr, iss)
		}
		return vErr
	}
	return nil
}

// claimsForAccount looks up the account on the server and returns the account claims.
func claimsForAccount(nc *nats.Conn, accountID string) (*jwt.AccountClaims, error) {
	// Lookup the account on the server.
	// If the account does not exist, the observed behaviour is that an empty JWT
	// is returned, without an error.
	// And re-creating the account using the same ID (public key) seems to work
	// just fine, so no need to handle an account that does not exist.
	lookupMessage, err := nc.Request(
		fmt.Sprintf("$SYS.REQ.ACCOUNT.%s.CLAIMS.LOOKUP", accountID),
		nil,
		time.Second,
	)
	if err != nil {
		return nil, fmt.Errorf("looking up account with ID %s: %w", accountID, err)
	}
	if lookupMessage.Data == nil {
		return nil, fmt.Errorf("account with ID %s does not exist", accountID)
	}
	// var ac jwt.AccountClaims
	ac, err := jwt.DecodeAccountClaims(string(lookupMessage.Data))
	if err != nil {
		return nil, fmt.Errorf("decoding account claims: %w", err)
	}
	return ac, nil
}
