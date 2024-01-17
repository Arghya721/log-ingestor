package config

import (
	"github.com/spf13/viper"
)

// KafkaConfiguration hold the values required to connect with the kafka
type KafkaConfiguration struct {
	KafkaConnectionURL string
	KafkaTopic string
}

// KafkaConfig holds the kafka configurations after reading it from config file
var KafkaConfig KafkaConfiguration

func GetKafkaConfig() {
	
	// Load the kafka configuration in the kafka struct
	KafkaConfig.KafkaConnectionURL = viper.GetString("KAFKA_CONNECTION_URL")
	KafkaConfig.KafkaTopic = viper.GetString("KAFKA_TOPIC")

}