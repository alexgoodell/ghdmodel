//test comment
package costanalysis

// import (
// 	"fmt"
// )

const (
	CYCLES = 40
)

type Inputs struct {
	spendings []Spending // InterventionSubpopulationNationCoverage
	costs     []Cost
	custom    CountryStats
}

type CountryStats struct {
	// the next too seem redundant
	diseaseStages                             []string
	diseaseAndTreatmentStages                 []string
	populationSize                            int
	populationSizeByGroup                     []int
	hivPrevalenceAdultsByGroup                []float64
	hivPrevalence15yoByGroup                  []float64
	proprtionDiseaseStage                     []float64
	infectiousnessByDiseaseStage              []float64
	hivDeathRateByDiseaseStage                []float64
	hivDeathRateByDiseaseStageTx              []float64
	initialTreatmentAccessByDiseaseStage      []float64
	treatmentRecuitingRateByDiseaseStage      []float64
	entryRateGenPop                           float64
	maturationRate                            float64
	deathRateGeneralCauses                    float64
	lifeExpectancy                            float64
	swInitiationRate                          float64
	swQuitRate                                float64
	entryRateMsm                              float64
	entryRateIdu                              float64
	iduInitiationRate                         float64
	iduSpontaneousQuitRate                    float64
	iduDeathRate                              float64
	increaseInInfectiousnessHomosexual        float64
	diseaseProgressionUntreatedAcuteToEarly   float64
	diseaseProgressionUntreatedEarlyToMedium  float64
	diseaseProgressionUntreatedMediumToLate   float64
	diseaseProgressionUntreatedLateToAdvanced float64
	diseaseProgressionUntreatedAdvancedToAids float64
	diseaseProgressionTreatedAcuteToEarly     float64
	diseaseProgressionTreatedEarlyToMedium    float64
	diseaseProgressionTreatedMediumToLate     float64
	diseaseProgressionTreatedLateToAdvanced   float64
	diseaseProgressionTreatedAdvancedToAids   float64
	generalNonSwPartnershipsYearly            float64
	generalCondomUse                          float64
	generalCondomEffectiveness                float64
	swProportionWhoUseServices                float64
	swPartnershipsYearly                      float64
	swCondomUseRate                           float64
	msmPartnershipsYearly                     float64
	smCondomUseRate                          float64
	treatmentReductionOfInfectiousness        float64
	treatmentQuitRate                         float64
	percentOfIduSexPartners                   float64
	iduPartnershipsYearly                     float64
	iduCondomUseRate                          float64
	annualNumberOfInjections                  float64
	percentSharedInjections                   float64
	percentMaleIdus                           float64
	infectiousnessInSharedInjection           float64
	circEffectiveness                         float64
}

type Cost struct {
	id            int
	nationID      int
	componentID   int
	costPerClient float64
	componentName string
}

type Spending struct {
	id               int
	interventionID   int
	subpopulationID  int
	nationID         int
	coverage         float64
	rrr              float64
	rrrTypeID        int
	hivStatus        int
	rrrStandardError float64
}

type Results struct {
	prevalenceByGroup               [][]int
	totalPlwa                       []int
	plwaByGroup                     [][]int
	incidenceRate                   [][]int
	totalNewInfectionsByGroup       [][]int
	totalPrevalence                 []int
	totalNewInfections              []int
	totalNewInfectionsPerPop        []int
	cumulativeTotalNewInfections    []int
	hivDeaths                       []int
	cumulativeHivDeathsfectionsByGroup   [][]int
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

func Predict(inputs *Inputs) *Inputs {
	return inputs;
}
