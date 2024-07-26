package singbox

import (
	"fmt"
	"sifu-clash/models"
)

func SetRulesets(serviceMap map[string][]models.Ruleset) []models.Ruleset {
	var newRulesets []models.Ruleset
	for _, rulesets := range serviceMap {
		newRulesets = append(newRulesets, rulesets...)
	}
	return newRulesets
}
func SetRules(serviceMap map[string][]models.Ruleset) []map[string]interface{}{
	var rules []map[string]interface{}
	for key,rulesets := range serviceMap{
		if key == ""{
			for _,ruleset := range(rulesets){
				if ruleset.China{
					rules = append(rules, map[string]interface{}{"rule_set":ruleset.Tag,"outbound":"direct"})
				}else{
					rules = append(rules, map[string]interface{}{"rule_set":ruleset.Tag,"outbound":fmt.Sprintf("select-%s",ruleset.Tag)})
				}
			}
		}else{
			var rulesetsList []string
			var china bool
			var label string
			for _,ruleset := range(rulesets){
				china = ruleset.China
				label = ruleset.Label
				rulesetsList = append(rulesetsList, ruleset.Tag)
			}
			if china{
				rules = append(rules, map[string]interface{}{"rule_set": rulesetsList,"outbound":"direct"})
			}else{
				rules = append(rules, map[string]interface{}{"rule_set": rulesetsList,"outbound":fmt.Sprintf("select-%s",label)})
			}
		}
		
	}
	return rules
}
