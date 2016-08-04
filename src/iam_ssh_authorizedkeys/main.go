package main

import (
	"os"
	"fmt"
	
	
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

var (
	version = "0.0.1"
)

const (
	usage = `
Usage:
  $ iam_ssh_authorizedkeys <username>
  
Options:
  -h, --help		Show usage
  --version			Show version
`
)

func main() {
	args := os.Args	
	if len(args) != 2 || args[1] == "-h" || args[1] == "--help" {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	if args[1] == "--version" {
		fmt.Printf("iam_ssh_authorizedkeys version %s\n", version)
		os.Exit(0)
	}

	user := args[1]	
	if len(user) > 32 {
		fmt.Fprintln(os.Stderr, "Username too long")
		os.Exit(1)
	}
	
	svc := iam.New(session.New(), &aws.Config{})
	
	resp, err := svc.ListSSHPublicKeys(&iam.ListSSHPublicKeysInput{ UserName: aws.String(user), })
	if err != nil {
		panic(err)
	}

	for _, pk := range resp.SSHPublicKeys {
		if *pk.Status == "Active" {
			pkresp, pkerr := svc.GetSSHPublicKey(&iam.GetSSHPublicKeyInput{
				Encoding:       aws.String("SSH"),
				SSHPublicKeyId: pk.SSHPublicKeyId,
				UserName:       aws.String(user),
			})
			if pkerr != nil {
				panic(pkerr)
			}
			
			fmt.Println(*pkresp.SSHPublicKey.SSHPublicKeyBody)
		}
	}
	
	os.Exit(0)
}

