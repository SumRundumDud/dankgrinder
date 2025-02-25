// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
)

func (in *Instance) balanceCheck(msg discord.Message) {
	if !strings.Contains(msg.Embeds[0].Title, in.Client.User.Username) {
		return
	}
	if !exp.bal.Match([]byte(msg.Embeds[0].Description)) {
		return
	}
	balstr := strings.Replace(exp.bal.FindStringSubmatch(msg.Embeds[0].Description)[1], ",", "", -1)
	balance, err := strconv.Atoi(balstr)
	if err != nil {
		in.Logger.Errorf("error while reading balance: %v", err)
		return
	}
	in.balance = balance
	in.Logger.Infof(
		"current wallet balance: %v coins",
		numFmt.Sprintf("%d", balance),
	)

	if balance > in.Features.AutoShare.MaximumBalance &&
		in.Features.AutoShare.Enable &&
		in.MasterID != "" &&
		in.Client.User.ID != in.MasterID {
		in.sdlr.Schedule(&scheduler.Command{
			Value: fmt.Sprintf(
				"pls share %v <@%v>",
				balance-in.Features.AutoShare.MinimumBalance,
				in.MasterID,
			),
			Log: "sharing all balance above minimum with master instance",
		})
	}

	if in.startingTime.IsZero() {
		in.initialBalance = balance
		in.startingTime = time.Now()
		return
	}
	inc := balance - in.initialBalance
	per := time.Now().Sub(in.startingTime)
	hourlyInc := int(math.Round(float64(inc) / per.Hours()))
	in.Logger.Infof(
		"average income: %v coins/h",
		numFmt.Sprintf("%d", hourlyInc),
	)
}
