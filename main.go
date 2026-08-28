package main

import (
	"flag"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"

	// Register image decoders.
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	// BMP and TIFF decoders.
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
)

var supportedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".bmp":  true,
	".tif":  true,
	".tiff": true,
}

func main() {
	quality := flag.Float64("q", 85, "WebP quality (0-100)")
	recursive := flag.Bool("r", false, "Process subdirectories recursively")
	deleteOriginal := flag.Bool("delete", false, "Delete original files after successful conversion")
	overwrite := flag.Bool("overwrite", false, "Overwrite existing WebP files")

	flag.Parse()

	if *quality < 0 || *quality > 100 {
		fmt.Println("Error: quality must be between 0 and 100")
		os.Exit(1)
	}

	root := "."

	fmt.Println("img2webp")
	fmt.Println("─────────")
	fmt.Printf("Directory: %s\n", root)
	fmt.Printf("Quality:   %.0f\n", *quality)
	fmt.Printf("Recursive: %v\n", *recursive)

	if *deleteOriginal {
		fmt.Println("Delete:    yes")
	} else {
		fmt.Println("Delete:    no")
	}

	fmt.Println()

	var files []string
	var err error

	if *recursive {
		files, err = findRecursive(root)
	} else {
		files, err = findCurrentDirectory(root)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No images found.")
		return
	}

	converted := 0
	skipped := 0
	failed := 0

	for _, filename := range files {
		result := convertImage(
			filename,
			*quality,
			*deleteOriginal,
			*overwrite,
		)

		switch result {
		case "converted":
			converted++
		case "skipped":
			skipped++
		case "failed":
			failed++
		}
	}

	fmt.Println()
	fmt.Println("─────────")
	fmt.Printf("Converted: %d\n", converted)
	fmt.Printf("Skipped:   %d\n", skipped)
	fmt.Printf("Failed:    %d\n", failed)

	if failed > 0 {
		os.Exit(1)
	}
}

func findCurrentDirectory(root string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var files []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if isSupportedImage(entry.Name()) {
			files = append(files, filepath.Join(root, entry.Name()))
		}
	}

	return files, nil
}

func findRecursive(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if isSupportedImage(path) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func isSupportedImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))

	return supportedExtensions[ext]
}

func convertImage(
	filename string,
	quality float64,
	deleteOriginal bool,
	overwrite bool,
) string {

	extension := filepath.Ext(filename)

	output := strings.TrimSuffix(filename, extension) + ".webp"

	// Prevent converting a WebP file to itself.
	if strings.EqualFold(extension, ".webp") {
		return "skipped"
	}

	// Check whether destination already exists.
	if _, err := os.Stat(output); err == nil {
		if !overwrite {
			fmt.Printf("⏭  %s → already exists\n", filename)
			return "skipped"
		}
	}

	input, err := os.Open(filename)
	if err != nil {
		fmt.Printf("❌ %s → %v\n", filename, err)
		return "failed"
	}
	defer input.Close()

	img, format, err := image.Decode(input)
	if err != nil {
		fmt.Printf("❌ %s → decode error: %v\n", filename, err)
		return "failed"
	}

	// Create temporary output first.
	// This prevents a broken WebP file if encoding fails.
	tempFile, err := os.CreateTemp(
		filepath.Dir(output),
		".img2webp-*.webp",
	)

	if err != nil {
		fmt.Printf("❌ %s → cannot create temporary file: %v\n", filename, err)
		return "failed"
	}

	tempName := tempFile.Name()

	defer func() {
		tempFile.Close()
		os.Remove(tempName)
	}()

	err = webp.Encode(tempFile, img, &webp.Options{
		Lossless: false,
		Quality: float32(quality),
	})

	if err != nil {
		fmt.Printf("❌ %s → encode error: %v\n", filename, err)
		return "failed"
	}

	if err := tempFile.Close(); err != nil {
		fmt.Printf("❌ %s → close error: %v\n", filename, err)
		return "failed"
	}

	// Remove existing output if overwrite was requested.
	if overwrite {
		if err := os.Remove(output); err != nil && !os.IsNotExist(err) {
			fmt.Printf("❌ %s → cannot replace existing file: %v\n", filename, err)
			return "failed"
		}
	}

	if err := os.Rename(tempName, output); err != nil {
		fmt.Printf("❌ %s → cannot create output: %v\n", filename, err)
		return "failed"
	}

	fmt.Printf(
		"✓  %s [%s] → %s\n",
		filename,
		format,
		output,
	)

	if deleteOriginal {
		if err := os.Remove(filename); err != nil {
			fmt.Printf(
				"⚠  %s → converted, but could not delete original: %v\n",
				filename,
				err,
			)
		}
	}

	return "converted"
}