package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fieldkit/cloud/server/common/logging"
)

type DescribeLocations struct {
	MapboxToken string
	metrics     *logging.Metrics
}

type LocationDescription struct {
	OtherLandName  *string
	NativeLandName *string
}

func NewDescribeLocations(mapboxToken string, metrics *logging.Metrics) (ls *DescribeLocations) {
	return &DescribeLocations{
		MapboxToken: mapboxToken,
		metrics:     metrics,
	}
}

type OtherLandFeature struct {
	Text      string `json:"text"`
	PlaceName string `json:"place_name"`
}

type OtherLandResponse struct {
	Features []*OtherLandFeature `json:"features"`
}

func (ls *DescribeLocations) queryOther(ctx context.Context, l *Location) (name *string, err error) {
	timing := ls.metrics.ThirdPartyLocation("mapbox")

	defer timing.Send()

	query := fmt.Sprintf("%f,%f", l.Longitude(), l.Latitude())
	url := "https://api.mapbox.com/geocoding/v5/mapbox.places/" + query + ".json?types=place&access_token=" + ls.MapboxToken
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	parsed := OtherLandResponse{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		log := Logger(ctx).Sugar()
		log.Warnw("location-other-parse-error", "body", string(body))
		return nil, err
	}

	for _, feature := range parsed.Features {
		return &feature.PlaceName, nil
	}

	return nil, nil
}

type NativeLandProperties struct {
	Name        string `json:"name"`
	FrenchName  string `json:"french_name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type NativeLandInfo struct {
	Properties *NativeLandProperties `json:"properties"`
}

type NativeLandResponse = []*NativeLandInfo

func (ls *DescribeLocations) queryNative(ctx context.Context, l *Location) (name *string, err error) {
	timing := ls.metrics.ThirdPartyLocation("nativeland")

	defer timing.Send()

	query := fmt.Sprintf("%f,%f", l.Latitude(), l.Longitude())
	url := "https://native-land.ca/api/index.php?maps=territories&position=" + query
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	parsed := NativeLandResponse{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		log := Logger(ctx).Sugar()
		log.Warnw("location-native-parse-error", "body", string(body))
		return nil, err
	}

	names := make([]string, 0)
	for _, feature := range parsed {
		names = append(names, feature.Properties.Name)
	}
	if len(names) > 0 {
		all := strings.Join(names, ", ")
		return &all, nil
	}

	return nil, nil
}

func (ls *DescribeLocations) Describe(ctx context.Context, l *Location) (ld *LocationDescription, err error) {
	log := Logger(ctx).Sugar()

	timing := ls.metrics.ThirdPartyLocationDescribe()

	defer timing.Send()

	if ls.MapboxToken == "" {
		return nil, fmt.Errorf("location description disabled")
	}

	other, err := ls.queryOther(ctx, l)
	if err != nil {
		log.Warnw("error getting other place names", "error", err)
	}

	native, err := ls.queryNative(ctx, l)
	if err != nil {
		log.Warnw("error getting native place names", "error", err)
	}

	ld = &LocationDescription{
		OtherLandName:  other,
		NativeLandName: native,
	}

	return
}
