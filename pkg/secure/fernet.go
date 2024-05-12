package secure

import (
	"APIGateWay/constants"
	"fmt"
	"github.com/fernet/fernet-go"
)

type fernetCypher struct {
	key *fernet.Key
}

func (f fernetCypher) Decrypt(bytes []byte) ([]byte, error) {
	bytes = fernet.VerifyAndDecrypt(bytes, 0, []*fernet.Key{f.key})
	if len(bytes) == 0 {
		return nil, fmt.Errorf("cannot decrypt fernet: %w", constants.ErrConfig)
	}
	return bytes, nil
}

func (f fernetCypher) Encrypt(bytes []byte) ([]byte, error) {
	return fernet.EncryptAndSign(bytes, f.key)
}

func NewFernetCypher(key string) (Cypher, error) {
	fkey, err := fernet.DecodeKey(key)
	if err != nil {
		return nil, err
	}
	return &fernetCypher{
		key: fkey,
	}, nil
}
