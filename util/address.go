package util

import (
	"regexp"
)

// Address constants
const (
	addressPrefixRepo                = "r/"
	addressPrefixAddressUser         = "a/"
	addressPrefixedIdentifierRegexp  = "^[ar]{1}/[a-zA-Z0-9_-]+$"              // e.g r/abc-xyz or a/abc-xyz
	AddressNamespaceDomainNameRegexp = "^[a-zA-Z0-9_-]+$"                      // e.g abc-xyz_
	addressNamespaceRegexp           = "^[a-zA-Z0-9_-]{3,}/[a-zA-Z0-9_-]{0,}$" // e.g namespace/abc-xyz_
)

// GetPrefixedAddressValue returns the address without the prefix
func GetPrefixedAddressValue(prefixedAddr string) string {
	return prefixedAddr[2:]
}

// IsPrefixedAddressRepo checks whether the given address is a prefixed repo address
func IsPrefixedAddressRepo(addr string) bool {
	if !IsPrefixedAddr(addr) {
		return false
	}
	return addr[:2] == addressPrefixRepo
}

// IsPrefixedAddressUserAccount checks whether the given address is a prefixed repo address
func IsPrefixedAddressUserAccount(addr string) bool {
	if !IsPrefixedAddr(addr) {
		return false
	}
	return addr[:2] == addressPrefixAddressUser &&
		IsValidAddr(addr[2:]) == nil
}

// IsPrefixedAddr checks whether the given address matches a prefixed address
func IsPrefixedAddr(addr string) bool {
	return regexp.MustCompile(addressPrefixedIdentifierRegexp).MatchString(addr)
}

// IsNamespaceURI checks whether the given address is a namespace URI
func IsNamespaceURI(addr string) bool {
	return regexp.MustCompile(addressNamespaceRegexp).MatchString(addr)
}

// Address represents an identifier for a resource
type Address string

func (a Address) String() string {
	return string(a)
}

// Equals checks whether a is equal to addr
func (a Address) Equal(addr Address) bool {
	return a == addr
}

// Empty checks whether the address is empty
func (a Address) Empty() bool {
	return a.String() == ""
}

// IsNamespaceURI checks whether the address is a namespace URI
func (a Address) IsNamespaceURI() bool {
	return IsNamespaceURI(string(a))
}

// IsPrefixed checks whether the address is prefixed with a/ or /r which
// indicates a repo and keystore address respectively
func (a Address) IsPrefixed() bool {
	return IsPrefixedAddr(string(a))
}

// IsPrefixedRepoAddress checks if the address is prefixed by `r/` which is used to
// identity a repo address
func (a Address) IsPrefixedRepoAddress() bool {
	if !a.IsPrefixed() {
		return false
	}
	return string(a)[:2] == addressPrefixRepo
}

// IsPrefixedUserAddress checks if the address is prefixed by
// `a/` which is used to identity an keystore address
func (a Address) IsPrefixedUserAddress() bool {
	if !a.IsPrefixed() {
		return false
	}
	return string(a)[:2] == addressPrefixAddressUser
}

// IsBech32MakerAddress checks whether the address is a
// bech32 address with the general HRP
func (a Address) IsBech32MakerAddress() bool {
	return IsValidAddr(string(a)) == nil
}
