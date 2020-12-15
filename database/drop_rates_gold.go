package database

import "github.com/ferdoran/go-sro-framework/db"

type DropRatesGold struct {
	MobLvl int
	Prob   float64
	Min    int
	Max    int
}

const (
	select_drop_rates_gold string = "SELECT MOB_LVL, PROB, MIN, MAX FROM DROP_RATES_GOLD"
)

func GetGoldDropRates() map[int]DropRatesGold {
	conn := db.OpenConnShard()
	defer conn.Close()

	queryHandle, err := conn.Query(select_drop_rates_gold)
	if err != nil {
		panic(err.Error())
	}

	var rates map[int]DropRatesGold
	for queryHandle.Next() {
		var mobLvl, min, max int
		var prob float64
		err = queryHandle.Scan(&mobLvl, &prob, &min, &max)
		if err != nil {
			panic(err.Error())
		}
		rate := DropRatesGold{
			MobLvl: mobLvl,
			Prob:   prob,
			Min:    min,
			Max:    max,
		}

		rates[rate.MobLvl] = rate
	}

	return rates
}
