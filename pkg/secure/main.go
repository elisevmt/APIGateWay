package secure

import (
	"reflect"
)

type Cypher interface {
	Decrypt(bytes []byte) ([]byte, error)
	Encrypt(bytes []byte) ([]byte, error)
}

func DecryptStruct(obj interface{}, cypher Cypher) error {
	return decryptStruct(reflect.Indirect(reflect.ValueOf(obj)), cypher)
}

func decryptStruct(obj reflect.Value, cypher Cypher) error {
	if obj.Kind() == reflect.Slice || obj.Kind() == reflect.Array {
		for i := 0; i < obj.Len(); i++ {
			err := decryptStruct(obj.Index(i), cypher)
			if err != nil {
				return err
			}
		}
		return nil
	}
	for i := 0; i < obj.Type().NumField(); i++ {
		if val, ok := obj.Type().Field(i).Tag.Lookup("encrypted"); val == "true" && ok {
			valueField := obj.Field(i)
			if valueField.Kind() == reflect.Ptr {
				valueField = valueField.Elem()
			}
			data, err := cypher.Decrypt([]byte(valueField.String()))
			if err != nil {
				return err
			}
			valueField.SetString(string(data))
		}
		if obj.Field(i).Kind() == reflect.Struct {
			err := decryptStruct(obj.Field(i), cypher)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func EncryptStruct(obj interface{}, cypher Cypher) error {
	return encryptStruct(reflect.Indirect(reflect.ValueOf(obj)), cypher)
}

func encryptStruct(obj reflect.Value, cypher Cypher) error {
	if obj.Kind() == reflect.Slice || obj.Kind() == reflect.Array {
		for i := 0; i < obj.Len(); i++ {
			err := encryptStruct(obj.Index(i), cypher)
			if err != nil {
				return err
			}
		}
		return nil
	}
	for i := 0; i < obj.Type().NumField(); i++ {
		if val, ok := obj.Type().Field(i).Tag.Lookup("encrypted"); val == "true" && ok {
			data, err := cypher.Encrypt([]byte(obj.Field(i).String()))
			if err != nil {
				return err
			}
			obj.Field(i).SetString(string(data))
		}
		if obj.Field(i).Kind() == reflect.Struct {
			err := encryptStruct(obj.Field(i), cypher)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
