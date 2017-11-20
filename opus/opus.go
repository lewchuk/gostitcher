package opus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

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
	colOpt := "cols=ringobsid,time1,time2,filter"
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

	fmt.Println(queryURL)

	client := cleanhttp.DefaultClient()

	resp, err := client.Get(queryURL)

	if (err != nil) { return err }

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if (err != nil) { return err }

	data := map[string]interface{}{}

	err = json.Unmarshal(body, &data)

	if (err != nil) { return err }

	fmt.Println(data)

	return nil
}
