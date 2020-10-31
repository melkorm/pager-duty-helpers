package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/melkorm/pagerduty-helpers/pd"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cmdOncall)
}

var cmdOncall = &cobra.Command{
	Use:   "oncall",
	Short: "Display current oncall",
	Run: func(cmd *cobra.Command, args []string) {
		pdToken1, _ := cmd.Flags().GetString("pdToken")
		pdClient := pd.NewClient(pdToken1, nil, "")
		fmt.Println(getCurrentOnCall(pdClient).String())
	},
}

func getCurrentOnCall(pdClient *pd.Client) *strings.Builder {
	var sb strings.Builder
	schedules := map[string]bool{}
	for _, oncalls := range getAllOncalls(pdClient, nil, nil) {
		for _, oncall := range oncalls.Items {
			if ok, _ := schedules[oncall.Schedule.ID]; ok {
				continue
			}
			if oncall.Schedule.Summary != "" && oncall.EscalationLevel == 1 && oncall.End != nil && oncall.Start != nil {
				schedules[oncall.Schedule.ID] = true
				sb.WriteString(
					fmt.Sprintf(
						"%s - @%s\n",
						oncall.Schedule.Summary,
						oncall.User.Summary,
					),
				)
			}
		}
	}

	return &sb
}

func getAllOncalls(pdClient *pd.Client, from, to *time.Time) []*pd.Oncalls {
	alloncalls := []*pd.Oncalls{}
	getOnCalls := func(offset int) (*pd.Oncalls, error) {
		oncalls, err := pdClient.GetOncalls(offset, from, to)

		if err != nil {
			return nil, err
		}

		return oncalls, nil
	}

	next := true
	curOffset := 0
	for next {
		oncalls, _ := getOnCalls(curOffset)
		alloncalls = append(alloncalls, oncalls)
		next = oncalls.More
		curOffset = oncalls.Offset
		curOffset = curOffset + 20
	}

	return alloncalls
}
