package costanalysis

// FIX ME: AMC is not integrated

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

const (
	NUMCYCLES              = 40
	NUMBERINTERVENTIONS    = 8
	NUMBERSUBINTERVENTIONS = 20
	NUMBEROFCOMPONENTS     = 60
)

// TODO: maybe some of these slices are fixed size and could just be arrays?
var (
	condomUse []float32
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

func (n NSlice) proportion(g, s int) float32 {
	return n.gs(g, s) / n.g(g)
}

// TODO: stop using hardcoded number of groups, stages. Use the data: len(groups), ..
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

func (n NSlice) setPopAtGS(pop float32, g int, s int) NSlice {

	var i int = g*13 + s
	n[i] = pop
	return n

}

type Cgs struct {
	C                     int
	G                     int
	S                     int
	Group                 string
	DiseaseStage          string
	Population            float32
	Proportion            float32
	Scr                   float32 //note this would be for PREVIOUS year
	DGeneral              float32
	DProgExits            float32
	DProgEntries          float32
	DTreatment            float32
	DIduSw                float32
	CompositePartnerships float32
	Infectiousness        float32
}

type CountryProfile struct {
	Groups                                    []string
	DiseaseStages                             []string
	PopulationSize                            float32
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
	CondomUseByGroup                      []float32
	PartnershipsByGroup                   []float32

	CondomUseByGroupAndHivStatus             [][]float32
	PartnershipsByGroupAndHivStatus          [][]float32
	AnnualNumberofRiskyInjectionsByHivStatus []float32
}

type Cost struct {
	InterventionId      int
	Id                  int
	NationId            int
	ComponentId         int
	CostPerClient       float32
	ComponentName       string
	SuperInterventionId int
}

type Spending struct {
	Id               int
	InterventionId   int
	SubpopulationId  int
	NationId         int
	Coverage         float32
	RRR              float32
	RRRTypeId        int
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
	PrevalenceByGroup               [][]float32
	TotalPlwa                       []float32
	PlwaByGroup                     [][]float32
	IncidenceRate                   []float32
	TotalNewInfectionsByGroup       [][]float32
	TotalNewInfectionsByGroupPerPop [][]float32
	HivDeathsByGroup                [][]float32
	TotalCostPerIntervention        []float32
	TotalCostPerComponent           []float32
	ComponentNames                  []string
	TotalCost                       float32
	TotalPopulation                 []float32
	PropOnArt                       []float32
	PercentOfTotalPopByGroup        [][]float32
}

type ChData struct {
	key   int
	value float32
}

func initResults() *Results {
	r := new(Results)
	//slice of slice
	r.PrevalenceByGroup = make([][]float32, 5, 5)
	r.PlwaByGroup = make([][]float32, 5, 5)
	r.TotalNewInfectionsByGroup = make([][]float32, 5, 5)
	r.TotalNewInfectionsByGroupPerPop = make([][]float32, 5, 5)
	r.HivDeathsByGroup = make([][]float32, 5, 5)
	r.PercentOfTotalPopByGroup = make([][]float32, 5, 5)
	return r
}

// Main entry point to the model
func Predict(inputs *Inputs) *Results {
	beginTime := time.Now()
	fmt.Println("Predicting results...")
	//buildCsvHeaders()
	results := initResults()
	p := &inputs.CountryProfile
	results = calculateCosts(inputs.Costs, inputs.Spendings, p, results)
	numPops := 65 //len(p.Groups) * len(p.DiseaseStages)
	currentCycle := make(NSlice, numPops, numPops)
	previousCycle := make(NSlice, numPops, numPops)

	// -- Apply intervention changes -- //

	spendings := inputs.Spendings

	// allCycles := make([]Cgs, 10000, 10000)
	p.DiseaseStages = []string{"Uninfected", "Acute", "Early", "Med", "Late", "Adv", "AIDS", "Acute Tx", "Early Tx", "Med Tx", "Late Tx", "Adv Tx", "AIDS Tx"}
	p.Step = 0.5
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
		0,
		p.DiseaseProgressionTreatedAcuteToEarly,
		p.DiseaseProgressionTreatedEarlyToMedium,
		p.DiseaseProgressionUntreatedMediumToLate,
		p.DiseaseProgressionTreatedLateToAdvanced,
		p.DiseaseProgressionTreatedAdvancedToAids,
		0}

	p.CondomUseByGroup = []float32{
		p.GeneralCondomUse,
		p.GeneralCondomUse,
		p.SwCondomUseRate,
		p.MsmCondomUseRate,
		p.IduCondomUseRate}

	p.PartnershipsByGroup = []float32{
		p.GeneralNonSwPartnershipsYearly,
		p.GeneralNonSwPartnershipsYearly,
		p.SwPartnershipsYearly,
		p.MsmPartnershipsYearly,
		p.IduPartnershipsYearly}

	p.AnnualNumberofRiskyInjectionsByHivStatus = make([]float32, 3, 3)
	p.AnnualNumberofRiskyInjectionsByHivStatus[0] = p.AnnualNumberOfInjections * p.PercentSharedInjections
	p.AnnualNumberofRiskyInjectionsByHivStatus[1] = p.AnnualNumberOfInjections * p.PercentSharedInjections
	p.AnnualNumberofRiskyInjectionsByHivStatus[2] = p.AnnualNumberOfInjections * p.PercentSharedInjections

	p = applyInterventions(p, spendings)

	// #############################################################################################################
	// ############ Step 3: Compute initual pops, vary parameters by group and disease stage #######################
	// #############################################################################################################

	//calculate initial populations, initialize other variables, matrices
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
	results.add(0, currentCycle, currentCycle, p)
	previousCycle = currentCycle

	//begin main loop
	ch := make(chan ChData)
	for c := 1; c < NUMCYCLES; c++ {
		for g, _ := range p.Groups {
			for s, _ := range p.DiseaseStages {
				go func(previousCycle NSlice, c int, g int, s int, p *CountryProfile, ch chan ChData) {
					totalDynamics := calculateTotalDynamics(previousCycle, g, s, p)
					newPopulation := previousCycle.gs(g, s) + totalDynamics
					// csvLine(previousCycle, c, g, s, p, newPopulation)
					ch <- ChData{g*13 + s, newPopulation}
				}(previousCycle, c, g, s, p, ch)
			} // end stage
		} //end group
		for g, _ := range p.Groups {
			for s, _ := range p.DiseaseStages {
				theChData := <-ch
				currentCycle[theChData.key] = theChData.value
				var _ int = g * s //just to avoid "not used error"
			} // end stage
		} //end group
		results = results.add(c, currentCycle, previousCycle, p)
		previousCycle = currentCycle
	} //end cycle

	// FIXME return results below
	fmt.Println("Done")
	fmt.Println("Time elapsed:", fmt.Sprint(time.Since(beginTime)))
	//fmt.Println(results)
	return results
} //end predict

func (r *Results) add(c int, currentCycle NSlice, previousCycle NSlice, p *CountryProfile) *Results {

	// --- --- Prepare data --- ---

	// Calculate basic variables
	totalPopulation := currentCycle.sum()
	totalInfected := totalPopulation - currentCycle.s(0)
	totalPrevalence := totalInfected / totalPopulation
	totalUninfedcted := totalPopulation - totalInfected
	totalOnArt := currentCycle.s(7) + currentCycle.s(8) + currentCycle.s(9) + currentCycle.s(10) + currentCycle.s(11) + currentCycle.s(12)

	totalNewInfectionsByGroup := make([]float32, 5, 5)
	totalNewInfectionsByGroup[0] = -dProgExits(previousCycle, 0, 0, p) // FIXME: Not the best way to od this, since these have already been calculated
	totalNewInfectionsByGroup[1] = -dProgExits(previousCycle, 1, 0, p)
	totalNewInfectionsByGroup[2] = -dProgExits(previousCycle, 2, 0, p)
	totalNewInfectionsByGroup[3] = -dProgExits(previousCycle, 3, 0, p)
	totalNewInfectionsByGroup[4] = -dProgExits(previousCycle, 4, 0, p)
	totalNewInfections := totalNewInfectionsByGroup[0] + totalNewInfectionsByGroup[1] + totalNewInfectionsByGroup[2] + totalNewInfectionsByGroup[3] + totalNewInfectionsByGroup[4]
	hivDeathsByGroup := make([]float32, 5, 5)
	hivDeathsByGroup[0] = dHivDeaths(previousCycle, 0, p)
	hivDeathsByGroup[1] = dHivDeaths(previousCycle, 1, p)
	hivDeathsByGroup[2] = dHivDeaths(previousCycle, 2, p)
	hivDeathsByGroup[3] = dHivDeaths(previousCycle, 3, p)
	hivDeathsByGroup[4] = dHivDeaths(previousCycle, 4, p)
	var hivDeaths float32 = hivDeathsByGroup[0] + hivDeathsByGroup[1] + hivDeathsByGroup[2] + hivDeathsByGroup[3] + hivDeathsByGroup[4]

	populationByGroup := make([]float32, 5, 5)
	populationByGroup[0] = currentCycle.g(0)
	populationByGroup[1] = currentCycle.g(1)
	populationByGroup[2] = currentCycle.g(2)
	populationByGroup[3] = currentCycle.g(3)
	populationByGroup[4] = currentCycle.g(4)

	plwaByGroup := make([]float32, 5, 5)
	plwaByGroup[0] = populationByGroup[0] - currentCycle.gs(0, 0)
	plwaByGroup[1] = populationByGroup[1] - currentCycle.gs(1, 0)
	plwaByGroup[2] = populationByGroup[2] - currentCycle.gs(2, 0)
	plwaByGroup[3] = populationByGroup[3] - currentCycle.gs(3, 0)
	plwaByGroup[4] = populationByGroup[4] - currentCycle.gs(4, 0)

	totalPlwa := plwaByGroup[0] + plwaByGroup[1] + plwaByGroup[2] + plwaByGroup[3] + plwaByGroup[4]

	prevalenceByGroup := make([]float32, 5, 5)
	prevalenceByGroup[0] = plwaByGroup[0] / populationByGroup[0]
	prevalenceByGroup[1] = plwaByGroup[1] / populationByGroup[1]
	prevalenceByGroup[2] = plwaByGroup[2] / populationByGroup[2]
	prevalenceByGroup[3] = plwaByGroup[3] / populationByGroup[3]
	prevalenceByGroup[4] = plwaByGroup[4] / populationByGroup[4]

	// --- --- Append to existing results --- ---

	// --- Basic ---
	r.TotalPrevalence = append(r.TotalPrevalence, totalPrevalence)
	r.TotalNewInfections = append(r.TotalNewInfections, totalNewInfections)
	r.TotalNewInfectionsPerPop = append(r.TotalNewInfectionsPerPop, totalNewInfections/totalPopulation)
	r.TotalPlwa = append(r.TotalPlwa, totalPlwa)
	r.HivDeaths = append(r.HivDeaths, hivDeaths)
	r.PropOnArt = append(r.PropOnArt, totalOnArt/totalInfected)
	r.TotalPopulation = append(r.TotalPopulation, totalPopulation)

	// -- Cumulative measurements --
	if c != 0 { // Not cycle zero
		r.CumulativeHivDeaths = append(r.CumulativeHivDeaths, float32(r.CumulativeHivDeaths[int(c)-1])+float32(hivDeaths))
		r.CumulativeTotalNewInfections = append(r.CumulativeTotalNewInfections, r.CumulativeTotalNewInfections[c-1]+totalNewInfections)
	} else { // Cycle zero
		r.CumulativeHivDeaths = append(r.CumulativeHivDeaths, hivDeaths)
		r.CumulativeTotalNewInfections = append(r.CumulativeTotalNewInfections, totalNewInfections)
	}

	// -- By Group  --
	for i := 0; i < 5; i++ {
		r.TotalNewInfectionsByGroup[i] = append(r.TotalNewInfectionsByGroup[i], totalNewInfectionsByGroup[i]) // FIXME: Not the best way to od this, since these have already been calculate
		r.PlwaByGroup[i] = append(r.PlwaByGroup[i], plwaByGroup[i])
		r.PrevalenceByGroup[i] = append(r.PrevalenceByGroup[i], prevalenceByGroup[i])
		r.HivDeathsByGroup[i] = append(r.HivDeathsByGroup[i], hivDeathsByGroup[i])
		r.TotalNewInfectionsByGroupPerPop[i] = append(r.TotalNewInfectionsByGroup[i], totalNewInfectionsByGroup[i]/populationByGroup[i])
		r.PercentOfTotalPopByGroup[i] = append(r.PercentOfTotalPopByGroup[i], populationByGroup[i]/totalPopulation)
	}

	// -- Only even --
	isOdd := math.Remainder(float64(c), 2)
	if isOdd != 0 && c != 0 { //odd year past cycle 0, calculate yearly values here
		totalNewInfectionsFromPrevAndCurrentCycle := totalNewInfections + r.TotalNewInfections[c-1]
		totalSusepctibleFromPrevAndCurrentCycle := totalUninfedcted + r.TotalPopulation[c-1] - r.TotalPlwa[c-1]
		r.IncidenceRate = append(r.IncidenceRate, totalNewInfectionsFromPrevAndCurrentCycle/totalSusepctibleFromPrevAndCurrentCycle)
	} else if isOdd == 0 && c != 0 { //even cycle, not cycle 0
		r.IncidenceRate = append(r.IncidenceRate, r.IncidenceRate[c-2])
	}

	// PrevalenceByGroup               []float32
	// TotalPlwa                       []float32
	// PlwaByGroup                     [][]float32
	// TotalNewInfectionsByGroup       [][]float32
	// TotalNewInfectionsByGroupPerPop [][]float32
	// HivDeathsByGroup                [][]float32
	// TotalCostPerIntervention        float32
	// TotalCostPerComponent           float32
	// ComponentNames                  []string
	// TotalCost                       float32
	// TotalPopulation                 []float32
	// PropOnArt                       []float32
	// PercentOfTotalPopByGroup        [][]float32

	return r
}

func applyInterventions(p *CountryProfile, spendings []Spending) *CountryProfile {

	for q := 0; q < len(spendings); q++ {
		//fmt.Println(spendings[q])
	}
	//fmt.Println("----------------------------------")

	p.CondomUseByGroupAndHivStatus = make([][]float32, 5, 5)
	p.PartnershipsByGroupAndHivStatus = make([][]float32, 5, 5)

	for g, _ := range p.Groups {

		p.CondomUseByGroupAndHivStatus[g] = make([]float32, 3, 3)
		p.PartnershipsByGroupAndHivStatus[g] = make([]float32, 3, 3)

		//foreach potential state in HIV status "negative, positive, treated"
		for hivStatus := 0; hivStatus < 3; hivStatus++ {

			//Compute total RRR from all interventions by group and calculate new condom use rate by group. This assumes interventions that do not affect a certain population group have RRR of 0.
			// Outcome type 1 = Condom use rate

			p.CondomUseByGroupAndHivStatus[g][hivStatus] = 1 - (1-p.CondomUseByGroup[g])*findCompositeRrr(spendings, g, 0, hivStatus)
			p.PartnershipsByGroupAndHivStatus[g][hivStatus] = p.PartnershipsByGroup[g] * findCompositeRrr(spendings, g, 1, hivStatus)
			p.AnnualNumberofRiskyInjectionsByHivStatus[hivStatus] = p.AnnualNumberOfInjections * p.PercentSharedInjections * findCompositeRrr(spendings, g, 2, hivStatus)

		} //end hiv status
	} // end groups

	//fmt.Println("con", p.CondomUseByGroupAndHivStatus)
	//fmt.Println("par", p.PartnershipsByGroupAndHivStatus)
	return p

} // end apply intervention

//find group (subpopulation), outcome, hiv status
func findCompositeRrr(spendings []Spending, g int, o int, h int) float32 {

	index_add := g*18 + o*3 + h
	var sum float32 = 1
	for i := 0; i < NUMBERINTERVENTIONS; i++ {
		index := i*90 + index_add
		spending := spendings[index]
		sum *= (1 - spending.RRR*spending.Coverage)
	}
	//fmt.Println("hiv status ", h, " group ", g+1, " outcome ", o+1, " sum_rrr ", sum)
	return sum

}

func srcSum(n NSlice, g int, s int, p *CountryProfile) float32 {
	var sum float32 = 0.0
	for ss, _ := range p.DiseaseStages {
		sum += src(n, g, ss, p)
		//fmt.Println(sum)
	}

	return sum
}

func getHivStatus(s int) int {

	var hivStatus int
	if s == 0 {
		hivStatus = 0
	} else if s > 0 && s < 7 {
		hivStatus = 1
	} else if s > 6 {
		hivStatus = 2
	}

	return hivStatus

}

//note costs are reported in intervention types, while rrr is reported in super-intervention
func calculateCosts(theCosts []Cost, theSpending []Spending, p *CountryProfile, theResults *Results) *Results {
	//create array of costs seperated into their interventions
	componentsByIntervention := make([][]Cost, NUMBERINTERVENTIONS+2) // FIXME the additional two is for ART and PMTCT, which have not yet been added to model
	costPerClientByIntervention := make([]float32, NUMBERINTERVENTIONS+2)
	totalCostByIntervention := make([]float32, NUMBERINTERVENTIONS+2)
	totalCostPerComponent := make([]float32, NUMBEROFCOMPONENTS, NUMBEROFCOMPONENTS) //unsure of size
	componentNames := make([]string, NUMBEROFCOMPONENTS, NUMBEROFCOMPONENTS)         //unsure of size
	var totalCost float32
	for i := 0; i < len(theCosts); i++ {
		fmt.Println(theCosts[i])
		componentsByIntervention[theCosts[i].SuperInterventionId-1] = append(componentsByIntervention[theCosts[i].SuperInterventionId-1], theCosts[i])
		costPerClientByIntervention[theCosts[i].SuperInterventionId-1] += theCosts[i].CostPerClient
	}
	//go through each spending, figure out how much will cost, add to approcpriate arrays
	for i := 0; i < len(theSpending)-18; i = i + 18 { //you move forward 18 because there are 18 grouping (outcomes * hivStatus) per intervention-subpop grouping
		clientsForThisSpending := theSpending[i].Coverage * p.PopulationSizeByGroup[theSpending[i].SubpopulationId-1]
		costForThisSpending := clientsForThisSpending * costPerClientByIntervention[theSpending[i].InterventionId-1]
		totalCostByIntervention[theSpending[i].InterventionId-1] += costForThisSpending
		totalCost += costForThisSpending
		interventionComponents := componentsByIntervention[theSpending[i].InterventionId-1]
		//now go through the components of this spending to itemize and fill in costsByComponent
		for p := 0; p < len(interventionComponents); p++ {
			totalCostPerComponent[interventionComponents[p].ComponentId-1] += clientsForThisSpending * interventionComponents[p].CostPerClient
			componentNames[interventionComponents[p].ComponentId-1] = interventionComponents[p].ComponentName
		}

		fmt.Println("For intervention ", theSpending[i].InterventionId, " there are ", len(interventionComponents), " comps totalling to ", totalCostByIntervention[theSpending[i].InterventionId-1])

	}

	//ART
	clientsForThisSpending := p.PopulationSize * (p.InitialTreatmentAccessByDiseaseStage[0]*p.ProprtionDiseaseStage[0] + p.InitialTreatmentAccessByDiseaseStage[1]*p.ProprtionDiseaseStage[1] + p.InitialTreatmentAccessByDiseaseStage[2]*p.ProprtionDiseaseStage[2] + p.InitialTreatmentAccessByDiseaseStage[3]*p.ProprtionDiseaseStage[3] + p.InitialTreatmentAccessByDiseaseStage[4]*p.ProprtionDiseaseStage[4] + p.InitialTreatmentAccessByDiseaseStage[5]*p.ProprtionDiseaseStage[5] + p.InitialTreatmentAccessByDiseaseStage[6]*p.ProprtionDiseaseStage[6])
	costForThisSpending := clientsForThisSpending * costPerClientByIntervention[8] // 8 = art
	totalCostByIntervention[8] = costForThisSpending
	totalCost += costForThisSpending
	interventionComponents := componentsByIntervention[8]
	//now go through the components of this spending to itemize and fill in costsByComponent
	for p := 0; p < len(interventionComponents); p++ {
		totalCostPerComponent[interventionComponents[p].ComponentId-1] += clientsForThisSpending * interventionComponents[p].CostPerClient
		componentNames[interventionComponents[p].ComponentId-1] = interventionComponents[p].ComponentName
	}

	fmt.Println("For ART there are ", len(interventionComponents), " comps totalling to ", totalCostByIntervention[8])

	theResults.TotalCost = totalCost
	theResults.ComponentNames = componentNames
	theResults.TotalCostPerComponent = totalCostPerComponent
	theResults.TotalCostPerIntervention = totalCostByIntervention
	return theResults
}

func src(n NSlice, g int, s int, p *CountryProfile) float32 { // FIX ME: still need to use the array version of condom use and # of partners, currently use general, pre-intervention levels

	hivStatus := getHivStatus(s)

	if g == 0 { // General population men

		// (
		// 	1-(
		// 		(
		// 			(1-PercentOfIduSexPartners) *
		// 			 (1-PercentMaleIdus) *
		// 			  subpopulations.g(4).c(c).compact.inject(0){|memo, n| memo + (n.Population) *(n.Partnerships)}
		// 		) / (
		// 			subpopulations.g(0).c(c).compact.inject(0){|memo, n| memo + (n.Population) * (n.Partnerships) }
		// 		)
		// 	)
		// ) *
		// genpop_female_mate.Probability *
		//  uninfected_men.CompositePartnerships *
		//   genpop_female_mate.Infectiousness

		// from genpop women
		scrGpW := (1 - (1-p.PercentOfIduSexPartners)*
			(1-p.PercentMaleIdus)*
			n.g(4)*
			p.IduPartnershipsYearly/(n.g(0)*
			p.GeneralNonSwPartnershipsYearly)) *
			n.proportion(1, s) *
			compositePartnerships(g, s, p) *
			infectiousness(s, p)

		// (
		// 	(
		// 		(1-PercentOfIduSexPartners) *
		// 		(1-PercentMaleIdus) *
		// 		 subpopulations.g(4).c(c).compact.inject(0){|memo, n| memo + (n.Population) *(n.Partnerships)}
		// 	) / (
		// 		subpopulations.g(0).c(c).compact.inject(0){|memo, n| memo + (n.Population) * (n.Partnerships) }
		// 	)
		// )  *
		// idu_mate.Probability *
		// uninfected_men.CompositePartnerships *
		// idu_mate.Infectiousness

		// from idu

		scrIdu := (1 - p.PercentOfIduSexPartners) *
			(1 - p.PercentMaleIdus) *
			n.g(4) *
			p.IduPartnershipsYearly /
			(n.g(0) *
				p.GeneralNonSwPartnershipsYearly) *
			n.proportion(4, s) *
			compositePartnerships(g, s, p) *
			infectiousness(s, p)

		//from sw

		// sw_mate.Probability *
		//  (
		//  	(subpopulations.g(2).c(c).compact.inject(0){|memo, n| memo + (n.Population) * (n.Partnerships) })
		//  	 / (total_men * SwProportionWhoUseServices)
		// )
		//  *(1-sw_mate.CondomUse
		//  * GeneralCondomEffectiveness)
		//  * sw_mate.Infectiousness
		//  * SwProportionWhoUseServices

		scrSw := (n.g(2) *
			p.SwPartnershipsYearly) /
			(p.SwProportionWhoUseServices *
				n.g(0)) *
			n.proportion(2, s) *
			p.SwProportionWhoUseServices *
			infectiousness(s, p) *
			(1 - p.CondomUseByGroupAndHivStatus[2][hivStatus]*p.GeneralCondomEffectiveness)

		//fmt.Println((n.gs(0, s)))
		//fmt.Println(scrGpW, scrIdu, scrSw)
		return scrGpW + scrIdu + scrSw
	}

	if g == 1 { // General population women

		//TotalMen := n.g(0)
		TotalWomen := n.g(1)
		TotalMenPartnerships := n.g(0) * p.GeneralNonSwPartnershipsYearly
		//TotalIdus := n.g(4)
		TotalMenPartnershipsLostToIdu := (1 - p.PercentOfIduSexPartners) * (1 - p.PercentMaleIdus) * n.g(4) * p.IduPartnershipsYearly
		TotalPartnershipsOfferedToGpWomenFromMen := TotalMenPartnerships - TotalMenPartnershipsLostToIdu
		AverageNumberOfParternshipsForWomenFromMen := TotalPartnershipsOfferedToGpWomenFromMen / TotalWomen
		TotalPartnershipsOfferedFromIdusToFemales := (1 - p.PercentOfIduSexPartners) * p.PercentMaleIdus * n.g(4) * p.IduPartnershipsYearly
		AveragePartnershipsFromIdusToFemale := TotalPartnershipsOfferedFromIdusToFemales / TotalWomen
		//TotalSws := n.g(2)
		//TotalMaleClients := TotalMen * p.SwProportionWhoUseServices
		//TotalSwPartnershipsOffered := n.g(2) * p.SwPartnershipsYearly
		//AverageSwPartnershipsPerUsingMale := TotalSwPartnershipsOffered / TotalMaleClients
		TotalMaleIduPartnerships := (1 - p.PercentOfIduSexPartners) * p.PercentMaleIdus * n.g(4) * p.IduPartnershipsYearly
		TotalFemaleIduPartnerships := (1 - p.PercentOfIduSexPartners) * (1 - p.PercentMaleIdus) * n.g(4) * p.IduPartnershipsYearly
		MaleIduPartnershipsOfferedToFemales := (1 - p.PercentOfIduSexPartners) * p.PercentMaleIdus * n.g(4) * p.IduPartnershipsYearly
		TotalFemalePartnerships := TotalMenPartnerships - TotalFemaleIduPartnerships + MaleIduPartnershipsOfferedToFemales
		ProbabilityOfChoosingIduPartner := TotalMaleIduPartnerships / TotalFemalePartnerships
		ProbabilityOfChoosingNonIduPartner := 1 - ProbabilityOfChoosingIduPartner
		SrcForFemalesFromMale := ProbabilityOfChoosingNonIduPartner * n.proportion(0, s) * AverageNumberOfParternshipsForWomenFromMen * (1 - (p.CondomUseByGroupAndHivStatus[0][hivStatus])*p.GeneralCondomEffectiveness) * infectiousness(s, p)
		SrcForFemalesFromIduMale := ProbabilityOfChoosingIduPartner * n.proportion(4, s) * AveragePartnershipsFromIdusToFemale * (1 - p.CondomUseByGroupAndHivStatus[4][hivStatus]*p.GeneralCondomEffectiveness) * infectiousness(s, p)
		TotalSrcForFemales := SrcForFemalesFromMale + SrcForFemalesFromIduMale
		return TotalSrcForFemales

	}
	if g == 2 { // Sex workers
		return n.proportion(0, s) * compositePartnerships(2, s, p) * infectiousness(s, p)
	}
	if g == 3 { // Men who have sex with men

		return n.proportion(3, s) * compositePartnerships(3, s, p) * infectiousness(s, p) * (1 + p.IncreaseInInfectiousnessHomosexual)
	}
	if g == 4 { // Injecting drug users

		//something here is slowing down program, added these lines and slowed by 10x...

		var treatmentModifier float32
		if s == 0 {
			treatmentModifier = 0
		} else if s > 6 {
			treatmentModifier = 0.5
		} else {
			treatmentModifier = 1.0
		}

		an := p.AnnualNumberofRiskyInjectionsByHivStatus[hivStatus]

		return p.PercentOfIduSexPartners*
			n.proportion(4, s)*
			compositePartnerships(4, s, p)*
			infectiousness(s, p) + n.proportion(4, s)*an*p.InfectiousnessInSharedInjection*treatmentModifier
	}

	// TODO log this situation?
	debug("Should not be here")
	return 0
}

func iduFemMate(n NSlice, p *CountryProfile) float32 {

	return ((1 - p.PercentOfIduSexPartners) * p.PercentMaleIdus * n.g(4) * p.IduPartnershipsYearly) / n.g(1)

}

func womenPartnership(n NSlice, p *CountryProfile) float32 {

	return (n.g(0)*p.GeneralNonSwPartnershipsYearly - (1-p.PercentOfIduSexPartners)*(1-p.PercentMaleIdus)*n.g(4)*p.IduPartnershipsYearly) / n.g(1)

}

func dHivDeaths(n NSlice, g int, p *CountryProfile) float32 {

	var sum float32 = 0.0
	for s := 0; s < 13; s++ {
		if s < 7 {
			sum += n.gs(g, s) * p.HivDeathRateByDiseaseStage[s]
		} else {
			sum += n.gs(g, s) * p.HivDeathRateByDiseaseStageTx[s-6]
		}

	}
	return sum

}

func dGeneral(n NSlice, g int, s int, p *CountryProfile) float32 {
	//fmt.Println("oh,", n.gs(g, s)*p.Step*-p.MaturationRate-p.DeathRateGeneralCauses-p.HivDeathRateByDiseaseStage[s]+p.Step*n.sum()*p.EntryRateByGroupAndStage[g][s])

	if g == 4 { //idu's have different death rate

		if s < 7 {
			return p.Step * (n.gs(g, s)*(-p.MaturationRate-p.IduDeathRate-p.HivDeathRateByDiseaseStage[s]) + n.sum()*p.EntryRateByGroupAndStage[g][s])
		} else {
			return p.Step * (n.gs(g, s)*(-p.MaturationRate-p.IduDeathRate-p.HivDeathRateByDiseaseStageTx[s-6]) + n.sum()*p.EntryRateByGroupAndStage[g][s])
		}
	} else {
		if s < 7 {
			return p.Step * (n.gs(g, s)*(-p.MaturationRate-p.DeathRateGeneralCauses-p.HivDeathRateByDiseaseStage[s]) + n.sum()*p.EntryRateByGroupAndStage[g][s])
		} else {
			return p.Step * (n.gs(g, s)*(-p.MaturationRate-p.DeathRateGeneralCauses-p.HivDeathRateByDiseaseStageTx[s-6]) + n.sum()*p.EntryRateByGroupAndStage[g][s])
		}
	}
	return 0 //shouldn't be here; just for compiler

}

func dProgExits(n NSlice, g int, s int, p *CountryProfile) float32 {
	if s == 0 {
		return -p.Step * n.gs(g, s) * srcSum(n, g, s, p)
	} else {
		return -p.Step * n.gs(g, s) * p.DiseaseProgressionExitsByDiseaseStage[s]
	}
	return 0 //shouldn't be here; just for compiler
}

func dProgEntries(n NSlice, g int, s int, p *CountryProfile) float32 {
	if s == 0 {
		return 0
	} else {
		return -dProgExits(n, g, s-1, p)
	}
	return 0 //shouldn't be here; just for compiler
}

func dTreatment(n NSlice, g int, s int, p *CountryProfile) float32 {
	if s == 0 {
		return 0
	} else if s > 0 && s < 7 {
		return p.Step * (n.gs(g, s)*-p.TreatmentRecuitingRateByDiseaseStage[s] + n.gs(g, s+6)*p.TreatmentQuitRate)
	} else {
		return p.Step * (n.gs(g, s-6)*p.TreatmentRecuitingRateByDiseaseStage[s-6] + n.gs(g, s)*-p.TreatmentQuitRate)
	}
	return 0 //shouldn't be here; just for compiler
}

func dIduSw(n NSlice, g int, s int, p *CountryProfile) float32 {
	if g == 0 { // Genpop men
		return p.Step * (n.gs(4, s)*p.PercentMaleIdus*p.IduSpontaneousQuitRate + n.gs(g, s)*-p.IduInitiationRate)
	} else if g == 1 { // Genpop women
		return p.Step * (n.gs(4, s)*(1-p.PercentMaleIdus)*p.IduSpontaneousQuitRate + n.gs(g, s)*-p.IduInitiationRate + n.gs(2, s)*p.SwQuitRate + n.gs(g, s)*-p.SwInitiationRate)
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

func compositePartnerships(g int, s int, p *CountryProfile) float32 {
	hivStatus := getHivStatus(s)
	return p.CondomUseByGroupAndHivStatus[g][hivStatus]*p.PartnershipsByGroupAndHivStatus[g][hivStatus]*(1-p.GeneralCondomEffectiveness) + (1-p.CondomUseByGroupAndHivStatus[g][hivStatus])*p.PartnershipsByGroupAndHivStatus[g][hivStatus]
}

func infectiousness(s int, p *CountryProfile) float32 {
	// replace calls to this function by direct acces to the matrix, when available
	if s < 7 {
		return p.InfectiousnessByDiseaseStage[s]
	} else {
		return p.InfectiousnessByDiseaseStage[s-6] * (1 - p.TreatmentReductionOfInfectiousness)
	}
	return 0 //shouldn't be here; just for compiler
}

func debug(s string) {
	fmt.Println(s)
}

func calculateTotalDynamics(previousCycle NSlice, g int, s int, p *CountryProfile) float32 {

	dGen := dGeneral(previousCycle, g, s, p)
	dProgEx := dProgExits(previousCycle, g, s, p)
	dProgEn := dProgEntries(previousCycle, g, s, p)
	dTx := dTreatment(previousCycle, g, s, p)
	dIduSw := dIduSw(previousCycle, g, s, p)

	//fmt.Println("Dynamics", "dGen", "dProgEx", "dProgEn", "dTx")
	//fmt.Println("Are:", dGen, dProgEx, dProgEn, dTx)
	return (dGen + dProgEx + dProgEn + dTx + dIduSw)
}

func csvLine(previousCycle NSlice, c int, g int, s int, p *CountryProfile, newPopulation float32) {

	var thisCgs Cgs
	thisCgs.C = c
	thisCgs.G = g
	thisCgs.S = s
	thisCgs.Group = p.Groups[g]
	thisCgs.DiseaseStage = p.DiseaseStages[s]
	thisCgs.Population = newPopulation
	//thisCgs.Proportion = currentCycle.proportion(g, s) //FIXME should be current cycle
	thisCgs.DGeneral = dGeneral(previousCycle, g, s, p)
	thisCgs.DProgExits = dProgExits(previousCycle, g, s, p)
	thisCgs.DProgEntries = dProgEntries(previousCycle, g, s, p)
	thisCgs.DTreatment = dTreatment(previousCycle, g, s, p)
	thisCgs.DIduSw = dIduSw(previousCycle, g, s, p)
	thisCgs.Scr = src(previousCycle, g, s, p)
	thisCgs.CompositePartnerships = compositePartnerships(g, s, p)
	thisCgs.Infectiousness = infectiousness(s, p)
	toCsvLine(thisCgs)

}

func toCsvLine(thisCgs Cgs) {
	//fmt.Println(".")
	file, error := os.OpenFile("output.csv", os.O_RDWR|os.O_APPEND, 0666)
	if error != nil {
		panic(error)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	thisCgsSlice := []string{strconv.Itoa(thisCgs.C),
		fmt.Sprint(thisCgs.G),
		fmt.Sprint(thisCgs.S),
		fmt.Sprint(thisCgs.Group),
		fmt.Sprint(thisCgs.DiseaseStage),
		fmt.Sprint(thisCgs.Population),
		fmt.Sprint(thisCgs.DGeneral),
		fmt.Sprint(thisCgs.DProgExits),
		fmt.Sprint(thisCgs.DProgEntries),
		fmt.Sprint(thisCgs.DTreatment),
		fmt.Sprint(thisCgs.DIduSw),
		fmt.Sprint(thisCgs.Scr),
		fmt.Sprint(thisCgs.CompositePartnerships),
		fmt.Sprint(thisCgs.Infectiousness)}
	returnError := writer.Write(thisCgsSlice)
	if returnError != nil {
		fmt.Println(returnError)
	}
	writer.Flush()
}

func buildCsvHeaders() {
	os.Create("output.csv")
	file, error := os.OpenFile("output.csv", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if error != nil {
		panic(error)
	}
	defer file.Close()
	// New Csv wriier
	writer := csv.NewWriter(file)
	var new_headers = []string{"cycle", "group", "stage", "groupName", "stageName", "population", "Dgeneral", "DProgExits", "DProgEntries", "D Treatment", "SWIDU dynamics", "Scr", "CompositePartnerships", "Infectiousness"}
	returnError := writer.Write(new_headers)
	if returnError != nil {
		fmt.Println("error")
	}
	writer.Flush()
}
