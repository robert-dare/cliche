package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/dc0d/argify"
	"github.com/hashicorp/hcl"
	"github.com/urfave/cli"
)

// build flags
var (
	BuildTime  string
	CommitHash string
	GoVersion  string
	GitTag     string
)

var conf struct {
	Info string `envvar:"APP_INFO" usage:"sample app info" value:"bare app structure"`

	Sample struct {
		SubCommand struct {
			Param string `envvar:"-"`
		}
	}
}

func defaultAppNameHandler() string {
	return filepath.Base(os.Args[0])
}

func defaultConfNameHandler() string {
	fp := fmt.Sprintf("%s.conf", defaultAppNameHandler())
	if _, err := os.Stat(fp); err != nil {
		fp = "app.conf"
	}
	return fp
}

func loadHCL(ptr interface{}, filePath ...string) error {
	var fp string
	if len(filePath) > 0 {
		fp = filePath[0]
	}
	if fp == "" {
		fp = defaultConfNameHandler()
	}
	cn, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}
	err = hcl.Unmarshal(cn, ptr)
	if err != nil {
		return err
	}

	return nil
}

func app() {
	if err := loadHCL(&conf); err != nil {
		log.Println("warn:", err)
		return
	}

	app := cli.NewApp()

	{
		app.Version = "0.0.1"
		app.Author = "__author__"
		app.Copyright = "__copyright__"
		now := time.Now()
		app.Description = fmt.Sprintf(
			"Build Time:  %v %v\n   Go:          %v\n   Commit Hash: %v\n   Git Tag:     %v",
			now.Weekday(),
			BuildTime,
			GoVersion,
			CommitHash,
			GitTag)
		app.Name = "__appname__"
		app.Usage = ""
	}

	{
		app.Action = cmdApp

		c := cli.Command{
			Name: `sample`,
		}
		c.Subcommands = append(c.Subcommands, cli.Command{
			Name:   "subcommand",
			Action: cmdSampleSubCommand,
		})
		app.Commands = append(app.Commands, c)
	}

	argify.NewArgify().Build(app, &conf)

	if err := app.Run(os.Args); err != nil {
		log.Fatalln("error:", err)
	}
}

func cmdApp(*cli.Context) error {
	defer finit(time.Second, true)
	fmt.Println(conf.Info, "ʕ⚆ϖ⚆ʔ")
	return nil
}

func cmdSampleSubCommand(*cli.Context) error {
	defer finit(time.Second, true)
	fmt.Println(conf.Sample.SubCommand.Param)
	return nil
}
