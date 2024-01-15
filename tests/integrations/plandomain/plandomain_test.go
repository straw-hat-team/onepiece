package plandomain_test

import (
	"github.com/straw-hat-team/onepiece-go/onepiece"
	"github.com/straw-hat-team/onepiece-go/onepiece/onepiecetesting"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestCreatePlanHandler(t *testing.T) {
	decider := onepiece.NewDecider(
		plandomain.Decide,
		plandomain.Evolve,
	)

	t.Run("creates a plan", func(t *testing.T) {
		onepiecetesting.NewTestCase(t, decider).
			When(&plandomain.Command{Command: &plandomain.Command_CreatePlan{CreatePlan: &plandomain.CreatePlan{
				PlanId: "d83a3744-0e53-4fb7-88f7-7ffc831f0090",
				Title:  "Vacation",
				Color:  "#FF0000",
				GoalAmount: &plandomain.Amount{
					Amount:       1000,
					Denomination: "USD",
				},
				Description:      "Plan for a vacation",
				Icon:             "https://some-url.com/icon.png",
				CreatedAt:        timestamppb.New(time.Date(1993, 7, 22, 7, 30, 0, 0, time.UTC)),
				DepositAccountId: "583448c0-696f-4ce5-a4c0-785a3b5c1603",
			}}}).
			Then(&plandomain.Event{Event: &plandomain.Event_PlanCreated{PlanCreated: &plandomain.PlanCreated{
				PlanId: "d83a3744-0e53-4fb7-88f7-7ffc831f0090",
				Title:  "Vacation",
				Color:  "#FF0000",
				GoalAmount: &plandomain.Amount{
					Amount:       1000,
					Denomination: "USD",
				},
				Description:      "Plan for a vacation",
				Icon:             "https://some-url.com/icon.png",
				CreatedAt:        timestamppb.New(time.Date(1993, 7, 22, 7, 30, 0, 0, time.UTC)),
				DepositAccountId: "583448c0-696f-4ce5-a4c0-785a3b5c1603",
			}}}).
			Run()
	})

	t.Run("fails to create a plan if the plan already exists", func(t *testing.T) {
		onepiecetesting.NewTestCase(t, decider).
			Given(&plandomain.Event{Event: &plandomain.Event_PlanCreated{PlanCreated: &plandomain.PlanCreated{
				PlanId: "d83a3744-0e53-4fb7-88f7-7ffc831f0090",
				Title:  "Vacation",
				Color:  "#FF0000",
				GoalAmount: &plandomain.Amount{
					Amount:       1000,
					Denomination: "USD",
				},
				Description:      "Plan for a vacation",
				Icon:             "https://some-url.com/icon.png",
				CreatedAt:        timestamppb.New(time.Date(1993, 7, 22, 7, 30, 0, 0, time.UTC)),
				DepositAccountId: "583448c0-696f-4ce5-a4c0-785a3b5c1603",
			}}}).
			When(&plandomain.Command{Command: &plandomain.Command_CreatePlan{CreatePlan: &plandomain.CreatePlan{
				PlanId: "d83a3744-0e53-4fb7-88f7-7ffc831f0090",
				Title:  "Vacation",
				Color:  "#FF0000",
				GoalAmount: &plandomain.Amount{
					Amount:       1000,
					Denomination: "USD",
				},
				Description:      "Plan for a vacation",
				Icon:             "https://some-url.com/icon.png",
				CreatedAt:        timestamppb.New(time.Date(1993, 7, 22, 7, 30, 0, 0, time.UTC)),
				DepositAccountId: "583448c0-696f-4ce5-a4c0-785a3b5c1603",
			}}}).
			Catch(plandomain.ErrAlreadyExists).
			Run()
	})
}
