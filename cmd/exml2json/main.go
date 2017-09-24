package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/konkers/nmsmine"
)

var dataDir = flag.String("data", "", "Directory containing exml files.")
var outFile = flag.String("out", "", "File to write json output.")

func main() {
	flag.Parse()

	if *dataDir == "" {
		fmt.Fprintf(os.Stderr, "-data must be specified.\n")
		flag.Usage()
		os.Exit(1)
	}
	if *outFile == "" {
		fmt.Fprintf(os.Stderr, "-out must be specified.\n")
		flag.Usage()
		os.Exit(1)
	}

	loader := nmsmine.NewExmlLoader()

	loader.LoadLocFile(
		path.Join(*dataDir, "LANGUAGE/NMS_LOC1_USENGLISH.exml"))
	loader.LoadLocFile(
		path.Join(*dataDir, "LANGUAGE/NMS_UPDATE3_USENGLISH.exml"))

	loader.LoadItemFile(
		path.Join(*dataDir, "METADATA/REALITY/TABLES/NMS_REALITY_GCPRODUCTTABLE.exml"))
	loader.LoadItemFile(
		path.Join(*dataDir, "METADATA/REALITY/TABLES/NMS_REALITY_GCSUBSTANCETABLE.exml"))
	loader.LoadItemFile(
		path.Join(*dataDir, "METADATA/REALITY/TABLES/NMS_REALITY_GCTECHNOLOGYTABLE.exml"))

	loader.LoadBuildingFile(
		path.Join(*dataDir, "METADATA/REALITY/TABLES/BASEBUILDINGTABLE.exml"))

	err := loader.Db.WriteToFile(*outFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
