package keystore

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"gitlab.com/makeos/mosdef/types/core"
)

// UIUnlockAccount renders a CLI UI to unlock a target keystore.
// addressOrIndex: The address or index of the keystore.
// passphrase: The user supplied passphrase. If not provided, an
// interactive session will be started to collect the passphrase
func (ks *Keystore) UIUnlockAccount(addressOrIndex, passphrase string) (core.StoredKey, error) {

	var err error

	// Get the keystore.
	storedAcct, err := ks.GetByIndexOrAddress(addressOrIndex)
	if err != nil {
		return nil, err
	}

	fmt.Println(color.HiBlackString("Chosen Account: ") + storedAcct.GetAddress())

	// Set the passphrase to the default passphrase if account
	// is encrypted with unsafe passphrase
	if storedAcct.IsUnsafe() {
		passphrase = DefaultPassphrase
	}

	// Ask for passphrase if unset
	if passphrase == "" {
		passphrase = ks.AskForPasswordOnce()
	}

	// If passphrase is not a path to a file, proceed to unlock the keystore
	if !strings.HasPrefix(passphrase, "./") &&
		!strings.HasPrefix(passphrase, "/") &&
		filepath.Ext(passphrase) == "" {
		goto unlock
	}

	// So, 'passphrase' contains a file path, read the passphrase from it
	passphrase, err = readPassFromFile(passphrase)
	if err != nil {
		return nil, err
	}

unlock:
	if err = storedAcct.Unlock(passphrase); err != nil {
		return nil, err
	}

	return storedAcct, nil
}