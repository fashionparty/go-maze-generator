package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/spf13/viper"
	"log"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Unable to read configuration: %v", err)
	}

	engine := Engine{}
	engine.InitEngine()
	ebiten.SetWindowSize(engine.width*25, engine.height*25)
	err := ebiten.RunGame(&engine)
	if err != nil {
		log.Fatal(err.Error())
	}
}
