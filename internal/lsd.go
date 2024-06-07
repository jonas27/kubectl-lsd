package internal

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ghodss/yaml"
)

type secret map[string]interface{}

type decodedSecret struct {
	Key   string
	Value string
}

func Run(stdin []byte) (string, error) {
	isjson := isJSON(stdin)
	if !isjson {
		injsonbytes, err := yaml.YAMLToJSON(stdin)
		if err != nil {
			return "", fmt.Errorf("error converting from yaml to json : %w", err)
		}
		stdin = injsonbytes
	}

	var orgSecret secret
	if err := unmarshal(stdin, &orgSecret, isjson); err != nil {
		return "", err
	}

	// Check if the object is a list.
	var err error
	if orgSecret["items"] != nil {
		if orgSecret, err = listSecrets(orgSecret); err != nil {
			return "", fmt.Errorf("could not convert list secrets: %w", err)
		}
	} else {
		if orgSecret, err = stringData(orgSecret); err != nil {
			return "", fmt.Errorf("could not convert secret to stringData: %w", err)
		}
	}

	var bs []byte
	if isjson {
		if bs, err = marshal(orgSecret); err != nil {
			return "", fmt.Errorf("can not marshal secret to JSON: %w", err)
		}
	} else {
		if bs, err = marshalYAML(orgSecret); err != nil {
			return "", fmt.Errorf("can not marshal secret to YAML: %w", err)
		}
	}

	return string(bs), nil
}

func listSecrets(orgSecret secret) (secret, error) {
	items, ok := orgSecret["items"].([]interface{})
	if !ok {
		return nil, errors.New("could not convert json object with key items to list")
	}
	for k, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return nil, errors.New("could not convert items to map")
		}
		plainSecret, err := stringData(itemMap)
		if err != nil {
			return nil, fmt.Errorf("could not convert item to stringData: %w", err)
		}

		items[k] = plainSecret
	}
	return orgSecret, nil
}

func isJSON(s []byte) bool {
	return json.Unmarshal(s, &json.RawMessage{}) == nil
}

func cast(data interface{}) (map[string]interface{}, bool) {
	d, ok := data.(map[string]interface{})
	return d, ok
}

func stringData(s secret) (secret, error) {
	data, ok := cast(s["data"])
	if !ok || len(data) == 0 {
		return s, nil
	}
	var err error
	s["stringData"], err = decode(data)
	if err != nil {
		return nil, fmt.Errorf("could not decode data: %w", err)
	}
	delete(s, "data")
	return s, nil
}

func unmarshal(in []byte, out interface{}, asJSON bool) error {
	if asJSON {
		return json.Unmarshal(in, out)
	}
	return yaml.Unmarshal(in, out)
}

func marshal(d interface{}) ([]byte, error) {
	return json.MarshalIndent(d, "", "    ")
}

func marshalYAML(d interface{}) ([]byte, error) {
	return yaml.Marshal(d)
}

func decodeSecret(key, secret string, secrets chan decodedSecret) {
	var value string
	// avoid wrong encoded secrets
	if decoded, err := base64.StdEncoding.DecodeString(secret); err == nil {
		value = string(decoded)
	} else {
		value = secret
	}
	secrets <- decodedSecret{Key: key, Value: value}
}

func decode(data map[string]interface{}) (map[string]string, error) {
	length := len(data)
	secrets := make(chan decodedSecret, length)
	decoded := make(map[string]string, length)
	for key, encoded := range data {
		encodedString, ok := encoded.(string)
		if !ok {
			return nil, fmt.Errorf("could not convert %v to string", key)
		}
		go decodeSecret(key, encodedString, secrets)
	}
	for range length {
		secret := <-secrets
		decoded[secret.Key] = secret.Value
	}
	return decoded, nil
}
