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

// Sums all cells in specified group
func (n NSlice) g(g int) float64 {
	var sum float64
	var start int = 13 * g
	for i := start; i < start+13; i++ {
		sum += n[i]
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

	step 									  float64
	entryRateByGroupAndStage				  []([]float64)
	diseaseProgressionExitsByDiseaseStage     []float64
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
	var p CountryInputs


	// ############################################################################################################# 
	// ####################################### Steps 1 and 2: prepare inputs #######################################
	// #############################################################################################################

	
	p.groups = []string { "Gen Pop Men", "Gen Pop Women", "SW Women", "MSM", "IDU" }
	p.diseaseAndTreatmentStages = []string { "Uninfected", "Acute", "Early" , "Med" ,"Late" ,"Adv", "AIDS", "Acute Tx", "Early Tx" , "Med Tx" ,"Late Tx" ,"Adv Tx", "AIDS Tx" }                                         
	p.populationSize = 35067464                                                    
	p.populationSizeByGroup = []float64 { 17390000.1, 16741000.1, 66964.1, 869500.1, 68262.1 }                                               
	p.hivPrevalenceAdultsByGroup = []float64 {  0.061, 0.083, 0.60, 0.039, 0.400 }                                         
	p.hivPrevalence15yoByGroup = []float64 { 0.02, 0.03 }                                           
	p.proprtionDiseaseStage = []float64 { 0.0, 0.05, 0.25, 0.20, 0.20, 0.20, 0.10 }                                               
	p.infectiousnessByDiseaseStage = []float64 { 0.16, 0.08, 0.09, 0.16, 0.5, 0.5 }                                          
	p.hivDeathRateByDiseaseStage = []float64 { 0.0, 0.0, 0.1, 0.2, 0.3, 0.4, 0.45, 0.0, 0.0, 0.02, 0.06, 0.06, 0.08, 0.11 }                                           
	p.hivDeathRateByDiseaseStageTx = []float64 {  }                                        
	p.initialTreatmentAccessByDiseaseStage = []float64 { 0.0, 0.0, 0.0, 0.0, 0.0, 0.2, 0.3 }                               
	p.treatmentRecuitingRateByDiseaseStage = []float64 { 0.0, 0.0, 0.0, 0.0, 0.1, 0.1, 0.1 } 
	p.entryRateGenPop =  0.09
	p.maturationRate =  0.07
	p.deathRateGeneralCauses =  0.01
	p.lifeExpectancy =  55
	p.swInitiationRate =  0.0006
	p.swQuitRate =  0.005
	p.entryRateMsm =  0.00225
	p.entryRateIdu =  0.00036
	p.iduInitiationRate =  0.0001
	p.iduSpontaneousQuitRate =  0.0003
	p.iduDeathRate =  0.05
	p.increaseInInfectiousnessHomosexual =  0.5
	p.diseaseProgressionUntreatedAcuteToEarly =  2
	p.diseaseProgressionUntreatedEarlyToMedium =  0.5
	p.diseaseProgressionUntreatedMediumToLate =  0.2
	p.diseaseProgressionUntreatedLateToAdvanced =  0.2
	p.diseaseProgressionUntreatedAdvancedToAids =  0.5
	p.diseaseProgressionTreatedAcuteToEarly =  0
	p.diseaseProgressionTreatedEarlyToMedium =  0.1
	p.diseaseProgressionTreatedMediumToLate =  0.1
	p.diseaseProgressionTreatedLateToAdvanced =  0.1
	p.diseaseProgressionTreatedAdvancedToAids =  0.1
	p.generalNonSwPartnershipsYearly =  1.5
	p.generalCondomUse =  0.3
	p.generalCondomEffectiveness =  0.9
	p.swProportionWhoUseServices =  0.1
	p.swPartnershipsYearly =  120
	p.swCondomUseRate =  0.66
	p.msmPartnershipsYearly =  3
	p.msmCondomUseRate =  0.49
	p.treatmentReductionOfInfectiousness =  0.95
	p.treatmentQuitRate =  0.05
	p.percentOfIduSexPartners =  0.4
	p.iduPartnershipsYearly =  4.5
	p.iduCondomUseRate =  0.2
	p.annualNumberOfInjections =  264
	p.percentSharedInjections =  0.4
	p.percentMaleIdus =  0.8
	p.infectiousnessInSharedInjection =  0.005
	p.circEffectiveness =  0.6 


	/// new creations

	p.entryRateByGroupAndStage = make([]([]float64), 5)  
	p.entryRateByGroupAndStage[0] = []float64{ p.entryRateGenPop/2 , 0, 0,  0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0  }
	p.entryRateByGroupAndStage[1] = []float64{ p.entryRateGenPop/2 , 0, 0,  0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0  }
	p.entryRateByGroupAndStage[2] = []float64{ 0 , 0, 0,  0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0  }
	p.entryRateByGroupAndStage[3] = []float64{ p.entryRateMsm , 0, 0,  0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0  }
	p.entryRateByGroupAndStage[4] = []float64{ p.entryRateIdu , 0, 0,  0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0  }

	p.diseaseProgressionExitsByDiseaseStage = []float64 { 0, p.diseaseProgressionUntreatedAcuteToEarly, p.diseaseProgressionUntreatedEarlyToMedium, p.diseaseProgressionUntreatedMediumToLate, p.diseaseProgressionUntreatedLateToAdvanced, p.diseaseProgressionUntreatedAdvancedToAids, p.diseaseProgressionTreatedAcuteToEarly, p.diseaseProgressionTreatedEarlyToMedium, p.diseaseProgressionUntreatedMediumToLate, p.diseaseProgressionTreatedLateToAdvanced, p.diseaseProgressionTreatedAdvancedToAids }





	// ############################################################################################################# 
	// ####################################### Step 3: Compute initual pops, vary parameters by group and disease stage #######################################
	// #############################################################################################################


	//calculate initial populations
	for g := 0; g < len(p.groups); g++ {
		for s := 0; s < len(p.diseaseAndTreatmentStages); s++ {
			var i int = g*13 + s
			if s == 0 {
				currentCycle[i] = p.populationSizeByGroup[g] * (1 - p.hivPrevalenceAdultsByGroup[g])
			}
			if s>0 && s<7 {
				currentCycle[i] = p.populationSizeByGroup[g] * p.hivPrevalenceAdultsByGroup[g] * p.proprtionDiseaseStage[s] * (1- p.initialTreatmentAccessByDiseaseStage[s])
			}
			if s>6 {
				var ds int = s - 6
				currentCycle[i] = p.populationSizeByGroup[g] * p.hivPrevalenceAdultsByGroup[g] * p.proprtionDiseaseStage[ds] * (p.initialTreatmentAccessByDiseaseStage[ds])
			}
		} // end disease stage
	} // end group
	previousCycle = currentCycle
	//begin main loop

	var s int = 0
	var g int = 0


	for c := 0; c < 40; c++ {

		// Step 6: develop dynamics


		// Step 6a: general dynamics
		var _ float64 = dGeneral(previousCycle,g,s,p) + dProgEntries(previousCycle,g,s,p) + dProgExits(previousCycle,g,s,p) + dTreatment(previousCycle,g,s,p)


		

		fmt.Println(g,s)
		fmt.Println(previousCycle)

		// Determine which group and stage you are in
		if s == 12 {
			s = 0
			g++
		}




	} //end cycle	







} //end predict



func src_sum(n NSlice,g int,s int) float64 {
	
	//src code here
	return 0
}


func dGeneral(n NSlice,g int,s int, p CountryInputs) float64 {
	return p.step * n.gs(g,s) * (-p.maturationRate - p.deathRateGeneralCauses - p.hivDeathRateByDiseaseStage[s]) + n.sum() * p.entryRateByGroupAndStage[g][s]
}

func dProgExits(n NSlice,g int,s int, p CountryInputs) float64 {
	if s == 0 {
		return p.step * n.gs(g,s) * src_sum(n,g,s)
	} else {
		return  p.step * n.gs(g,s) * p.diseaseProgressionExitsByDiseaseStage[s]
	}
}

func dProgEntries(n NSlice,g int,s int, p CountryInputs) float64 {
	if s == 0 {
		return p.step * n.gs(g,s) * src_sum(n,g,s)
	} else  {
		return p.step * n.gs(g,s) * p.diseaseProgressionExitsByDiseaseStage[s]
	}
}

func dTreatment(n NSlice,g int,s int, p CountryInputs) float64 {
	if s == 0 {
		return 0
	} else  if s > 0 && s < 7 {
		return p.step * n.gs(g,s) * p.treatmentRecuitingRateByDiseaseStage[s] + n.gs(g,s+6) * p.treatmentQuitRate
	} else {
		return p.step * n.gs(g,s-6) * p.treatmentRecuitingRateByDiseaseStage[s] + n.gs(g,s) * -p.treatmentQuitRate
	}
}

func dIduSw(n NSlice,g int,s int, p CountryInputs) float64 {
	if g == 0 {
		return p.step * ( n.gs(4,s) * p.percentMaleIdus * p.iduSpontaneousQuitRate + n.gs(g,s) * - p.iduInitiationRate
	} else if g == 1 {
		return p.step * ( n.gs(4,s) * (1-p.percentMaleIdus) * p.iduSpontaneousQuitRate + n.gs(g,s) * - p.iduInitiationRate + n.gs(3,s) * p.swQuitRate + n.gs(g,s) * - p.swInitiationRate
	} else if g == 2 {
		return n.gs(g,s) * -p.swQuitRate + n.gs(1,s) * p.swInitiationRate
	} else if g == 0 {
		return 0
	} else if g == 4 {
		return n.gs(g,s) * -p.iduSpontaneousQuitRate + 
	}
}



		



func main() {


	Predict()
	fmt.Println("DOne")


}






