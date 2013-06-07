package costanalysis

// Test comment

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	NUMCYCLES = 40
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
	CondomUseByGroup                      []float32
	PartnershipsByGroup                   []float32
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
	IncidenceRate                   []float32
	TotalNewInfectionsByGroup       [][]float32
	TotalNewInfectionsByGroupPerPop [][]float32
	HivDeathsByGroup                [][]float32
	TotalCostPerIntervention        float32
	TotalCostPerComponent           float32
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

// Main entry point to the model
func Predict(inputs *Inputs) *Results {
	beginTime := time.Now()
	fmt.Println("Predicting results...")
	buildCsvHeaders()
	results := new(Results)
	//	fmt.Println("Initializted results: ", results)
	p := &inputs.CountryProfile
	numPops := 65 //len(p.Groups) * len(p.DiseaseStages)
	currentCycle := make(NSlice, numPops, numPops)
	previousCycle := make(NSlice, numPops, numPops)
	secondPreviousCycle := make(NSlice, numPops, numPops)
	allCycles := make([]Cgs, 10000, 10000)
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

	// #############################################################################################################
	// ############ Step 3: Compute initual pops, vary parameters by group and disease stage #######################
	// #############################################################################################################

	//calculate initial populations, initialize other variables, matrices
	for g, _ := range p.Groups {
		// condomUse[g] = condomUse
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

			var thisCgs Cgs
			thisCgs.C = 0
			thisCgs.G = g
			thisCgs.S = s
			thisCgs.Group = p.Groups[g]
			thisCgs.DiseaseStage = p.DiseaseStages[s]
			thisCgs.Population = currentCycle.gs(g, s)
			thisCgs.Population = currentCycle.gs(g, s)
			thisCgs.Proportion = currentCycle.proportion(g, s)
			thisCgs.DGeneral = dGeneral(previousCycle, g, s, p)
			thisCgs.DProgExits = dProgExits(previousCycle, g, s, p)
			thisCgs.DProgEntries = dProgEntries(previousCycle, g, s, p)
			thisCgs.DTreatment = dTreatment(previousCycle, g, s, p)
			thisCgs.DIduSw = dIduSw(previousCycle, g, s, p)
			thisCgs.Scr = src(previousCycle, g, s, p)
			thisCgs.CompositePartnerships = compositePartnerships(g, p)
			thisCgs.Infectiousness = infectiousness(s, p)
			toCsvLine(thisCgs)
			var _ = append(allCycles, thisCgs)

		} // end disease stage
	} // end group
	previousCycle = currentCycle
	//fmt.Println(previousCycle[:10])
	//begin main loop

	ch := make(chan ChData)

	for c := 1; c < NUMCYCLES; c++ {

		for g, _ := range p.Groups {
			for s, _ := range p.DiseaseStages {
				go func(previousCycle NSlice, c int, g int, s int, p *CountryProfile, ch chan ChData) {
					totalDynamics := calculateTotalDynamics(previousCycle, g, s, p)
					newPopulation := previousCycle.gs(g, s) + totalDynamics
					csvLine(previousCycle, c, g, s, p, newPopulation)
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

		fmt.Println("cycle: ", c)
		fmt.Println(currentCycle)

		secondPreviousCycle = previousCycle
		previousCycle = currentCycle
		currentCycle := make(NSlice, numPops, numPops)

		if 0 == 1 {
			fmt.Print(secondPreviousCycle)
			fmt.Print(currentCycle)
		}

	} //end cycle

	// FIXME return results below

	fmt.Println("Done")
	fmt.Println("Time elapsed:", fmt.Sprint(time.Since(beginTime)))
	return results
} //end predict

func srcSum(n NSlice, g int, s int, p *CountryProfile) float32 {
	var sum float32 = 0.0
	for ss, _ := range p.DiseaseStages {
		sum += src(n, g, ss, p)
		//fmt.Println(sum)
	}

	return sum
}

func src(n NSlice, g int, s int, p *CountryProfile) float32 {
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
			compositePartnerships(g, p) *
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
			compositePartnerships(g, p) *
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
			(1 - p.SwCondomUseRate*p.GeneralCondomEffectiveness)

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
		SrcForFemalesFromMale := ProbabilityOfChoosingNonIduPartner * n.proportion(0, s) * AverageNumberOfParternshipsForWomenFromMen * (1 - (p.GeneralCondomUse)*p.GeneralCondomEffectiveness) * infectiousness(s, p)
		SrcForFemalesFromIduMale := ProbabilityOfChoosingIduPartner * n.proportion(4, s) * AveragePartnershipsFromIdusToFemale * (1 - (p.IduCondomUseRate)*p.GeneralCondomEffectiveness) * infectiousness(s, p)
		TotalSrcForFemales := SrcForFemalesFromMale + SrcForFemalesFromIduMale
		return TotalSrcForFemales

	}
	if g == 2 { // Sex workers
		return n.proportion(0, s) * compositePartnerships(2, p) * infectiousness(s, p)
	}
	if g == 3 { // Men who have sex with men

		return n.proportion(3, s) * compositePartnerships(3, p) * infectiousness(s, p) * (1 + p.IncreaseInInfectiousnessHomosexual)
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

		return p.PercentOfIduSexPartners*
			n.proportion(4, s)*
			compositePartnerships(4, p)*
			infectiousness(s, p) +

			n.proportion(4, s)*
				p.AnnualNumberOfInjections*
				p.PercentSharedInjections*
				p.InfectiousnessInSharedInjection*
				treatmentModifier
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

func dGeneral(n NSlice, g int, s int, p *CountryProfile) float32 {
	//fmt.Println("oh,", n.gs(g, s)*p.Step*-p.MaturationRate-p.DeathRateGeneralCauses-p.HivDeathRateByDiseaseStage[s]+p.Step*n.sum()*p.EntryRateByGroupAndStage[g][s])

	if g == 4 { //idu's have different death rate

		return p.Step *
			(n.gs(g, s)*(-p.MaturationRate-p.IduDeathRate-p.HivDeathRateByDiseaseStage[s]) +
				n.sum()*p.EntryRateByGroupAndStage[g][s])

	} else {

		return p.Step *
			(n.gs(g, s)*(-p.MaturationRate-p.DeathRateGeneralCauses-p.HivDeathRateByDiseaseStage[s]) +
				n.sum()*p.EntryRateByGroupAndStage[g][s])

	}

}

func dProgExits(n NSlice, g int, s int, p *CountryProfile) float32 {
	if s == 0 {
		return -p.Step * n.gs(g, s) * srcSum(n, g, s, p)
	} else {
		return -p.Step * n.gs(g, s) * p.DiseaseProgressionExitsByDiseaseStage[s]
	}
}

func dProgEntries(n NSlice, g int, s int, p *CountryProfile) float32 {
	if s == 0 {
		return 0
	} else {
		return -dProgExits(n, g, s-1, p)
	}
}

func dTreatment(n NSlice, g int, s int, p *CountryProfile) float32 {
	if s == 0 {
		return 0
	} else if s > 0 && s < 7 {
		return p.Step * (n.gs(g, s)*-p.TreatmentRecuitingRateByDiseaseStage[s] + n.gs(g, s+6)*p.TreatmentQuitRate)
	} else {
		return p.Step * (n.gs(g, s-6)*p.TreatmentRecuitingRateByDiseaseStage[s-6] + n.gs(g, s)*-p.TreatmentQuitRate)
	}
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

func compositePartnerships(g int, p *CountryProfile) float32 {
	return p.CondomUseByGroup[g]*p.PartnershipsByGroup[g]*(1-p.GeneralCondomEffectiveness) + (1-p.CondomUseByGroup[g])*p.PartnershipsByGroup[g]
}

func infectiousness(s int, p *CountryProfile) float32 {
	// replace calls to this function by direct acces to the matrix, when available
	if s < 7 {
		return p.InfectiousnessByDiseaseStage[s]
	} else {
		return p.InfectiousnessByDiseaseStage[s-6] * (1 - p.TreatmentReductionOfInfectiousness)
	}
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
	thisCgs.CompositePartnerships = compositePartnerships(g, p)
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
