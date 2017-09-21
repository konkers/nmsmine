package main

import (
	"encoding/json"
	"fmt"

	"github.com/konkers/nmsmine"
)

func main() {
	loader := nmsmine.NewExmlLoader()

	loader.LoadLocFile("data/LANGUAGE/NMS_LOC1_USENGLISH.exml")
	loader.LoadLocFile("data/LANGUAGE/NMS_UPDATE3_USENGLISH.exml")

	loader.LoadItemFile(
		"data/METADATA/REALITY/TABLES/NMS_REALITY_GCPRODUCTTABLE.exml")
	loader.LoadItemFile(
		"data/METADATA/REALITY/TABLES/NMS_REALITY_GCSUBSTANCETABLE.exml")
	loader.LoadItemFile(
		"data/METADATA/REALITY/TABLES/NMS_REALITY_GCTECHNOLOGYTABLE.exml")

	loader.LoadBuildingFile(
		"data/METADATA/REALITY/TABLES/BASEBUILDINGTABLE.exml")

	b, err := json.MarshalIndent(loader.Db, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))

}
