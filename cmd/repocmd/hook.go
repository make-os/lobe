package repocmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	restclient "github.com/make-os/lobe/api/remote/client"
	"github.com/make-os/lobe/api/rpc/client"
	"github.com/make-os/lobe/api/utils"
	"github.com/make-os/lobe/cmd/common"
	"github.com/make-os/lobe/cmd/signcmd"
	"github.com/make-os/lobe/config"
	plumbing2 "github.com/make-os/lobe/remote/plumbing"
	"github.com/make-os/lobe/remote/server"
	"github.com/make-os/lobe/remote/types"
	"github.com/make-os/lobe/util/colorfmt"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type HookArgs struct {
	*signcmd.SignCommitArgs

	// Args is the command arguments
	Args []string

	// AuthMode outputs credentials
	AskPass bool

	// RpcClient is the RPC client
	RPCClient client.Client

	// RemoteClients is the remote server API client.
	RemoteClients []restclient.Client

	// KeyUnlocker is a function for getting and unlocking a push key from keystore
	KeyUnlocker common.KeyUnlocker

	// GetNextNonce is a function for getting the next nonce of the owner account of a pusher key
	GetNextNonce utils.NextNonceGetter

	// SetRemotePushToken is a function for setting push tokens on a git remote config
	SetRemotePushToken server.RemotePushTokenSetter

	CommitSigner signcmd.SignCommitFunc
	TagSigner    signcmd.SignTagFunc
	NoteSigner   signcmd.SignNoteFunc

	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

// HookCmd handles pre-push calls by git
func HookCmd(cfg *config.AppConfig, repo types.LocalRepo, args *HookArgs) error {

	if args.AskPass {
		return HandleAskPass(args.Stdout, args.Stderr, repo, args.Args)
	}

	updates, err := ioutil.ReadAll(args.Stdin)
	if err != nil {
		return err
	}

	// Read the references to be updated
	var references []plumbing.ReferenceName
	for _, line := range strings.Split(strings.TrimSpace(string(updates)), "\n") {
		refname := strings.Split(strings.TrimSpace(line), " ")[0]
		references = append(references, plumbing.ReferenceName(refname))
	}

	rcfg, err := repo.Config()
	if err != nil {
		return err
	}

	// Set hook.origin config var to be used
	// by HandleAskPass to determine which remote tokens to use
	rcfg.Raw.Section("hook").SetOption("curRemote", args.Args[0])
	if err = repo.SetConfig(rcfg); err != nil {
		return errors.Wrap(err, "failed to set `hook.curRemote` value")
	}

	// Sign each reference
	for _, ref := range references {
		if ref.IsBranch() {
			if err := args.CommitSigner(cfg, repo, &signcmd.SignCommitArgs{
				Branch:             ref.String(),
				ForceCheckout:      false,
				AmendCommit:        plumbing2.IsMergeRequestReference(ref.String()) || plumbing2.IsIssueReference(ref.String()),
				Remote:             args.Args[0],
				NoPrompt:           true,
				ResetTokens:        false,
				RPCClient:          args.RPCClient,
				RemoteClients:      args.RemoteClients,
				KeyUnlocker:        args.KeyUnlocker,
				GetNextNonce:       args.GetNextNonce,
				SetRemotePushToken: args.SetRemotePushToken,
			}); err != nil {
				return err
			}
		}

		if ref.IsTag() {
			name := strings.Replace(ref.String(), "refs/tags/", "", 1)
			if err := args.TagSigner(cfg, []string{name}, repo, &signcmd.SignTagArgs{
				Remote:             args.Args[0],
				NoPrompt:           true,
				Force:              true,
				ResetTokens:        false,
				RPCClient:          args.RPCClient,
				RemoteClients:      args.RemoteClients,
				KeyUnlocker:        args.KeyUnlocker,
				GetNextNonce:       args.GetNextNonce,
				SetRemotePushToken: args.SetRemotePushToken,
			}); err != nil {
				return err
			}
		}

		if ref.IsNote() {
			if err := args.NoteSigner(cfg, repo, &signcmd.SignNoteArgs{
				Name:               strings.Replace(ref.String(), "refs/notes/", "", 1),
				Remote:             args.Args[0],
				NoPrompt:           true,
				ResetTokens:        false,
				RPCClient:          args.RPCClient,
				RemoteClients:      args.RemoteClients,
				KeyUnlocker:        args.KeyUnlocker,
				GetNextNonce:       args.GetNextNonce,
				SetRemotePushToken: args.SetRemotePushToken,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

// HandleAskPass handles git's askpass request. It collects the push token
// created by SignCmd and passes it to git as the password. It also resets
// fields
func HandleAskPass(stdout, stderr io.Writer, repo types.LocalRepo, args []string) error {

	// Output nothing for password request
	if strings.HasPrefix(args[0], "Password") {
		fmt.Fprint(stdout, "")
		return nil
	}

	// Get the remote currently being pushed by the push hook.
	// HookCmd will store it at `hook.curRemote`.
	// Remove the 'hook' section afterwards.
	rcfg, err := repo.Config()
	if err != nil {
		return errors.Wrap(err, "failed to get repo config")
	}
	curRemote := rcfg.Raw.Section("hook").Option("curRemote")
	rcfg.Raw.RemoveSection("hook")

	// Get the remote's push token
	tokens := rcfg.Raw.Section("remote").Subsection(curRemote).Option("tokens")
	if tokens == "" {
		fmt.Fprintln(stderr, colorfmt.RedString("Push token was not found for remote (%)", curRemote))
		return fmt.Errorf("push token was not found")
	}

	// Output push token as username
	fmt.Fprintf(stdout, tokens)

	// Clear the remote.*.tokens and sign.mergeID fields since these
	// fields' values are supposed to be for one-time use
	rcfg.Raw.Section("remote").Subsection(curRemote).RemoveOption("tokens")
	rcfg.Raw.Section("sign").RemoveOption("mergeID")
	return repo.SetConfig(rcfg)
}
