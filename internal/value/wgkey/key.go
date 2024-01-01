package wgkey

import (
	"fmt"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type (
	Key[T Type] wgtypes.Key
	Type        interface{ keyType() }
	Private     struct{}
	Public      struct{}
	PreShared   struct{}
)

func (Public) keyType()    {}
func (Private) keyType()   {}
func (PreShared) keyType() {}

type ErrInvalidKeyType[T Type] struct{}

func (ErrInvalidKeyType[T]) Error() string { return fmt.Sprintf("invalid key type %T", *new(T)) }
