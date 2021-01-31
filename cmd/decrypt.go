package main

import (
	"os"

	"github.com/go-clix/cli"
	"github.com/xvzf/htw-crypto-project/pkg/crypt"
	"github.com/xvzf/htw-crypto-project/pkg/image"
)

func decryptCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "decrypt",
		Short: "Decrypt a file using the cipher",
		Args:  cli.ArgsExact(2),
	}

	key := cmd.Flags().StringP("key-file", "k", "", "Key File (Image) used for encryption")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		k, err := os.Open(*key)
		if err != nil {
			return err
		}
		defer k.Close()

		s, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer s.Close()

		t, err := os.Create(args[1])
		if err != nil {
			return err
		}
		defer t.Close()

		img, err := image.Read(k)
		if err != nil {
			return err
		}

		c, err := crypt.New(img)
		if err != nil {
			return err
		}

		// Read source
		enc, err := crypt.Read(s)
		if err != nil {
			return err
		}

		// Decrypt
		dec, err := c.Decrypt(enc)
		if err != nil {
			return err
		}

		// Write plaintext to file
		t.Write([]byte(dec))

		return nil
	}
	return cmd
}
