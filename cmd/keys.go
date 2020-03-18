// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	path "path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/makeos/mosdef/config"
	"gitlab.com/makeos/mosdef/keystore"
	"gitlab.com/makeos/mosdef/types/core"
)

// keysCmd represents the keystore command
var keysCmd = &cobra.Command{
	Use:   "key command [flags]",
	Short: "Create and manage your account and push keys.",
	Long: `Description:
This command provides the ability to create, list, import and update 
keys. Keys are stored in an encrypted format using a passphrase. 
Please understand that if you forget the password, it is IMPOSSIBLE to 
unlock your key. 

During creation, if a passphrase is not provided, the key is still encrypted using
a default (unsafe) passphrase and marked as 'unsafe'. You can change the passphrase 
at any time. (not recommended)

Keys are stored under <DATADIR>/` + config.KeystoreDirName + `. It is safe to transfer the 
directory or individual accounts to another node. 

Always backup your keeps regularly.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// keyCreateCmd represents the keystore command
var keyCreateCmd = &cobra.Command{
	Use:   "create [flags]",
	Short: "Create an key.",
	Long: `This command creates a key and encrypts it using a passphrase
you provide. Do not forget your passphrase, you will not be able 
to unlock your key if you do.

Password will be stored under <DATADIR>/` + config.KeystoreDirName + `. 
It is safe to transfer the directory or individual accounts to another node. 

Use --pass to directly specify a password without going interactive mode. You 
can also provide a path to a file containing a password. If a path is provided,
password is fetched with leading and trailing newline character removed. 

Always backup your keeps regularly.`,
	Run: func(cmd *cobra.Command, args []string) {
		seed, _ := cmd.Flags().GetInt64("seed")
		pass, _ := cmd.Flags().GetString("pass")
		nopass, _ := cmd.Flags().GetBool("nopass")
		pushType, _ := cmd.Flags().GetBool("push")

		ks := keystore.New(path.Join(cfg.DataDir(), config.KeystoreDirName))
		kt := core.KeyTypeAccount
		if pushType {
			kt = core.KeyTypePush
		}
		_, err := ks.CreateCmd(kt, seed, pass, nopass)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

var keyListCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: "List all accounts.",
	Long: `Description:
This command lists all accounts existing under <DATADIR>/` + config.KeystoreDirName + `.

Given that keys in the keystore directory are prefixed with their creation timestamp, the 
list is lexicographically sorted such that the oldest keystore will be at the top on the list
`,
	Run: func(cmd *cobra.Command, args []string) {
		ks := keystore.New(path.Join(cfg.DataDir(), config.KeystoreDirName))
		if err := ks.ListCmd(); err != nil {
			log.Fatal(err.Error())
		}
	},
}

var keyUpdateCmd = &cobra.Command{
	Use:   "update [flags] <address>",
	Short: "Update a key",
	Long: `Description:
This command allows you to update the password of a key and to
convert a key encrypted in an old format to a new one.
`,
	Run: func(cmd *cobra.Command, args []string) {

		var address string
		if len(args) >= 1 {
			address = args[0]
		}

		pass, _ := cmd.Flags().GetString("pass")

		ks := keystore.New(path.Join(cfg.DataDir(), config.KeystoreDirName))
		if err := ks.UpdateCmd(address, pass); err != nil {
			log.Fatal(err.Error())
		}
	},
}

var keyImportCmd = &cobra.Command{
	Use:   "import [flags] <keyfile>",
	Short: "Import an existing, unencrypted private key.",
	Long: `Description:
This command allows you to create a new key by importing a private key from a <keyfile>. 
You will be prompted to provide your password. Your key is saved in an encrypted format.

The keyfile is expected to contain an unencrypted private key in Base58 format.

You can skip the interactive mode by providing your password via the '--pass' flag. 
Also, a path to a file containing a password can be provided to the flag.

You must not forget your password, otherwise you will not be able to unlock your
key.
`,
	Run: func(cmd *cobra.Command, args []string) {

		var keyFile string
		if len(args) >= 1 {
			keyFile = args[0]
		}

		pass, _ := cmd.Flags().GetString("pass")
		pushType, _ := cmd.Flags().GetBool("push")
		kt := core.KeyTypeAccount
		if pushType {
			kt = core.KeyTypePush
		}

		ks := keystore.New(path.Join(cfg.DataDir(), config.KeystoreDirName))
		if err := ks.ImportCmd(keyFile, kt, pass); err != nil {
			log.Fatal(err.Error())
		}
	},
}

var keyRevealCmd = &cobra.Command{
	Use:   "reveal [flags] <address>",
	Short: "Reveal the private key of a key.",
	Long: `Description:
This command reveals the private key of a key. You will be prompted to 
provide your password. 
	
You can skip the interactive mode by providing your password via the '--pass' flag. 
Also, the flag accepts a path to a file containing a password.
`,
	Run: func(cmd *cobra.Command, args []string) {

		var address string
		if len(args) >= 1 {
			address = args[0]
		}

		_ = viper.BindPFlag("node.passphrase", cmd.Flags().Lookup("pass"))
		pass := viper.GetString("node.passphrase")

		ks := keystore.New(path.Join(cfg.DataDir(), config.KeystoreDirName))
		if err := ks.RevealCmd(address, pass); err != nil {
			log.Fatal(err.Error())
		}
	},
}

func setKeyCmdAndFlags() {
	keysCmd.AddCommand(keyCreateCmd)
	keysCmd.AddCommand(keyListCmd)
	keysCmd.AddCommand(keyUpdateCmd)
	keysCmd.AddCommand(keyImportCmd)
	keysCmd.AddCommand(keyRevealCmd)
	keysCmd.PersistentFlags().String("pass", "", "Password to unlock the target key and skip interactive mode")
	keyCreateCmd.Flags().Int64P("seed", "s", 0, "Provide a strong seed (not recommended)")
	keyCreateCmd.Flags().Bool("nopass", false, "Force key to be created with no passphrase")
	keyCreateCmd.Flags().Bool("push", false, "Create a push key")
	keyImportCmd.Flags().Bool("push", false, "Create a push key")
	rootCmd.AddCommand(keysCmd)
}
