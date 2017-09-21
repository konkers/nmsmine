package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/konkers/nmsmine"
)

var dataDir = flag.String("data", "", "Directory containing exml files.")
var outFile = flag.String("out", "", "File to write json output.")

func main() {
	flag.Parse()

	if *dataDir == "" {
		fmt.Errorf("-data must be specified.\n")
		flag.Usage()
		os.Exit(1)
	}
	if *outFile == "" {
		fmt.Errorf("-out must be specified.\n")
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

	b, err := json.MarshalIndent(loader.Db, "", "  ")
	if err != nil {
		fmt.Errorf("Error encoding JSON: %s\n", err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(*outFile, b, 0644)
	if err != nil {
		fmt.Errorf("Error writing to %s: %s\n", *outFile, err)
		os.Exit(1)
	}
}
