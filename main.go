package main

import (
	"time"

	logging "github.com/op/go-logging"

	"encoding/json"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
	"io/ioutil"
	"os"
)

type MinerConfig struct {
	Name string // NAME
	Pin  string // PIN-NUMBER OF GPIO
	Ip   string // IP ADDRESS
	Info string // ADDITIONAL INFO
}

var log = logging.MustGetLogger("auto-hard-reset-log")

func main() {
	log.Notice("Reading file config.json...")
	configFileContent, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Error("Trying to read file config.json, but:", err)
		os.Exit(1)
	}

	log.Notice("Parsing configuration file...")
	var minerConfigs []MinerConfig
	err = json.Unmarshal(configFileContent, &minerConfigs)
	if err != nil {
		log.Error("Parsing JSON content, but:", err)
		os.Exit(2)
	}

	totalMinerConfigs := len(minerConfigs)
	log.Notice("Found miner configurations:", totalMinerConfigs)
	
	r := raspi.NewAdaptor()
	///MINING RIGS CONFIGURATION///
	miningRigs := make([]Rig, 0)
	for _, m := range minerConfigs {
		log.Notice("minerConfig:", m)
		miningRigs = append(miningRigs, Rig{m.Name, gpio.NewRelayDriver(r, m.Pin), m.Ip, m.Info})
	}
	log.Notice("Configured rigs: ", len(miningRigs))

	LogMachines()

	work := func() {
		timer := 33 * time.Minute
		log.Notice("HELLO! I WILL KEEP YOUR MONEY MAKING MACHINES ONLINE!")
		log.Notice("Starting timer: ", timer)

		//Check the machines every 33 minutes
		gobot.Every(timer, func() {
			log.Notice("Checking machines: ")
			for i := 0; i < len(miningRigs); i++ {
				log.Notice("Ping miner: ", i, "name: ", miningRigs[i].name, "ip: ", miningRigs[i].ip)
				if !miningRigs[i].Ping() {
					miningRigs[i].Restarter()
				}
			}

			log.Notice("Checking machines DONE")
			log.Notice("Restarting timer")
		})
	}

	robot := gobot.NewRobot("RPiMinerHardReset",
		[]gobot.Connection{r},
		[]gobot.Device{miningRigs[0].pin},
		[]gobot.Device{miningRigs[1].pin},
		[]gobot.Device{miningRigs[2].pin},
		[]gobot.Device{miningRigs[3].pin},
		[]gobot.Device{miningRigs[4].pin},
		[]gobot.Device{miningRigs[5].pin},
		[]gobot.Device{miningRigs[6].pin},
		[]gobot.Device{miningRigs[7].pin},
		[]gobot.Device{miningRigs[8].pin},
		[]gobot.Device{miningRigs[9].pin},
		[]gobot.Device{miningRigs[10].pin},
		[]gobot.Device{miningRigs[11].pin},
		//IF YOU ADD MORE MACHINES ADD INFO HERE
		//[]gobot.Device{miningRigs[12].pin},
		work,
	)

	robot.Start()
}
