package nmsmine

import (
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

type ExmlLoader struct {
	Db      ItemDb
	strings map[string]string
}

func NewExmlLoader() *ExmlLoader {
	return &ExmlLoader{
		Db:      NewItemDb(),
		strings: make(map[string]string),
	}
}

func (e *ExmlLoader) loadDataFile(fileName string) *Data {
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

func (e *ExmlLoader) handleString(entry *Property) string {
	return entry.Value
}

func (e *ExmlLoader) handleIndirectString(entry *Property) string {
	return entry.Properties[0].Value
}

func (e *ExmlLoader) handleIndirectStringArray(entry *Property) []string {
	var strings []string
	for _, data := range entry.Properties {
		strings = append(strings, e.handleIndirectString(data))
	}
	return strings
}

func (e *ExmlLoader) handleBool(entry *Property) bool {
	return entry.Value == "True"
}

func (e *ExmlLoader) handleInt(entry *Property) int64 {
	val, _ := strconv.ParseInt(entry.Value, 10, 64)
	return val
}

func (e *ExmlLoader) handleFloat(entry *Property) float64 {
	val, _ := strconv.ParseFloat(entry.Value, 64)
	return val
}

func (e *ExmlLoader) handleLoc(entry *Property) {
	var name string
	var value string

	for _, data := range entry.Properties {
		switch data.Name {
		case "Id":
			name = e.handleString(data)
		case "USEnglish":
			value = e.handleIndirectString(data)
		}
	}

	if name != "" && value != "" {
		e.strings[name] = value
	}
}

func (e *ExmlLoader) LoadLocFile(fileName string) {
	data := e.loadDataFile(fileName)
	for _, data := range data.Properties[0].Properties {
		e.handleLoc(data)
	}
}

func (e *ExmlLoader) handleItemIcon(entry *Property) string {
	path := entry.Properties[0].Value
	path = strings.TrimPrefix(path, "TEXTURES/UI/FRONTEND/ICONS/")
	path = strings.TrimSuffix(path, ".DDS")
	path = strings.ToLower(path)
	return path
}

func (e *ExmlLoader) handleItemCost(entry *Property) ItemCost {
	var cost ItemCost
	for _, data := range entry.Properties {
		switch data.Name {
		case "SpaceStationMarkup":
			cost.SpaceStationMarkup = e.handleFloat(data)
		case "LowPriceMod":
			cost.LowPriceMod = e.handleFloat(data)
		case "HighPriceMod":
			cost.HighPriceMod = e.handleFloat(data)
		case "BuyBaseMarkup":
			cost.BuyBaseMarkup = e.handleFloat(data)
		case "BuyMarkupMod":
			cost.BuyMarkupMod = e.handleFloat(data)
		}
	}
	return cost
}

func (e *ExmlLoader) handleItemColor(entry *Property) string {
	var color ItemColor
	for _, data := range entry.Properties {
		switch data.Name {
		case "R":
			color.R = e.handleFloat(data)
		case "G":
			color.G = e.handleFloat(data)
		case "B":
			color.B = e.handleFloat(data)
		case "A":
			color.A = e.handleFloat(data)
		}
	}
	return fmt.Sprintf("#%02x%02x%02x",
		uint8(color.R*255), uint8(color.G*255), uint8(color.B*255))
}

func (e *ExmlLoader) handleItemRequirements(entry *Property) []ItemRequirement {
	var reqs []ItemRequirement
	for _, line := range entry.Properties {
		var req ItemRequirement
		for _, data := range line.Properties {
			switch data.Name {
			case "ID":
				req.Id = e.handleString(data)
			case "Amount":
				req.Quantity = e.handleInt(data)
			}
		}
		reqs = append(reqs, req)
	}
	return reqs
}

func (e *ExmlLoader) handleItemStatBonus(entry *Property) ItemStatBonus {
	var bonus ItemStatBonus
	for _, data := range entry.Properties {
		switch data.Name {
		case "StatsTypes":
			bonus.Type = e.handleIndirectString(data)
		case "Bonus":
			bonus.Bonus = e.handleInt(data)
		case "Level":
			bonus.Level = e.handleInt(data)
		}
	}
	return bonus
}

func (e *ExmlLoader) handleItemStatBonuses(entry *Property) []ItemStatBonus {
	var bonuses []ItemStatBonus
	for _, data := range entry.Properties {
		bonuses = append(bonuses, e.handleItemStatBonus(data))
	}
	return bonuses
}

func (e *ExmlLoader) handleItem(entry *Property) {
	item := &Item{}

	for _, data := range entry.Properties {
		switch data.Name {
		case "ID":
			item.Id = e.handleString(data)
		case "Id":
			item.Id = e.handleString(data)
		case "Name":
			item.Name = e.strings[e.handleString(data)]
		case "NameLower":
			item.NameLower = e.strings[e.handleString(data)]
		case "Symbol":
			item.Symbol = e.strings[e.handleString(data)]
		case "Subtitle":
			item.Subtitle = e.strings[e.handleIndirectString(data)]
		case "Description":
			item.Description = e.strings[e.handleIndirectString(data)]
		case "Teach":
			item.Teach = e.handleBool(data)
		case "Hint":
			item.Hint = e.strings[e.handleString(data)]
		case "BaseValue":
			item.BaseValue = e.handleFloat(data)
		case "Level":
			item.Level = e.handleInt(data)
		case "Icon":
			item.Icon = e.handleItemIcon(data)
		case "Colour":
			item.Color = e.handleItemColor(data)
		case "WorldColour":
			item.WorldColor = e.handleItemColor(data)
		case "Chargeable":
			item.Charegeable = e.handleBool(data)
		case "ChargeAmount":
			item.CharegeAmount = e.handleInt(data)
		case "SubstanceCategory":
			item.SubstanceCategory = e.handleIndirectString(data)
		case "Category":
			item.ProductCategory = e.handleIndirectString(data)
		case "ChargeBy":
			item.ChargeBy = e.handleIndirectStringArray(data)
		case "BuildFullyCharged":
			item.BuildFullyCharged = e.handleBool(data)
		case "Upgrade":
			item.Upgrade = e.handleBool(data)
		case "Core":
			item.Core = e.handleBool(data)
		case "TechnologyCategory":
			item.TechnologyCategory = e.handleIndirectString(data)
		case "TechnologyRarity":
			item.TechnologyRarity = e.handleIndirectString(data)
		case "Value":
			item.Value = e.handleInt(data)
		case "BaseStat":
			item.BaseStat = e.handleIndirectString(data)
		case "StatBonuses":
			item.StatBonuses = e.handleItemStatBonuses(data)
		case "RequiredTech":
			item.RequiredTech = e.handleString(data)
		case "RequiredLevel":
			item.RequiredLevel = e.handleInt(data)
		case "UpgradeColour":
			item.UpgradeColor = e.handleItemColor(data)
		case "LinkColour":
			item.LinkColor = e.handleItemColor(data)
		case "RewardGroup":
			item.RewardGroup = e.handleString(data)
		case "Rarity":
			item.Rarity = e.handleIndirectString(data)
		case "Legality":
			item.Legality = e.handleIndirectString(data)
		case "Consumable":
			item.Consumable = e.handleBool(data)
		case "ChargeValue":
			item.ChargeValue = e.handleInt(data)
		case "Requirements":
			item.Requirements = e.handleItemRequirements(data)
		case "Cost":
			item.Cost = e.handleItemCost(data)
		case "RequiredRank":
			item.RequiredRank = e.handleInt(data)
		case "DispensingRace":
			item.DispensingRace = e.handleIndirectString(data)
		case "TechShopRarity":
			item.TechShopRarity = e.handleIndirectString(data)
		case "SpecificChargeOnly":
			item.SpecificChargeOnly = e.handleBool(data)
		case "NormalisedValueOnWorld":
			item.NormalisedValueOnWorld = e.handleFloat(data)
		case "NormalisedValueOffWorld":
			item.NormalisedValueOffWorld = e.handleFloat(data)

			// TODO(TradingCategory)
		case "WikiEanabled":
			item.WikiEanabled = e.handleBool(data)
		case "IsCraftable":
			item.IsCraftable = e.handleBool(data)

		case "EconomyInfluenceMultiplier":
			item.EconomyInfluenceMultiplier = e.handleFloat(data)
		}
	}

	e.Db[item.Id] = item
}

func (e *ExmlLoader) LoadItemFile(fileName string) {
	data := e.loadDataFile(fileName)
	for _, data := range data.Properties[0].Properties {
		e.handleItem(data)
	}
}

func (e *ExmlLoader) handleBuildingObject(entry *Property) {
	var info ItemBuildingInfo
	var name string

	for _, data := range entry.Properties {
		switch data.Name {
		case "ID":
			name = e.handleString(data)
		case "BuildableOnBase":
			info.BuildableOnBase = e.handleBool(data)
		case "BuildableOnFreighter":
			info.BuildableOnFreighter = e.handleBool(data)
		case "BuildableOnPlanet":
			info.BuildableOnPlanet = e.handleBool(data)
		case "ComplexityCost":
			info.ComplexityCost = e.handleInt(data)
		case "Group":
			info.Group = e.handleString(data)
		case "CanChangeColour":
			info.CanChangeColor = e.handleBool(data)
		case "CanChangeMaterial":
			info.CanChangeMaterial = e.handleBool(data)
		}
	}

	if item, ok := e.Db[name]; ok {
		item.BuildingInfo = &info
	} else {
		fmt.Errorf("can't find %s in items", name)
	}
}

func (e *ExmlLoader) handleBuildingObjects(entry *Property) {
	for _, data := range entry.Properties {
		e.handleBuildingObject(data)
	}
}

func (e *ExmlLoader) LoadBuildingFile(fileName string) {
	data := e.loadDataFile(fileName)
	for _, entry := range data.Properties {
		if entry.Name == "Objects" {
			e.handleBuildingObjects(entry)
		}
	}
}
