package purge

import (
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

const (
	ctnCreated = "Created"
	ctnExited = "Exited"
)

var cmdCtn = &cobra.Command{
	Use: "container",
	Short: "Purge stopped containers",
	Long: "Purge stopped containers",
	RunE: RunCmdContainer,
}

type CtnFilter func(ctn docker.APIContainers) bool


type CtnValidator struct {
	Filters []CtnFilter
}

func (c *CtnValidator) Satisfied(ctn docker.APIContainers) bool {
	if len(c.Filters) == 0 {
		return false
	}
	for _, Func := range c.Filters {
		if !Func(ctn) {
			return false
		}
	}
	return true
}

func NewCtnValidator(filters ...Filter) (c *CtnValidator, err error) {
	c = new(CtnValidator)
	var filter CtnFilter
	for _, f := range filters {
		switch f.Field {
		case "created":
			filter, err = GenFilterCreated(f)
		case "exited":
			filter, err = GenFilterExited(f)
		default:
			continue
		}
		if err != nil {
			return
		}
		c.Filters = append(c.Filters, filter)
	}
	return
}

// GenFilterCreated creates a filter that filter containers with created timestamp
func GenFilterCreated(f Filter) (filter CtnFilter, err error) {
	ago, err := parseDuration(f.Value)
	if err != nil {
		return
	}
	cmp, ok := int64Comparator[f.Comparator]
	if !ok {
		tips := fmt.Sprintf("unsupported filter: %s, field: %s", f.Source, f.Comparator)
		return  nil, errors.New(tips)
	}
	filter = func(ctn docker.APIContainers) bool {
		return cmp(ago.Timestamp(), ctn.Created)
	}
	return
}


// GenFilterExited creates a filter that filter containers with exited time
func GenFilterExited(f Filter) (filter CtnFilter, err error) {
	ago, err := parseDuration(f.Value)
	if err != nil {
		return
	}
	cmp, ok := int64Comparator[f.Comparator]
	if !ok {
		tips := fmt.Sprintf("unsupported filter: %s, field: %s", f.Source, f.Comparator)
		return  nil, errors.New(tips)
	}
	filter = func(ctn docker.APIContainers) bool {
		status, err := parseContainerStatus(ctn.Status)
		if status.Status != ctnExited {
			return false
		}
		if err != nil {
			return false
		}
		exts, _ := status.ExitedTimestamp()
		return cmp(ago.Timestamp(), exts)
	}
	return
}

// CtnStatus stores container status
type CtnStatus struct {
	Status string
	Code int
	Num int
	Unit string
}

// ExitedTimestamp returns the exited timestamp of a container
func (s *CtnStatus) ExitedTimestamp() (int64, error) {
	if s.Status != ctnExited {
		return 0, errors.New("status is not Exited")
	}
	now := time.Now()
	ago := Ago{}
	hours := 0
	switch s.Unit {
	case "years":
		ago.Years = s.Num
	case "months":
		ago.Months = s.Num
	case "weeks":
		ago.Days = s.Num * 7
	case "days":
		ago.Days = s.Num
	case "hours":
		hours = s.Num
	}
	then := now.AddDate(-ago.Years, -ago.Months, -ago.Days)
	then = then.Add(-time.Duration(hours)*time.Hour)
	return then.Unix(), nil
}



func RunCmdContainer(cmd *cobra.Command, args []string) error {
	var filters []Filter
	for _, f := range filter {
		fmt.Println("Filter: ", f)
		parsed, err := parseFilter(f)
		if err != nil {
			return err
		}
		filters = append(filters, parsed)
	}
	return RemoveContainers(filters...)
}


func RemoveContainers(filters ...Filter) (err error) {
	var cli *docker.Client
	if dockerUri == "" {
		cli, err = docker.NewClientFromEnv()
	} else {
		cli, err = docker.NewClient(dockerUri)
	}
	if err != nil {
		return
	}
	validator, err := NewCtnValidator(filters...)
	if err != nil {
		return
	}

	containers, err := cli.ListContainers(docker.ListContainersOptions{All:true})
	if err != nil {
		return
	}
	for _, ctn := range containers {
		if validator.Satisfied(ctn) {
			e := cli.RemoveContainer(docker.RemoveContainerOptions{ID:ctn.ID})
			if e != nil {
				fmt.Printf("Can not remove container: %s, reason: %s", ctn.ID, e)
			}
			fmt.Println("removing container ", ctn.ID)
		}
	}
	return
}

// parseContainerStatus parses container status with statusPtn
func parseContainerStatus(s string) (a *CtnStatus, err error) {
	a = new(CtnStatus)
	m := statusPtn.FindStringSubmatch(s)
	for i, name := range statusPtn.SubexpNames() {
		if i != 0 && name != "" && m[i] != "" {
			switch name {
			case "status":
				a.Status = m[i]
			case "code":
				a.Code, err = strconv.Atoi(m[i])
			case "num":
				a.Num, err = strconv.Atoi(m[i])
			case "unit":
				a.Unit = m[i]
			}
			if err != nil {
				return
			}
		}
	}

	return a, nil
}
