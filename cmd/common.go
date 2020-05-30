package cmd

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	restclient "gitlab.com/makeos/mosdef/api/rest/client"
	"gitlab.com/makeos/mosdef/api/rpc/client"
	"gitlab.com/makeos/mosdef/remote/plumbing"
	rr "gitlab.com/makeos/mosdef/remote/repo"
	"gitlab.com/makeos/mosdef/remote/types"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
)

// getJSONRPCClient returns a JSON-RPCclient or error if unable to
// create one. It will return nil client and nil error if --no.rpc
// is true.
func getJSONRPCClient(cmd *cobra.Command) (*client.RPCClient, error) {
	noRPC, _ := cmd.Flags().GetBool("no.rpc")
	if noRPC {
		return nil, nil
	}

	rpcAddress, _ := cmd.Flags().GetString("rpc.address")
	rpcUser, _ := cmd.Flags().GetString("rpc.user")
	rpcPassword, _ := cmd.Flags().GetString("rpc.password")
	rpcSecured, _ := cmd.Flags().GetBool("rpc.https")

	host, port, err := net.SplitHostPort(rpcAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed parse rpc address")
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, errors.Wrap(err, "failed convert rpc port")
	}

	c := client.NewClient(&client.Options{
		Host:     host,
		Port:     portInt,
		User:     rpcUser,
		Password: rpcPassword,
		HTTPS:    rpcSecured,
	})

	return c, nil
}

// getRemoteAPIClients gets REST clients for every  http(s) remote
// URL set on the given repository. Immediately returns nothing if
// --no.remote is true.
func getRemoteAPIClients(cmd *cobra.Command, repo types.LocalRepo) (clients []restclient.RestClient) {
	noRemote, _ := cmd.Flags().GetBool("no.remote")
	if noRemote {
		return
	}

	for _, url := range repo.GetRemoteURLs() {
		ep, _ := transport.NewEndpoint(url)
		if !funk.ContainsString([]string{"http", "https"}, ep.Protocol) {
			continue
		}

		apiURL := fmt.Sprintf("%s://%s", ep.Protocol, ep.Host)
		if ep.Port != 0 {
			apiURL = fmt.Sprintf("%s:%d", apiURL, ep.Port)
		}

		clients = append(clients, restclient.NewClient(apiURL))
	}
	return
}

// getClients returns RPCClient and Remote API clients
func getRepoAndClients(cmd *cobra.Command) (types.LocalRepo,
	*client.RPCClient, []restclient.RestClient) {

	// Get the repository
	targetRepo, err := rr.GetAtWorkingDir(cfg.Node.GitBinPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Get JSON RPCClient client
	var rpcClient *client.RPCClient
	rpcClient, err = getJSONRPCClient(cmd)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Get remote APIs from the repository
	remoteClients := getRemoteAPIClients(cmd, targetRepo)

	return targetRepo, rpcClient, remoteClients
}

// rejectFlagCombo rejects unwanted flag combination
func rejectFlagCombo(cmd *cobra.Command, flags ...string) {
	var found = []string{}
	for _, f := range flags {
		if len(found) > 0 && cmd.Flags().Changed(f) {
			str := "--" + f
			if fShort := cmd.Flags().Lookup(f).Shorthand; fShort != "" {
				str += "|-" + fShort
			}
			found = append(found, str)
			log.Fatal(fmt.Sprintf("flags %s can't be used together", strings.Join(found, ", ")))
		}
		if cmd.Flags().Changed(f) {
			str := "--" + f
			if fShort := cmd.Flags().Lookup(f).Shorthand; fShort != "" {
				str += "|-" + fShort
			}
			found = append(found, str)
		}
	}
}

// requireFlag enforces flag requirement
func requireFlag(cmd *cobra.Command, flags ...string) {
	for _, f := range flags {
		if !cmd.Flags().Changed(f) {
			log.Fatal(fmt.Sprintf("flag (--%s) is required", f))
		}
	}
}

func getRef(curRepo types.LocalRepo, args []string) string {
	var ref string
	var err error

	if len(args) > 0 {
		ref = args[0]
	}

	if ref == "" {
		ref, err = curRepo.Head()
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to get HEAD").Error())
		}
	} else {
		ref = strings.ToLower(ref)
		if strings.HasPrefix(ref, plumbing.IssueBranchPrefix) {
			ref = fmt.Sprintf("refs/heads/%s", ref)
		}
		if !plumbing.IsMergeRequestReferencePath(ref) {
			ref = plumbing.MakeMergeRequestReference(ref)
		}
		if !plumbing.IsMergeRequestReference(ref) {
			log.Fatal(fmt.Sprintf("invalid issue path (%s)", ref))
		}
	}

	return ref
}
