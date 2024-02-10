package cmd

import (
	"archiver/lib/comperssion"
	"archiver/lib/comperssion/vlc"
	"archiver/lib/comperssion/vlc/table/shannon_fano"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var unpackCmd = &cobra.Command{
	Use:   "unpack",
	Short: "Unpack file",
	Run:   unpack,
}

const unpackedExtensions = "txt"

func unpack(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		handleErr(ErrEmptyPath)
	}

	var decoder comperssion.Decoder
	method := cmd.Flag("method").Value.String()

	switch method {
	case "vlc":
		decoder = vlc.New(shannon_fano.NewGenerator())
	default:
		cmd.PrintErr("unknown method")
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

	unpacked := decoder.Decode(data)

	// TODO: add file path to output file.
	err = os.WriteFile(unpackedFileName("out"), []byte(unpacked), 0644)
	if err != nil {
		handleErr(err)
	}
}

// TODO: refactor this
func unpackedFileName(path string) string {
	fileName := filepath.Base(path)
	ext := filepath.Ext(fileName)
	baseName := strings.TrimSuffix(fileName, ext)

	return baseName + "." + unpackedExtensions
}

func init() {
	rootCmd.AddCommand(unpackCmd)

	unpackCmd.Flags().StringP("method", "m", "",
		"compression method: vlc")

	if err := packCmd.MarkFlagRequired("method"); err != nil {
		panic(err)
	}
}
