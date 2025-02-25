// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"fmt"
	"time"

	"github.com/dankgrinder/dankgrinder/instance/scheduler"
)

// commands returns a command pointer slice with all commands that should be
// executed periodically. It contains all commands as configured.
func (in *Instance) newCmds() []*scheduler.Command {
	var cmds []*scheduler.Command
	if in.Features.Commands.Beg {
		cmds = append(cmds, &scheduler.Command{
			Value:    "pls beg",
			Interval: time.Duration(in.Compat.Cooldown.Beg) * time.Second,
		})
	}
	if in.Features.Commands.Postmeme {
		cmds = append(cmds, &scheduler.Command{
			Value:       "pls pm",
			Interval:    time.Duration(in.Compat.Cooldown.Postmeme) * time.Second,
			AwaitResume: true,
		})
	}
	if in.Features.Commands.Search {
		cmds = append(cmds, &scheduler.Command{
			Value:       "pls search",
			Interval:    time.Duration(in.Compat.Cooldown.Search) * time.Second,
			AwaitResume: true,
		})
	}
	if in.Features.Commands.Highlow {
		cmds = append(cmds, &scheduler.Command{
			Value:       "pls hl",
			Interval:    time.Duration(in.Compat.Cooldown.Highlow) * time.Second,
			AwaitResume: true,
		})
	}
	if in.Features.Commands.Fish {
		cmds = append(cmds, &scheduler.Command{
			Value:       "pls fish",
			Interval:    time.Duration(in.Compat.Cooldown.Fish) * time.Second,
			AwaitResume: true,
		})
	}
	if in.Features.Commands.Hunt {
		cmds = append(cmds, &scheduler.Command{
			Value:       "pls hunt",
			Interval:    time.Duration(in.Compat.Cooldown.Hunt) * time.Second,
			AwaitResume: true,
		})
	}
	if in.Features.BalanceCheck.Enable {
		cmds = append(cmds, &scheduler.Command{
			Value:    "pls bal",
			Interval: time.Duration(in.Features.BalanceCheck.Interval) * time.Second,
		})
	}
	if in.Features.AutoTidepod.Enable {
		cmds = append(cmds, &scheduler.Command{
			Value:       "pls use tidepod",
			Interval:    time.Duration(in.Features.AutoTidepod.Interval) * time.Second,
			AwaitResume: true,
		})
	}
	if in.Features.AutoBlackjack.Enable {
		cmds = append(cmds, in.newAutoBlackjackCmd())
	}

	for _, cmd := range in.Features.CustomCommands {
		// cmd.Value and cmd.Amount are not checked for correct values here
		// because they were checked when the application started using
		// cfg.Validate().
		cmds = append(cmds, &scheduler.Command{
			Value:    cmd.Value,
			Interval: time.Duration(cmd.Interval) * time.Second,
			Amount:   uint(cmd.Amount),
			CondFunc: func() bool {
				return cmd.PauseBelowBalance == 0 || in.balance >= cmd.PauseBelowBalance
			},
		})
	}
	return cmds
}

func (in *Instance) newAutoSellChain() *scheduler.Command {
	var cmds []*scheduler.Command
	for _, item := range in.Features.AutoSell.Items {
		cmds = append(cmds, &scheduler.Command{
			Value:    fmt.Sprintf("pls sell %v max", item),
			Interval: time.Duration(in.Compat.Cooldown.Sell) * time.Second,
		})
	}
	return in.newCmdChain(
		cmds,
		time.Duration(in.Features.AutoSell.Interval)*time.Second,
	)
}

func (in *Instance) newAutoGiftChain() *scheduler.Command {
	var cmds []*scheduler.Command
	for _, item := range in.Features.AutoGift.Items {
		cmds = append(cmds, &scheduler.Command{
			Value:       fmt.Sprintf("pls shop %v", item),
			Interval:    time.Duration(in.Compat.Cooldown.Gift) * time.Second,
			AwaitResume: true,
		})
	}
	return in.newCmdChain(
		cmds,
		time.Duration(in.Features.AutoGift.Interval)*time.Second,
	)
}

func (in *Instance) newAutoBlackjackCmd() *scheduler.Command {
	cmd := &scheduler.Command{
		Value:    fmt.Sprintf("pls blackjack %v", in.Features.AutoBlackjack.Amount),
		Interval: time.Duration(in.Compat.Cooldown.Blackjack) * time.Second,
		CondFunc: func() bool {
			correctBalance := in.Features.AutoBlackjack.PauseBelowBalance == 0 || in.balance >= in.Features.AutoBlackjack.PauseBelowBalance
			return correctBalance && in.balance < 10000000
		},
		AwaitResume:          true,
		RescheduleAsPriority: in.Features.AutoBlackjack.Priority,
	}
	if in.Features.AutoBlackjack.Amount == 0 {
		cmd.Value = fmt.Sprintf("pls bet max")
	}
	return cmd
}

func (in *Instance) newCmdChain(cmds []*scheduler.Command, chainInterval time.Duration) *scheduler.Command {
	for i := 0; i < len(cmds); i++ {
		if i != 0 {
			cmds[i-1].Next = cmds[i]
		}
		if i == len(cmds)-1 {
			cmds[i].Next = cmds[0]
			cmds[i].Interval = chainInterval
		}
	}
	return cmds[0]
}
