package keystore

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"gitlab.com/makeos/mosdef/crypto"
	"gitlab.com/makeos/mosdef/types/core"
	"gitlab.com/makeos/mosdef/util"
)

const (
	DefaultPassphrase = "passphrase"
)

// CreateKey creates a new privKey
func (ks *Keystore) CreateKey(key *crypto.Key, keyType core.KeyType, passphrase string) error {

	// Check whether the privKey already exists. Return error if true.
	exist, err := ks.Exist(key.Addr().String())
	if err != nil {
		return err
	} else if exist {
		return fmt.Errorf("key already exists")
	}

	// When no passphrase is provided, we use a default, publicly known
	// and completely unsafe passphrase just so we don't leave the keys
	// in a non-encrypted state and be forced to write special logic.
	if passphrase == "" {
		passphrase = DefaultPassphrase
	}

	// Harden the passphrase
	passphraseHardened := hardenPassword([]byte(passphrase))

	// Create the serialized privKey data
	keyData := util.ToBytes(core.KeyPayload{
		SecretKey:     key.PrivKey().Base58(),
		Type:          int(keyType),
		FormatVersion: Version,
	})

	// Encode the privKey data with base58 checksum enabled and encrypt
	b58AcctBs := base58.CheckEncode(keyData, 1)
	ct, err := util.Encrypt([]byte(b58AcctBs), passphraseHardened[:])
	if err != nil {
		return errors.Wrap(err, "privKey encryption failed")
	}

	addr := key.Addr()
	if keyType == core.KeyTypePush {
		addr = key.PushAddr()
	}

	// Save the privKey to disk
	now := time.Now().Unix()
	fileName := createKeyFileName(now, addr.String(), passphrase)
	err = ioutil.WriteFile(filepath.Join(ks.dir, fileName), ct, 0644)
	if err != nil {
		return err
	}

	return nil
}

func createKeyFileName(timeNow int64, addr, passphrase string) string {
	fn := fmt.Sprintf("%d_%s", timeNow, addr)
	if passphrase == DefaultPassphrase {
		fn = fn + "_unsafe"
	}
	return fn
}

// CreateCmd creates a new privKey in the keystore.
// It will prompt the user to obtain encryption passphrase if one is not provided.
// If seed is non-zero, it is used as the seed for privKey generation, otherwise,
// one will be randomly generated.
// If passphrase is a file path, the file is read and its content is used as the
// passphrase.
// If nopass is true, the default encryption passphrase is used and the privKey will be marked 'unsafe'
func (ks *Keystore) CreateCmd(keyType core.KeyType, seed int64, passphrase string, nopass bool) (*crypto.Key, error) {

	var passFromPrompt string
	var err error

	// If no passphrase is provided, start an interactive session to
	// collect the passphrase
	if !nopass && strings.TrimSpace(passphrase) == "" {
		fmt.Println("Your new privKey needs to be locked with a passphrase. Please enter a passphrase.")
		passFromPrompt, err = ks.AskForPassword()
		if err != nil {
			return nil, err
		}
	}

	// But if the passphrase is set and is a valid file, read it and use as passphrase
	if len(passphrase) > 0 && (os.IsPathSeparator(passphrase[0]) || (len(passphrase) >= 2 && passphrase[:2] == "./")) {
		content, err := ioutil.ReadFile(passphrase)
		if err != nil {
			return nil, errors.Wrap(err, "unable to read passphrase file")
		}
		passFromPrompt = strings.TrimSpace(strings.Trim(string(content), "/n"))
	} else if len(passphrase) > 0 {
		passFromPrompt = passphrase
	}

	// Generate a privKey
	key, err := crypto.NewKey(nil)
	if seed != 0 {
		key, err = crypto.NewKey(&seed)
	}
	if err != nil {
		return nil, err
	}

	// Create the privKey
	if err := ks.CreateKey(key, keyType, passFromPrompt); err != nil {
		return nil, err
	}

	fmt.Println("New privKey created, encrypted and stored.")
	if keyType == core.KeyTypeAccount {
		fmt.Println("Address:", color.CyanString(key.Addr().String()))
	} else if keyType == core.KeyTypePush {
		fmt.Println("Address:", color.CyanString(key.PushAddr().String()))
	}

	return key, nil
}