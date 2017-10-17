// imgops project imgops.go
package imgops

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"

	"github.com/koyachi/go-nude"
	"github.com/nfnt/resize"
)

type ImageProc struct {
	buffer   []byte
	i        int64 // current reading index
	prevRune int
}

func New(buff []byte) ImageProc {
	ip := ImageProc{}
	ip.buffer = buff
	ip.i = 0
	ip.prevRune = -1

	return ip
}

//obsolete since there is error in Image.Decode. //This isssue is caused by the image api.
func (r ImageProc) Read(p []byte) (n int, err error) {
	if r.i >= int64(len(r.buffer)) {
		return 0, io.EOF
	}
	r.prevRune = -1
	n = copy(p, r.buffer[r.i:])
	r.i += int64(n)
	return
}

//Reads buffer and create a new io.Reader
func Read(buf []byte) io.Reader {
	return bytes.NewReader(buf)
}

func (imp ImageProc) IsNude() (isnude bool, err error) {
	fmt.Println("Working")
	img, _, err := image.Decode(Read(imp.buffer))
	isnude, err = nude.IsImageNude(img)
	if err != nil {
		return false, err
	}
	if isnude == true {
		return true, nil
	} else {
		return false, nil
	}
}

//Checkes whether a image is nude or not
//inmemory nude check
//this method is not working accurately, hence just being used.
/*func IsNude(reader io.Reader) (isnude bool, err error) {
	img, _, err := image.Decode(reader)
	isnude, err = nude.IsImageNude(img)
	if err != nil {
		return false, err
	}
	if isnude == true {
		return true, nil
	} else {
		return false, nil
	}
}*/

//takes image from src and resize with given width and height,copies to the destination.
//this mehthod has been created for goroutine since ResizeAndMove returns err
func ResizeAndMoveGo(src string, dest string, w, h uint) {
	err := ResizeAndMove(src, dest, w, h)
	if err != nil {
		//logger.WriteLog(err)
	}
}

//takes image from src and resize with given width and height,copies to the destination.
func ResizeAndMove(src string, dest string, w, h uint) (err error) {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}
	file.Close()

	m := resize.Resize(w, h, img, resize.Lanczos3)

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// write new image to file

	err = jpeg.Encode(out, m, nil)

	if err != nil {
		return err
	}
	return nil
}

//returns width and height of an image.
//dimensions are measured on the stored file
func (imp ImageProc) GetImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return image.Width, image.Height
}

//returns width and height of an image.
//this is an inmemory operation
func (imp ImageProc) GetImageDimensionBy() (int, int) {

	image, _, err := image.DecodeConfig(Read(imp.buffer))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	return image.Width, image.Height
}
