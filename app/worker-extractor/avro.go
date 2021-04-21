package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/linkedin/goavro"
)

func newCodecFromSchema(schemaPath string) (*goavro.Codec, error) {
	schemaBytes, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read schema file %s: %v", schemaPath, err)
	}

	codec, err := goavro.NewCodec(string(schemaBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create new codec %s: %v", schemaPath, err)
	}

	return codec, nil
}

type Message interface{}

func avroEncode(m Message, codec *goavro.Codec) ([]byte, error) {
	txt, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	ntv, _, err := codec.NativeFromTextual(txt)
	if err != nil {
		return nil, err
	}

	bin, err := codec.BinaryFromNative(nil, ntv)
	if err != nil {
		return nil, err
	}

	return bin, nil
}
