// test comment
package costanalysis

import (
	"fmt"
)

const (
	CYCLES = 40
)

type Inputs struct {
	spendings []Spending // InterventionSubpopulationNationCoverage
	costs     []Cost
	custom    Custom
}

// CountryStats?
type Custom struct {
	groups []string
	// the next too seem redundant
	diseaseStages                             []string
	diseaseAndTreatmentStages                 []string
	populationSize                            int
	populationSizeByGroup                     []int
	HivPrevalenceAdultsByGroup                []float
	HivPrevalence15yoByGroup                  []float
	ProprtionDiseaseStage                     []float
	InfectiousnessByDiseaseStage              []float
	HivDeathRateByDiseaseStage                []float
	HivDeathRateByDiseaseStageTx              []float
	InitialTreatmentAccessByDiseaseStage      []float
	TreatmentRecuitingRateByDiseaseStage      []float
	EntryRateGenPop                           float
	MaturationRate                            float
	DeathRateGeneralCauses                    float
	LifeExpectancy                            float
	SwInitiationRate                          float
	SwQuitRate                                float
	EntryRateMsm                              float
	EntryRateIdu                              float
	IduInitiationRate                         float
	IduSpontaneousQuitRate                    float
	IduDeathRate                              float
	IncreaseInInfectiousnessHomosexual        float
	DiseaseProgressionUntreatedAcuteToEarly   float
	DiseaseProgressionUntreatedEarlyToMedium  float
	DiseaseProgressionUntreatedMediumToLate   float
	DiseaseProgressionUntreatedLateToAdvanced float
	DiseaseProgressionUntreatedAdvancedToAids float
	DiseaseProgressionTreatedAcuteToEarly     float
	DiseaseProgressionTreatedEarlyToMedium    float
	DiseaseProgressionTreatedMediumToLate     float
	DiseaseProgressionTreatedLateToAdvanced   float
	DiseaseProgressionTreatedAdvancedToAids   float
	GeneralNonSwPartnershipsYearly            float
	GeneralCondomUse                          float
	GeneralCondomEffectiveness                float
	SwProportionWhoUseServices                float
	SwPartnershipsYearly                      float
	SwCondomUseRate                           float
	MsmPartnershipsYearly                     float
	MsmCondomUseRate                          float
	TreatmentReductionOfInfectiousness        float
	TreatmentQuitRate                         float
	PercentOfIduSexPartners                   float
	IduPartnershipsYearly                     float
	IduCondomUseRate                          float
	AnnualNumberOfInjections                  float
	PercentSharedInjections                   float
	PercentMaleIdus                           float
	InfectiousnessInSharedInjection           float
	CircEffectiveness                         float
}

type Cost struct {
	id            int
	nationID      int
	componentID   int
	costPerClient float
	componentName string
}

type Spending struct {
	id               int
	interventionID   int
	subpopulationID  int
	nationID         int
	coverage         float
	RRR              float
	RRRTypeID        int
	HIVStatus        int
	RRRStandardError float
}

type Results struct {
	totalPrevalence                 []int
	totalNewInfections              []int
	totalNewInfectionsPerPop        []int
	cumulativeTotalNewInfections    []int
	hivDeaths                       []int
	cumulativeHivDeaths             []int
	prevalenceByGroup               []int
	prevalenceByGroup               []int
	prevalenceByGroup               []int
	prevalenceByGroup               []int
	prevalenceByGroup               []int
	prevalenceByGroup               []int
	totalPlwa                       []int
	plwaByGroup                     [][]int
	incidenceRate                   [][]int
	totalNewInfectionsByGroup       [][]int
	totalNewInfectionsByGroupPerPop [][]int
	hivDeathsByGroup                [][]int
	totalCostPerIntervention        int
	totalCostPerComponent           int
	componentNames                  int
	totalCost                       int
	totalPopulation                 []int
	propOnArt                       []int
	percentOfTotalPopByGroup        [][]int
}

func Predict(inputs *Inputs) {

}
