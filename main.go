package main

import (
	"github.com/yfujita/slackutil"
	"os/exec"
	"encoding/json"
	"fmt"
	"strconv"
)

type Tag struct {
	Key		string
	Value	string
}

type Instance struct {
	InstanceId		string
	InstanceType	string
	Tags			[]Tag
}

type Reservation struct {
	Instances		[]Instance
}

type Resp struct {
	Reservations []Reservation
}

func main() {
	reservations := getReservations()

	bot := slackutil.NewBot("{url}", "#bot_test", "ec2-reminder", ":ghost:")

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
			message += "\n"
			instanceNum++
		}
	}

	if instanceNum > 0 {
		message += "```"
	}

	title := "ec2(Running)インスタンス数:" + strconv.Itoa(instanceNum)
	bot.Message(title, message)
}

func getReservations() []Reservation {
	region := "ap-northeast-1a"
	jsonStr := executeCmd("aws", "ec2", "describe-instances", "--filters", "Name=instance-state-code,Values=16", "Name=availability-zone,Values=" + region)
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