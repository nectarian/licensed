package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	rsautil "github.com/nectarian/licensed/pkg/utils/rsa"
)

var root = &cobra.Command{
	Use:   "keymaker",
	Short: "keymaker is a tool used for generate ras key.",
	Long:  "keymaker is a tool used for generate ras key.",
	Run:   Keymaker,
}

func init() {

	root.PersistentFlags().StringP("output", "o", ".", "output folder")
	root.PersistentFlags().StringP("key", "k", "rsa", "key file name prefix")
	root.PersistentFlags().IntP("bits", "b", 4096, "key length")
	root.PersistentFlags().StringP("metadata", "m", "", "metadata , text key-value pairs with json format")

	if err := viper.BindPFlag("output", root.PersistentFlags().Lookup("output")); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("key", root.PersistentFlags().Lookup("key")); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("bits", root.PersistentFlags().Lookup("bits")); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("metadata", root.PersistentFlags().Lookup("metadata")); err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Keymaker execute
func Keymaker(cmd *cobra.Command, args []string) {
	var err error
	var output = viper.GetString("output")
	var filename = viper.GetString("key")
	var bits = viper.GetInt("bits")
	var metadataPar = viper.GetString("metadata")

	var metadata map[string]string
	if metadataPar != "" {
		if err = json.Unmarshal([]byte(metadataPar), &metadata); err != nil {
			log.Fatal("unmarshal metadta", err)
		}
	}

	// Generate Private Key
	var privateKey *rsa.PrivateKey
	if privateKey, err = rsa.GenerateKey(rand.Reader, bits); err != nil {
		log.Fatal("generate private key", err)
	}

	// save private key
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyFilePath := filepath.Join(output, filename+"_private.pem")
	if err = rsautil.WritePEMBlock(privateKeyFilePath, "RSA Private Key", metadata, X509PrivateKey); err != nil {
		log.Fatal("write key file", err)
	}

	fmt.Printf("private key generated : %s\n", privateKeyFilePath)

	// save public key
	publicKey := privateKey.PublicKey
	var X509PublicKey = x509.MarshalPKCS1PublicKey(&publicKey)

	publicKeyFilePath := filepath.Join(output, filename+"_public.pem")
	if err = rsautil.WritePEMBlock(publicKeyFilePath, "RSA Public Key", metadata, X509PublicKey); err != nil {
		log.Fatal("write key file", err)
	}
	fmt.Printf(" public key generated : %s\n", publicKeyFilePath)
}
