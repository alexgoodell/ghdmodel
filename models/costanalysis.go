package costanalysis

import (
	"fmt"
)

const (
	NUMCYCLES = 40
)

// Main data input struct
type Inputs struct {
	Spendings      []Spending // InterventionSubpopulationNationCoverage
	Costs          []Cost
	CountryProfile CountryProfile
}

// Step populations slice
type NSlice []float32

// Sums all parameters in passed slice
func (n NSlice) sum() float32 {
	var sum float32
	for i := 0; i < len(n); i++ {
		sum += n[i]
	}
	return sum
}

// Sums all cells in specified disease stage
func (n NSlice) s(s int) float32 {
	var sum float32
	var q = make([]int, 5, 5)
	q[0], q[1], q[2], q[3], q[4] = s, s+13, s+13*2, s+13*3, s+13*4
	for i := 0; i < len(q); i++ {
		sum += n[q[i]]
	}
	return sum
}

// Sums all cells in specified group
func (n NSlice) g(g int) float32 {
	var sum float32
	var start int = 13 * g
	for i := start; i < start+13; i++ {
		sum += n[i]
	}
	return sum
}

// Finds cell that matches group and disease stage
func (n NSlice) gs(g int, s int) float32 {
	var i int = g*13 + s
	return n[i]
}

type CountryProfile struct {
	Groups                                    []string
	DiseaseStages                             []string
	PopulationSize                            int
	PopulationSizeByGroup                     []float32
	HivPrevalenceAdultsByGroup                []float32
	HivPrevalence15yoByGroup                  []float32
	ProprtionDiseaseStage                     []float32
	InfectiousnessByDiseaseStage              []float32
	HivDeathRateByDiseaseStage                []float32
	HivDeathRateByDiseaseStageTx              []float32
	InitialTreatmentAccessByDiseaseStage      []float32
	TreatmentRecuitingRateByDiseaseStage      []float32
	EntryRateGenPop                           float32
	MaturationRate                            float32
	DeathRateGeneralCauses                    float32
	LifeExpectancy                            float32
	SwInitiationRate                          float32
	SwQuitRate                                float32
	EntryRateMsm                              float32
	EntryRateIdu                              float32
	IduInitiationRate                         float32
	IduSpontaneousQuitRate                    float32
	IduDeathRate                              float32
	IncreaseInInfectiousnessHomosexual        float32
	DiseaseProgressionUntreatedAcuteToEarly   float32
	DiseaseProgressionUntreatedEarlyToMedium  float32
	DiseaseProgressionUntreatedMediumToLate   float32
	DiseaseProgressionUntreatedLateToAdvanced float32
	DiseaseProgressionUntreatedAdvancedToAids float32
	DiseaseProgressionTreatedAcuteToEarly     float32
	DiseaseProgressionTreatedEarlyToMedium    float32
	DiseaseProgressionTreatedMediumToLate     float32
	DiseaseProgressionTreatedLateToAdvanced   float32
	DiseaseProgressionTreatedAdvancedToAids   float32
	GeneralNonSwPartnershipsYearly            float32
	GeneralCondomUse                          float32
	GeneralCondomEffectiveness                float32
	SwProportionWhoUseServices                float32
	SwPartnershipsYearly                      float32
	SwCondomUseRate                           float32
	MsmPartnershipsYearly                     float32
	MsmCondomUseRate                          float32
	TreatmentReductionOfInfectiousness        float32
	TreatmentQuitRate                         float32
	PercentOfIduSexPartners                   float32
	IduPartnershipsYearly                     float32
	IduCondomUseRate                          float32
	AnnualNumberOfInjections                  float32
	PercentSharedInjections                   float32
	PercentMaleIdus                           float32
	InfectiousnessInSharedInjection           float32
	CircEffectiveness                         float32

	Step                                  float32
	EntryRateByGroupAndStage              []([]float32)
	DiseaseProgressionExitsByDiseaseStage []float32
}

type Cost struct {
	Id            int
	NationID      int
	ComponentID   int
	CostPerClient float32
	ComponentName string
}

type Spending struct {
	Id               int
	InterventionID   int
	SubpopulationID  int
	NationID         int
	Coverage         float32
	RRR              float32
	RRRTypeID        int
	HIVStatus        int
	RRRStandardError float32
}

// You get this back
type Results struct {
	TotalPrevalence                 []float32
	TotalNewInfections              []float32
	TotalNewInfectionsPerPop        []float32
	CumulativeTotalNewInfections    []float32
	HivDeaths                       []float32
	CumulativeHivDeaths             []float32
	PrevalenceByGroup               []float32
	TotalPlwa                       []float32
	PlwaByGroup                     [][]float32
	IncidenceRate                   [][]float32
	TotalNewInfectionsByGroup       [][]float32
	TotalNewInfectionsByGroupPerPop [][]float32
	HivDeathsByGroup                [][]float32
	TotalCostPerIntervention        float32
	TotalCostPerComponent           float32
	ComponentNames                  float32
	TotalCost                       float32
	TotalPopulation                 []float32
	PropOnArt                       []float32
	PercentOfTotalPopByGroup        [][]float32
}

// Main entry point to the model
func Predict(inputs *Inputs) *Results {
	p := &inputs.CountryProfile
	results := new(Results)
	numPops := 65 //len(p.Groups) * len(p.DiseaseStages)
	currentCycle := make(NSlice, numPops, numPops)
	previousCycle := make(NSlice, numPops, numPops)
	secondPreviousCycle := make(NSlice, numPops, numPops)

	// #############################################################################################################
	// ####################################### Steps 1 and 2: prepare inputs #######################################
	// #############################################################################################################

	// // p.Groups = []string { "Gen Pop Men", "Gen Pop Women", "SW Women", "MSM", "IDU" }
	// p.DiseaseAndTreatmentStages = []string { "Uninfected", "Acute", "Early" , "Med" ,"Late" ,"Adv", "AIDS", "Acute Tx", "Early Tx" , "Med Tx" ,"Late Tx" ,"Adv Tx", "AIDS Tx" }
	// p.PopulationSize = 35067464
	// p.PopulationSizeByGroup = []float32 { 17390000.1, 16741000.1, 66964.1, 869500.1, 68262.1 }
	// p.HivPrevalenceAdultsByGroup = []float32 {  0.061, 0.083, 0.60, 0.039, 0.400 }
	// p.HivPrevalence15yoByGroup = []float32 { 0.02, 0.03 }
	// p.ProprtionDiseaseStage = []float32 { 0.0, 0.05, 0.25, 0.20, 0.20, 0.20, 0.10 }
	// p.InfectiousnessByDiseaseStage = []float32 { 0.16, 0.08, 0.09, 0.16, 0.5, 0.5 }
	// p.HivDeathRateByDiseaseStage = []float32 { 0.0, 0.0, 0.1, 0.2, 0.3, 0.4, 0.45, 0.0, 0.0, 0.02, 0.06, 0.06, 0.08, 0.11 }
	// p.HivDeathRateByDiseaseStageTx = []float32 {  }
	// p.InitialTreatmentAccessByDiseaseStage = []float32 { 0.0, 0.0, 0.0, 0.0, 0.0, 0.2, 0.3 }
	// p.TreatmentRecuitingRateByDiseaseStage = []float32 { 0.0, 0.0, 0.0, 0.0, 0.1, 0.1, 0.1 }
	// p.EntryRateGenPop =  0.09
	// p.MaturationRate =  0.07
	// p.DeathRateGeneralCauses =  0.01
	// p.LifeExpectancy =  55
	// p.SwInitiationRate =  0.0006
	// p.SwQuitRate =  0.005
	// p.EntryRateMsm =  0.00225
	// p.EntryRateIdu =  0.00036
	// p.IduInitiationRate =  0.0001
	// p.IduSpontaneousQuitRate =  0.0003
	// p.IduDeathRate =  0.05
	// p.IncreaseInInfectiousnessHomosexual =  0.5
	// p.DiseaseProgressionUntreatedAcuteToEarly =  2
	// p.DiseaseProgressionUntreatedEarlyToMedium =  0.5
	// p.DiseaseProgressionUntreatedMediumToLate =  0.2
	// p.DiseaseProgressionUntreatedLateToAdvanced =  0.2
	// p.DiseaseProgressionUntreatedAdvancedToAids =  0.5
	// p.DiseaseProgressionTreatedAcuteToEarly =  0
	// p.DiseaseProgressionTreatedEarlyToMedium =  0.1
	// p.DiseaseProgressionTreatedMediumToLate =  0.1
	// p.DiseaseProgressionTreatedLateToAdvanced =  0.1
	// p.DiseaseProgressionTreatedAdvancedToAids =  0.1
	// p.GeneralNonSwPartnershipsYearly =  1.5
	// p.GeneralCondomUse =  0.3
	// p.GeneralCondomEffectiveness =  0.9
	// p.SwProportionWhoUseServices =  0.1
	// p.SwPartnershipsYearly =  120
	// p.SwCondomUseRate =  0.66
	// p.MsmPartnershipsYearly =  3
	// p.MsmCondomUseRate =  0.49
	// p.TreatmentReductionOfInfectiousness =  0.95
	// p.TreatmentQuitRate =  0.05
	// p.PercentOfIduSexPartners =  0.4
	// p.IduPartnershipsYearly =  4.5
	// p.IduCondomUseRate =  0.2
	// p.AnnualNumberOfInjections =  264
	// p.PercentSharedInjections =  0.4
	// p.PercentMaleIdus =  0.8
	// p.InfectiousnessInSharedInjection =  0.005
	// p.CircEffectiveness =  0.6

	// /// new creations
	// NOTE: should be made private?
	p.EntryRateByGroupAndStage = make([]([]float32), 5)
	p.EntryRateByGroupAndStage[0] = []float32{p.EntryRateGenPop / 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	p.EntryRateByGroupAndStage[1] = []float32{p.EntryRateGenPop / 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	p.EntryRateByGroupAndStage[2] = []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	p.EntryRateByGroupAndStage[3] = []float32{p.EntryRateMsm, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	p.EntryRateByGroupAndStage[4] = []float32{p.EntryRateIdu, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	// NOTE: should be made private?
	p.DiseaseProgressionExitsByDiseaseStage = []float32{
		0,
		p.DiseaseProgressionUntreatedAcuteToEarly,
		p.DiseaseProgressionUntreatedEarlyToMedium,
		p.DiseaseProgressionUntreatedMediumToLate,
		p.DiseaseProgressionUntreatedLateToAdvanced,
		p.DiseaseProgressionUntreatedAdvancedToAids,
		p.DiseaseProgressionTreatedAcuteToEarly,
		p.DiseaseProgressionTreatedEarlyToMedium,
		p.DiseaseProgressionUntreatedMediumToLate,
		p.DiseaseProgressionTreatedLateToAdvanced,
		p.DiseaseProgressionTreatedAdvancedToAids}

	// #############################################################################################################
	// ####################################### Step 3: Compute initual pops, vary parameters by group and disease stage #######################################
	// #############################################################################################################

	//calculate initial populations
	for g, _ := range p.Groups {
		for s, _ := range p.DiseaseStages {
			var i int = g*13 + s
			if s == 0 {
				currentCycle[i] = p.PopulationSizeByGroup[g] * (1 - p.HivPrevalenceAdultsByGroup[g])
			}
			if s > 0 && s < 7 {
				currentCycle[i] = p.PopulationSizeByGroup[g] *
					p.HivPrevalenceAdultsByGroup[g] *
					p.ProprtionDiseaseStage[s] *
					(1 - p.InitialTreatmentAccessByDiseaseStage[s])
			}
			if s > 6 {
				var ds int = s - 6
				currentCycle[i] = p.PopulationSizeByGroup[g] *
					p.HivPrevalenceAdultsByGroup[g] *
					p.ProprtionDiseaseStage[ds] *
					p.InitialTreatmentAccessByDiseaseStage[ds]
			}
		} // end disease stage
	} // end group

	previousCycle = currentCycle

	//begin main loop
	s := 0 // stage
	g := 0 // group

	for c := 0; c < NUMCYCLES; c++ {

		// TODO
		// Step 6: develop dynamics

		// Step 6a: general dynamics
		var _ float32 = dGeneral(previousCycle, g, s, p) +
			dProgEntries(previousCycle, g, s, p) +
			dProgExits(previousCycle, g, s, p) +
			dTreatment(previousCycle, g, s, p) +
			dIduSw(previousCycle, g, s, p)

		fmt.Println(g, s)
		fmt.Println(previousCycle)

		// At the end of the group cycle, prepare the s and g variables
		if s == 12 {
			s = 0
			g++
		}

		//cycle is over, replace cycle variables
		secondPreviousCycle = previousCycle
		previousCycle = currentCycle
		currentCycle := make(NSlice, numPops, numPops)

		fmt.Println(previousCycle, secondPreviousCycle, currentCycle)

	} //end cycle
	return results
} //end predict

func srcSum(n NSlice, g int, s int) float32 {
	// TODO src code here
	return 0
}

func dGeneral(n NSlice, g int, s int, p *CountryProfile) float32 {
	return p.Step * (n.gs(g, s)*(-p.MaturationRate-p.DeathRateGeneralCauses-p.HivDeathRateByDiseaseStage[s]) + n.sum()*p.EntryRateByGroupAndStage[g][s])
}

func dProgExits(n NSlice, g int, s int, p *CountryProfile) float32 {
	if s == 0 {
		return p.Step * n.gs(g, s) * srcSum(n, g, s)
	} else {
		return p.Step * n.gs(g, s) * p.DiseaseProgressionExitsByDiseaseStage[s]
	}
}

func dProgEntries(n NSlice, g int, s int, p *CountryProfile) float32 {
	if s == 0 {
		return p.Step * n.gs(g, s) * srcSum(n, g, s)
	} else {
		return p.Step * n.gs(g, s) * p.DiseaseProgressionExitsByDiseaseStage[s]
	}
}

func dTreatment(n NSlice, g int, s int, p *CountryProfile) float32 {
	if s == 0 {
		return 0
	} else if s > 0 && s < 7 {
		return p.Step * (n.gs(g, s)*p.TreatmentRecuitingRateByDiseaseStage[s] + n.gs(g, s+6)*p.TreatmentQuitRate)
	} else {
		return p.Step * (n.gs(g, s-6)*p.TreatmentRecuitingRateByDiseaseStage[s] + n.gs(g, s)*-p.TreatmentQuitRate)
	}
}

func dIduSw(n NSlice, g int, s int, p *CountryProfile) float32 {
	if g == 0 { // Genpop men
		return p.Step * (n.gs(4, s)*p.PercentMaleIdus*p.IduSpontaneousQuitRate + n.gs(g, s)*-p.IduInitiationRate)
	} else if g == 1 { // Genpop women
		return p.Step * (n.gs(4, s)*(1-p.PercentMaleIdus)*p.IduSpontaneousQuitRate + n.gs(g, s)*-p.IduInitiationRate + n.gs(3, s)*p.SwQuitRate + n.gs(g, s)*-p.SwInitiationRate)
	} else if g == 2 { // Sew workers
		return p.Step * (n.gs(g, s)*-p.SwQuitRate + n.gs(1, s)*p.SwInitiationRate)
	} else if g == 3 { // MSM
		return 0
	} else if g == 4 { // Drug users
		return p.Step * (n.gs(g, s)*-p.IduSpontaneousQuitRate + n.gs(0, s)*p.IduInitiationRate + n.gs(1, s)*p.IduInitiationRate)
	}
	// TODO log this situation?
	debug("Should not be here")
	return 0
}

func debug(s string) {
	fmt.Println(s)
}
