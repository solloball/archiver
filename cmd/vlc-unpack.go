package cmd

import (
	"archiver/lib/vlc"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var vlcUnpackCmd = &cobra.Command{
	Use:   "vlc",
	Short: "Unpack file using variable-length code",
	Run:   unpack,
}

const unpackedExtensions = "txt"

func unpack(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		handleErr(ErrEmptyPath)
	}
	filePath := args[0]

	r, err := os.Open(filePath)
	if err != nil {
		handleErr(err)
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		handleErr(err)
	}

	unpacked := vlc.Decode(string(data))

	// TODO: add file path to output file.
	err = os.WriteFile(unpackedFileName("out"), []byte(unpacked), 0644)
	if err != nil {
		handleErr(err)
	}
}

func unpackedFileName(path string) string {
	fileName := filepath.Base(path)
	ext := filepath.Ext(fileName)
	baseName := strings.TrimSuffix(fileName, ext)

	return baseName + "." + packedExtensions
}

func init() {
	unpackCmd.AddCommand(vlcUnpackCmd)
}
