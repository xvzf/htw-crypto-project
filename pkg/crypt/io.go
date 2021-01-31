package crypt

import (
	"encoding/json"
	"os"
)

func Read(f *os.File) ([]PixelPosition, error) {

	var out []PixelPosition

	dec := json.NewDecoder(f)

	if err := dec.Decode(&out); err != nil {
		return nil, err
	}

	return out, nil
}

func Write(f *os.File, in []PixelPosition) error {

	enc := json.NewEncoder(f)

	if err := enc.Encode(in); err != nil {
		return err
	}

	return nil
}
