package wgapi

import (
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// NewPrivatePublic generates a new private key and also returns its corresponding public key.
func NewPrivatePublic() (private PrivateKey, public PublicKey, err error) {
	private, err = NewPrivate()
	if err != nil {
		return
	}
	p, _ := private.Public()
	public = PublicKey(p)
	return
}

// NewPrivate generates a new private key.
func NewPrivate() (PrivateKey, error) {
	k, err := wgtypes.GeneratePrivateKey()
	return PrivateKey(k), err
}

// NewPreshared generates a new preshared key.
func NewPreshared() (PresharedKey, error) {
	k, err := wgtypes.GenerateKey()
	return PresharedKey(k), err
}
