package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/labels"
)

// PREFIX is config prefix
const PREFIX = "kuberule"

var (
	AppName         string
	PodNamespace    string
	Namespace       string
	CertDir         string
	SecretName      string
	ServiceName     string
	ServiceSelector labels.Set
)

func init() {
	Initialize()
}

// Initialize set configurations from
func Initialize() {

	viper.SetEnvPrefix(PREFIX)
	viper.SetConfigName(PREFIX)
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("%s", err)
	}

	viper.SetDefault("app.name", PREFIX)
	AppName = viper.GetString("app.name")

	PodNamespace, ok := os.LookupEnv("POD_NAMESPACE")
	if !ok {
		PodNamespace = "default"
	}

	viper.SetDefault("namespace", PodNamespace)
	Namespace = viper.GetString("namespace")

	viper.SetDefault("cert.dir", "/tmp/cert")
	CertDir = viper.GetString("cert.dir")

	viper.SetDefault("service.name", AppName)
	ServiceName = viper.GetString("service.name")

	viper.SetDefault("secret.name", AppName)
	SecretName = viper.GetString("secret.name")

	viper.SetDefault("service.selector", "app="+AppName)
	selectorString := viper.GetString("service.selector")

	if selector, err := labels.ConvertSelectorToLabelsMap(selectorString); err != nil {
		panic(fmt.Errorf("service.selector=\"%s\"\n%s", selectorString, err))
	} else {
		ServiceSelector = selector
	}
}

// Debug returns the entire config map
func Debug() map[string]interface{} {
	return viper.AllSettings()
}

// JSON returns JSON representation of the config
func JSON() string {
	if result, err := json.Marshal(Debug()); err == nil {
		return string(result)
	}

	return "{}"
}
