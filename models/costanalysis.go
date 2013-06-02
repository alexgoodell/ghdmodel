package costanalysis

import (
	"fmt"
)

const (
	CYCLES = 40
)

type Inputs struct {
	Spendings 		[]Spending // InterventionSubpopulationNationCoverage
	Costs     		[]Cost
	CountryProfile  CountryProfile
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

type CountryProfile struct {
	Groups                                    []string
	DiseaseAndTreatmentStages                 []string
	PopulationSize                            int
	PopulationSizeByGroup                     []float64
	HivPrevalenceAdultsByGroup                []float64
	HivPrevalence15yoByGroup                  []float64
	ProprtionDiseaseStage                     []float64
	InfectiousnessByDiseaseStage              []float64
	HivDeathRateByDiseaseStage                []float64
	HivDeathRateByDiseaseStageTx              []float64
	InitialTreatmentAccessByDiseaseStage      []float64
	TreatmentRecuitingRateByDiseaseStage      []float64
	EntryRateGenPop                           float64
	MaturationRate                            float64
	DeathRateGeneralCauses                    float64
	LifeExpectancy                            float64
	SwInitiationRate                          float64
	SwQuitRate                                float64
	EntryRateMsm                              float64
	EntryRateIdu                              float64
	IduInitiationRate                         float64
	IduSpontaneousQuitRate                    float64
	IduDeathRate                              float64
	IncreaseInInfectiousnessHomosexual        float64
	DiseaseProgressionUntreatedAcuteToEarly   float64
	DiseaseProgressionUntreatedEarlyToMedium  float64
	DiseaseProgressionUntreatedMediumToLate   float64
	DiseaseProgressionUntreatedLateToAdvanced float64
	DiseaseProgressionUntreatedAdvancedToAids float64
	DiseaseProgressionTreatedAcuteToEarly     float64
	DiseaseProgressionTreatedEarlyToMedium    float64
	DiseaseProgressionTreatedMediumToLate     float64
	DiseaseProgressionTreatedLateToAdvanced   float64
	DiseaseProgressionTreatedAdvancedToAids   float64
	GeneralNonSwPartnershipsYearly            float64
	GeneralCondomUse                          float64
	GeneralCondomEffectiveness                float64
	SwProportionWhoUseServices                float64
	SwPartnershipsYearly                      float64
	SwCondomUseRate                           float64
	MsmPartnershipsYearly                     float64
	MsmCondomUseRate                          float64
	TreatmentReductionOfInfectiousness        float64
	TreatmentQuitRate                         float64
	PercentOfIduSexPartners                   float64
	IduPartnershipsYearly                     float64
	IduCondomUseRate                          float64
	AnnualNumberOfInjections                  float64
	PercentSharedInjections                   float64
	PercentMaleIdus                           float64
	InfectiousnessInSharedInjection           float64
	CircEffectiveness                         float64

	Step 									  float64
	EntryRateByGroupAndStage				  []([]float64)
	DiseaseProgressionExitsByDiseaseStage     []float64
}



type Cost struct {
	Id            int
	NationID      int
	ComponentID   int
	CostPerClient float64
	ComponentName string
}

type Spending struct {
	Id               int
	InterventionID   int
	SubpopulationID  int
	NationID         int
	Coverage         float64
	RRR              float64
	RRRTypeID        int
	HIVStatus        int
	RRRStandardError float64
}

type Results struct {
	TotalPrevalence                 []float64
	TotalNewInfections              []float64
	TotalNewInfectionsPerPop        []float64
	CumulativeTotalNewInfections    []float64
	HivDeaths                       []float64
	CumulativeHivDeaths             []float64
	PrevalenceByGroup               []float64
	TotalPlwa                       []float64
	PlwaByGroup                     [][]float64
	IncidenceRate                   [][]float64
	TotalNewInfectionsByGroup       [][]float64
	TotalNewInfectionsByGroupPerPop [][]float64
	HivDeathsByGroup                [][]float64
	TotalCostPerIntervention        float64
	TotalCostPerComponent           float64
	ComponentNames                  float64
	TotalCost                       float64
	TotalPopulation                 []float64
	PropOnArt                       []float64
	PercentOfTotalPopByGroup        [][]float64
}



func Predict() {

	//var theResults Results
	var currentCycle = make(NSlice, 65, 65)
	var previousCycle = make(NSlice, 65, 65)
	//var secondPreviousCycle []float64
	var p CountryProfile


	// ############################################################################################################# 
	// ####################################### Steps 1 and 2: prepare inputs #######################################
	// #############################################################################################################

	
	p.Groups = []string { "Gen Pop Men", "Gen Pop Women", "SW Women", "MSM", "IDU" }
	p.DiseaseAndTreatmentStages = []string { "Uninfected", "Acute", "Early" , "Med" ,"Late" ,"Adv", "AIDS", "Acute Tx", "Early Tx" , "Med Tx" ,"Late Tx" ,"Adv Tx", "AIDS Tx" }                                         
	p.PopulationSize = 35067464                                                    
	p.PopulationSizeByGroup = []float64 { 17390000.1, 16741000.1, 66964.1, 869500.1, 68262.1 }                                               
	p.HivPrevalenceAdultsByGroup = []float64 {  0.061, 0.083, 0.60, 0.039, 0.400 }                                         
	p.HivPrevalence15yoByGroup = []float64 { 0.02, 0.03 }                                           
	p.ProprtionDiseaseStage = []float64 { 0.0, 0.05, 0.25, 0.20, 0.20, 0.20, 0.10 }                                               
	p.InfectiousnessByDiseaseStage = []float64 { 0.16, 0.08, 0.09, 0.16, 0.5, 0.5 }                                          
	p.HivDeathRateByDiseaseStage = []float64 { 0.0, 0.0, 0.1, 0.2, 0.3, 0.4, 0.45, 0.0, 0.0, 0.02, 0.06, 0.06, 0.08, 0.11 }                                           
	p.HivDeathRateByDiseaseStageTx = []float64 {  }                                        
	p.InitialTreatmentAccessByDiseaseStage = []float64 { 0.0, 0.0, 0.0, 0.0, 0.0, 0.2, 0.3 }                               
	p.TreatmentRecuitingRateByDiseaseStage = []float64 { 0.0, 0.0, 0.0, 0.0, 0.1, 0.1, 0.1 } 
	p.EntryRateGenPop =  0.09
	p.MaturationRate =  0.07
	p.DeathRateGeneralCauses =  0.01
	p.LifeExpectancy =  55
	p.SwInitiationRate =  0.0006
	p.SwQuitRate =  0.005
	p.EntryRateMsm =  0.00225
	p.EntryRateIdu =  0.00036
	p.IduInitiationRate =  0.0001
	p.IduSpontaneousQuitRate =  0.0003
	p.IduDeathRate =  0.05
	p.IncreaseInInfectiousnessHomosexual =  0.5
	p.DiseaseProgressionUntreatedAcuteToEarly =  2
	p.DiseaseProgressionUntreatedEarlyToMedium =  0.5
	p.DiseaseProgressionUntreatedMediumToLate =  0.2
	p.DiseaseProgressionUntreatedLateToAdvanced =  0.2
	p.DiseaseProgressionUntreatedAdvancedToAids =  0.5
	p.DiseaseProgressionTreatedAcuteToEarly =  0
	p.DiseaseProgressionTreatedEarlyToMedium =  0.1
	p.DiseaseProgressionTreatedMediumToLate =  0.1
	p.DiseaseProgressionTreatedLateToAdvanced =  0.1
	p.DiseaseProgressionTreatedAdvancedToAids =  0.1
	p.GeneralNonSwPartnershipsYearly =  1.5
	p.GeneralCondomUse =  0.3
	p.GeneralCondomEffectiveness =  0.9
	p.SwProportionWhoUseServices =  0.1
	p.SwPartnershipsYearly =  120
	p.SwCondomUseRate =  0.66
	p.MsmPartnershipsYearly =  3
	p.MsmCondomUseRate =  0.49
	p.TreatmentReductionOfInfectiousness =  0.95
	p.TreatmentQuitRate =  0.05
	p.PercentOfIduSexPartners =  0.4
	p.IduPartnershipsYearly =  4.5
	p.IduCondomUseRate =  0.2
	p.AnnualNumberOfInjections =  264
	p.PercentSharedInjections =  0.4
	p.PercentMaleIdus =  0.8
	p.InfectiousnessInSharedInjection =  0.005
	p.CircEffectiveness =  0.6 


	/// new creations

	p.EntryRateByGroupAndStage = make([]([]float64), 5)  
	p.EntryRateByGroupAndStage[0] = []float64{ p.EntryRateGenPop/2 , 0, 0,  0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0  }
	p.EntryRateByGroupAndStage[1] = []float64{ p.EntryRateGenPop/2 , 0, 0,  0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0  }
	p.EntryRateByGroupAndStage[2] = []float64{ 0 , 0, 0,  0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0  }
	p.EntryRateByGroupAndStage[3] = []float64{ p.EntryRateMsm , 0, 0,  0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0  }
	p.EntryRateByGroupAndStage[4] = []float64{ p.EntryRateIdu , 0, 0,  0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0  }

	p.DiseaseProgressionExitsByDiseaseStage = []float64 { 0, p.DiseaseProgressionUntreatedAcuteToEarly, p.DiseaseProgressionUntreatedEarlyToMedium, p.DiseaseProgressionUntreatedMediumToLate, p.DiseaseProgressionUntreatedLateToAdvanced, p.DiseaseProgressionUntreatedAdvancedToAids, p.DiseaseProgressionTreatedAcuteToEarly, p.DiseaseProgressionTreatedEarlyToMedium, p.DiseaseProgressionUntreatedMediumToLate, p.DiseaseProgressionTreatedLateToAdvanced, p.DiseaseProgressionTreatedAdvancedToAids }





	// ############################################################################################################# 
	// ####################################### Step 3: Compute initual pops, vary parameters by group and disease stage #######################################
	// #############################################################################################################


	//calculate initial populations
	for g := 0; g < len(p.Groups); g++ {
		for s := 0; s < len(p.DiseaseAndTreatmentStages); s++ {
			var i int = g*13 + s
			if s == 0 {
				currentCycle[i] = p.PopulationSizeByGroup[g] * (1 - p.HivPrevalenceAdultsByGroup[g])
			}
			if s>0 && s<7 {
				currentCycle[i] = p.PopulationSizeByGroup[g] * p.HivPrevalenceAdultsByGroup[g] * p.ProprtionDiseaseStage[s] * (1- p.InitialTreatmentAccessByDiseaseStage[s])
			}
			if s>6 {
				var ds int = s - 6
				currentCycle[i] = p.PopulationSizeByGroup[g] * p.HivPrevalenceAdultsByGroup[g] * p.ProprtionDiseaseStage[ds] * (p.InitialTreatmentAccessByDiseaseStage[ds])
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



func srcSum(n NSlice,g int,s int) float64 {
	
	//src code here
	return 0
}


func dGeneral(n NSlice,g int,s int, p CountryProfile) float64 {
	return p.Step * n.gs(g,s) * (-p.MaturationRate - p.DeathRateGeneralCauses - p.HivDeathRateByDiseaseStage[s]) + n.sum() * p.EntryRateByGroupAndStage[g][s]
}

func dProgExits(n NSlice,g int,s int, p CountryProfile) float64 {
	if s == 0 {
		return p.Step * n.gs(g,s) * srcSum(n,g,s)
	} else {
		return  p.Step * n.gs(g,s) * p.DiseaseProgressionExitsByDiseaseStage[s]
	}
}

func dProgEntries(n NSlice,g int,s int, p CountryProfile) float64 {
	if s == 0 {
		return p.Step * n.gs(g,s) * srcSum(n,g,s)
	} else  {
		return p.Step * n.gs(g,s) * p.DiseaseProgressionExitsByDiseaseStage[s]
	}
}

func dTreatment(n NSlice,g int,s int, p CountryProfile) float64 {
	if s == 0 {
		return 0
	} else  if s > 0 && s < 7 {
		return p.Step * n.gs(g,s) * p.TreatmentRecuitingRateByDiseaseStage[s] + n.gs(g,s+6) * p.TreatmentQuitRate
	} else {
		return p.Step * n.gs(g,s-6) * p.TreatmentRecuitingRateByDiseaseStage[s] + n.gs(g,s) * -p.TreatmentQuitRate
	}
}

func dIduSw(n NSlice,g int,s int, p CountryProfile) float64 {
	if g == 0 {
		return p.Step * ( n.gs(4,s) * p.PercentMaleIdus * p.IduSpontaneousQuitRate + n.gs(g,s) * - p.IduInitiationRate
	} else if g == 1 {
		return p.Step * ( n.gs(4,s) * (1-p.PercentMaleIdus) * p.IduSpontaneousQuitRate + n.gs(g,s) * - p.IduInitiationRate + n.gs(3,s) * p.SwQuitRate + n.gs(g,s) * - p.SwInitiationRate
	} else if g == 2 {
		return n.gs(g,s) * -p.SwQuitRate + n.gs(1,s) * p.SwInitiationRate
	} else if g == 0 {
		return 0
	} else if g == 4 {
		return n.gs(g,s) * -p.IduSpontaneousQuitRate + 
	}
}



		



func main() {


	Predict()
	fmt.Println("DOne")
}