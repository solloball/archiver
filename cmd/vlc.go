package cmd

import (
	"archiver/lib/vlc"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var vlcCmd = &cobra.Command{
	Use:   "vlc",
	Short: "Pack file using variable-length code",
	Run:   pack,
}

const packedExtensions = "vlc"

func pack(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		handleErr(errors.New("path to file is not declarated"))
	}
	filePath := args[0]

	r, err := os.Open(filePath)
	if err != nil {
		handleErr(err)
	}

	data, err := io.ReadAll(r)
	if err != nil {
		handleErr(err)
	}
	defer r.Close()

	packed := vlc.Encode(string(data))

	// TODO: add file path to output file.
	err = os.WriteFile(packedFileName("out"), []byte(packed), 0644)
	if err != nil {
		handleErr(err)
	}
}

func packedFileName(path string) string {
	fileName := filepath.Base(path)
	ext := filepath.Ext(fileName)
	baseName := strings.TrimSuffix(fileName, ext)

	return baseName + "." + packedExtensions
}

func init() {
	packCmd.AddCommand(vlcCmd)
}
