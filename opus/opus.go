package opus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

type OpusAPIResponse struct {
	PageNo  int        `json:"page_no"`
	Columns []string   `json:"columns"`
	Images  [][]string `json:"page"`
}

type OpusImage struct {
	RingObsId string
	Time      time.Time
	Filter    string
}

func getAPIResponse(request *http.Request) (OpusAPIResponse, error) {
	client := cleanhttp.DefaultClient()

	var data OpusAPIResponse

	resp, err := client.Do(request)

	if err != nil {
		return data, fmt.Errorf("requesting from %s: %s", request.URL, err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return data, fmt.Errorf("loading api response %s: %s", request.URL, err)
	}

	data = OpusAPIResponse{}

	if err := json.Unmarshal(body, &data); err != nil {
		return data, fmt.Errorf("parsing api response %s: %s", body, err)
	}

	return data, nil
}

func findIndex(val string, content []string) int {
	index := -1

	for i, v := range content {
		if val == v {
			index = i
			break
		}
	}

	return index
}

func translateAPIResonse(data OpusAPIResponse) ([]OpusImage, error) {
	idIndex := findIndex("Ring Observation ID", data.Columns)
	timeIndex := findIndex("Observation Time 1 (UTC)", data.Columns)
	filterIndex := findIndex("Filter", data.Columns)

	if idIndex == -1 || timeIndex == -1 || filterIndex == -1 {
		return nil, fmt.Errorf(
			"Getting column indexes from %s (id: %d, time: %d, filter: %d)",
			data.Columns, idIndex, timeIndex, filterIndex)
	}

	images := make([]OpusImage, len(data.Images))
	for i, imgArray := range data.Images {
		timeS := imgArray[timeIndex]
		// Replace the day of year with a dummy value singe golang time doesn't support day of year.
		dayOfYear, err := strconv.Atoi(timeS[5:8])
		if err != nil {
			return nil, fmt.Errorf("parsing day of year from %s: %s", timeS, err)
		}
		dayOfYearDuration := time.Duration(dayOfYear*24) * time.Hour
		dayStripped := strings.Replace(timeS, timeS[4:8], "-01", 1)
		timeTaken, err := time.Parse("2006-02T15:04:05.000", dayStripped)
		dateTimeTaken := timeTaken.Add(dayOfYearDuration)

		if err != nil {
			return nil, fmt.Errorf("parsing time %s: %s", imgArray[timeIndex], err)
		}

		images[i] = OpusImage{
			RingObsId: imgArray[idIndex],
			Time:      dateTimeTaken,
			Filter:    imgArray[filterIndex],
		}
	}

	return images, nil
}

func CombineImages() error {

	// OPUS API components for a single page of images of Enceladus from https://tools.pds-rings.seti.org/opus/api/.
	// Result metadata filtered to just narrow angle RED, BL1 and GRN images. Ultimately the planet and target opts could
	// be more configurable and this could load multiple pages.
	apiRoot := "https://tools.pds-rings.seti.org/opus/api"
	cassiniISSOpts := "instrumentid=Cassini+ISS&typeid=Image"
	saturnOpt := "planet=Saturn"
	enceladusTargetOpt := "target=ENCELADUS"
	filterOpts := "FILTER=BL1,GRN,RED"
	orderOpt := "order=time1"
	colOpt := "cols=ringobsid,time1,filter"
	narrowCameraOpt := "camera=Narrow+Angle"
	pageSizeOpt := "limit=99"

	// 2 minute range for correlating images.
	// timeRange := 2 * 60 * 1000

	queryURL := fmt.Sprintf(
		"%s/data.json?%s&%s&%s&%s&%s&%s&%s&%s&page=1",
		apiRoot,
		cassiniISSOpts,
		saturnOpt,
		enceladusTargetOpt,
		filterOpts,
		narrowCameraOpt,
		orderOpt,
		colOpt,
		pageSizeOpt)

	request, err := http.NewRequest("GET", queryURL, nil)

	if err != nil {
		return fmt.Errorf("creating request for %s: %s", queryURL, err)
	}

	data, err := getAPIResponse(request)

	if err != nil {
		return err
	}

	images, err := translateAPIResonse(data)

	if err != nil {
		return err
	}

	fmt.Println(images)

	return nil
}
