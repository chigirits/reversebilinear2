package main

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func demixTL(img image.Image, x, y int) *Pixel {
	return AddPixel4(
		NewPixel(img.At(x+1, y+1)).Mul(9.0/4),
		NewPixel(img.At(x+2, y+1)).Mul(-3.0/4),
		NewPixel(img.At(x+1, y+2)).Mul(-3.0/4),
		NewPixel(img.At(x+2, y+2)).Mul(1.0/4),
	)
}

func demixTR(img image.Image, x, y int) *Pixel {
	return AddPixel4(
		NewPixel(img.At(x+1, y+1)).Mul(-3.0/4),
		NewPixel(img.At(x+2, y+1)).Mul(9.0/4),
		NewPixel(img.At(x+1, y+2)).Mul(1.0/4),
		NewPixel(img.At(x+2, y+2)).Mul(-3.0/4),
	)
}

func demixBL(img image.Image, x, y int) *Pixel {
	return AddPixel4(
		NewPixel(img.At(x+1, y+1)).Mul(-3.0/4),
		NewPixel(img.At(x+2, y+1)).Mul(1.0/4),
		NewPixel(img.At(x+1, y+2)).Mul(9.0/4),
		NewPixel(img.At(x+2, y+2)).Mul(-3.0/4),
	)
}

func demixBR(img image.Image, x, y int) *Pixel {
	return AddPixel4(
		NewPixel(img.At(x+1, y+1)).Mul(1.0/4),
		NewPixel(img.At(x+2, y+1)).Mul(-3.0/4),
		NewPixel(img.At(x+1, y+2)).Mul(-3.0/4),
		NewPixel(img.At(x+2, y+2)).Mul(9.0/4),
	)
}

func processFile(inFile, outFile string) error {
	fmt.Println(inFile + " -> " + outFile)

	inReader, err := os.Open(inFile)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", inFile, err)
	}
	defer inReader.Close()

	img, err := png.Decode(inReader)
	if err != nil {
		return fmt.Errorf("failed to decode PNG file %s: %w", inFile, err)
	}
	bounds := img.Bounds()
	x0 := bounds.Min.X
	y0 := bounds.Min.Y
	sizeX := (bounds.Max.X - x0) / 2
	sizeY := (bounds.Max.Y - y0) / 2
	dest := image.NewRGBA(image.Rect(x0, y0, sizeX, sizeY))

	for y := 0; y < sizeY-1; y++ {
		for x := 0; x < sizeX; x++ {
			var p *Pixel
			if y < sizeY-1 {
				if x < sizeX-1 {
					p = demixTL(img, x0+x*2, y0+y*2)
				} else {
					p = demixTR(img, x0+(x-1)*2, y0+y*2)
				}
			} else {
				if x < sizeX-1 {
					p = demixBL(img, x0+x*2, y0+(y-1)*2)
				} else {
					p = demixBR(img, x0+(x-1)*2, y0+(y-1)*2)
				}
			}
			dest.Set(x, y, p.ToColor())
		}
	}

	outWriter, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outFile, err)
	}
	defer outWriter.Close()
	if err := png.Encode(outWriter, dest); err != nil {
		return fmt.Errorf("failed to encode image into PNG file %s: %w", outFile, err)
	}

	return nil
}

func pause() {
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func raiseError(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	pause()
}

func main() {
	if len(os.Args) < 2 {
		raiseError(fmt.Errorf("USAGE: reversebilinear2.exe <IN> [<OUT>]"))
		return
	}
	inFile := os.Args[1]
	outFile := strings.TrimSuffix(inFile, ".png") + ".half"
	if 3 <= len(os.Args) {
		outFile = os.Args[2]
	}

	info, err := os.Stat(inFile)
	if err != nil {
		raiseError(err)
		return
	}
	if info.IsDir() {
		if err := os.Mkdir(outFile, 0777); err != nil {
			raiseError(err)
			return
		}
		fn := func(inFile string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			name := filepath.Base(inFile)
			processFile(inFile, path.Join(outFile, name))
			return nil
		}
		if err := filepath.Walk(inFile, fn); err != nil {
			raiseError(err)
			return
		}
	} else {
		if err := processFile(inFile, outFile+".png"); err != nil {
			raiseError(err)
			return
		}
	}
	fmt.Println("Ok")
}
