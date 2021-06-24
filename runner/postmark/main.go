package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nectarian/licensed/pkg/plugin"
	"github.com/nectarian/licensed/pkg/plugin/byted"
	rsautil "github.com/nectarian/licensed/pkg/utils/rsa"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var root = &cobra.Command{
	Use:   "postmark",
	Short: "postmark is a tool used for generate license.",
	Long:  "postmark is a tool used for generate license.",
	Run:   Postmark,
}

func init() {
	root.PersistentFlags().StringP("output", "o", "license.key", "the output path of the license file")
	root.PersistentFlags().StringP("input", "i", "clientID.pem", "the path of input file which include client ID data")
	root.PersistentFlags().StringP("key", "k", "rsa_private.pem", "rsa private key file path")
	root.PersistentFlags().StringP("metadata", "m", "", "metadata json format")
	root.PersistentFlags().StringP("title", "t", "license", "license file title")
	root.PersistentFlags().StringP("expiration", "e", "", "expiration data , format : yyyy-MM-dd , this flag can be empty, it will not expire if it is empty")

	if err := viper.BindPFlag("output", root.PersistentFlags().Lookup("output")); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("input", root.PersistentFlags().Lookup("input")); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("key", root.PersistentFlags().Lookup("key")); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("metadata", root.PersistentFlags().Lookup("metadata")); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("title", root.PersistentFlags().Lookup("title")); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("expiration", root.PersistentFlags().Lookup("expiration")); err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Postmark generate license
func Postmark(cmd *cobra.Command, args []string) {
	var err error
	var CreateClientFunc = byted.CreateClientID
	var output = viper.GetString("output")
	var input = viper.GetString("input")
	var key = viper.GetString("key")
	var metadataPar = viper.GetString("metadata")
	var title = viper.GetString("title")
	var expirationPar = viper.GetString("expiration")

	var metadata map[string]string = make(map[string]string)
	if metadataPar != "" {
		if err = json.Unmarshal([]byte(metadataPar), &metadata); err != nil {
			log.Fatal("unmarshal metadta", err)
		}
	}

	var expiration int64
	if expirationPar != "" {
		var expirationTime time.Time
		if expirationTime, err = time.Parse("2006-01-02", expirationPar); err != nil {
			log.Fatal("bad expiration fomat", err)
		}
		// add one day , to the end of date
		expiration = expirationTime.UnixNano()
		metadata["expiration"] = expirationPar
	}

	// read private key
	var privateKey *rsa.PrivateKey
	if privateKey, err = rsautil.LoadPrivateKey(key); err != nil {
		log.Fatal(err)
	}

	// read client raw data
	var raw *pem.Block
	if raw, err = rsautil.ReadPEMBlock(input); err != nil {
		log.Fatal("ready raw pem failed : ", err)
	}

	// create ClientID and deserialize
	var clientID plugin.ClientID = CreateClientFunc()

	if err = clientID.Deserialize(raw.Bytes); err != nil {
		log.Fatal("deserialize ClientID failed : ", err)
	}

	var msg []byte

	if msg, err = plugin.Wrap(clientID, expiration); err != nil {
		log.Fatal(err)
	}

	// signate
	hash := sha512.New()
	if _, err = hash.Write(msg); err != nil {
		log.Fatal("sha512sum(data) : ", err)
	}
	sum := hash.Sum(nil)
	var sign []byte
	if sign, err = rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA512, sum); err != nil {
		log.Fatal("SignPKCS1v15 : ", err)
	}

	// write pem file
	if err = rsautil.WritePEMBlock(output, title, metadata, sign); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("license generated : %s\n", output)
}
