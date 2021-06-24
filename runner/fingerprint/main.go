package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nectarian/licensed/pkg/plugin"
	"github.com/nectarian/licensed/pkg/plugin/factory"
	"github.com/nectarian/licensed/pkg/plugin/mac"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	rsautil "github.com/nectarian/licensed/pkg/utils/rsa"
)

var root = &cobra.Command{
	Use:   "fingerprint",
	Short: "fingerprint is a tool used for collect client unique identification information.",
	Long: `fingerprint is a tool used for collect client unique identification information.
identification information can be machine-id in linux or network card MAC address`,
	Run: Fingerprint,
}

func init() {
	root.PersistentFlags().StringP("output", "o", "clientID.pem", "output path")
	root.PersistentFlags().StringP("network", "n", "", "network card name")
	root.PersistentFlags().StringP("type", "t", "baseboard", `client information type, should be one of : machine-id, mac, text, baseboard
machine-id : machine-id , linux and windows only
mac : network card MAC address
baseboard : base board serialNumber
text : plaintext "license", just for opeartion which do not validate client
`)

	if err := viper.BindPFlag("output", root.PersistentFlags().Lookup("output")); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("network", root.PersistentFlags().Lookup("network")); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("type", root.PersistentFlags().Lookup("type")); err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Fingerprint execute
func Fingerprint(cmd *cobra.Command, args []string) {
	var err error
	var output = viper.GetString("output")
	var ciType = viper.GetString("type")
	var network = viper.GetString("network")

	var CreateClientFunc plugin.CreateClientIDFunc
	CreateClientFunc, err = factory.CreateClientIDFuncFactory(ciType)
	if err != nil {
		log.Fatalf("bad type , type must be one of : mac machine-id text")
	}

	opts := []plugin.Option{}
	if network != "" {
		opts = append(opts, mac.NetworkName(network))
	}

	var clientID plugin.ClientID = CreateClientFunc(opts...)
	if clientID.Load() != nil {
		log.Fatalf("load ClientID failed : %v", err)
	}
	var msg []byte
	if msg, err = clientID.Serialize(); err != nil {
		log.Fatalf("serialize client ID failed : %v", err)
	}
	if err = rsautil.WritePEMBlock(output, "Client ID", nil, msg); err != nil {
		log.Fatalf("write client ID failed : %v", err)
	}
	fmt.Printf("client ID generated : %s\n", output)
}
