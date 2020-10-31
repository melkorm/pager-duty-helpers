package cmd

import (
	"fmt"
	"time"

	"github.com/melkorm/pagerduty-helpers/pd"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cmdReport)
}

var cmdReport = &cobra.Command{
	Use:   "report",
	Short: "Display users oncall report",
	Run: func(cmd *cobra.Command, args []string) {
		pdToken1, _ := cmd.Flags().GetString("pdToken")
		pdClient := pd.NewClient(pdToken1, nil, "")
		countHours(pdClient)
	},
}

func countHours(pdClient *pd.Client) {
	location, _ := time.LoadLocation("EET")
	onCallStart := time.Date(2020, 9, 24, 0, 0, 0, 0, location)
	onCallEnd := time.Date(2020, 10, 22, 0, 0, 0, 0, location)

	summary := map[string]map[time.Time]time.Time{}
	for _, oncalls := range getAllOncalls(pdClient, &onCallStart, &onCallEnd) {
		for _, oncall := range oncalls.Items {
			if oncall.Schedule.Summary != "" && oncall.EscalationLevel == 1 && oncall.End != nil && oncall.Start != nil {
				if summary[oncall.User.Summary] == nil {
					summary[oncall.User.Summary] = map[time.Time]time.Time{}
				}

				start := time.Date(
					(*oncall.Start).Year(),
					(*oncall.Start).Month(),
					(*oncall.Start).Day(),
					(*oncall.Start).Hour(),
					0,
					0,
					0,
					location,
				)

				hours := (*oncall.End).Hour()
				if (*oncall.End).Minute() > 0 {
					hours++
				}
				end := time.Date(
					(*oncall.End).Year(),
					(*oncall.End).Month(),
					(*oncall.End).Day(),
					hours,
					0,
					0,
					0,
					location,
				)

				summary[oncall.User.Summary][start] = end
			}
		}
	}

	for username, timing := range summary {
		sum := map[string]float64{"weekend": 0.0, "normal": 0.0}
		for start, end := range timing {
			for i := 0; i < int(end.Sub(start).Hours()); i++ {
				current := start.Add(time.Duration(i * int(time.Hour)))
				if current.Before(onCallStart) || current.After(onCallEnd) {
					continue
				}
				switch current.Weekday() {
				case 0:
					sum["weekend"]++
				case 6:
					sum["weekend"]++
				default:
					if current.Hour() < 8 || current.Hour() >= 16 {
						sum["normal"]++
					}
				}
			}
		}
		fmt.Printf("%s: %d\n", username, int(sum["weekend"]+sum["normal"]))
	}
}
