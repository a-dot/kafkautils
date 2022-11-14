package main

import (
	"fmt"
	"time"

	"github.com/a-dot/kafkautils/pkg/client"
	"github.com/alecthomas/kong"
	"github.com/segmentio/kafka-go"
)

var CLI struct {
	Brokers []string `help:"Broker to connect to" required:""`
	Timeout int      `help:"Kafka client timeout" default:"10"`

	Security string `help:"Security mode to connect to kafka (default: plaintext)" default:"plaintext" enum:"plaintext,ssl"`
	// JKS
	TrustStore    string `help:"JKS File for the truststore"`
	KeyStore      string `help:"JKS File for the keystore"`
	KeyStoreAlias string `help:"Alias for entry to use in JKS keystore"`
	// PEM
	CAPEM          string `help:"CA certificate for verifying the broker's key"`
	CertificatePEM string `help:"Client's public key used for authentication"`
	KeyPEM         string `help:"Client's private key used for authentication"`
	KeyPassword    string `help:"Private key passphrase"`

	List struct {
		Topics struct {
		} `cmd:""`

		Groups struct {
		} `cmd:""`
	} `cmd:""`

	Describe struct {
		Group struct {
			Group string `help:"Consumer group name" arg:""`
		} `cmd:""`
	} `cmd:""`
}

func main() {
	k := kong.Parse(&CLI)

	c := &kafka.Client{
		Addr:      kafka.TCP(CLI.Brokers...),
		Timeout:   time.Duration(CLI.Timeout) * time.Second,
		Transport: &kafka.Transport{},
	}

	fmt.Println(CLI)
	fmt.Println(k.Command())

	switch k.Command() {
	case "list topics":
		z := client.ListTopics(c)
		z.Print()
	case "list groups":
		z := client.ListGroups(c)
		z.Print()
	case "describe group <group>":
		z := client.DescribeGroup(c, CLI.Describe.Group.Group)
		z.Print()
	default:
		panic(k.Command())
	}

	// list consumer groups offsets for a topic
	// delete offsets for consumer group on a topic
}
