package collect_data

import (
	"fmt"
	"strings"
	"time"

	"github.com/laidbackware/cf-nsx-rule-usage/internal/nsx_client"
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
	UnusedRules map[string]map[string][]Rule
	AllRules		map[string]map[string][]Rule
}

var ruleUsage RuleUsage

func CollectData(nsxApi, nsxUsername, nsxPassword string, skipVerify bool) (RuleUsage, error) {
	client, err := nsx_client.SetupClient(nsxApi, nsxUsername, nsxPassword)
	if err != nil {return ruleUsage, err}

	sections, err := client.GetSgSections()
	if err != nil {return ruleUsage, err}

	_, err = processSections(client, sections)
	if err != nil {return ruleUsage, err}

	return ruleUsage, nil
}

func processSections(client *nsx_client.Client, sections []nsx_client.Section) (RuleUsage, error) {
	// var ruleUsage RuleUsage
	var rule Rule
	ruleUsage.AllRules = make(map[string]map[string][]Rule)
	ruleUsage.UnusedRules = make(map[string]map[string][]Rule)

	for _, section := range(sections) {
		foundationName, err := findTag("ncp/cluster", section.DisplayName, section.Tags)
		if err != nil {return ruleUsage, err}

		asgName, err := findTag("ncp/cf_asg_name", section.DisplayName, section.Tags)
		if err != nil {return ruleUsage, err}

		sectionRules, err := client.GetSectionRules(section.ID)
		if err != nil {return ruleUsage, err}

		sectionStats, err := client.GetSectionStats(section.ID)
		if err != nil {return ruleUsage, err}

		for idx, sectionRule := range(sectionRules) {
			rule = Rule{
				// Assume a single target in all ASG rules
				Target:				sectionRule.Destinations[0].TargetID,
				Ports:				strings.Join(sectionRule.Services[0].Service.Destination_ports[:], ","),
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