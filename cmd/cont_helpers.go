package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/influxdata/influxdb/client/v2"

	"gopkg.in/yaml.v2"
)

var (
	i icinga
)

type conf struct {
	Icinga icingaConf
}

type bps []bp

func feed() {
	c, err := readConf()
	if err != nil {
		log.Fatal(err)
	}

	b, err := readBPs()
	if err != nil {
		log.Fatal(err)
	}

	i = newIcinga(c.Icinga)
	infl, _ := NewInflux(client.HTTPConfig{
		Addr: "http://***REMOVED***:8086",
	})
	for _, bp := range b {
		rs := bp.Status()
		err = infl.writeResultSet(rs)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func run() {
	c, err := readConf()
	if err != nil {
		log.Fatal(err)
	}

	b, err := readBPs()
	if err != nil {
		log.Fatal(err)
	}

	i = newIcinga(c.Icinga)
	//infl, _ := NewInflux(client.HTTPConfig{
	//	Addr: "http://***REMOVED***:8086",
	//})
	for _, bp := range b {
		rs := bp.Status()
		fmt.Println(rs.PrettyPrint(0))
		//err = infl.writeResultSet(rs)
		//if err != nil {
		//	fmt.Println(err)
		//}
	}
}

func readConf() (conf, error) {
	conf := conf{}
	file, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error while reading %s: %s", cfgFile, err.Error()))
		return conf, err
	}

	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error while parsing %s: %s", cfgFile, err.Error()))
		return conf, err
	}

	return conf, nil
}

func readBPs() (bps, error) {
	bps := bps{}
	files, err := ioutil.ReadDir(bpPath)
	if err != nil {
		return bps, err
	}

	for _, f := range files {
		match, err := filepath.Match(bpPattern, f.Name())
		if err != nil {
			return bps, err
		}
		if !match {
			continue
		}
		bp := bp{}
		file, err := ioutil.ReadFile(bpPath + "/" + f.Name())
		if err != nil {
			err = errors.New(fmt.Sprintf("Error while reading %s/%s: %s", bpPath, f.Name(), err.Error()))
			return bps, err
		}

		err = yaml.Unmarshal(file, &bp)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error while parsing %s/%s: %s", bpPath, f.Name(), err.Error()))
			return bps, err
		}

		bps = append(bps, bp)
	}
	return bps, nil
}
