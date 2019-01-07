package purge

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strconv"
	"time"
)

var (
	Mismatched = errors.New("pattern mismatched")
)

var (
	dockerUri string

	// filter stores filter strings from CMD
	filter []string

	// dryRun with it set to true, only print operations without actually applying them
	dryRun bool

	// durationPtn is responsible for matching duration string from CMD
	durationPtn = regexp.MustCompile(`((?P<years>\d+?)y)?((?P<months>\d+?)m)?((?P<days>\d+?)d)?`)

	// statusPtn is responsible for matching docker resource, mainly container status
	statusPtn = regexp.MustCompile(`(?P<status>\w+) ?(\((?P<code>\d+)\))? ?(?P<num>\d+)? ?(?P<unit>\w+)? ?(ago)?`)

	// sizePtn matches human readable size. "500m", "2G", etc
	sizePtn = regexp.MustCompile(`(?P<amount>\d+)(?P<unit>[k|m|g|K|M|G])`)

	// filterPtn matches a whole filter string
	filterPtn = regexp.MustCompile(`(?P<field>\w+)(?P<op>=|!=|>|>=|<|<=)(?P<value>[^(=|\s)]+)`)
)

// rootCmd the entry of dkp
var rootCmd = &cobra.Command{
	Use: "dkp",
	Short: "purge resources",
	Long: "purge allows you to clean images, containers, swarm services with filters",
	// Run: Purge,
}


func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}


func Version(cmd *cobra.Command, args []string) {

}

// Ago stores duration from now
type Ago struct {
	Years, Months, Days int
}

// Timestamp convert ago to a specific timestamp
func (a *Ago) Timestamp() int64 {
	return time.Now().AddDate(-a.Years, -a.Months, -a.Days).Unix()
}


// parseDuration parses strings that are like "10m", "1d", "5y"
// see durationPtn for pattern
func parseDuration(d string) (a *Ago, err error) {
	a = new(Ago)
	m := durationPtn.FindStringSubmatch(d)
	if m != nil {
		return a, Mismatched
	}
	for i, name := range durationPtn.SubexpNames() {
		if i != 0 && name != "" && m[i] != "" {
			switch name {
			case "years":
				a.Years, err = strconv.Atoi(m[i])
			case "months":
				a.Months, err = strconv.Atoi(m[i])
			case "days":
				a.Days, err = strconv.Atoi(m[i])
			}
			if err != nil {
				return
			}
		}
	}

	return a, nil
}


// parseSize parses strings that are formed of "12m", "2G", etc
// see sizePtn for the pattern
func parseSize(s string) (size int64, err error) {
	var amount int64
	var unit string
	m := sizePtn.FindStringSubmatch(s)
	if m == nil {
		return 0, Mismatched
	}
	for i, name := range sizePtn.SubexpNames() {
		if i != 0 && name != "" && m[i] != "" {
			switch name {
			case "amount":
				amount, err = strconv.ParseInt(m[i], 10, 64)
			case "unit":
				unit = m[i]
			}
			if err != nil {
				return
			}
		}
	}
	return ByteSize(amount, unit)
}


func init() {
	rootCmd.AddCommand(cmdImg)
	rootCmd.AddCommand(cmdCtn)
	rootCmd.AddCommand(cmdSvc)
	rootCmd.PersistentFlags().StringSliceVarP(
		&filter, "filter", "f", nil, "filter conditions")
	rootCmd.PersistentFlags().StringVarP(
		&dockerUri,
		"docker",
		"d",
		"",
		"docker uri. if not provided, use the default")
	rootCmd.PersistentFlags().BoolVarP(
		&dryRun,
		"dry-run",
		"p",
		false,
		"Only prints actions but not actually apply them")
}
