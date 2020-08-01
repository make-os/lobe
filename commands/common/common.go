package common

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/pkg/errors"
	client2 "github.com/themakeos/lobe/api/remote/client"
	"github.com/themakeos/lobe/api/rpc/client"
	"github.com/themakeos/lobe/api/utils"
	"github.com/themakeos/lobe/config"
	"github.com/themakeos/lobe/keystore"
	"github.com/themakeos/lobe/keystore/types"
	"github.com/themakeos/lobe/modules"
	remotetypes "github.com/themakeos/lobe/remote/types"
	"github.com/themakeos/lobe/util/colorfmt"
)

var (
	ErrBodyRequired  = fmt.Errorf("body is required")
	ErrTitleRequired = fmt.Errorf("title is required")
)

// pagerWriter describes a function for writing a specified content to a pager program
type PagerWriter func(pagerCmd string, content io.Reader, stdOut, stdErr io.Writer)

// WriteToPager spawns the specified page, passing the given content to it
func WriteToPager(pagerCmd string, content io.Reader, stdOut, stdErr io.Writer) {
	args := strings.Split(pagerCmd, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	cmd.Stdin = content
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(stdOut, err.Error())
		fmt.Fprint(stdOut, content)
		return
	}
}

// UnlockKeyArgs contains arguments for UnlockKey
type UnlockKeyArgs struct {
	KeyAddrOrIdx string
	AskPass      bool
	Passphrase   string
	TargetRepo   remotetypes.LocalRepo
	Prompt       string
	Stdout       io.Writer
}

// KeyUnlocker describes a function for unlocking a keystore key.
type KeyUnlocker func(cfg *config.AppConfig, args *UnlockKeyArgs) (types.StoredKey, error)

// UnlockKey takes a key address or index, unlocks it and returns the key.
// - It will using the given passphrase if set, otherwise
// - if the target repo is set, it will try to get it from the git config (user.passphrase).
// - If passphrase is still unknown, it will attempt to get it from an environment variable.
// - On success, args.Passphrase is updated with the passphrase used to unlock the key.
func UnlockKey(cfg *config.AppConfig, args *UnlockKeyArgs) (types.StoredKey, error) {

	// Get the key from the key store
	ks := keystore.New(cfg.KeystoreDir())
	if args.Stdout != nil {
		ks.SetOutput(args.Stdout)
	}

	// Get the key by address or index
	key, err := ks.GetByIndexOrAddress(args.KeyAddrOrIdx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find key (%s)", args.KeyAddrOrIdx)
	}

	// If passphrase is unset and target repo is set, attempt to get the
	// passphrase from 'user.passphrase' config.
	unprotected := key.IsUnprotected()
	if !unprotected && args.Passphrase == "" && args.TargetRepo != nil {
		repoCfg, err := args.TargetRepo.Config()
		if err != nil {
			return nil, errors.Wrap(err, "unable to get repo config")
		}
		args.Passphrase = repoCfg.Raw.Section("user").Option("passphrase")

		// If we still don't have a passphrase, get it from the repo scoped env variable.
		if args.Passphrase == "" {
			args.Passphrase = os.Getenv(MakeRepoScopedPassEnvVar(config.AppName, args.TargetRepo.GetName()))
		}
	}

	// If key is protected and still no passphrase,
	// try to get it from the general passphrase env variable
	if !unprotected && args.Passphrase == "" {
		args.Passphrase = os.Getenv(MakePassEnvVar(config.AppName))
	}

	// If key is protected and still no passphrase, exit with error
	if !unprotected && args.Passphrase == "" && !args.AskPass {
		return nil, fmt.Errorf("passphrase of signing key is required")
	}

	key, passphrase, err := ks.UnlockKeyUI(args.KeyAddrOrIdx, args.Passphrase, args.Prompt)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unlock key (%s)", args.KeyAddrOrIdx)
	}

	// Index the passphrase used to unlock the key.
	key.GetMeta()["passphrase"] = passphrase

	return key, nil
}

// MakeRepoScopedPassEnvVar returns a repo-specific env variable
// expected to contain passphrase for unlocking an account.
func MakeRepoScopedPassEnvVar(appName, repoName string) string {
	return strings.ToUpper(fmt.Sprintf("%s_%s_PASS", appName, repoName))
}

// MakePassEnvVar is the name of the env variable expected to contain a key's passphrase.
func MakePassEnvVar(appName string) string {
	return strings.ToUpper(fmt.Sprintf("%s_PASS", appName))
}

type TxStatusTrackerFunc func(stdout io.Writer, hash string, rpcClient client.Client,
	remoteClients []client2.Client) error

// ShowTxStatusTracker tracks transaction status and displays updates to stdout.
func ShowTxStatusTracker(stdout io.Writer, hash string, rpcClient client.Client, remoteClients []client2.Client) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Writer = stdout
	s.Prefix = " "
	s.Start()
	lastStatus := ""
	for {
		time.Sleep(1 * time.Second)
		resp, err := utils.GetTransaction(hash, rpcClient, remoteClients)
		if err != nil {
			s.Stop()
			return err
		}
		if lastStatus == resp.Status {
			continue
		}
		lastStatus = resp.Status
		if resp.Status == modules.TxStatusInMempool {
			s.Suffix = colorfmt.YellowString(" In mempool")
		} else if resp.Status == modules.TxStatusInPushpool {
			s.Suffix = colorfmt.YellowString(" In pushpool")
		} else {
			s.FinalMSG = colorfmt.GreenString("   Confirmed!\n")
			s.Stop()
			break
		}
	}
	return nil
}
