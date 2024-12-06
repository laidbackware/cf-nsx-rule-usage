package collect_data

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/laidbackware/cf-nsx-rule-usage/internal/nsx_client"
	"github.com/schollz/progressbar/v3"
)

type Rule struct {
	Target 				string	`json:"target"`
	Ports 				string	`json:"prots"`
	Protocol 			string	`json:"protocol"`
	Created 			string	`json:"created"`
	LastUpdated 	string	`json:"last_updated"`
	HitCount 			int			`json:"hit_count"`
	PacketCount 	int			`json:"packet_count"`

}

type RuleUsage struct {
	// map[foundation_name]map[asg_name][]Rule
	UnusedRules				 	map[string]map[string][]Rule
	UnusedRulesMonths 	map[string]map[string][]Rule
	AllRules						map[string]map[string][]Rule

}

var ruleUsage RuleUsage

func CollectData(nsxApi, nsxUsername, nsxPassword string, unusedMonths int, skipVerify, debug bool, log Logger) (RuleUsage, error) {
	client, err := nsx_client.SetupClient(nsxApi, nsxUsername, nsxPassword, skipVerify, debug, log)
	if err != nil {return ruleUsage, err}

	sections, err := client.GetSgSections(debug, log)
	if err != nil {return ruleUsage, err}

	_, err = processSections(client, sections, unusedMonths, debug, log)
	if err != nil {return ruleUsage, err}

	return ruleUsage, nil
}

func processSections(client *nsx_client.Client, sections []nsx_client.Section, unusedMonths int, debug bool,log Logger) (RuleUsage, error) {
	var rule Rule
	var ports string
	var bar *progressbar.ProgressBar
	ruleUsage.AllRules = make(map[string]map[string][]Rule)
	ruleUsage.UnusedRules = make(map[string]map[string][]Rule)
	ruleUsage.UnusedRulesMonths = make(map[string]map[string][]Rule)

	log.Printf("Processing sections...")
	if !debug{
		bar = progressbar.Default(int64(len(sections)))
	}
	for _, section := range(sections) {
		startTime := time.Now()
		foundationName, err := findTag("ncp/cluster", section.DisplayName, section.Tags)
		if err != nil {return ruleUsage, err}

		asgName, err := findTag("ncp/cf_asg_name", section.DisplayName, section.Tags)
		if err != nil {return ruleUsage, err}

		sectionRules, err := client.GetSectionRules(section.ID, debug, log)
		if err != nil {return ruleUsage, err}

		sectionStats, err := client.GetSectionStats(section.ID, debug, log)
		if err != nil {return ruleUsage, err}

		for idx, sectionRule := range(sectionRules) {
			if debug {log.Printf("Building struct: " + section.DisplayName + ":" + sectionRule.DisplayName + ":" + strconv.Itoa(idx))}
			oldSection := false
			if time.UnixMilli(section.LastModifiedTime).Before(time.Now().AddDate(0, -unusedMonths,0)) {
				oldSection = true
			}
			
			if sectionRule.DisplayName[0:8] == "all_all" {
				ports = "all"
			} else {
				ports = strings.Join(sectionRule.Services[0].Service.Destination_ports[:], ",")
			}

			// If ASG destination is 0.0.0.0/0, no destination is set in NSX, meanning any
			var destination string
			if sectionRule.Destinations != nil {
				destination = sectionRule.Destinations[0].TargetID
			} else {
				destination = "ANY"
			}

			rule = Rule{
				// Assume a single target in all ASG rules
				Target:				destination,
				Ports:				ports,
				Protocol:			sectionRule.Services[0].Service.L4Protocol,
				Created:			time.UnixMilli(section.CreateTime).Format(time.DateTime),
				LastUpdated:	time.UnixMilli(section.LastModifiedTime).Format(time.DateTime),
				HitCount:			sectionStats.Results[idx].HitCount,
				PacketCount:	sectionStats.Results[idx].PacketCount,
			}
			addRule(foundationName, asgName, ruleUsage.AllRules, rule)
			if rule.HitCount == 0 {
				addRule(foundationName, asgName, ruleUsage.UnusedRules, rule)
			}
			if oldSection {
				addRule(foundationName, asgName, ruleUsage.UnusedRulesMonths, rule)
			}
		}

		elapsedMillis := time.Since(startTime).Milliseconds()
		// Ensure that no more than 100 requests per second can be made to prevent NSX API rate limiting
		if elapsedMillis < 20 {
			time.Sleep(time.Duration(21 - elapsedMillis) * time.Millisecond)
		}
		if !debug{
			bar.Add(1)
		}
	}

	return ruleUsage, nil
}

func findTag(scope, section string, tags []nsx_client.Tags) (string, error) {
	for _, tag := range(tags) {
		if tag.Scope == scope {
			return tag.Tag, nil
		}
	}
	return "", fmt.Errorf("tag scope [%s] not found for section %s", scope, section)
}

func addRule(foundaitonName, asgName string, targetMap map[string]map[string][]Rule, rule Rule) {
	if _, ok := targetMap[foundaitonName]; !ok {
		targetMap[foundaitonName] = make(map[string][]Rule)
	}
	if _, ok := targetMap[foundaitonName][asgName]; !ok {
		targetMap[foundaitonName][asgName] = []Rule{rule}
	} else {
		targetMap[foundaitonName][asgName] = append(targetMap[foundaitonName][asgName], rule)
	}
}