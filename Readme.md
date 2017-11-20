# gostitcher

This project is my attempt to learn the go language by building an application that can stitch together the monochrome color photos taken by space probes into visible color photos.

## Iterations

My iterations are based on the following manual process that is outlined by Emily Ladawalla in [Tutorial: Making RGB Images in Photoshop](http://www.planetary.org/explore/space-topics/space-imaging/tutorial_rgb_ps.html). I am using the same images of Reha against Saturn from the Cassini space probe.

### 1. Color Masking

The manual tutorial uses a layers and channels tool of Photoshop. My first attempt to replicate the color combining in an automated fashion involved converting the gray scale images to an alpha mask and then layer the images on top of each other. This attempt produced a very red image and applied the red layer last. Thus I tried changing the order to have the blue image last and it produced a very blue result. This indicates that this approach is not replicating the channels approach and is not blending the colors together.

Red Image|Blue Image
----------|----------
![Red Heavy Image of Rhea](images/rhea/output_v1_alpha.jpg)|![Blue Heavy Image of Rhea](images/rhea/output_v1_beta.jpg)

### 2. Color Blending

To replicate the channel approach of the tutorial, this approach uses the Gray value from the three images as the RGB values of the resulting color pixel. This approach was much more successful than masking and generated an appropriate image.

![Blended Cropped image of Rhea](images/rhea/output_v2_alpha.jpg)

### 3. Alignment

This iteration attempted to align images by minimizing the subtracted images. The initial attempt at this did not seem to produce fundamentally better images with the minimization selecting the original alignment for blue and red, and picking an alignment with a clear arc for blue and green. I have not gone any further with this current approach.

Blue Green Original Subtraction|Blue Green "Optimal" Subtraction|Blue Red Original & "Optimal" Subtraction
--------|----------|--------------
![Blue Green Original Subtraction](images/rhea/output_v3_bg_align_00.jpg)|![Blue Green "Optimal" Subtraction](images/rhea/output_v3_bg_align_02.jpg)|![Blue Red Original & "Optimal" Subtraction](images/rhea/output_v3_br_align_00.jpg)
