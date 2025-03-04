package main

import (
	"Ada/api"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// Load the configuration using viper
	viper.SetConfigName("config")  // name of the config file (without extension)
	viper.SetConfigType("json")    // specify type. Use "yaml" if your config is in YAML format
	viper.AddConfigPath("configs") // optionally look for config in the working directory

	if err := viper.ReadInConfig(); err != nil {
		panic("Failed to read configuration file: " + err.Error())
	}

	router := gin.Default()

	api.RegisterRoutes(router)

	router.Run(":8080")
}
