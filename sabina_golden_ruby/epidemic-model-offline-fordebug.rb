### the shebang headers are passed in the "header" file, which varies by environment

puts "Content-Type: text/html\n\n"

=begin
########################################################################################################################################################
########################################################## Step 1: Comments ############################################################
########################################################################################################################################################

This epidemic model is writen in Ruby and is based on an excel model written by Sabina Alistar. The program tracks the estimated HIV epidemic of one country 
given a set of parameters for that country. It breaks the nation up into five groups: general population males, general population females, sex workers, 
men who have sex with men, and intervenous drug users. Each of these five groups is accounted for throughout the modelling process. Within each of the five groups, 
there are 13 different disease stages: uninfected, acute, early, medium, late, aids, acute with treatment, early w tx, medium w tx, late w tx, and aids w tx.

The model is more or less a markov process. It uses a half-year cycle and examines how these populations change over 40 half-year cycles, or twenty years. 

To best represent this process, we create an set of ojects called a CycleGroupDiseaseStages, or CGDs. A CGD represents a single disease stage in a single 
group in a single year. For example, there are 65 total CGDs in the first cycle. There are 5 groups, 13 disease stages in one cycle = 65 CGDs. 
In the first and second cycle, there are a total of 65x2 =  130 CGDs.  

The model has a somewhat complex system which i will explain here. The first few steps take imported data from a JSON object, and turn them into global variables.
These variables generally concern epidemic parameters, such as the size of a cycle, or demographic or epidemiological data. The second major step is to assign all the relevant  

=end




########################################################################################################################################################
########################################################## Step 1: Programming Housekeeping ############################################################
########################################################################################################################################################

#this step sets up debugging and library support. Non-technical users do not need to understand what is going on in this step, except the addressing function

# ------------ needed libraries  ------------  #


require 'rubygems' #gems are ruby-speak for libraires (sets of useful variables and functions)
require 'irb' #allows for debugging
require 'csv' #allows for export to csv
require 'json' #json, a javascript data format, is a universal data format between many languages
require 'cgi'#common gateway interface, opens a port to talk to the client-side, essentially
require 'narray' #numeric array library


# ------------ debugging set up  ------------  #


module IRB # :nodoc:
  def self.start_session(binding)
 unless @__initialized
   args = ARGV
   ARGV.replace(ARGV.dup)
   IRB.setup(nil)
   ARGV.replace(args)
   @__initialized = true
 end

 workspace = WorkSpace.new(binding)

 irb = Irb.new(workspace)

 @CONF[:IRB_RC].call(irb.context) if @CONF[:IRB_RC]
 @CONF[:MAIN_CONTEXT] = irb.context

 catch(:IRB_EXIT) do
   irb.eval_input
 end
  end
end

# ------------  allows hash access like objects, through dot format  ------------  #

#hashes are associaitive arrays (arrays with string keys). this function allows us to access hashes through "dot format" as if they were methods

class ::Hash
 def method_missing(name)
    return self[name] if key? name
    self.each { |k,v| return v if k.to_s.to_sym == name }
    super.method_missing name
  end
end

# ------------  Addressing function ------------  #

#these functions were built to allow us to avoid using the ".select" method, and directly access the item we're looking for.
#essentially, it uses a "address" array to store all the ids of different CDG's

class ::Array
  
  ##Need to compute next generation for cycle 41 so that cycle 40 updates and is identical to period 20 in excel. 
  
  def get_cycle(c)
    address = NArray.int(13,5,41).indgen!
    subpops = address[true,true,c]
    new_subpops = Array.new
    subpops.each{|s| new_subpops[s] = self[s] unless self[s].nil? }
    new_subpops
  end


  def get_group(g)
    address = NArray.int(13,5,41).indgen!
    subpops = address[true,g,true]
    new_subpops = Array.new
    subpops.each{|s| new_subpops[s] = self[s] unless self[s].nil? }
    new_subpops
  end

  def get_disease_stage(d)
    address = NArray.int(13,5,41).indgen!
    subpops = address[d,true,true]
    new_subpops = Array.new
    subpops.each{|s| new_subpops[s] = self[s] unless self[s].nil? }
    new_subpops
  end
  
  alias :d :get_disease_stage
  alias :g :get_group
  alias :c :get_cycle
  
end


################################################################################################################################################################
########################################################## Step 2: Pull key inputs from JSON and model housekeeping variables ############################################################
################################################################################################################################################################

# ------------  json data version (online only)  ------------  #

#cgi = CGI.new('html5')
#json = cgi.params["custom"][0]
#json2 = cgi.params["InterventionSubpopulationNationCoverage"][0]
#json3 = cgi.params["RrrInterventionSubpopulation"][0] # currently not used

#custom = JSON.parse(json)
#intervention_subpopulation_coverage = JSON.parse(json2)

#intervention_subpopulation_coverage.collect{|item| item['SubpopulationId'] = item['SubpopulationId'] -1}


# ------------  file data version (for offline)  ------------  #

file = File.open('custom.json', 'r')
json = file.readlines                      #to rb1.9
custom = JSON.parse(json[0]) #to rb1.9

file2 = File.open('intervention_subpopulation_coverage.json', 'r')
json2 = file2.readlines  #to rb1.9
intervention_subpopulation_coverage = JSON.parse(json2[0]) #to rb1.9

intervention_subpopulation_coverage.collect{|item| item['SubpopulationId'] = item['SubpopulationId'] -1}


#IRB.start_session(binding)



# ------------  pre-build arrays  ------------  #

step= 0.5
cycles = Array(0..40)

#arrays

groups = custom.groups
disease_stages = custom.disease_stages
disease_and_treatment_stages = custom.disease_and_treatment_stages
PopulationSizeByGroup = custom.PopulationSizeByGroup
HivPrevalenceAdultsByGroup = custom.HivPrevalenceAdultsByGroup
HivPrevalence15yoByGroup = custom.HivPrevalence15yoByGroup
ProprtionDiseaseStage = custom.ProprtionDiseaseStage
InfectiousnessByDiseaseStage = custom.InfectiousnessByDiseaseStage
HivDeathRateByDiseaseStage= custom.HivDeathRateByDiseaseStage
HivDeathRateByDiseaseStageTx = custom.HivDeathRateByDiseaseStageTx
InitialTreatmentAccessByDiseaseStage = custom.InitialTreatmentAccessByDiseaseStage
TreatmentRecuitingRateByDiseaseStage = custom.TreatmentRecuitingRateByDiseaseStage

PopulationSizeByGroup.each{|item| item = item.to_f } #convert to float
HivPrevalenceAdultsByGroup.each{|item| item = item.to_f }
HivPrevalence15yoByGroup.each{|item| item = item.to_f } #convert to float
ProprtionDiseaseStage.each{|item| item = item.to_f } #convert to float
InfectiousnessByDiseaseStage.each{|item| item = item.to_f } #convert to float
HivDeathRateByDiseaseStage.each{|item| item = item.to_f } #convert to float
HivDeathRateByDiseaseStageTx.each{|item| item = item.to_f } #convert to float
InitialTreatmentAccessByDiseaseStage.each{|item| item = item.to_f } #convert to float
TreatmentRecuitingRateByDiseaseStage.each{|item| item = item.to_f } #convert to float

# ------------  asign variables from custom to global variables  ------------  #
 
EntryRateGenPop =                               custom.EntryRateGenPop.to_f 
MaturationRate =                                custom.MaturationRate.to_f  
DeathRateGeneralCauses =                        custom.DeathRateGeneralCauses.to_f  
LifeExpectancy =                                custom.LifeExpectancy.to_f 
SwInitiationRate =                              custom.SwInitiationRate.to_f  
SwQuitRate =                                    custom.SwQuitRate.to_f 
EntryRateMsm =                                  custom.EntryRateMsm.to_f 
EntryRateIdu =                                  custom.EntryRateIdu.to_f  
IduInitiationRate =	                            custom.IduInitiationRate.to_f 	
IduSpontaneousQuitRate	=                       custom.IduSpontaneousQuitRate.to_f	 
IduDeathRate	=                                 custom.IduDeathRate.to_f	
IncreaseInInfectiousnessHomosexual =            custom.IncreaseInInfectiousnessHomosexual.to_f 
DiseaseProgressionUntreatedAcuteToEarly =       custom.DiseaseProgressionUntreatedAcuteToEarly.to_f  
DiseaseProgressionUntreatedEarlyToMedium =      custom.DiseaseProgressionUntreatedEarlyToMedium.to_f  
DiseaseProgressionUntreatedMediumToLate =       custom.DiseaseProgressionUntreatedMediumToLate.to_f  
DiseaseProgressionUntreatedLateToAdvanced =     custom.DiseaseProgressionUntreatedLateToAdvanced.to_f  
DiseaseProgressionUntreatedAdvancedToAids =     custom.DiseaseProgressionUntreatedAdvancedToAids.to_f  
DiseaseProgressionTreatedAcuteToEarly =         custom.DiseaseProgressionTreatedAcuteToEarly.to_f 
DiseaseProgressionTreatedEarlyToMedium =        custom.DiseaseProgressionTreatedEarlyToMedium.to_f  
DiseaseProgressionTreatedMediumToLate =         custom.DiseaseProgressionTreatedMediumToLate.to_f  
DiseaseProgressionTreatedLateToAdvanced =       custom.DiseaseProgressionTreatedLateToAdvanced.to_f  
DiseaseProgressionTreatedAdvancedToAids =       custom.DiseaseProgressionTreatedAdvancedToAids.to_f  

GeneralNonSwPartnershipsYearlyUninfected =                custom.GeneralNonSwPartnershipsYearly.to_f  
GeneralNonSwPartnershipsYearlyInfected =                custom.GeneralNonSwPartnershipsYearly.to_f 
GeneralNonSwPartnershipsYearlyTreated =                custom.GeneralNonSwPartnershipsYearly.to_f 

GeneralCondomUseUninfected =                              custom.GeneralCondomUse.to_f  
GeneralCondomUseInfected =		custom.GeneralCondomUse.to_f 
GeneralCondomUseTreated =			custom.GeneralCondomUse.to_f 

GeneralCondomEffectiveness =                    custom.GeneralCondomEffectiveness.to_f  
SwProportionWhoUseServices =                    custom.SwProportionWhoUseServices.to_f  

SwPartnershipsYearlyUninfected =                          custom.SwPartnershipsYearly.to_f  
SwPartnershipsYearlyInfected =                          custom.SwPartnershipsYearly.to_f  
SwPartnershipsYearlyTreated =                          custom.SwPartnershipsYearly.to_f  

SwCondomUseRateUninfected =                               custom.SwCondomUseRate.to_f  
SwCondomUseRateInfected =                               custom.SwCondomUseRate.to_f  
SwCondomUseRateTreated =                               custom.SwCondomUseRate.to_f  

MsmPartnershipsYearlyUninfected =                         custom.MsmPartnershipsYearly.to_f 
MsmPartnershipsYearlyInfected =                         custom.MsmPartnershipsYearly.to_f 
MsmPartnershipsYearlyTreated =                         custom.MsmPartnershipsYearly.to_f 

MsmCondomUseRateUninfected =                              custom.MsmCondomUseRate.to_f  
MsmCondomUseRateInfected =                              custom.MsmCondomUseRate.to_f  
MsmCondomUseRateTreated =                              custom.MsmCondomUseRate.to_f  

TreatmentReductionOfInfectiousness =            custom.TreatmentReductionOfInfectiousness.to_f  
TreatmentQuitRate =                             custom.TreatmentQuitRate.to_f  
PercentOfIduSexPartners =	                      custom.PercentOfIduSexPartners.to_f 

IduPartnershipsYearlyUninfected	=                         custom.IduPartnershipsYearly.to_f	
IduPartnershipsYearlyInfected	=                         custom.IduPartnershipsYearly.to_f
IduPartnershipsYearlyTreated	=                         custom.IduPartnershipsYearly.to_f

IduCondomUseRateUninfected =                              custom.IduCondomUseRate.to_f  
IduCondomUseRateInfected =                              custom.IduCondomUseRate.to_f  
IduCondomUseRateTreated =                              custom.IduCondomUseRate.to_f  

AnnualNumberOfInjections =                      custom.AnnualNumberOfInjections.to_f  

PercentSharedInjectionsUninfected =                       custom.PercentSharedInjections.to_f  
PercentSharedInjectionsInfected =                       custom.PercentSharedInjections.to_f  
PercentSharedInjectionsTreated =                       custom.PercentSharedInjections.to_f  

PercentMaleIdus =                               custom.PercentMaleIdus.to_f  
InfectiousnessInSharedInjection =               custom.InfectiousnessInSharedInjection.to_f 
CircEffectiveness = 				custom.CircEffectiveness.to_f


################################################################################################################################################################
######################## Step 2.9: Adjust behavioral and mortality rates by *intervention* for *key group parameters*   ########################################
################################################################################################################################################################


## Use choo to select intervention (for testing purposes): 0= No interventions, 1 = circumcision, 2 = VCT, 3 = Media, 4 = SWP, 5 = NEP, 6 = MMT, 7 = CD4 Monitoring, 8 = VL Monitoring; look right before sections 3.5 and 10 to find circumcision formulas and remove choo
##choo=[1,2,3,4,5,6,7,8]
   choo=[0] 
groups.each_with_index do |group, g|

for disstage in 0..2

    ##Compute total RRR from all interventions by group and calculate new condom use rate by group. This assumes interventions that do not affect a certain population group have RRR of 0. 
    ## Outcome type 1 = Condom use rate
    
    TotalGroupRrr = intervention_subpopulation_coverage.select{|item| item.SubpopulationId == g && item.HivStatus==disstage && item.RrrTypeId ==1 && choo.include?(item.InterventionId)}.inject(1){|memo, n| memo *(1- n.Rrr.to_f * n.Coverage.to_f) }
  
if disstage==0  
    if g==0 
      GeneralCondomUseUninfected = 1 - (1-GeneralCondomUseUninfected) * TotalGroupRrr
    elsif g==2
      SwCondomUseRateUninfected = 1 - (1-SwCondomUseRateUninfected) * TotalGroupRrr
    elsif g==3
      MsmCondomUseRateUninfected = 1 - (1-MsmCondomUseRateUninfected) * TotalGroupRrr
    elsif g==4
      IduCondomUseRateUninfected = 1 - (1-IduCondomUseRateUninfected) * TotalGroupRrr
    end
    
    elsif disstage==1
	   if g==0 
      GeneralCondomUseInfected = 1 - (1-GeneralCondomUseInfected) * TotalGroupRrr
    elsif g==2
      SwCondomUseRateInfected = 1 - (1-SwCondomUseRateInfected) * TotalGroupRrr
    elsif g==3
      MsmCondomUseRateInfected = 1 - (1-MsmCondomUseRateInfected) * TotalGroupRrr
    elsif g==4
      IduCondomUseRateInfected = 1 - (1-IduCondomUseRateInfected) * TotalGroupRrr
    end 
 
else
	if g==0 ##|| g==1
      GeneralCondomUseTreated = 1 - (1-GeneralCondomUseTreated) * TotalGroupRrr
    elsif g==2
      SwCondomUseRateTreated = 1 - (1-SwCondomUseRateTreated) * TotalGroupRrr
    elsif g==3
      MsmCondomUseRateTreated = 1 - (1-MsmCondomUseRateTreated) * TotalGroupRrr
    elsif g==4
      IduCondomUseRateTreated = 1 - (1-IduCondomUseRateTreated) * TotalGroupRrr
    end
    end
    
     ##Compute total RRR from all interventions by group and calculate new number of partnerships by group. Not computed for Gen Pop Females, since their partnerships are the "absorbing bucket" for partnerships. This assumes interventions that do not affect a certain population group have RRR of 0. 
      ## Outcome type 2 = # partnerships
      
    TotalGroupRrr = intervention_subpopulation_coverage.select{|item| item.SubpopulationId == g && item.HivStatus==disstage&& item.RrrTypeId ==2 && choo.include?(item.InterventionId)}.inject(1){|memo, n| memo *(1- n.Rrr.to_f * n.Coverage.to_f) }
    
if disstage==0   
   if g==0 
      GeneralNonSwPartnershipsYearlyUninfected = GeneralNonSwPartnershipsYearlyUninfected* TotalGroupRrr
    elsif g==2
      SwPartnershipsYearlyUninfected = SwPartnershipsYearlyUninfected * TotalGroupRrr
    elsif g==3
      MsmPartnershipsYearlyUninfected = MsmPartnershipsYearlyUninfected * TotalGroupRrr
    elsif g==4
      IduPartnershipsYearlyUninfected = IduPartnershipsYearlyUninfected * TotalGroupRrr
    end  
    
    elsif disstage ==1
	if g==0 
      GeneralNonSwPartnershipsYearlyInfected = GeneralNonSwPartnershipsYearlyInfected* TotalGroupRrr
    elsif g==2
      SwPartnershipsYearlyInfected = SwPartnershipsYearlyInfected * TotalGroupRrr
    elsif g==3
      MsmPartnershipsYearlyInfected = MsmPartnershipsYearlyInfected * TotalGroupRrr
    elsif g==4
      IduPartnershipsYearlyInfected = IduPartnershipsYearlyInfected * TotalGroupRrr
    end  
else 
	if g==0 
      GeneralNonSwPartnershipsYearlyTreated = GeneralNonSwPartnershipsYearlyTreated* TotalGroupRrr
    elsif g==2
      SwPartnershipsYearlyTreated = SwPartnershipsYearlyTreated * TotalGroupRrr
    elsif g==3
      MsmPartnershipsYearlyTreated = MsmPartnershipsYearlyTreated * TotalGroupRrr
    elsif g==4
      IduPartnershipsYearlyTreated = IduPartnershipsYearlyTreated * TotalGroupRrr
    end  
end


 ##Compute total RRR from all interventions by group and calculate new number of annual injections for IDUs. This assumes interventions that do not affect a certain population group have RRR of 0. 
    ## Outcome type 3 = # injections
    
   TotalGroupRrr = intervention_subpopulation_coverage.select{|item| item.SubpopulationId == g && item.HivStatus==disstage&& item.RrrTypeId ==3 && choo.include?(item.InterventionId)}.inject(1){|memo, n| memo *(1- n.Rrr.to_f * n.Coverage.to_f) }
    
    if g==4 && disstage==0
      AnnualNumberOfRiskyInjectionsUninfected = AnnualNumberOfInjections * PercentSharedInjectionsUninfected * TotalGroupRrr
      elsif g==4&&disstage==1
	       AnnualNumberOfRiskyInjectionsInfected = AnnualNumberOfInjections * PercentSharedInjectionsInfected * TotalGroupRrr
	       else
		        AnnualNumberOfRiskyInjectionsTreated = AnnualNumberOfInjections * PercentSharedInjectionsTreated * TotalGroupRrr
    end 

end

     ##Compute total RRR from all interventions and progression rates by group. All disease stages are affected similarly. This assumes interventions that do not affect a certain population group have RRR of 0. 
      ## Outcome type 4 = Disease progression rates
      
    TotalGroupRrr = intervention_subpopulation_coverage.select{|item| item.SubpopulationId == g && item.HivStatus==2&& item.RrrTypeId ==4 && choo.include?(item.InterventionId)}.inject(1){|memo, n| memo *(1- n.Rrr.to_f * n.Coverage.to_f) }
    
    ## This is a simple fix - if we want differential RRR by group, we should assign a progression rate to each group and modify if accordingly. 
if g==0
	DiseaseProgressionTreatedAcuteToEarly = DiseaseProgressionTreatedAcuteToEarly * TotalGroupRrr
	DiseaseProgressionTreatedEarlyToMedium = DiseaseProgressionTreatedEarlyToMedium * TotalGroupRrr 
	DiseaseProgressionTreatedMediumToLate = DiseaseProgressionTreatedMediumToLate * TotalGroupRrr  
	DiseaseProgressionTreatedLateToAdvanced = DiseaseProgressionTreatedLateToAdvanced * TotalGroupRrr  
	DiseaseProgressionTreatedAdvancedToAids = DiseaseProgressionTreatedAdvancedToAids * TotalGroupRrr 
end

     ##Compute total RRR from all interventions and HIV mortality rates by group. All disease stages are affected similarly. This assumes interventions that do not affect a certain population group have RRR of 0. 
      ## Outcome type 5 = Disease mortality rate
      
    TotalGroupRrr = intervention_subpopulation_coverage.select{|item| item.SubpopulationId == g && item.HivStatus==2&& item.RrrTypeId ==5 && choo.include?(item.InterventionId)}.inject(1){|memo, n| memo *(1- n.Rrr.to_f * n.Coverage.to_f) }

if g==0

##this formulation doesn;t work
##HivDeathRateByDiseaseStageTx.each{|item| item = (item.to_f * TotalGroupRrr )}

for j in 0..6
	HivDeathRateByDiseaseStageTx[j]=HivDeathRateByDiseaseStageTx[j] * TotalGroupRrr
	
end

end

end




################################################################################################################################################################
################################# Step 3: Compute intial populations, Vary parameters by group and disease stage. ##############################################
################################################################################################################################################################



class EpiSubpopulationCycle
  
  #the below specifies the different parameters a EpiSubpopulation object can hold
  attr_accessor :DynamicsGeneral, :DynamicsArray, :TotalDynamics, :SwAndIduDynamics, :TreatmentDynamics, :DiseaseProgressionEntries, :DiseaseProgressionExits, :TotalDiseaseProgressionExits, :Cycle, :InfectiousnessInSharedInjection, :Scr, :Probability, :GroupId, :DiseaseStageId, :DiseaseProgressionUntreatedAcuteToEarly, :DiseaseProgressionTreatedAcuteToEarly, :AnnualNumberOfInjections, :PercentSharedInjections, :IduSpontaneousQuitRate, :Id, :Group, :DiseaseStage, :Population, :Cycle, :Partnerships, :CondomUse, :CompositePartnerships, :Infectiousness, :TreatmentRecuitingRate, :HivDeathRate, :EntryRate, :PopulationSize, :PopulationSizeGenPopMen, :PopulationSizeGenPopWomen, :PopulationSizeSwWomen, :PopulationSizeMsm, :EntryRateGenPop, :MaturationRate, :DeathRateGeneralCauses, :LifeExpectancy, :SwInitiationRate, :SwQuitRate, :EntryRateMsm, :HivPrevalenceAdultsGenPopMen, :HivPrevalenceAdultsGenPopWomen, :HivPrevalenceAdultsSwWomen, :HivPrevalenceAdultsMsm, :HivPrevalence15yoGenPopMen, :HivPrevalence15yoGenPopWomen, :ProprtionDiseaseStageEarly, :ProprtionDiseaseStageMedium, :ProprtionDiseaseStageLate, :ProprtionDiseaseStageAdvanced, :ProprtionDiseaseStageAids, :InfectiousnessEarly, :InfectiousnessMedium, :InfectiousnessLate, :InfectiousnessLate, :InfectiousnessAids, :IncreaseInInfectiousnessHomosexual, :DiseaseProgressionUntreatedEarlyToMedium, :DiseaseProgressionUntreatedMediumToLate, :DiseaseProgressionUntreatedLateToAdvanced, :DiseaseProgressionUntreatedAdvancedToAids, :DiseaseProgressionTreatedEarlyToMedium, :DiseaseProgressionTreatedMediumToLate, :DiseaseProgressionTreatedLateToAdvanced, :DiseaseProgressionTreatedAdvancedToAids, :HivDeathRateUntreatedEarly, :HivDeathRateUntreatedMedium, :HivDeathRateUntreatedLate, :HivDeathRateUntreatedAdvanced, :HivDeathRateUntreatedAids, :HivDeathRateTreatedEarly, :HivDeathRateTreatedMedium, :HivDeathRateTreatedLate, :HivDeathRateTreatedAdvanced, :HivDeathRateTreatedAids, :GeneralNonSwPartnershipsYearly, :GeneralCondomUse, :GeneralCondomEffectiveness, :SwProportionWhoUseServices, :SwPartnershipsYearly, :SwCondomUseRate, :MsmPartnershipsYearly, :MsmCondomUseRate, :TreatmentReductionOfInfectiousness, :InitialTreatmentAccessEarly, :InitialTreatmentAccessMedium, :InitialTreatmentAccessLate, :InitialTreatmentAccessAdvanced, :InitialTreatmentAccessAids, :TreatmentRecuitingRateEarly, :TreatmentRecuitingRateMedium, :TreatmentRecuitingRateLate, :TreatmentRecuitingRateAdvanced, :TreatmentRecuitingRateAids, :TreatmentQuitRate, :PropMaleCirc

end

subpopulations = Array.new
s = 0

groups.each_with_index do |group,g| 
  
  disease_and_treatment_stages.each_with_index do |disease_stage,d|
    
    this_subpopulation = EpiSubpopulationCycle.new
    
    this_subpopulation.Cycle = 0 #this first array, subpopulations, will hold all the data for the 0th cycle
    
    # --- fill in their populations ---
    
    #uninfected
    
    if d==0 
      this_subpopulation.Population = PopulationSizeByGroup[g]*(1-HivPrevalenceAdultsByGroup[g])
    end
    
    #acute to aids, no treatment 
    
    if  d>0 and d<7
      this_subpopulation.Population  = PopulationSizeByGroup[g]*HivPrevalenceAdultsByGroup[g]*ProprtionDiseaseStage[d]*(1-InitialTreatmentAccessByDiseaseStage[d])
    end
    
    #treatment
      
    if d>6
      ds = d - 6;
      this_subpopulation.Population = PopulationSizeByGroup[g]*HivPrevalenceAdultsByGroup[g]*ProprtionDiseaseStage[ds]*(InitialTreatmentAccessByDiseaseStage[ds])
    end
      
    # ------------------------------- basic info -------------------------------
    
    this_subpopulation.Id = s
    this_subpopulation.Group = group
    this_subpopulation.GroupId = g
    this_subpopulation.DiseaseStage = disease_stage
    this_subpopulation.DiseaseStageId = d
    
    # ------------------------------- these do not vary  -------------------------------
      
    this_subpopulation.MaturationRate = MaturationRate
    this_subpopulation.DeathRateGeneralCauses = DeathRateGeneralCauses  #idus override below
    this_subpopulation.InfectiousnessInSharedInjection = 0.0 #idus override below
    this_subpopulation.GeneralCondomEffectiveness = GeneralCondomEffectiveness
    
     # ------------------------------- these vary by GROUP and Uninfected/Infected/Treated-------------------------------
    
    #gen pop male
    
    if g==0
	    if d==0
      this_subpopulation.Partnerships = GeneralNonSwPartnershipsYearlyUninfected
      this_subpopulation.CondomUse = GeneralCondomUseUninfected
        elsif d>0 && d<7
		this_subpopulation.Partnerships = GeneralNonSwPartnershipsYearlyInfected
		this_subpopulation.CondomUse = GeneralCondomUseInfected
		elsif d>6
			this_subpopulation.Partnerships = GeneralNonSwPartnershipsYearlyTreated
			this_subpopulation.CondomUse = GeneralCondomUseTreated
		end
		
    end
    
    #gen pop female
    
    if g==1
        this_subpopulation.SwInitiationRate = SwInitiationRate
	
        	    if d==0
      this_subpopulation.Partnerships = GeneralNonSwPartnershipsYearlyUninfected
      this_subpopulation.CondomUse = GeneralCondomUseUninfected
        elsif d>0 && d<7
		this_subpopulation.Partnerships = GeneralNonSwPartnershipsYearlyInfected
		this_subpopulation.CondomUse = GeneralCondomUseInfected
		elsif d>6
			this_subpopulation.Partnerships = GeneralNonSwPartnershipsYearlyTreated
			this_subpopulation.CondomUse = GeneralCondomUseTreated
		end
        #genpop females are not assigned partners -- they get assigned the "remaining" partnerships from males
    end
    
    #sws
    
    if g==2
      this_subpopulation.SwQuitRate = SwQuitRate
      
     	    if d==0
      this_subpopulation.Partnerships = SwPartnershipsYearlyUninfected
      this_subpopulation.CondomUse = SwCondomUseRateUninfected
        elsif d>0 && d<7
		this_subpopulation.Partnerships = SwPartnershipsYearlyInfected
		this_subpopulation.CondomUse = SwCondomUseRateInfected
		elsif d>6
			this_subpopulation.Partnerships = SwPartnershipsYearlyTreated
			this_subpopulation.CondomUse = SwCondomUseRateTreated
		end
    end

    #msm

    if g==3

       	    if d==0
      this_subpopulation.Partnerships = MsmPartnershipsYearlyUninfected
      this_subpopulation.CondomUse = MsmCondomUseRateUninfected
        elsif d>0 && d<7
		this_subpopulation.Partnerships = MsmPartnershipsYearlyInfected
		this_subpopulation.CondomUse = MsmCondomUseRateInfected
		elsif d>6
			this_subpopulation.Partnerships = MsmPartnershipsYearlyTreated
			this_subpopulation.CondomUse = MsmCondomUseRateTreated
		end
      
    end
    
    #idu
    
    if g==4

      this_subpopulation.IduSpontaneousQuitRate	= IduSpontaneousQuitRate
      
  	    if d==0
      this_subpopulation.Partnerships = IduPartnershipsYearlyUninfected
      this_subpopulation.CondomUse = IduCondomUseRateUninfected
            #needle sharing
      this_subpopulation.AnnualNumberOfInjections = AnnualNumberOfRiskyInjectionsUninfected
      this_subpopulation.PercentSharedInjections = PercentSharedInjectionsUninfected
      
        elsif d>0 && d<7
		this_subpopulation.Partnerships = IduPartnershipsYearlyInfected
		this_subpopulation.CondomUse = IduCondomUseRateInfected
		
		      #needle sharing
      this_subpopulation.AnnualNumberOfInjections = AnnualNumberOfRiskyInjectionsInfected
      this_subpopulation.PercentSharedInjections = PercentSharedInjectionsInfected
      
		elsif d>6
			this_subpopulation.Partnerships = IduPartnershipsYearlyTreated
			this_subpopulation.CondomUse = IduCondomUseRateTreated
			
			      #needle sharing
      this_subpopulation.AnnualNumberOfInjections = AnnualNumberOfRiskyInjectionsTreated
      this_subpopulation.PercentSharedInjections = PercentSharedInjectionsTreated
      
		end
		
      this_subpopulation.DeathRateGeneralCauses = IduDeathRate
        
    end
    
    this_subpopulation.CompositePartnerships = (this_subpopulation.CondomUse) * (this_subpopulation.Partnerships) * (1-GeneralCondomEffectiveness) + (1-this_subpopulation.CondomUse) * (this_subpopulation.Partnerships)
    

    # ------------------------------- these vary by DISEASE STAGE -------------------------------
    
    # uninfected
    
    if d==0
        this_subpopulation.Infectiousness = 0
        this_subpopulation.HivDeathRate = 0
        this_subpopulation.InfectiousnessInSharedInjection = 0.0
    end

    #acute to aids, no treatment
    
    if d>0 && d<7
        this_subpopulation.Infectiousness = InfectiousnessByDiseaseStage[d]
        this_subpopulation.TreatmentRecuitingRate = TreatmentRecuitingRateByDiseaseStage[d]
        this_subpopulation.HivDeathRate = HivDeathRateByDiseaseStage[d]
        this_subpopulation.InfectiousnessInSharedInjection = InfectiousnessInSharedInjection
    end

    # early to aids, with treatment
    
    if d>6
        ds = d - 6;
        this_subpopulation.Infectiousness = InfectiousnessByDiseaseStage[ds]*(1-TreatmentReductionOfInfectiousness)
        this_subpopulation.TreatmentQuitRate = TreatmentQuitRate
        this_subpopulation.HivDeathRate = HivDeathRateByDiseaseStageTx[ds]
        this_subpopulation.InfectiousnessInSharedInjection = InfectiousnessInSharedInjection / 2.0
    end

    #acute

    if d==0
      this_subpopulation.DiseaseProgressionExits = 0
    end

    
    if d==1
      this_subpopulation.DiseaseProgressionUntreatedAcuteToEarly = DiseaseProgressionUntreatedAcuteToEarly
      this_subpopulation.DiseaseProgressionExits = DiseaseProgressionUntreatedAcuteToEarly
    end
    
    #early
    
    if d==2
      this_subpopulation.DiseaseProgressionUntreatedEarlyToMedium = DiseaseProgressionUntreatedEarlyToMedium
      this_subpopulation.DiseaseProgressionExits = DiseaseProgressionUntreatedEarlyToMedium
    end
    
    #medium

    if d==3 
      this_subpopulation.DiseaseProgressionUntreatedMediumToLate = DiseaseProgressionUntreatedMediumToLate
      this_subpopulation.DiseaseProgressionExits = DiseaseProgressionUntreatedMediumToLate
    end
    
    #late
    
    if d==4
      this_subpopulation.DiseaseProgressionUntreatedLateToAdvanced = DiseaseProgressionUntreatedLateToAdvanced
      this_subpopulation.DiseaseProgressionExits = DiseaseProgressionUntreatedLateToAdvanced
    end
    
    #advanced
    
    if d==5
      this_subpopulation.DiseaseProgressionUntreatedAdvancedToAids = DiseaseProgressionUntreatedAdvancedToAids
      this_subpopulation.DiseaseProgressionExits = DiseaseProgressionUntreatedAdvancedToAids
    end
    
    if d==6
      # aids, 6, has no transition
      this_subpopulation.DiseaseProgressionExits = 0.0
    end
    
    if d==7
      this_subpopulation.DiseaseProgressionTreatedAcuteToEarly = DiseaseProgressionTreatedAcuteToEarly
      this_subpopulation.DiseaseProgressionExits = DiseaseProgressionTreatedAcuteToEarly
    end

    if d==8
      this_subpopulation.DiseaseProgressionTreatedEarlyToMedium = DiseaseProgressionTreatedEarlyToMedium
      this_subpopulation.DiseaseProgressionExits = DiseaseProgressionTreatedEarlyToMedium
    end
    
    if d==9
      this_subpopulation.DiseaseProgressionTreatedMediumToLate = DiseaseProgressionTreatedMediumToLate
      this_subpopulation.DiseaseProgressionExits = DiseaseProgressionTreatedMediumToLate
    end
    
    if d==10
      this_subpopulation.DiseaseProgressionTreatedLateToAdvanced = DiseaseProgressionTreatedLateToAdvanced
      this_subpopulation.DiseaseProgressionExits = DiseaseProgressionTreatedLateToAdvanced
    end
    
    if d==11
       this_subpopulation.DiseaseProgressionTreatedAdvancedToAids = DiseaseProgressionTreatedAdvancedToAids
       this_subpopulation.DiseaseProgressionExits = DiseaseProgressionTreatedAdvancedToAids
     end
     
     if d==12
         this_subpopulation.DiseaseProgressionExits = 0.0
         # treated aids, 12, has no transition
     end
    
     # --- these vary by both group and disease stage, they are one-off ones ---
    
    
     # MSM men enter through uninfected group

     if g == 3 && d == 0
         this_subpopulation.EntryRateMsm = EntryRateMsm
         this_subpopulation.EntryRate = EntryRateMsm #note ask Sabina why there are two MSM entry numbers
    elsif g == 4 && d == 0
        this_subpopulation.EntryRate = EntryRateIdu
    else
        this_subpopulation.EntryRate = 0    
    end

     # msm more infectios (anal sex)

     if g == 3
         this_subpopulation.Infectiousness = this_subpopulation.Infectiousness * ( 1 + IncreaseInInfectiousnessHomosexual )
     end

     #entry rates for genpop men and women

     if g==0 && d==0 || g==1 && d==0
         this_subpopulation.EntryRate = EntryRateGenPop/2
     end

     # determine proportion of circumcised men in initial step

if choo.include?(1)
   this_subpopulation.PropMaleCirc= step *  intervention_subpopulation_coverage.select{|item| item.Id == 1}[0].Coverage.to_f 
   else
	   this_subpopulation.PropMaleCirc=0
   end
   
    # add subpopulation to subpopulations array
      
    subpopulations << this_subpopulation
    s = s + 1
      
      
    subpopulations.g(0).compact.inject(0){|memo, n| memo + n.Population }
     
  end #disease and treatment stages
  
end #groups


################################################################################################################################################################
#################################################### Step 3.5: Adjust behavioral and mortality rates by *intervention* ###########################################
################################################################################################################################################################

#discuss with Jim and Sabina about where these modifications should go



#intervention_subpopulation_coverage.each_with_index do |isc,i|
#  
#  testRrr = isc.TestRrr
#  testRrr = testRrr.to_f
#
#  if isc.InterventionId == 1
#    subpopulations.select{|item| item.GroupId == isc.SubpopulationId}.each{|anotheritem| anotheritem.CondomUse = 1 - ((1-anotheritem.CondomUse) * ( (1- testRrr* isc.Coverage.to_f/100) ) ) }
#  end
#end #ICS


################################################################################################################################################################
############################################# Step 4: Calculate partnerships per female and per male client of sex #############################################
################################################################################################################################################################

cycles.each_with_index do |cycle,c|

  #determine average risky partnerships for women from men

  total_men = subpopulations.g(0).c(c).compact.inject(0){|memo, n| memo + n.Population }
  total_women = subpopulations.g(1).c(c).compact.inject(0){|memo, n| memo + n.Population }
  total_men_partnerships = subpopulations.g(0).c(c).compact.inject(0){|memo, n| memo + (n.Population) * (n.Partnerships) }
  total_idus = subpopulations.g(4).c(c).compact.inject(0){|memo, n| memo + n.Population }
  ##total_male_idus = total_idus * PercentMaleIdus
  ##total_female_idus = total_idus * (1-PercentMaleIdus)
  total_men_partnerships_lost_to_idu = (1-PercentOfIduSexPartners) * (1-PercentMaleIdus) * subpopulations.g(4).c(c).compact.inject(0){|memo, n| memo + (n.Population) * (n.Partnerships)}
  total_partnerships_offered_to_gp_women_from_men = total_men_partnerships - total_men_partnerships_lost_to_idu
  average_number_of_parternships_for_women_from_men = total_partnerships_offered_to_gp_women_from_men / total_women
  ##average_number_of_risky_parternships_for_women_from_men = average_number_of_parternships_for_women_from_men * (1-GeneralCondomUse + GeneralCondomUse*(1-GeneralCondomEffectiveness)) 
  ##Include condom use in SCR computation (varies by uninfected/infected/treated)

  #determine average risky partnerships for women from idu

   total_partnerships_offered_from_idus_to_females = (1-PercentOfIduSexPartners) * PercentMaleIdus * subpopulations.g(4).c(c).compact.inject(0){|memo, n| memo + (n.Population) * (n.Partnerships)}
   average_partnerships_from_idus_to_female = total_partnerships_offered_from_idus_to_females / total_women
   ##average_risky_partnerships_from_idus_to_female = average_partnerships_from_idus_to_female * (1-IduCondomUseRate+IduCondomUseRate*(1-GeneralCondomEffectiveness))
  ##Include condom use in SCR computation (varies by uninfected/infected/treated)
  
   #determine average sex worker partnerships per sw client

   total_sws = subpopulations.g(2).c(c).compact.inject(0){|memo, n| memo + n.Population }
   total_male_clients = total_men * SwProportionWhoUseServices
   total_sw_partnerships_offered = subpopulations.g(2).c(c).compact.inject(0){|memo, n| memo + (n.Population) * (n.Partnerships) }
   average_sw_partnerships_per_using_male = total_sw_partnerships_offered / total_male_clients

   #determine average number of risky sw partnerships per sw client - Not needed,  condom rate included in SCR

   ##average_risky_sw_partnerships_per_using_male = average_sw_partnerships_per_using_male * (1-SwCondomUseRate + SwCondomUseRate*(1-GeneralCondomEffectiveness))

   ################################################################################################################################################################
   ############################################# # Step 5: Calculate proportion of people in each disease state per group, ie "probability tab" ###################
   ################################################################################################################################################################


   subpopulations.c(c).compact.each_with_index do |this_subpopulation,s|
     this_subpopulation.Probability = this_subpopulation.Population / subpopulations.g(this_subpopulation.GroupId).c(c).compact.inject(0){|memo, n| memo + n.Population }
   end

  ################################################################################################################################################################
  ############################################# # Step 6: Calculate sufficient contact rate from each disease state to each ######################################
  ################################################################################################################################################################

  # in order to calculate the SCR for a certain subpop, you're finding the chance that the uninfected subpopulation in question might be infected by the particular disease state from a different
  # group. For example, the SCR assigned under acute male gen pop is the chance that an uninfected genpop men will be infected by an acute SW, IDU, or genpop woman. Therefore, I've assigned
  # "mate" subpopulations designating which subpops the SCR refers to. In the above sample, acute GP M has a SW, IDU and GP W mate groups.


  subpopulations.select{|item| item.Cycle == c}.each_with_index do |this_subpopulation,s|
  
    if this_subpopulation.GroupId == 0
    
      # SCR GenPopMen

      ##Not needed ##total_female_idus = total_idus * (1-PercentMaleIdus) 
	total_men_partnerships = subpopulations.g(0).c(c).compact.inject(0){|memo, n| memo + (n.Population) * (n.Partnerships) }
      total_female_idu_partnerships = (1-PercentOfIduSexPartners) * (1-PercentMaleIdus) * subpopulations.g(4).c(c).compact.inject(0){|memo, n| memo + (n.Population) *(n.Partnerships)}
      probability_of_choosing_idu_partner = total_female_idu_partnerships / total_men_partnerships
      probability_of_choosing_non_idu_partner = 1 - probability_of_choosing_idu_partner
    
      genpop_female_mate = subpopulations.g(1).d(this_subpopulation.DiseaseStageId).c(c).compact[0]
      idu_mate = subpopulations.g(4).d(this_subpopulation.DiseaseStageId).c(c).compact[0]
      sw_mate = subpopulations.g(2).d(this_subpopulation.DiseaseStageId).c(c).compact[0]
    
    #Need to multiply by partnerships for uninfected men
    uninfected_men=subpopulations.g(0).d(0).c(c).compact[0]
    
      scr_for_males_from_genpop_females = probability_of_choosing_non_idu_partner * genpop_female_mate.Probability * uninfected_men.CompositePartnerships * genpop_female_mate.Infectiousness
      scr_for_males_from_idu_females = probability_of_choosing_idu_partner  * idu_mate.Probability * uninfected_men.CompositePartnerships * idu_mate.Infectiousness
      scr_for_males_from_sws = sw_mate.Probability * average_sw_partnerships_per_using_male *(1-sw_mate.CondomUse * GeneralCondomEffectiveness) * sw_mate.Infectiousness * SwProportionWhoUseServices
	
      total_scr_male = scr_for_males_from_genpop_females + scr_for_males_from_idu_females + scr_for_males_from_sws
    
    #include circumcision effects
      this_subpopulation.Scr = total_scr_male * (1-this_subpopulation.PropMaleCirc.to_f * CircEffectiveness.to_f)
    
    end

      if this_subpopulation.GroupId == 1

        #scr genpop women

        #from idu males + genpop males

        ##total_male_idus = total_idus * PercentMaleIdus
        total_male_idu_partnerships =  (1-PercentOfIduSexPartners) * PercentMaleIdus * subpopulations.g(4).c(c).compact.inject(0){|memo, n| memo + (n.Population) *(n.Partnerships)}
        total_female_idu_partnerships =  (1-PercentOfIduSexPartners)*(1-PercentMaleIdus )* subpopulations.g(4).c(c).compact.inject(0){|memo, n| memo + (n.Population) * (n.Partnerships)}
        male_idu_partnerships_offered_to_females = (1-PercentOfIduSexPartners) * PercentMaleIdus * subpopulations.g(4).c(c).compact.inject(0){|memo, n| memo + (n.Population) * (n.Partnerships)}
        total_female_partnerships = total_men_partnerships - total_female_idu_partnerships + male_idu_partnerships_offered_to_females
        probability_of_choosing_idu_partner = total_male_idu_partnerships / total_female_partnerships
        probability_of_choosing_non_idu_partner = ( 1 - probability_of_choosing_idu_partner )

        male_mate = subpopulations.g(0).d(this_subpopulation.DiseaseStageId).c(c).compact[0]
        idu_mate = subpopulations.g(4).d(this_subpopulation.DiseaseStageId).c(c).compact[0]

        src_for_females_from_male = probability_of_choosing_non_idu_partner * male_mate.Probability * average_number_of_parternships_for_women_from_men *(1-(this_subpopulation.CondomUse) *GeneralCondomEffectiveness) * male_mate.Infectiousness
        src_for_females_from_idu_male = probability_of_choosing_idu_partner * idu_mate.Probability * average_partnerships_from_idus_to_female *(1-(idu_mate.CondomUse)*GeneralCondomEffectiveness)* idu_mate.Infectiousness 
        total_src_for_females = src_for_females_from_male + src_for_females_from_idu_male

        this_subpopulation.Scr = total_src_for_females

        if s == 14
            #IRB.start_session(binding)
        end


      end

      if this_subpopulation.GroupId == 2

       #=Probability!B4*('By group'!AB$25+'By group'!AB$26)*'By group'!B$23

       male_mate = subpopulations.g(0).d(this_subpopulation.DiseaseStageId).c(c).compact[0]
       
       #need to use partnerships per uninfected SW
      uninfected_sw=subpopulations.g(2).d(0).c(c).compact[0]
       
       src_for_sw_from_male = male_mate.Probability * uninfected_sw.CompositePartnerships * male_mate.Infectiousness

       this_subpopulation.Scr = src_for_sw_from_male

      end

      if this_subpopulation.GroupId == 3 

        #=Probability!AO5*('By group'!$AO$26+'By group'!$AO$25)*'By group'!AO$23
	#Need to use number of partnerships for uninfected MSM
	uninfected_msm=subpopulations.g(3).d(0).c(c).compact[0]
        msm_mate = subpopulations.g(3).d(this_subpopulation.DiseaseStageId).c(c).compact[0]
       src_for_msm_from_msm = msm_mate.Probability * uninfected_msm.CompositePartnerships * msm_mate.Infectiousness
        this_subpopulation.Scr = src_for_msm_from_msm  

      end 

      if this_subpopulation.GroupId == 4 

        #=prop_idu_part*Probability!BB4*'By group'!$BB$24*'By group'!BB$23+  Probability!BB4*'By group'!$BB$30*'By group'!BB$31

        #Probability!BC4*('By group'!$BB$26+'By group'!$BB$25)*'By group'!BC$23 + Probability!BC4*'By group'!$BB$30*'By group'!BC$31
	uninfected_idu=subpopulations.g(4).d(0).c(c).compact[0]
	
        scr_from_sex = PercentOfIduSexPartners * this_subpopulation.Probability * uninfected_idu.CompositePartnerships * this_subpopulation.Infectiousness 
        scr_from_needle_sharing = this_subpopulation.Probability * uninfected_idu.AnnualNumberOfInjections * this_subpopulation.InfectiousnessInSharedInjection
# Annualnumberofinjecitons already includes percentake shared (computed in section 3 as risky injections)

        this_subpopulation.Scr = scr_from_sex + scr_from_needle_sharing

      end

    end #subpopulation

  ################################################################################################################################################################
  ########################################################### Step 8: Develop dynamics for year 0 to year 1 ######################################################
  ################################################################################################################################################################

  subpopulations.c(c).compact.each_with_index do |this_subpopulation,s|
  
    #General/IDU death rate, maturation, HIV death rate, entry
  
    total_population = subpopulations.c(c).compact.inject(0){|memo, n| memo + n.Population }
    general_deaths = step * ((-this_subpopulation.MaturationRate + -this_subpopulation.DeathRateGeneralCauses + -this_subpopulation.HivDeathRate) * this_subpopulation.Population + this_subpopulation.EntryRate * total_population)
    
    this_subpopulation.DynamicsGeneral = general_deaths
    #subpopulations.select{|item| item.GroupId == this_subpopulation.GroupId}.each{|thing| print thing.Id, thing.DiseaseProgressionExits}


      #-------Disease progression exits

      if this_subpopulation.DiseaseStageId == 0.0
        total_scr = subpopulations.g(this_subpopulation.GroupId).c(c).compact.inject(0){|memo, n| memo + n.Scr.to_f }
		disease_progression_exits = step * -total_scr * this_subpopulation.Population
      else
        #=step*-SUM('By group'!C$11:C$21)*Population!C4
        disease_progression_exits = step * -this_subpopulation.DiseaseProgressionExits * this_subpopulation.Population
      end

      if s == 1 && c==1
        #IRB.start_session(binding)
      end

      this_subpopulation.TotalDiseaseProgressionExits = disease_progression_exits

      #--------Disease progression entries


      if this_subpopulation.DiseaseStageId == 0 
         disease_progression_entries = 0.0
      else
        from = subpopulations.g(this_subpopulation.GroupId).d(this_subpopulation.DiseaseStageId - 1).c(c).compact[0]
        disease_progression_entries = -from.TotalDiseaseProgressionExits
      end


      this_subpopulation.DiseaseProgressionEntries = disease_progression_entries

  
        #--------Treatment dynamics


        if this_subpopulation.DiseaseStageId == 0 
           treatment_dynamics = 0.0
        elsif this_subpopulation.DiseaseStageId > 0 && this_subpopulation.DiseaseStageId < 7
          treated_pop = subpopulations.g(this_subpopulation.GroupId).d(this_subpopulation.DiseaseStageId + 6).c(c).compact.first
          treatment_dynamics = step * (-this_subpopulation.Population * TreatmentRecuitingRateByDiseaseStage[this_subpopulation.DiseaseStageId] + treated_pop.Population * TreatmentQuitRate)
        else
          untreated_pop = subpopulations.g(this_subpopulation.GroupId).d(this_subpopulation.DiseaseStageId + -6).c(c).compact.first
          treatment_dynamics = -untreated_pop.TreatmentDynamics
        end

        this_subpopulation.TreatmentDynamics = treatment_dynamics

        #-----SW and IDU dynamics		

        if this_subpopulation.GroupId == 0

          idu_mate_total = subpopulations.g(4).d(this_subpopulation.DiseaseStageId).c(c).compact[0]
          total_unif_male_idus = idu_mate_total.Population * PercentMaleIdus
          sw_and_idu_dynamics = step * (total_unif_male_idus * IduSpontaneousQuitRate - this_subpopulation.Population * IduInitiationRate)

        elsif this_subpopulation.GroupId == 1

          idu_mate_total = subpopulations.g(4).d(this_subpopulation.DiseaseStageId).c(c).compact[0]
          sw_mate = subpopulations.g(2).d(this_subpopulation.DiseaseStageId).c(c).compact.first 
          total_unif_female_idus = idu_mate_total.Population * (1-PercentMaleIdus)
          sw_and_idu_dynamics = step * (total_unif_female_idus * IduSpontaneousQuitRate - this_subpopulation.Population * IduInitiationRate + sw_mate.Population * SwQuitRate - this_subpopulation.Population * SwInitiationRate)

        elsif this_subpopulation.GroupId == 2

          genpop_female = subpopulations.g(1).d(this_subpopulation.DiseaseStageId).c(c).compact[0]
          sw_and_idu_dynamics = step * (genpop_female.Population * SwInitiationRate - this_subpopulation.Population * SwQuitRate) 

        elsif this_subpopulation.GroupId == 3

          sw_and_idu_dynamics = 0.0

        elsif this_subpopulation.GroupId == 4

          genpop_male = subpopulations.g(0).d(this_subpopulation.DiseaseStageId).c(c).compact[0]
          genpop_female = subpopulations.g(1).d(this_subpopulation.DiseaseStageId).c(c).compact[0]
          sw_and_idu_dynamics =  step * (IduInitiationRate * genpop_female.Population + IduInitiationRate * genpop_male.Population - this_subpopulation.Population * IduSpontaneousQuitRate)

        end #if groups

        this_subpopulation.SwAndIduDynamics = sw_and_idu_dynamics

        dynamics_array = [general_deaths, disease_progression_exits, disease_progression_entries, treatment_dynamics, sw_and_idu_dynamics]

        this_subpopulation.DynamicsArray = dynamics_array

        this_subpopulation.TotalDynamics = sw_and_idu_dynamics + treatment_dynamics + disease_progression_entries + disease_progression_exits + general_deaths


      end #subpopulations

      #Step 9, build next cycle, populate each subpop

 
 
     subpopulations.c(c).compact.each_with_index do |this_subpopulation,s| 

       new_subpopulation = this_subpopulation.clone
       new_subpopulation.Cycle = this_subpopulation.Cycle + 1
       new_subpopulation.Id = this_subpopulation.Id + 65
       new_subpopulation.Population = this_subpopulation.Population + this_subpopulation.TotalDynamics
       
       if choo.include?(1)
	       new_subpopulation.PropMaleCirc= this_subpopulation.PropMaleCirc.to_f + step *  intervention_subpopulation_coverage.select{|item| item.Id == 1}[0].Coverage.to_f * (1-this_subpopulation.PropMaleCirc.to_f) 
	       else 
		    new_subpopulation.PropMaleCirc= this_subpopulation.PropMaleCirc.to_f 
	    end
	    
       subpopulations = subpopulations.push(new_subpopulation)

    end
  
  #subpopulations = subpopulations + new_subpopulations

end #cycles


CSV.open("file.csv", "wb") do |csv|
  
  csv << [ "Id" , "Cycle", "Population" , "Group", "DiseaseStage",  "Scr", "Probability", "Dgeneral", "DProgExits", "DProgEntries", "D Treatment", "SWIDU dynamics", "CondomUse", "Partnerships","PropMaleCirc" ] 

  subpopulations.each_with_index do |this_subpopulation,s| 
    
    csv << [ this_subpopulation.Id , this_subpopulation.Cycle,this_subpopulation.Population , this_subpopulation.Group, this_subpopulation.DiseaseStage,  this_subpopulation.Scr, this_subpopulation.Probability, this_subpopulation.DynamicsGeneral, this_subpopulation.DiseaseProgressionExits, this_subpopulation.DiseaseProgressionEntries, this_subpopulation.TreatmentDynamics, this_subpopulation.SwAndIduDynamics,  this_subpopulation.CondomUse, this_subpopulation.Partnerships, this_subpopulation.PropMaleCirc]
 
 end
 
end

 
################################################################################################################################################################
############################################### Step 10: Build chart functions to format data into correct arrays (ie, ########################################
################################################################################################################################################################

#class EpiOutputs
  
  #the below specifies the different parameters a EpiSubpopulation object can hold
  #attr_accessor :total_prevalence, :prevalence_by_group, :total_new_infections, :total_new_infections_per_pop, :total_new_infections_by_group, :total_new_infections_by_group_per_pop, :cumulative_total_new_infections, :hiv_deaths, :cumulative_hiv_deaths, :hiv_deaths_by_group
  
#end

epi_outputs = Hash.new
epi_outputs['total_prevalence'] = Array.new
epi_outputs['total_new_infections']= Array.new
epi_outputs['total_new_infections_per_pop'] = Array.new
epi_outputs['cumulative_total_new_infections'] = Array.new
epi_outputs['hiv_deaths'] = Array.new
epi_outputs['cumulative_hiv_deaths'] = Array.new
epi_outputs['prevalence_by_group'] = Array.new
epi_outputs['prevalence_by_group'][0] = Array.new
epi_outputs['prevalence_by_group'][1] = Array.new
epi_outputs['prevalence_by_group'][2] = Array.new
epi_outputs['prevalence_by_group'][3] = Array.new
epi_outputs['prevalence_by_group'][4] = Array.new
epi_outputs['total_plwa'] = Array.new
epi_outputs['plwa_by_group'] = Array.new
epi_outputs['plwa_by_group'][0] = Array.new
epi_outputs['plwa_by_group'][1] = Array.new
epi_outputs['plwa_by_group'][2] = Array.new
epi_outputs['plwa_by_group'][3] = Array.new
epi_outputs['plwa_by_group'][4] = Array.new
epi_outputs['total_new_infections_by_group'] = Array.new
epi_outputs['total_new_infections_by_group'][0] = Array.new
epi_outputs['total_new_infections_by_group'][1] = Array.new
epi_outputs['total_new_infections_by_group'][2] = Array.new
epi_outputs['total_new_infections_by_group'][3] = Array.new
epi_outputs['total_new_infections_by_group'][4] = Array.new
epi_outputs['total_new_infections_by_group_per_pop'] = Array.new
epi_outputs['total_new_infections_by_group_per_pop'][0] = Array.new
epi_outputs['total_new_infections_by_group_per_pop'][1] = Array.new
epi_outputs['total_new_infections_by_group_per_pop'][2] = Array.new
epi_outputs['total_new_infections_by_group_per_pop'][3] = Array.new
epi_outputs['total_new_infections_by_group_per_pop'][4] = Array.new
epi_outputs['hiv_deaths_by_group'] = Array.new
epi_outputs['hiv_deaths_by_group'][0] = Array.new
epi_outputs['hiv_deaths_by_group'][1] = Array.new
epi_outputs['hiv_deaths_by_group'][2] = Array.new
epi_outputs['hiv_deaths_by_group'][3] = Array.new
epi_outputs['hiv_deaths_by_group'][4] = Array.new


cycles.each_with_index do |this_cycle,c|
  
  
  #------- Total prevalence and plwa -------#
  
  total_infected = subpopulations.select{|item| item.Cycle == c && item.DiseaseStageId != 0}.inject(0){|memo, n| memo + n.Population }
  total_pop = subpopulations.select{|item| item.Cycle == c }.inject(0){|memo, n| memo + n.Population }
  epi_outputs.total_plwa[c] = total_infected
  epi_outputs.total_prevalence[c] = total_infected / total_pop 
  
  #------- Total incidence -------#
  
  if c.even? && c !=0
    new_infections = subpopulations.select{|item| item.Cycle == c-1 && item.DiseaseStageId == 0}.inject(0){|memo, n| memo - n.TotalDiseaseProgressionExits } + subpopulations.select{|item| item.Cycle == c-2 && item.DiseaseStageId == 0}.inject(0){|memo, n| memo - n.TotalDiseaseProgressionExits }
    total_old_pop = subpopulations.select{|item| item.Cycle == c-2 }.inject(0){|memo, n| memo + n.Population }
    epi_outputs.total_new_infections.push new_infections
    epi_outputs.total_new_infections_per_pop.push new_infections / total_old_pop * 100_000
    if epi_outputs.cumulative_total_new_infections.empty?
      epi_outputs.cumulative_total_new_infections.push new_infections
    else
      epi_outputs.cumulative_total_new_infections.push epi_outputs.cumulative_total_new_infections.last + new_infections
    end
  end
  
  #------- Total HIV deaths -------#
    
  if c.even? && c !=0
    deaths = subpopulations.select{|item| item.Cycle == c-1}.inject(0){|memo, n| memo + n.Population * n.HivDeathRate * step } + subpopulations.select{|item| item.Cycle == c-2}.inject(0){|memo, n| memo + n.Population * n.HivDeathRate * step } 
    epi_outputs.hiv_deaths.push deaths
    if epi_outputs.cumulative_hiv_deaths.empty?
      epi_outputs.cumulative_hiv_deaths.push deaths
    else
      epi_outputs.cumulative_hiv_deaths.push epi_outputs.cumulative_hiv_deaths.last + deaths
    end
  end

  groups.each_with_index do |group,g|
    
    #------- Prevalence and plwa by group -------#
    
    total_infected = subpopulations.select{|item| item.Cycle == c && item.DiseaseStageId != 0 && item.GroupId == g}.inject(0){|memo, n| memo + n.Population }
    total_pop = subpopulations.select{|item| item.Cycle == c && item.GroupId == g }.inject(0){|memo, n| memo + n.Population }
    epi_outputs.plwa_by_group[g].push total_infected
    epi_outputs.prevalence_by_group[g].push total_infected / total_pop

    
    #------- Incidence by group -------#
    
    if c.even? && c != 0
      new_infections_in_group = subpopulations.select{|item| item.Cycle == c-1 && item.DiseaseStageId == 0 && item.GroupId == g}.inject(0){|memo, n| memo - n.TotalDiseaseProgressionExits } + subpopulations.select{|item| item.Cycle == c-2 && item.DiseaseStageId == 0 && item.GroupId == g}.inject(0){|memo, n| memo - n.TotalDiseaseProgressionExits }
      total_pop_in_group = subpopulations.select{|item| item.Cycle == c  && item.GroupId == g }.inject(0){|memo, n| memo + n.Population }
      epi_outputs.total_new_infections_by_group[g].push new_infections_in_group
      epi_outputs.total_new_infections_by_group_per_pop[g].push new_infections_in_group / total_pop * 100_000
    end
    
    #------- HIV deaths by group -------#
    
    if c.even? && c != 0
      deaths_in_group = subpopulations.select{|item| item.Cycle == c-1 && item.GroupId == g}.inject(0){|memo, n| memo + n.Population * n.HivDeathRate * step } + subpopulations.select{|item| item.Cycle == c-2 && item.GroupId == g}.inject(0){|memo, n| memo + n.Population * n.HivDeathRate * step } 
      epi_outputs.hiv_deaths_by_group[g].push deaths_in_group
    end
    
    
  end #groups
  
  
  
end #each subpopulation



# Step 11: Export chart data via JSON to website


#IRB.start_session(binding)

output = JSON.generate(epi_outputs)
##puts output

#IRB.start_session(binding)