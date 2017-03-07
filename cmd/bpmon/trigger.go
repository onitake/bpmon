package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/spf13/cobra"
	"github.com/unprofession-al/bpmon"
)

var triggerCmd = &cobra.Command{
	Use:   "trigger",
	Short: "Run all business process checks and trigger temploted command on BP issues",
	Run: func(cmd *cobra.Command, args []string) {
		c, b, err := bpmon.Configure(cfgFile, cfgSection, bpPath, bpPattern)
		if err != nil {
			log.Fatal(err)
		}

		t := template.Must(template.New("t1").Parse(c.Trigger.Template))

		i, err := bpmon.NewIcinga(c.Icinga, c.Rules)
		if err != nil {
			log.Fatal(err)
		}
		stripBy := []bpmon.Status{bpmon.StatusUnknown, bpmon.StatusOK}
		var sets []bpmon.ResultSet
		for _, bp := range b {
			rs := bp.Status(i)
			set, stripped := rs.StripByStatus(stripBy)
			if !stripped {
				sets = append(sets, set)
			}
		}
		var command bytes.Buffer
		t.Execute(&command, sets)
		if len(sets) > 0 {
			fmt.Println(command.String())
		}
	},
}

func init() {
	RootCmd.AddCommand(triggerCmd)
}
