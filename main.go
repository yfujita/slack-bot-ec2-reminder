package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/yfujita/slackutil"
	"os/exec"
	"strconv"
)

type Tag struct {
	Key   string
	Value string
}

type Instance struct {
	InstanceId   string
	InstanceType string
	LaunchTime   string
	Tags         []Tag
}

type Reservation struct {
	Instances []Instance
}

type Resp struct {
	Reservations []Reservation
}

func main() {
	var slackUrl string
	var slackChannel string
	var slackBotName string
	var slackBotIcon string
	var region string
	flag.StringVar(&slackUrl, "slackUrl", "blank", "config file path")
	flag.StringVar(&region, "region", "ap-northeast-1a", "ec2 region")
	flag.StringVar(&slackChannel, "slackChannel", "#bot_test", "slack channel")
	flag.StringVar(&slackBotName, "slackBotName", "ec2-reminder", "slack bot name")
	flag.StringVar(&slackBotIcon, "slackBotIcon", ":ghost:", "slack bot name")

	flag.Parse()
	if slackUrl == "blank" {
		panic("Invalid url parameter")
	}

	reservations := getReservations(region)

	bot := slackutil.NewBot(slackUrl, slackChannel, slackBotName, slackBotIcon)

	reservationNum := len(reservations)

	message := ""
	if reservationNum > 0 {
		message += "```\n"
	}

	instanceNum := 0
	for _, reservation := range reservations {
		instances := reservation.Instances
		for _, instance := range instances {
			message += "id:" + instance.InstanceId
			message += " type:" + instance.InstanceType
			message += " tags:["
			for i, tag := range instance.Tags {
				if i > 0 {
					message += ", "
				}
				message += "{key:" + tag.Key + ", value:" + tag.Value + "}"
			}
			message += "]"
			message += " launchTime:" + instance.LaunchTime
			message += "\n"
			instanceNum++
		}
	}

	if instanceNum > 0 {
		message += "```"
	}

	title := "ec2(Running)インスタンス数:" + strconv.Itoa(instanceNum)
	fmt.Println("Send message. " + title + "\n" + message)
	err := bot.Message(title, message)
	if err != nil {
		panic(err.Error())
	}
}

func getReservations(region string) []Reservation {
	jsonStr := executeCmd("aws", "ec2", "describe-instances", "--output", "json", "--filters", "Name=instance-state-code,Values=16", "Name=availability-zone,Values="+region)
	fmt.Println(jsonStr)
	var resp Resp
	json.Unmarshal([]byte(jsonStr), &resp)
	fmt.Println("instance num=" + strconv.Itoa(len(resp.Reservations)))
	return resp.Reservations
}

func executeCmd(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	if err != nil {
		panic(err.Error())
	}
	return string(out)
}
