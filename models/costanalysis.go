package main


import (

	"fmt"


)


const (
	CYCLES = 40
)

type Inputs struct {
	spendings 		[]Spending // InterventionSubpopulationNationCoverage
	costs     		[]Cost
	countryInputs   CountryInputs
}

type NSlice []float64


// Sums all parameters in passed slice
func (n NSlice) sum() float64 {

	var sum float64
	for i := 0; i < len(n); i++ {
		sum += n[i]
	}
	return sum

}

// Sums all cells in specified disease stage
func (n NSlice) s(s int) float64 {
	var sum float64
	var q = make([]int,5,5)
	q[0], q[1], q[2], q[3], q[4] = s , s + 13, s+13*2, s+13*3, s+13*4
	for i := 0; i < len(q); i++ {
		sum += n[q[i]]
	}
	return sum
}


// Finds cell that matches group and disease stage
func (n NSlice) gs(g int, s int) float64 {
	var i int = g*13 + s
	return n[i]
}


type CountryInputs struct {
	groups                                    []string
	diseaseAndTreatmentStages                 []string
	populationSize                            int
	populationSizeByGroup                     []float64
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
	msmCondomUseRate                          float64
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
	RRR              float64
	RRRTypeID        int
	HIVStatus        int
	RRRStandardError float64
}

type Results struct {
	totalPrevalence                 []float64
	totalNewInfections              []float64
	totalNewInfectionsPerPop        []float64
	cumulativeTotalNewInfections    []float64
	hivDeaths                       []float64
	cumulativeHivDeaths             []float64
	prevalenceByGroup               []float64
	totalPlwa                       []float64
	plwaByGroup                     [][]float64
	incidenceRate                   [][]float64
	totalNewInfectionsByGroup       [][]float64
	totalNewInfectionsByGroupPerPop [][]float64
	hivDeathsByGroup                [][]float64
	totalCostPerIntervention        float64
	totalCostPerComponent           float64
	componentNames                  float64
	totalCost                       float64
	totalPopulation                 []float64
	propOnArt                       []float64
	percentOfTotalPopByGroup        [][]float64
}



func Predict() {

	//var theResults Results
	var currentCycle = make(NSlice, 65, 65)
	var previousCycle = make(NSlice, 65, 65)
	//var secondPreviousCycle []float64
	var theCountryInputs CountryInputs

	// ############################################################################################################# 
	// ####################################### Steps 1 and 2: prepare inputs #######################################
	// #############################################################################################################

	
	theCountryInputs.groups = []string { "Gen Pop Men", "Gen Pop Women", "SW Women", "MSM", "IDU" }
	theCountryInputs.diseaseAndTreatmentStages = []string { "Uninfected", "Acute", "Early" , "Med" ,"Late" ,"Adv", "AIDS", "Acute Tx", "Early Tx" , "Med Tx" ,"Late Tx" ,"Adv Tx", "AIDS Tx" }                                         
	theCountryInputs.populationSize = 35067464                                                    
	theCountryInputs.populationSizeByGroup = []float64 { 17390000.1, 16741000.1, 66964.1, 869500.1, 68262.1 }                                               
	theCountryInputs.hivPrevalenceAdultsByGroup = []float64 {  0.061, 0.083, 0.60, 0.039, 0.400 }                                         
	theCountryInputs.hivPrevalence15yoByGroup = []float64 { 0.02, 0.03 }                                           
	theCountryInputs.proprtionDiseaseStage = []float64 { 0.0, 0.05, 0.25, 0.20, 0.20, 0.20, 0.10 }                                               
	theCountryInputs.infectiousnessByDiseaseStage = []float64 { 0.16, 0.08, 0.09, 0.16, 0.5, 0.5 }                                          
	theCountryInputs.hivDeathRateByDiseaseStage = []float64 { 0.0, 0.0, 0.1, 0.2, 0.3, 0.4, 0.45 }                                           
	theCountryInputs.hivDeathRateByDiseaseStageTx = []float64 { 0.0, 0.0, 0.02, 0.06, 0.06, 0.08, 0.11 }                                        
	theCountryInputs.initialTreatmentAccessByDiseaseStage = []float64 { 0.0, 0.0, 0.0, 0.0, 0.0, 0.2, 0.3 }                               
	theCountryInputs.treatmentRecuitingRateByDiseaseStage = []float64 { 0.0, 0.0, 0.0, 0.0, 0.1, 0.1, 0.1 } 
	theCountryInputs.entryRateGenPop =  0.09
	theCountryInputs.maturationRate =  0.07
	theCountryInputs.deathRateGeneralCauses =  0.01
	theCountryInputs.lifeExpectancy =  55
	theCountryInputs.swInitiationRate =  0.0006
	theCountryInputs.swQuitRate =  0.005
	theCountryInputs.entryRateMsm =  0.00225
	theCountryInputs.entryRateIdu =  0.00036
	theCountryInputs.iduInitiationRate =  0.0001
	theCountryInputs.iduSpontaneousQuitRate =  0.0003
	theCountryInputs.iduDeathRate =  0.05
	theCountryInputs.increaseInInfectiousnessHomosexual =  0.5
	theCountryInputs.diseaseProgressionUntreatedAcuteToEarly =  2
	theCountryInputs.diseaseProgressionUntreatedEarlyToMedium =  0.5
	theCountryInputs.diseaseProgressionUntreatedMediumToLate =  0.2
	theCountryInputs.diseaseProgressionUntreatedLateToAdvanced =  0.2
	theCountryInputs.diseaseProgressionUntreatedAdvancedToAids =  0.5
	theCountryInputs.diseaseProgressionTreatedAcuteToEarly =  0
	theCountryInputs.diseaseProgressionTreatedEarlyToMedium =  0.1
	theCountryInputs.diseaseProgressionTreatedMediumToLate =  0.1
	theCountryInputs.diseaseProgressionTreatedLateToAdvanced =  0.1
	theCountryInputs.diseaseProgressionTreatedAdvancedToAids =  0.1
	theCountryInputs.generalNonSwPartnershipsYearly =  1.5
	theCountryInputs.generalCondomUse =  0.3
	theCountryInputs.generalCondomEffectiveness =  0.9
	theCountryInputs.swProportionWhoUseServices =  0.1
	theCountryInputs.swPartnershipsYearly =  120
	theCountryInputs.swCondomUseRate =  0.66
	theCountryInputs.msmPartnershipsYearly =  3
	theCountryInputs.msmCondomUseRate =  0.49
	theCountryInputs.treatmentReductionOfInfectiousness =  0.95
	theCountryInputs.treatmentQuitRate =  0.05
	theCountryInputs.percentOfIduSexPartners =  0.4
	theCountryInputs.iduPartnershipsYearly =  4.5
	theCountryInputs.iduCondomUseRate =  0.2
	theCountryInputs.annualNumberOfInjections =  264
	theCountryInputs.percentSharedInjections =  0.4
	theCountryInputs.percentMaleIdus =  0.8
	theCountryInputs.infectiousnessInSharedInjection =  0.005
	theCountryInputs.circEffectiveness =  0.6   

	// ############################################################################################################# 
	// ####################################### Step 3: Compute initual pops, vary parameters by group and disease stage #######################################
	// #############################################################################################################


	//calculate initial populations
	for g := 0; g < len(theCountryInputs.groups); g++ {
		for s := 0; s < len(theCountryInputs.diseaseAndTreatmentStages); s++ {
			var i int = g*13 + s
			if s == 0 {
				currentCycle[i] = theCountryInputs.populationSizeByGroup[g] * (1 - theCountryInputs.hivPrevalenceAdultsByGroup[g])
			}
			if s>0 && s<7 {
				currentCycle[i] = theCountryInputs.populationSizeByGroup[g] * theCountryInputs.hivPrevalenceAdultsByGroup[g] * theCountryInputs.proprtionDiseaseStage[s] * (1- theCountryInputs.initialTreatmentAccessByDiseaseStage[s])
			}
			if s>6 {
				var ds int = s - 6
				currentCycle[i] = theCountryInputs.populationSizeByGroup[g] * theCountryInputs.hivPrevalenceAdultsByGroup[g] * theCountryInputs.proprtionDiseaseStage[ds] * (theCountryInputs.initialTreatmentAccessByDiseaseStage[ds])
			}
		} // end disease stage
	} // end group
	previousCycle = currentCycle
	//begin main loop
} //end predict




func main() {


	Predict()
	fmt.Println("DOne")


}






