package image_manger

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imageorient"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	_ "golang.org/x/image/webp"
	"os"
	"path"
)

type ImageManager struct {
	targetDir string
}

type ImageProcessJob struct {
	Width  int
	Height int
	Anchor imaging.Anchor
}

func NewImageProcessJob(width int, height int, anchor imaging.Anchor) *ImageProcessJob {
	return &ImageProcessJob{Width: width, Height: height, Anchor: anchor}
}

func NewImageManager(targetDir string) *ImageManager {
	return &ImageManager{targetDir: targetDir}
}

func (im ImageManager) ProcessImage(source []byte, processingJobs ...*ImageProcessJob) (guids []uuid.UUID, err error) {
	img, _, err := imageorient.Decode(bytes.NewReader(source))
	if err == nil {
		guids = make([]uuid.UUID, len(processingJobs))
		for index, pJob := range processingJobs {
			newImage := imaging.Fill(img, pJob.Width, pJob.Height, pJob.Anchor, imaging.Lanczos)
			newUUID := uuid.New()
			err = imaging.Save(newImage, im.getImageFilePath(newUUID))
			if err == nil {
				guids[index] = newUUID
			} else {
				return
			}
		}
	}
	return
}

func (im ImageManager) getImageFilePath(uuid uuid.UUID) string {
	jpgExt := "jpg"
	return im.getImageFilePathExtension(uuid, &jpgExt)
}

func (im ImageManager) getImageFilePathExtension(uuid uuid.UUID, extension *string) string {
	extensionStr := "jpg"
	if extension != nil {
		extensionStr = *extension
	}
	imgPath := path.Clean(fmt.Sprintf("%s/%s.%s", im.targetDir, uuid, extensionStr))
	dir := path.Dir(imgPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
	}
	return imgPath
}

func (im ImageManager) ProcessImageAs16by9(source []byte) (guids []uuid.UUID, err error) {
	return im.ProcessImage(
		source,
		NewImageProcessJob(295, 166, imaging.Center),
		NewImageProcessJob(960, 540, imaging.Center),
		NewImageProcessJob(1920, 1080, imaging.Center),
	)
}

func (im ImageManager) ProcessImageAsSquare(source []byte) (guids []uuid.UUID, err error) {
	return im.ProcessImage(
		source,
		NewImageProcessJob(96, 96, imaging.Center),
		NewImageProcessJob(512, 512, imaging.Center),
		NewImageProcessJob(1080, 1080, imaging.Center),
	)
}
