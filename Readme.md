# gostitcher

This program uses the [OPUS API](https://tools.pds-rings.seti.org/opus/about/) to stitch images taken with red, blue and green (RGB) filters together to create "colour" images from space craft, specifically the Cassini ISS, see the sections below for more details about the use of the OPUS API and the image stitching algorithms.

To use this program:

1. Install the go programming language: https://golang.org/doc/install.
1. Fetch this program: `go get github.com/lewchuk/gostitcher`.
1. Navigate to the location of this program, e.g. `cd ~/go/src/github.com/lewchuk/gostitcher` if you installed go into your home directory.
1. Build the project `go build`.
1. Run the program to show the different options `./gostitcher --help`.

Some example commands:

- `./gostitcher --api ~/rhea --observation ISS_046RH_LIMB270L001_PRIME`. Pulls images from a single observation of Rhea.
- `./gostitcher --api ~/rhea --target rhea`. Pulls all RGB image metadata from the API, group by observation and filter to those with a full compliment of RGB images.


## Iterations

This project was both an attempt to learn Go as well as become more familiar with space craft images. Thus I have implemented a few different algorithms for stiching together images. Currently the one implemented for the OPUS API integration is the V2 algorithm described below.

My iterations are based on the following manual process that is outlined by Emily Lakdawalla in [Tutorial: Making RGB Images in Photoshop](http://www.planetary.org/explore/space-topics/space-imaging/tutorial_rgb_ps.html). I am using the same images of Reha against Saturn from the Cassini space probe.

These iterations can be run by specifying a `--path <path>` parameter where the path points to a folder containing source images and a config.json file matching the images to the filter used to take them. Note that downloading images from the API will cache the source images and generate a config.json file appropriate for using `--path` mode of this program.

### 1. Colour Masking

The manual tutorial uses a layers and channels tool of Photoshop. My first attempt to replicate the colour combining in an automated fashion involved converting the gray scale images to an alpha mask and then layer the images on top of each other. This attempt produced a very red image and applied the red layer last. Thus I tried changing the order to have the blue image last and it produced a very blue result. This indicates that this approach is not replicating the channels approach and is not blending the colors together.

Red Image|Blue Image
----------|----------
![Red Heavy Image of Rhea](images/rhea/output_v1_alpha.jpg)|![Blue Heavy Image of Rhea](images/rhea/output_v1_beta.jpg)

### 2. Color Blending

To replicate the channel approach of the tutorial, this approach uses the Gray value from the three images as the RGB values of the resulting color pixel. This approach was much more successful than masking and generated an appropriate image.

![Blended Cropped image of Rhea](images/rhea/output_v2_alpha.jpg)

### 3. Alignment

This iteration attempts to align images by minimizing the subtracted images. The initial attempt at this did not seem to produce fundamentally better images with the minimization selecting the original alignment for blue and red, and picking an alignment with a clear arc for blue and green as seen below:

Blue Green Original Subtraction|Blue Green "Optimal" Subtraction|Blue Red Original & "Optimal" Subtraction
--------|----------|--------------
![Blue Green Original Subtraction](images/rhea/output_v3_bg_align_00.jpg)|![Blue Green "Optimal" Subtraction](images/rhea/output_v3_bg_align_02.jpg)|![Blue Red Original & "Optimal" Subtraction](images/rhea/output_v3_br_align_00.jpg)

However, after seeing a highly shifted set of images of Enceladus and applying the same algorithm again, it performed much better. I suspect that the the background colour of the Rhea images and thus the difference in intensity between the various filters overwhelms the differences in intensity caused by misalignment.

Original Blended Image|Aligned Blended Image
----|----
![Original Blended Image](images/opus/enceladus/output_v2_alpha.jpg)|![Aligned Blended Image](images/opus/enceladus/output_v3.jpg)
----|----
Original Blue Green Subtraction|Aligned Blue Green Subtraction (-1, -44)
----|----
![Original Blue Green Subtraction](images/opus/enceladus/output_v3_bg_align_00.jpg)|![Aligned Blue Green Subtraction (-1, -44)](images/opus/enceladus/output_v3_bg_align_-1-44.jpg)
----|----
Original Blue Red Subtraction|Aligned Blue Red Subtraction (-1, -44)
----|----
![Original Blue Red Subtraction](images/opus/enceladus/output_v3_br_align_00.jpg)|![Aligned Blue Red Subtraction (-4, -89)](images/opus/enceladus/output_v3_br_align_-4-89.jpg)

The current alignment algorithm is pure brute force. However, there are likely some significant possible improvements:

- Since the effectiveness of each alignment can be summarized in a single value, algorithms for efficiently finding minimum points in 2D matrix of values can be applied to this problem to reduce the number of alignments to consider.
- It seems likely that for many images like the Enceladus image the "geometry" of the effectiveness matrix will be simple so greedy path finding algorithms could be used to quickly find the minimum location.
- When combined with other OPUS metadata on space craft and target positions to estimate an alignment.

## OPUS API

[OPUS](https://tools.pds-rings.seti.org/opus/about/) is a data search tool for NASA outer planets missions, a project of the [Planetary Rings Node](http://pds-rings.seti.org/). It provides a way to search for images across most of the metadata available for those images. Currently, of the data available on OPUS, only images from the Cassini Imaging Science Subsystem (ISS) contain the necessary Filter metadata to programatically select and combine images so only images from that mission are supported.

This mode of gostitcher uses the search tools to identify observations with a full set of RGB images and then combines them with the V2 algorithm above. Since this does no alignment results are not publication ready but this can be an effective tool to preview which observations are likely to have promising images.

This mode is run by specifying a `--api <output>` option identifying where to put the images and then some set of filtering parameters such as `--target` or `--observation` to filter images to a manageable result.

### Output

Using the API mode will select three images (one of each filter) from each observation and download those images into folders named for the observation. It will then combine the images and write the result into the observation folder as well as another `result` folder that will only include the "color" images.

### Examples

Some promising images of Enceladus and Rhea:

#### Rhea

![Rhea ISS_046RH_LIMB270L001_PRIME](images/opus/rhea/results/ISS_046RH_LIMB270L001_PRIME.jpg)
![Rhea ISS_048RH_LIMB270L001_PRIME](images/opus/rhea/results/ISS_048RH_LIMB270L001_PRIME.jpg)
![Rhea ISS_121RH_MUTUALEVE001_PRIME](images/opus/rhea/results/ISS_121RH_MUTUALEVE001_PRIME.jpg)
![Rhea ISS_219RH_RHEA001_CIRS](images/opus/rhea/results/ISS_219RH_RHEA001_CIRS.jpg)

#### Enceladus

![Enceladus ISS_019EN_FP3HOTSPT020_CIRS](images/opus/enceladus/results/ISS_019EN_FP3HOTSPT020_CIRS.jpg)
![Enceladus ISS_080EN_FP1ECLSCN001_CIRS](images/opus/enceladus/results/ISS_080EN_FP1ECLSCN001_CIRS.jpg)
![Enceladus ISS_219EN_ENCEL001_CIRS](images/opus/enceladus/results/ISS_219EN_ENCEL001_CIRS.jpg)
![Enceladus ISS_141EN_PLMHPHR001_PIE](images/opus/enceladus/results/ISS_141EN_PLMHPHR001_PIE.jpg)

#### Titan

![Titan ISS_128TI_M180R2HZ075_PRIME](images/opus/titan/results/ISS_128TI_M180R2HZ075_PRIME.jpg)
![Titan ISS_130TI_MUTUALEVE006_PRIME](images/opus/titan/results/ISS_130TI_MUTUALEVE006_PRIME.jpg)
![Titan ISS_148TI_MUTUALEVE003_PRIME](images/opus/titan/results/ISS_148TI_MUTUALEVE003_PRIME.jpg)

### Improvements

Some potential improvements:

- Being more selective about images when multiple options for each filter are available. For example observation ISS_021RH_GLOCOL001_PRIME is a great example where there are two green filtered images but the last one (which this algorithm currently picks) is ten minutes later and has a much longer exposure, making the resulting image much more green. A first alternative is to select three images that minimize the time between them. Another might be to minimize on exposure length to better match them.
- Hosting as a web service to explore OPUS images in an interactive manner.
- Figure out if images other than the preview (e.g. the calibrated images) could be better used. Note that the two data sources created images of Rhea against Saturn of different colours.
- Consider leveraging additional information about the images to improve the image stitching process. For example image metadata contains and exposure value, if images have differing exposure levels images could be adjusted.
