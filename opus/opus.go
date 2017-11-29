package opus

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	cleanhttp "github.com/hashicorp/go-cleanhttp"

	"github.com/lewchuk/gostitcher/algv2blending"
	"github.com/lewchuk/gostitcher/common"
)

type OpusDataAPIResponse struct {
	PageNo  int        `json:"page_no"`
	Columns []string   `json:"columns"`
	Images  [][]string `json:"page"`
}

type OpusImage struct {
	RingObsId string
	ObsKey    string
	Filter    string
	Time      time.Time
}

type OpusCountAPIResponse struct {
	Data []OpusCountResultAPIResponse `json:"data"`
}

type OpusCountResultAPIResponse struct {
	ResultCount int `json:"result_count"`
}

type OpusFilesAPIResponse struct {
	Data map[string]OpusFilesAPIImageResponse `json:"data"`
}

type OpusFilesAPIImageResponse struct {
	RawImages []string `json:"RAW_IMAGE"`
	CalibratedImages []string `json:"CALIBRATED"`
	PreviewImages []string `json:"preview_image"`
}

// getAPIQuery requests a URL and returns a reader of the response
func getAPIQuery(url string) (io.ReadCloser, error) {
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("creating request for %s: %s", url, err)
	}

	client := cleanhttp.DefaultClient()

	resp, err := client.Do(request)

	if err != nil {
		return nil, fmt.Errorf("requesting from %s: %s", url, err)
	}

	return resp.Body, nil
}

// getAPIQueryAsBytes requests a URL and returns the bytes of the response
func getAPIQueryAsBytes(url string) ([]byte, error) {
	reader, err := getAPIQuery(url)
	defer reader.Close()

	if err != nil {
		return nil, fmt.Errorf("loading api response %s: %s", url, err)
	}

	body, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, fmt.Errorf("loading api response %s: %s", url, err)
	}

	return body, nil

}

// findIdndex Finds the index of a string in an array of strings
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

// getDataAPIResponse given a query to the Opus data.json api, fetches and parses the JSON response.
func getDataAPIResponse(queryURL string) (OpusDataAPIResponse, error) {
	var data OpusDataAPIResponse

	body, err := getAPIQueryAsBytes(queryURL)

	if err != nil {
		return data, fmt.Errorf("loading api response %s: %s", queryURL, err)
	}

	data = OpusDataAPIResponse{}

	if err := json.Unmarshal(body, &data); err != nil {
		return data, fmt.Errorf("parsing api response %s: %s", body, err)
	}

	return data, nil
}

// getCountAPIResponse given a query to the Opus meta/result_count.json api, fetches and parses the JSON response.
func getCountAPIResponse(queryURL string) (OpusCountAPIResponse, error) {
	var data OpusCountAPIResponse

	body, err := getAPIQueryAsBytes(queryURL)

	if err != nil {
		return data, fmt.Errorf("loading api response %s: %s", queryURL, err)
	}

	data = OpusCountAPIResponse{}

	if err := json.Unmarshal(body, &data); err != nil {
		return data, fmt.Errorf("parsing api response %s: %s", body, err)
	}

	return data, nil
}

// translateDataAPIResonse translates a raw JSON response from the Opus data.json endpoint into an array of image metadata.
func translateDataAPIResonse(data OpusDataAPIResponse) ([]OpusImage, error) {
	idIndex := findIndex("Ring Observation ID", data.Columns)
	obsIndex := findIndex("Observation Name", data.Columns)
	timeIndex := findIndex("Observation Time 1 (UTC)", data.Columns)
	filterIndex := findIndex("Filter", data.Columns)

	if idIndex == -1 || obsIndex == -1 || filterIndex == -1 || timeIndex == -1 {
		return nil, fmt.Errorf(
			"Getting column indexes from %s (id: %d, obs: %d, time: %s, filter: %d)",
			data.Columns, idIndex, obsIndex, timeIndex, filterIndex)
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
			ObsKey:    imgArray[obsIndex],
			Filter:    imgArray[filterIndex],
			Time: dateTimeTaken,
		}
	}

	return images, nil
}

// groupImage groups images by the observation name and discards any observation without a
// full RGB image set.
func groupImages(images []OpusImage) map[string]common.ImageFilenameMap {
	lastObs := ""
	imageGroups := make(map[string]common.ImageFilenameMap)
	imageGroupIndex := -1

	for _, image := range images {
		// We are on a new group of images.
		if lastObs != image.ObsKey {
			// There is a previous group we just finished.
			if lastObs != "" {
				if err := common.ValidateImageMap(imageGroups[lastObs]); err != nil {
					fmt.Println("Group is not valid:", err)
					// Delete group so we don't try to process it more.
					delete(imageGroups, lastObs)
				}
			}

			lastObs = image.ObsKey
			fmt.Println("Starting new group:", imageGroupIndex)
			imageGroups[lastObs] = make(common.ImageFilenameMap)
		}
		fmt.Println(imageGroupIndex, image)
		// TODO: Take the set of Filter/Time tuples and pick the three images
		// with the least time deltas, probably using blue as a start point
		// since those images seem to be more rare in clusters of images.
		imageGroups[lastObs][image.Filter] = image.RingObsId
	}

	if err := common.ValidateImageMap(imageGroups[lastObs]); err != nil {
		fmt.Println("Group is not valid:", err)
		// Delete group so we don't try to process it more.
		delete(imageGroups, lastObs)
	}

	return imageGroups
}

// loadImage loads and caches the full sized JPEG preview image from OPUS for an observation id.
func loadImage(obsName, imageId, outputFolder string) (*image.Gray, error) {
	cacheFolder := fmt.Sprintf("%s/%s/", outputFolder, obsName)
	cachePath := fmt.Sprintf("%s/%s.jpg", cacheFolder, imageId)

	if err := os.MkdirAll(cacheFolder, os.ModePerm); err != nil {
		return nil, fmt.Errorf("cannot create cache folder %s: %s", cacheFolder, err)
	}

	if _, err := os.Stat(cachePath); !os.IsNotExist(err) {
		return common.LoadImageFromPath(cachePath)
	}

	apiRoot := "https://tools.pds-rings.seti.org/opus/api"

	queryURL := fmt.Sprintf(
		"%s/files/%s.json",
		apiRoot,
		imageId,
	)

	body, err := getAPIQueryAsBytes(queryURL)

	if err != nil {
		return nil, fmt.Errorf("loading api response %s: %s", queryURL, err)
	}

	data := OpusFilesAPIResponse{}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("parsing api response %s: %s", body, err)
	}

	files := data.Data[imageId].PreviewImages

	var fullImage string

	for _, path := range files {
		if strings.HasSuffix(path, "full.jpg") {
			fullImage = path
			break
		}
	}

	if fullImage == "" {
		return nil, fmt.Errorf("no full preview image in %s", files)
	}

	fmt.Println("Loading", fullImage)

	reader, err := getAPIQuery(fullImage)
	defer reader.Close()
	if err != nil {
		return nil, fmt.Errorf("error loading image from %s: %s", fullImage, err)
	}

	image, err := common.LoadImage(reader)
	if err != nil {
		return nil, fmt.Errorf("error loading image from %s: %s", fullImage, err)
	}

	err = common.WriteImage(cachePath, image)

	if err != nil {
		return nil, fmt.Errorf("error caching image at %s: %s", cachePath, err)
	}

	return image, nil
}

// combineImages combines a set of images representing a single observation.
func combineImages(obsName string, idMap common.ImageFilenameMap, outputFolder string) error {
	imageMap := make(common.ImageMap)

	for _, filter := range common.Filters {
		image, err := loadImage(obsName, idMap[filter], outputFolder)
		if err != nil {
			return err
		}
		imageMap[filter] = *image
	}

	outpuImage := algv2blending.BlendImage(imageMap)

	outputPath := fmt.Sprintf("%s/%s/%s.jpg", outputFolder, obsName, obsName)

	if err := common.WriteImage(outputPath, outpuImage); err != nil {
		return fmt.Errorf("error writing image to %s: %s", outputPath, err)
	}

	outputPath = fmt.Sprintf("%s/results/%s.jpg", outputFolder, obsName)

	if err := common.WriteImage(outputPath, outpuImage); err != nil {
		return fmt.Errorf("error writing image to %s: %s", outputPath, err)
	}

	return nil
}

func CombineImages(rawTarget string) error {

	target := strings.ToUpper(rawTarget)
	outputLocation := fmt.Sprintf("images/opus/%s", rawTarget)

	if err := os.MkdirAll(fmt.Sprintf("%s/results", outputLocation), os.ModePerm); err != nil {
		return fmt.Errorf("cannot create resulsts folder: %s", err)
	}

	// OPUS API components for a page of images from the Cassini ISS instrument
	// from https://tools.pds-rings.seti.org/opus/api/.
	// Result metadata filtered to just narrow angle RED, BL1 and GRN images.
	apiRoot := "https://tools.pds-rings.seti.org/opus/api"
	cassiniISSOpts := "instrumentid=Cassini+ISS&typeid=Image"
	saturnOpt := "planet=Saturn"
	enceladusTargetOpt := fmt.Sprintf("target=%s", target)
	filterOpts := "FILTER=BL1,GRN,RED"
	orderOpt := "order=time1"
	colOpt := "cols=ringobsid,obsname,filter,time1"
	narrowCameraOpt := "camera=Narrow+Angle"
	pageSizeOpt := "limit=100"

	searchParams := fmt.Sprintf(
		"%s&%s&%s&%s&%s&%s&%s&%s",
		cassiniISSOpts,
		saturnOpt,
		enceladusTargetOpt,
		filterOpts,
		narrowCameraOpt,
		orderOpt,
		colOpt,
		pageSizeOpt)

	countURL := fmt.Sprintf(
		"%s/meta/result_count.json?%s",
		apiRoot,
		searchParams)

	countData, err := getCountAPIResponse(countURL)

	if err != nil {
		return err
	}

	count := countData.Data[0].ResultCount

	baseURL := fmt.Sprintf(
		"%s/data.json?%s",
		apiRoot,
		searchParams)

	fmt.Println(baseURL)

	images := make([]OpusImage, 0, count)

	for i := 1; i <= count/100; i++ {
	    queryURL := fmt.Sprintf("%s&page=%d", baseURL, i)
	    fmt.Println(queryURL)
	    data, err := getDataAPIResponse(queryURL)

		if err != nil {
			return err
		}

		imagePage, err := translateDataAPIResonse(data)

		if err != nil {
			return err
		}

		images = append(images, imagePage...)
	}

	groups := groupImages(images)

	for obsName, imageMap := range groups {
		err := combineImages(obsName, imageMap, outputLocation)
		if err != nil {
			return err
		}
	}

	return nil
}
