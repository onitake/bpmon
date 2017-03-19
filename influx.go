package bpmon

import (
	"fmt"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

type InfluxConf struct {
	Connection struct {
		Server string `yaml:"server"`
		Port   int    `yaml:"port"`
		Pass   string `yaml:"pass"`
		User   string `yaml:"user"`
		Proto  string `yaml:"proto"`
	} `yaml:"connection"`
	SaveOK   []string `yaml:"save_ok"`
	Database string   `yaml:"database"`
}

type Influx struct {
	cli      client.Client
	saveOK   []string
	database string
}

type Influxable interface {
	AsInflux([]string) []Point
}

func NewInflux(conf InfluxConf) (Influx, error) {
	addr := fmt.Sprintf("%s://%s:%d", conf.Connection.Proto, conf.Connection.Server, conf.Connection.Port)
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: conf.Connection.User,
		Password: conf.Connection.Pass,
	})
	cli := Influx{
		cli:      c,
		saveOK:   conf.SaveOK,
		database: conf.Database,
	}
	return cli, err
}

func (i Influx) Write(in Influxable) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  i.database,
		Precision: "s",
	})
	if err != nil {
		return err
	}

	points := in.AsInflux(i.saveOK)

	for _, p := range points {
		pt, _ := client.NewPoint(p.Series, p.Tags, p.Fields, p.Timestamp)
		bp.AddPoint(pt)
	}
	err = i.cli.Write(bp)

	return err
}

type Point struct {
	Timestamp time.Time
	Series    string
	Tags      map[string]string
	Fields    map[string]interface{}
}
