package purge

import (
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
	"strings"
)

var cmdImg = &cobra.Command{
	Use: "image",
	Short: "Clean images",
	Long: "Clean images",
	RunE: RunCmdImage,
}


type ImgFilter func(img docker.APIImages) bool

type ImageValidator struct {
	Validators []ImgFilter
}

// Satisfied checks if an image can pass all filters of the validator
func (i *ImageValidator) Satisfied(img docker.APIImages) bool {
	if len(i.Validators) == 0 {
		return false
	}
	for _, Func := range i.Validators {
		if !Func(img) {
			return false
		}
	}
	return true
}


func NewImageValidator(filters ...Filter) (iv *ImageValidator, err error) {
	iv = new(ImageValidator)
	var filter ImgFilter
	for _, f := range filters {
		switch f.Field {
		case "created":
			filter, err = ImgCreatedFilter(f)
		case "name":
			filter, err = ImgNameFilter(f)
		case "tag":
			filter, err = ImgTagFilter(f)
		case "size":
			filter, err = ImgSizeFilter(f)
		}
		if err != nil {
			return
		}
		iv.Validators = append(iv.Validators, filter)
	}
	return
}

// ImgCreatedFilter creates a filter that filters image with created timestamp
func ImgCreatedFilter(f Filter) (filter ImgFilter, err error) {
	ago, err := parseDuration(f.Value)
	if err != nil {
		return
	}
	cmp, ok := int64Comparator[f.Comparator]
	if !ok {
		tips := fmt.Sprintf("unsupported filter: %s, field: %s", f.Source, f.Comparator)
		return  nil, errors.New(tips)
	}
	filter = func(img docker.APIImages) bool {
		return cmp(ago.Timestamp(), img.Created)
	}
	return
}

// ImgNameFilter creates a filter that filters image with name
func ImgNameFilter(f Filter) (filter ImgFilter, err error) {
	op, ok := stringComparator[f.Comparator]
	if !ok {
		tips := fmt.Sprintf("unsupported filter: %s, field: %s", f.Source, f.Comparator)
		return  nil, errors.New(tips)
	}
	filter = func(img docker.APIImages) bool {
		// tags form: repo/name:tag
		for _, tags := range img.RepoTags {
			name := strings.Split(tags, ":")[0]
			if op(name, f.Value) {
				return true
			}
		}
		return false
	}
	return
}

// ImgTagFilter creates a filter that filters image with tag
func ImgTagFilter(f Filter) (filter ImgFilter, err error) {
	op, ok := stringComparator[f.Comparator]
	if !ok {
		tips := fmt.Sprintf("unsupported filter: %s, field: %s", f.Source, f.Comparator)
		return  nil, errors.New(tips)
	}
	filter = func(img docker.APIImages) bool {
		// tags form: repo/name:tag
		for _, tags := range img.RepoTags {
			parts := strings.Split(tags, ":")
			if len(parts) < 2 {
				return false
			}

			if op(parts[1], f.Value) {
				return true
			}
		}
		return false
	}
	return
}

// ImgSizeFilter creates a filter that filters image with size
func ImgSizeFilter(f Filter) (filter ImgFilter, err error) {
	size, err := parseSize(f.Value)
	if err != nil {
		return
	}
	op, ok := int64Comparator[f.Comparator]
	if !ok {
		tips := fmt.Sprintf("unsupported filter: %s, field: %s", f.Source, f.Comparator)
		return  nil, errors.New(tips)
	}
	filter = func(img docker.APIImages) bool {
		return op(img.Size, size)
	}
	return
}

func RunCmdImage(cmd *cobra.Command, args []string) error {
	var filters []Filter
	for _, f := range filter {
		fmt.Println("Filter: ", f)
		parsed, err := parseFilter(f)
		if err != nil {
			return err
		}
		filters = append(filters, parsed)
	}
	return RemoveImages(filters...)
}


func RemoveImages (filters ...Filter) (err error) {
	var cli *docker.Client
	if dockerUri == "" {
		cli, err = docker.NewClientFromEnv()
	} else {
		cli, err = docker.NewClient(dockerUri)
	}
	if err != nil {
		return
	}
	iv, err := NewImageValidator(filters...)
	if err != nil {
		return
	}
	images, err := cli.ListImages(docker.ListImagesOptions{All:true})
	if err != nil {
		return
	}
	for _, img := range images {
		if iv.Satisfied(img) {
			if dryRun {
				fmt.Println("[DryRun]Removing image:", img.ID, img.RepoTags)
				continue
			}
			er := cli.RemoveImage(img.ID)
			if er != nil {
				fmt.Printf("can not remove image %s, reason: %s\n", img.ID, er)
			}
			fmt.Println("removed:", img.ID, img.RepoTags)
		}
	}
	return
}


// byteSize returns the size in Bytes against the given amount and unit
func ByteSize(amount int64, unit string) (int64, error) {
	var base int64
	switch strings.ToLower(unit) {
	case "k":
		base = 1 << 10
	case "m":
		base = 1 << 20
	case "g":
		base = 1 << 30
	default:
		return amount, errors.New("unsupported unit")
	}
	return amount * base, nil
}

