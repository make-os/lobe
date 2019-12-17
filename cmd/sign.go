// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"net"
	"strconv"

	"github.com/makeos/mosdef/config"
	"github.com/makeos/mosdef/repo"
	"github.com/makeos/mosdef/rpc/client"
	"github.com/makeos/mosdef/types"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

func getRPCClient(cmd *cobra.Command) (*client.RPCClient, error) {
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

	c := client.NewClient(&types.Options{
		Host:     host,
		Port:     portInt,
		User:     rpcUser,
		Password: rpcPassword,
		HTTPS:    rpcSecured,
	})

	return c, nil
}

// signCmd represents the commit command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Add transaction information and sign git commits, tags and notes",
	Long:  `Add transaction information and sign git commits, tags and notes`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var signCommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Add transaction information and sign the last commit",
	Long:  `Add transaction information and sign the last commit`,
	Run: func(cmd *cobra.Command, args []string) {
		fee, _ := cmd.Flags().GetString("fee")
		nonce, _ := cmd.Flags().GetString("nonce")
		sk, _ := cmd.Flags().GetString("signingKey")

		var client *client.RPCClient
		var err error
		if nonce == "" || nonce == "0" {
			client, err = getRPCClient(cmd)
			if err != nil {
				log.Fatal(err.Error())
			}
		}

		if err := repo.SignCommitCmd(cfg.Node.GitBinPath, fee,
			nonce, sk, client); err != nil {
			cfg.G().Log.Fatal(err.Error())
		}
	},
}

var signTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Create an annotated tag with transaction information and sign it",
	Long:  `Create an annotated tags with transaction information and sign it`,
	Run: func(cmd *cobra.Command, args []string) {
		fee, _ := cmd.Flags().GetString("fee")
		nonce, _ := cmd.Flags().GetString("nonce")
		sk, _ := cmd.Flags().GetString("signingKey")

		var client *client.RPCClient
		var err error
		if nonce == "" || nonce == "0" {
			client, err = getRPCClient(cmd)
			if err != nil {
				log.Fatal(err.Error())
			}
		}

		args = cmd.Flags().Args()
		if err := repo.SignTagCmd(args, cfg.Node.GitBinPath, fee,
			nonce, sk, client); err != nil {
			cfg.G().Log.Fatal(err.Error())
		}
	},
}

var signNoteCmd = &cobra.Command{
	Use:   "note",
	Short: "Create a note, add transaction information and sign it",
	Long:  `Create a note, add transaction information and sign it`,
	Run: func(cmd *cobra.Command, args []string) {
		fee, _ := cmd.Flags().GetString("fee")
		nonce, _ := cmd.Flags().GetString("nonce")
		sk, _ := cmd.Flags().GetString("signingKey")

		var client *client.RPCClient
		var err error
		if nonce == "" || nonce == "0" {
			client, err = getRPCClient(cmd)
			if err != nil {
				log.Fatal(err.Error())
			}
		}

		if len(args) == 0 {
			log.Fatal("name is required")
		}

		if err := repo.SignNoteCmd(
			cfg.Node.GitBinPath,
			fee,
			nonce,
			sk,
			args[0],
			client); err != nil {
			log.Fatal(err.Error())
		}
	},
}

func initCommit() {
	rootCmd.AddCommand(signCmd)
	signCmd.AddCommand(signTagCmd)
	signCmd.AddCommand(signCommitCmd)
	signCmd.AddCommand(signNoteCmd)

	signCmd.PersistentFlags().StringP("fee", "f", "0", "Set the transaction fee")
	signCmd.PersistentFlags().StringP("nonce", "n", "0", "Set the transaction nonce")
	signCmd.PersistentFlags().StringP("signingKey", "s", "", "Set the GPG signing key ID")
	signCmd.PersistentFlags().String("rpc.user", "", "Set the RPC username")
	signCmd.PersistentFlags().String("rpc.password", "", "Set the RPC password")
	signCmd.PersistentFlags().String("rpc.address", config.DefaultRPCAddress,
		"Set the RPC listening address")
	signCmd.PersistentFlags().String("rpc.https", "", "Force the client to use https:// protocol")
}