package config

import (
	"errors"
	"os"
	"reflect"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	cfgDefaultTag = "cfgDefault"
	cfgNameTag    = "json"
)

// LoadConfig [ConfigType] loading config with json tag as env var
// and cfgDefault tag as default value
func LoadConfig[T any]() (*T, error) {
	_ = godotenv.Load()

	config := new(T)
	val := reflect.ValueOf(config).Elem()

	err := loadConfig(val)
	if err != nil {
		return nil, err // TODO: wrap error
	}

	return config, nil
}

func loadConfig(val reflect.Value) error {
	if val.Kind() != reflect.Struct {
		return errors.New("only struct can be parsed")
	}

	typeof := val.Type()

	for i := 0; i < typeof.NumField(); i++ {
		field := typeof.Field(i)
		x := getEnv(field)

		switch field.Type.Kind() {
		case reflect.String:
			val.Field(i).SetString(x)

		case reflect.Int:
			x, err := strconv.ParseInt(x, 10, 64)
			if err != nil {
				return err // TODO: wrap error
			}

			val.Field(i).SetInt(x)

		case reflect.Bool:
			x, err := strconv.ParseBool(x)
			if err != nil {
				return err // TODO: wrap error
			}

			val.Field(i).SetBool(x)

		case reflect.Struct:
			err := loadConfig(val.Field(i))
			if err != nil {
				return err
			}

		default:
			// skip
		}
	}

	return nil
}

func getEnv(field reflect.StructField) string {
	key, ok := field.Tag.Lookup(cfgNameTag)
	if !ok {
		key = field.Name
	}

	value := os.Getenv(key)
	if value == "" {
		value, _ = field.Tag.Lookup(cfgDefaultTag)
	}

	return value
}
