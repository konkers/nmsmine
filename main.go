package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Property struct {
	Name       string      `xml:"name,attr"`
	Value      string      `xml:"value,attr"`
	Properties []*Property `xml:"Property"`
}

type Data struct {
	Properties []*Property `xml:"Property"`
}

type ItemColor struct {
	R float64 `json:"r"`
	G float64 `json:"g"`
	B float64 `json:"b"`
	A float64 `josn:"a"`
}

type ItemCost struct {
	SpaceStationMarkup float64 `json:"space_station_markup"`
	LowPriceMod        float64 `json:"low_price_mod"`
	HighPriceMod       float64 `json:"high_price_mod"`
	BuyBaseMarkup      float64 `json:"buy_base_markup"`
	BuyMarkupMod       float64 `json:"buy_markup_mod"`
}

type ItemRequirement struct {
	Id       string `json:"id"`
	Quantity int64  `json:"quantity"`
}

type ItemStatBonus struct {
	Type  string `json:"type"`
	Bonus int64  `json:"bonus"`
	Level int64  `json:"level"`
}

type ItemBuildingInfo struct {
	BuildableOnBase      bool   `json:"buildable_on_base"`
	BuildableOnFreighter bool   `json:"buildable_on_freighter"`
	BuildableOnPlanet    bool   `json:"buildable_on_planet"`
	ComplexityCost       int64  `json:"complexity_cost"`
	Group                string `json:"group"`
	CanChangeColor       bool   `json:"can_change_color"`
	CanChangeMaterial    bool   `json:"can_change_material"`
}

type Item struct {
	Id                      string            `json:"id"`
	Name                    string            `json:"name"`
	NameLower               string            `json:"name_lower"`
	Symbol                  string            `json:"symbol"`
	Subtitle                string            `json:"subtitle"`
	Description             string            `json:"description"`
	Teach                   bool              `json:"teach"`
	Hint                    string            `json:"hint"`
	BaseValue               float64           `json:"base_value"`
	Level                   int64             `json:"level"`
	Icon                    string            `json:"icon"`
	Color                   string            `json:"color"`
	WorldColor              string            `json:"world_color"`
	Charegeable             bool              `json:"chargeable"`
	CharegeAmount           int64             `json:"charge_amount"`
	SubstanceCategory       string            `json:"substance_category"`
	ProductCategory         string            `json:"product_category"`
	ChargeBy                []string          `json:"charge_by"`
	BuildFullyCharged       bool              `json:"build_fully_charged"`
	Upgrade                 bool              `json:"upgrade"`
	Core                    bool              `json:"core"`
	TechnologyCategory      string            `json:"technology_category"`
	TechnologyRarity        string            `json:"technology_rarity"`
	Value                   int64             `json:"value"`
	BaseStat                string            `json:"base_stat"`
	StatBonuses             []ItemStatBonus   `json:"stat_bonuses"`
	RequiredTech            string            `json:"required_tech"`
	RequiredLevel           int64             `json:"required_level"`
	UpgradeColor            string            `json:"upgrade_color"`
	LinkColor               string            `json:"link_color"`
	RewardGroup             string            `json:"reward_group"`
	Rarity                  string            `json:"rarity"`
	Legality                string            `json:"legality"`
	Consumable              bool              `json:"consumable"`
	ChargeValue             int64             `json:"charge_value"`
	Requirements            []ItemRequirement `json:"requirements"`
	Cost                    ItemCost          `json:"cost"`
	RequiredRank            int64             `json:"required_rank"`
	DispensingRace          string            `json:"dispensing_race"`
	TechShopRarity          string            `json:"tech_shop_rarity"`
	SpecificChargeOnly      bool              `json:"specific_charge_only"`
	NormalisedValueOnWorld  float64           `json:"normalised_value_on_world"`
	NormalisedValueOffWorld float64           `json:"normalised_value_off_world"`
	// TODO(TradingCategory)
	WikiEanabled               bool    `json:"wiki_enabled"`
	IsCraftable                bool    `json:"is_craftable"`
	EconomyInfluenceMultiplier float64 `json:"economy_influence_multiplier"`

	BuildingInfo *ItemBuildingInfo `json:"building_info"`
}

func loadDataFile(fileName string) *Data {
	file, err := os.Open(fileName) // For read access.
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)

	var data Data
	xml.Unmarshal(byteValue, &data)

	return &data
}

func handleString(entry *Property) string {
	return entry.Value
}

func handleIndirectString(entry *Property) string {
	return entry.Properties[0].Value
}

func handleIndirectStringArray(entry *Property) []string {
	var strings []string
	for _, data := range entry.Properties {
		strings = append(strings, handleIndirectString(data))
	}
	return strings
}

func handleBool(entry *Property) bool {
	return entry.Value == "True"
}

func handleInt(entry *Property) int64 {
	val, _ := strconv.ParseInt(entry.Value, 10, 64)
	return val
}

func handleFloat(entry *Property) float64 {
	val, _ := strconv.ParseFloat(entry.Value, 64)
	return val
}

func handleLoc(strings map[string]string, entry *Property) {
	var name string
	var value string

	for _, data := range entry.Properties {
		switch data.Name {
		case "Id":
			name = data.Value
		case "USEnglish":
			value = data.Properties[0].Value
		}
	}

	if name != "" && value != "" {
		strings[name] = value
	}
}

func loadLocFile(strings map[string]string, fileName string) {
	data := loadDataFile(fileName)
	for _, data := range data.Properties[0].Properties {
		handleLoc(strings, data)
	}
}

func handleItemIcon(entry *Property) string {
	path := entry.Properties[0].Value
	path = strings.TrimPrefix(path, "TEXTURES/UI/FRONTEND/ICONS/")
	path = strings.TrimSuffix(path, ".DDS")
	path = strings.ToLower(path)
	return path
}

func handleItemCost(entry *Property) ItemCost {
	var cost ItemCost
	for _, data := range entry.Properties {
		switch data.Name {
		case "SpaceStationMarkup":
			cost.SpaceStationMarkup, _ = strconv.ParseFloat(data.Value, 64)
		case "LowPriceMod":
			cost.LowPriceMod, _ = strconv.ParseFloat(data.Value, 64)
		case "HighPriceMod":
			cost.HighPriceMod, _ = strconv.ParseFloat(data.Value, 64)
		case "BuyBaseMarkup":
			cost.BuyBaseMarkup, _ = strconv.ParseFloat(data.Value, 64)
		case "BuyMarkupMod":
			cost.BuyMarkupMod, _ = strconv.ParseFloat(data.Value, 64)
		}
	}
	return cost
}

func handleItemColor(entry *Property) string {
	var color ItemColor
	for _, data := range entry.Properties {
		switch data.Name {
		case "R":
			color.R, _ = strconv.ParseFloat(data.Value, 64)
		case "G":
			color.G, _ = strconv.ParseFloat(data.Value, 64)
		case "B":
			color.B, _ = strconv.ParseFloat(data.Value, 64)
		case "A":
			color.A, _ = strconv.ParseFloat(data.Value, 64)
		}
	}
	return fmt.Sprintf("#%02x%02x%02x",
		uint8(color.R*255), uint8(color.G*255), uint8(color.B*255))
}

func handleItemRequirements(entry *Property) []ItemRequirement {
	var reqs []ItemRequirement
	for _, line := range entry.Properties {
		var req ItemRequirement
		for _, data := range line.Properties {
			switch data.Name {
			case "ID":
				req.Id = data.Value
			case "Amount":
				req.Quantity, _ = strconv.ParseInt(data.Value, 10, 64)
			}
		}
		reqs = append(reqs, req)
	}
	return reqs
}

func handleItemStatBonus(entry *Property) ItemStatBonus {
	var bonus ItemStatBonus
	for _, data := range entry.Properties {
		switch data.Name {
		case "StatsTypes":
			bonus.Type = handleIndirectString(data)
		case "Bonus":
			bonus.Bonus = handleInt(data)
		case "Level":
			bonus.Level = handleInt(data)
		}
	}
	return bonus
}

func handleItemStatBonuses(entry *Property) []ItemStatBonus {
	var bonuses []ItemStatBonus
	for _, data := range entry.Properties {
		bonuses = append(bonuses, handleItemStatBonus(data))
	}
	return bonuses
}

func handleItem(items map[string]*Item, entry *Property, strings map[string]string) {
	item := &Item{}

	for _, data := range entry.Properties {
		switch data.Name {
		case "ID":
			item.Id = data.Value
		case "Id":
			item.Id = data.Value
		case "Name":
			item.Name = strings[data.Value]
		case "NameLower":
			item.NameLower = strings[data.Value]
		case "Symbol":
			item.Symbol = strings[data.Value]
		case "Subtitle":
			item.Subtitle = strings[data.Properties[0].Value]
		case "Description":
			item.Description = strings[data.Properties[0].Value]
		case "Teach":
			item.Teach = handleBool(data)
		case "Hint":
			item.Hint = strings[data.Value]
		case "BaseValue":
			item.BaseValue, _ = strconv.ParseFloat(data.Value, 64)
		case "Level":
			item.Level, _ = strconv.ParseInt(data.Value, 10, 64)
		case "Icon":
			item.Icon = handleItemIcon(data)
		case "Colour":
			item.Color = handleItemColor(data)
		case "WorldColour":
			item.WorldColor = handleItemColor(data)
		case "Chargeable":
			item.Charegeable = handleBool(data)
		case "ChargeAmount":
			item.CharegeAmount = handleInt(data)
		case "SubstanceCategory":
			item.SubstanceCategory = data.Properties[0].Value
		case "Category":
			item.ProductCategory = data.Properties[0].Value
		case "ChargeBy":
			item.ChargeBy = handleIndirectStringArray(data)
		case "BuildFullyCharged":
			item.BuildFullyCharged = handleBool(data)
		case "Upgrade":
			item.Upgrade = handleBool(data)
		case "Core":
			item.Core = handleBool(data)
		case "TechnologyCategory":
			item.TechnologyCategory = handleIndirectString(data)
		case "TechnologyRarity":
			item.TechnologyRarity = handleIndirectString(data)
		case "Value":
			item.Value = handleInt(data)
		case "BaseStat":
			item.BaseStat = handleIndirectString(data)
		case "StatBonuses":
			item.StatBonuses = handleItemStatBonuses(data)
		case "RequiredTech":
			item.RequiredTech = handleString(data)
		case "RequiredLevel":
			item.RequiredLevel = handleInt(data)
		case "UpgradeColour":
			item.UpgradeColor = handleItemColor(data)
		case "LinkColour":
			item.LinkColor = handleItemColor(data)
		case "RewardGroup":
			item.RewardGroup = handleString(data)
		case "Rarity":
			item.Rarity = data.Properties[0].Value
		case "Legality":
			item.Legality = data.Properties[0].Value
		case "Consumable":
			item.Consumable = handleBool(data)
		case "ChargeValue":
			item.ChargeValue, _ = strconv.ParseInt(data.Value, 10, 64)
		case "Requirements":
			item.Requirements = handleItemRequirements(data)
		case "Cost":
			item.Cost = handleItemCost(data)
		case "RequiredRank":
			item.RequiredRank = handleInt(data)
		case "DispensingRace":
			item.DispensingRace = handleIndirectString(data)
		case "TechShopRarity":
			item.TechShopRarity = handleIndirectString(data)
		case "SpecificChargeOnly":
			item.SpecificChargeOnly = handleBool(data)
		case "NormalisedValueOnWorld":
			item.NormalisedValueOnWorld, _ = strconv.ParseFloat(data.Value, 64)
		case "NormalisedValueOffWorld":
			item.NormalisedValueOffWorld, _ = strconv.ParseFloat(data.Value, 64)

			// TODO(TradingCategory)
		case "WikiEanabled":
			item.WikiEanabled = data.Value == "True"
		case "IsCraftable":
			item.IsCraftable = data.Value == "True"

		case "EconomyInfluenceMultiplier":
			item.EconomyInfluenceMultiplier, _ = strconv.ParseFloat(data.Value, 64)

		}
	}

	items[item.Id] = item
}

func loadItemFile(items map[string]*Item, fileName string, strings map[string]string) {
	data := loadDataFile(fileName)
	for _, data := range data.Properties[0].Properties {
		handleItem(items, data, strings)
	}
}

func handleBuildingObject(items map[string]*Item, entry *Property) {
	var info ItemBuildingInfo
	var name string

	for _, data := range entry.Properties {
		switch data.Name {
		case "ID":
			name = handleString(data)
		case "BuildableOnBase":
			info.BuildableOnBase = handleBool(data)
		case "BuildableOnFreighter":
			info.BuildableOnFreighter = handleBool(data)
		case "BuildableOnPlanet":
			info.BuildableOnPlanet = handleBool(data)
		case "ComplexityCost":
			info.ComplexityCost = handleInt(data)
		case "Group":
			info.Group = handleString(data)
		case "CanChangeColour":
			info.CanChangeColor = handleBool(data)
		case "CanChangeMaterial":
			info.CanChangeMaterial = handleBool(data)
		}
	}

	if item, ok := items[name]; ok {
		item.BuildingInfo = &info
	} else {
		fmt.Errorf("can't find %s in items", name)
	}
}

func handleBuildingObjects(items map[string]*Item, entry *Property) {
	for _, data := range entry.Properties {
		handleBuildingObject(items, data)
	}
}

func loadBuildingFile(items map[string]*Item, fileName string) {
	data := loadDataFile(fileName)
	for _, entry := range data.Properties {
		if entry.Name == "Objects" {
			handleBuildingObjects(items, entry)
		}
	}
}

func main() {
	strings := make(map[string]string)

	loadLocFile(strings, "data/LANGUAGE/NMS_LOC1_USENGLISH.exml")
	loadLocFile(strings, "data/LANGUAGE/NMS_UPDATE3_USENGLISH.exml")

	items := make(map[string]*Item)
	loadItemFile(items,
		"data/METADATA/REALITY/TABLES/NMS_REALITY_GCPRODUCTTABLE.exml",
		strings)
	loadItemFile(items,
		"data/METADATA/REALITY/TABLES/NMS_REALITY_GCSUBSTANCETABLE.exml",
		strings)
	loadItemFile(items,
		"data/METADATA/REALITY/TABLES/NMS_REALITY_GCTECHNOLOGYTABLE.exml",
		strings)

	loadBuildingFile(items,
		"data/METADATA/REALITY/TABLES/BASEBUILDINGTABLE.exml")

	b, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))

}
