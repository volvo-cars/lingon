package bla

import (
	"errors"
	"fmt"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nkeys"
)

// User represents a NATS user.
type User struct {
	// ID of the user, which for NATS is the public key of the user
	// and the subject of the user's JWT.
	ID string
	// NKey of the user.
	// The NKey (or "seed") can be converted into the user public
	// and private keys. The public key must match the user ID.
	NKey []byte
	// JWT of the account.
	// The JWT contains the user claims (i.e. name, permissions, limits, etc.)
	// and is signed using an account NKey (that defines which account this user
	// belongs to).
	//
	// In a way, the JWT *is* the user definition.
	// We do not need to send the JWT to the NATS server to *create* the user.
	// The user exists once the JWT is created.
	// A client can connect by providing the JWT, to which the server sends a nonce,
	// which is signed using the user's private key.
	// This verifies that the user has access to the NKeys of the user (which is the
	// public/private key pair).
	JWT string
}

// UserRequest defines the options available for creating a NATS user.
type UserRequest struct {
	// Name of the user
	Name string
}

// SyncUser takes a [UserRequest] and an optional [User], and
// uses the provided [nats.Conn] to create or update the user.
//
// The user JWT is signed using the account NKey to which this user belongs.
//
// The [User] is optional.
// In case it is not provided, a new NATS user is created, with a new key pair.
func SyncUser(
	accountNKey []byte,
	user *User,
	req UserRequest,
) (*User, error) {

	// If user is nil, we need to create a new user
	if user == nil {
		user, err := newUser(accountNKey, req)
		if err != nil {
			return nil, fmt.Errorf("creating user: %w", err)
		}
		// We do not need to send the JWT to the server to create the user.
		return user, nil
	}

	userKeyPair, err := nkeys.FromSeed(user.NKey)
	if err != nil {
		return nil, fmt.Errorf("getting user key pair from nkey: %w", err)
	}
	userPublicKey, err := userKeyPair.PublicKey()
	if err != nil {
		return nil, fmt.Errorf("getting user public key from key pair: %w", err)
	}

	if user.ID != userPublicKey {
		return nil, fmt.Errorf("user ID does not match public key from user NKey")
	}

	// Create a user JWT based on the request
	userJWT, err := generateUserJWT(accountNKey, userPublicKey, req)
	if err != nil {
		return nil, fmt.Errorf("creating user JWT: %w", err)
	}
	user.JWT = userJWT

	return user, nil
}

// newUser creates a whole new user and returns the key pair and JWT.
// It is used when a brand new user is requested, and it does not
// communicate with the NATS server (only returns the new [User]).
func newUser(accountNKey []byte, req UserRequest) (*User, error) {
	userKeyPair, err := nkeys.CreateUser()
	if err != nil {
		return nil, fmt.Errorf("creating user nkeys: %w", err)
	}
	userPublicKey, err := userKeyPair.PublicKey()
	if err != nil {
		return nil, fmt.Errorf("getting user public key: %w", err)
	}
	userSeed, err := userKeyPair.Seed()
	if err != nil {
		return nil, fmt.Errorf("getting user seed from key pair: %w", err)
	}
	userJWT, err := generateUserJWT(accountNKey, userPublicKey, req)
	if err != nil {
		return nil, fmt.Errorf("creating user JWT: %w", err)
	}
	return &User{
		ID:   userPublicKey,
		NKey: userSeed,
		JWT:  userJWT,
	}, nil
}

// generateUserJWT creates a new user JWT with the given subject (user ID / public key)
// using the given [UserRequest] to populate the claims, and finally signing
// it with the account NKey.
//
// It does not talk with the server, but it does validate the claims locally.
func generateUserJWT(accountNKey []byte, userPublicKey string, req UserRequest) (string, error) {
	accountKeyPair, err := nkeys.FromSeed(accountNKey)
	if err != nil {
		return "", fmt.Errorf("getting account key pair from nkey: %w", err)
	}
	accountPublicKey, err := accountKeyPair.PublicKey()
	if err != nil {
		return "", fmt.Errorf("getting account public key from key pair: %w", err)
	}

	userClaims := jwt.NewUserClaims(userPublicKey)
	userClaims.Name = req.Name
	// Set the issuer_account which is the link saying which account the user
	// belongs to.
	userClaims.IssuerAccount = accountPublicKey

	// Validate user claims
	vr := jwt.ValidationResults{}
	userClaims.Validate(&vr)
	if vr.IsBlocking(true) {
		var vErr error
		for _, iss := range vr.Issues {
			vErr = errors.Join(vErr, iss)
		}
		return "", fmt.Errorf("invalid user claims: %w", err)
	}

	userJWT, err := userClaims.Encode(accountKeyPair)
	if err != nil {
		return "", fmt.Errorf("encoding user claims: %w", err)
	}
	return userJWT, nil
}
