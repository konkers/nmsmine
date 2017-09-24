package nmsmine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

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

type ItemDb map[string]*Item

func NewItemDb() ItemDb {
	return make(ItemDb)
}

func LoadItemDb(filename string) (ItemDb, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Error reading from %s: %s\n", filename, err.Error())
	}

	db := NewItemDb()
	err = json.Unmarshal(b, &db)
	if err != nil {
		return nil, fmt.Errorf("Error decoding JSON: %s\n", err.Error())
	}
	return db, nil
}

func (db ItemDb) WriteToFile(filename string) error {
	b, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return fmt.Errorf("Error encoding JSON: %s\n", err.Error())
	}

	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		return fmt.Errorf("Error writing to %s: %s\n", filename, err.Error())
	}

	return nil
}
