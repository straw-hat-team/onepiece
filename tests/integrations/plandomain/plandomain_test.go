package plandomain_test

import (
	"github.com/straw-hat-team/onepiece-go/onepiece/onepiecetesting"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain"
	"github.com/straw-hat-team/onepiece-go/tests/integrations/plandomain/planproto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestCreatePLan(t *testing.T) {
	t.Run("creates a plan", func(t *testing.T) {
		onepiecetesting.NewTestCase(t, plandomain.Decider).
			When(&planproto.Command{Command: &planproto.Command_CreatePlan{CreatePlan: &planproto.CreatePlan{
				PlanId: "d83a3744-0e53-4fb7-88f7-7ffc831f0090",
				Title:  "Vacation",
				Color:  "#FF0000",
				GoalAmount: &planproto.Amount{
					Amount:       1000,
					Denomination: "USD",
				},
				Description:      "Plan for a vacation",
				Icon:             "https://some-url.com/icon.png",
				CreatedAt:        timestamppb.New(time.Date(1993, 7, 22, 7, 30, 0, 0, time.UTC)),
				DepositAccountId: "583448c0-696f-4ce5-a4c0-785a3b5c1603",
			}}}).
			Then(&planproto.Event{Event: &planproto.Event_PlanCreated{PlanCreated: &planproto.PlanCreated{
				PlanId: "d83a3744-0e53-4fb7-88f7-7ffc831f0090",
				Title:  "Vacation",
				Color:  "#FF0000",
				GoalAmount: &planproto.Amount{
					Amount:       1000,
					Denomination: "USD",
				},
				Description:      "Plan for a vacation",
				Icon:             "https://some-url.com/icon.png",
				CreatedAt:        timestamppb.New(time.Date(1993, 7, 22, 7, 30, 0, 0, time.UTC)),
				DepositAccountId: "583448c0-696f-4ce5-a4c0-785a3b5c1603",
			}}}).
			Assert()
	})

	t.Run("fails to create a plan if the plan already exists", func(t *testing.T) {
		onepiecetesting.NewTestCase(t, plandomain.Decider).
			Given(&planproto.Event{Event: &planproto.Event_PlanCreated{PlanCreated: &planproto.PlanCreated{
				PlanId: "d83a3744-0e53-4fb7-88f7-7ffc831f0090",
				Title:  "Vacation",
				Color:  "#FF0000",
				GoalAmount: &planproto.Amount{
					Amount:       1000,
					Denomination: "USD",
				},
				Description:      "Plan for a vacation",
				Icon:             "https://some-url.com/icon.png",
				CreatedAt:        timestamppb.New(time.Date(1993, 7, 22, 7, 30, 0, 0, time.UTC)),
				DepositAccountId: "583448c0-696f-4ce5-a4c0-785a3b5c1603",
			}}}).
			When(&planproto.Command{Command: &planproto.Command_CreatePlan{CreatePlan: &planproto.CreatePlan{
				PlanId: "d83a3744-0e53-4fb7-88f7-7ffc831f0090",
				Title:  "Vacation",
				Color:  "#FF0000",
				GoalAmount: &planproto.Amount{
					Amount:       1000,
					Denomination: "USD",
				},
				Description:      "Plan for a vacation",
				Icon:             "https://some-url.com/icon.png",
				CreatedAt:        timestamppb.New(time.Date(1993, 7, 22, 7, 30, 0, 0, time.UTC)),
				DepositAccountId: "583448c0-696f-4ce5-a4c0-785a3b5c1603",
			}}}).
			Catch(plandomain.ErrPlanExists).
			Assert()
	})
}

func TestArchivePlan(t *testing.T) {
	t.Run("fails to archive a plan if the plan does not exist", func(t *testing.T) {
		onepiecetesting.NewTestCase(t, plandomain.Decider).
			When(&planproto.Command{Command: &planproto.Command_ArchivePlan{ArchivePlan: &planproto.ArchivePlan{
				PlanId:     "d83a3744-0e53-4fb7-88f7-7ffc831f0090",
				ArchivedBy: "583448c0-696f-4ce5-a4c0-785a3b5c1603",
				ArchivedAt: timestamppb.New(time.Date(1993, 7, 22, 7, 30, 0, 0, time.UTC)),
			}}}).
			Catch(plandomain.ErrPlanNotFound).
			Assert()
	})

	t.Run("fails to archive a plan if the plan is already archived", func(t *testing.T) {
		onepiecetesting.NewTestCase(t, plandomain.Decider).
			Given(
				&planproto.Event{Event: &planproto.Event_PlanCreated{PlanCreated: &planproto.PlanCreated{
					PlanId: "d83a3744-0e53-4fb7-88f7-7ffc831f0090",
					Title:  "Vacation",
					Color:  "#FF0000",
					GoalAmount: &planproto.Amount{
						Amount:       1000,
						Denomination: "USD",
					},
					Description:      "Plan for a vacation",
					Icon:             "https://some-url.com/icon.png",
					CreatedAt:        timestamppb.New(time.Date(1993, 7, 22, 7, 30, 0, 0, time.UTC)),
					DepositAccountId: "583448c0-696f-4ce5-a4c0-785a3b5c1603",
				}}},
				&planproto.Event{Event: &planproto.Event_PlanArchived{PlanArchived: &planproto.PlanArchived{
					PlanId:     "d83a3744-0e53-4fb7-88f7-7ffc831f0090",
					ArchivedBy: "583448c0-696f-4ce5-a4c0-785a3b5c1603",
					ArchivedAt: timestamppb.New(time.Date(1993, 7, 22, 7, 30, 0, 0, time.UTC)),
				}}},
			).
			When(&planproto.Command{Command: &planproto.Command_ArchivePlan{ArchivePlan: &planproto.ArchivePlan{
				PlanId:     "d83a3744-0e53-4fb7-88f7-7ffc831f0090",
				ArchivedBy: "583448c0-696f-4ce5-a4c0-785a3b5c1603",
				ArchivedAt: timestamppb.New(time.Date(1993, 7, 22, 7, 30, 0, 0, time.UTC)),
			}}}).Catch(plandomain.ErrPlanArchived).Assert()
	})

	t.Run("archives a plan", func(t *testing.T) {
		onepiecetesting.NewTestCase(t, plandomain.Decider).
			Given(&planproto.Event{Event: &planproto.Event_PlanCreated{PlanCreated: &planproto.PlanCreated{
				PlanId: "d83a3744-0e53-4fb7-88f7-7ffc831f0090",
				Title:  "Vacation",
				Color:  "#FF0000",
				GoalAmount: &planproto.Amount{
					Amount:       1000,
					Denomination: "USD",
				},
				Description:      "Plan for a vacation",
				Icon:             "https://some-url.com/icon.png",
				CreatedAt:        timestamppb.New(time.Date(1993, 7, 22, 7, 30, 0, 0, time.UTC)),
				DepositAccountId: "583448c0-696f-4ce5-a4c0-785a3b5c1603",
			}}}).
			When(&planproto.Command{Command: &planproto.Command_ArchivePlan{ArchivePlan: &planproto.ArchivePlan{
				PlanId:     "d83a3744-0e53-4fb7-88f7-7ffc831f0090",
				ArchivedBy: "583448c0-696f-4ce5-a4c0-785a3b5c1603",
				ArchivedAt: timestamppb.New(time.Date(1993, 7, 22, 7, 30, 0, 0, time.UTC)),
			}}}).
			Then(&planproto.Event{Event: &planproto.Event_PlanArchived{PlanArchived: &planproto.PlanArchived{
				PlanId:     "d83a3744-0e53-4fb7-88f7-7ffc831f0090",
				ArchivedBy: "583448c0-696f-4ce5-a4c0-785a3b5c1603",
				ArchivedAt: timestamppb.New(time.Date(1993, 7, 22, 7, 30, 0, 0, time.UTC)),
			}}}).Assert()
	})
}

func TestUpdatePlan(t *testing.T) {
	t.Run("fails to update a plan if the plan does not exist", func(t *testing.T) {
		onepiecetesting.NewTestCase(t, plandomain.Decider).
			When(&planproto.Command{Command: &planproto.Command_UpdatePlan{UpdatePlan: &planproto.UpdatePlan{
				PlanId: "d83a3744-0e53-4fb7-88f7-7ffc831f0090",
				Title:  "Vacation",
				Color:  "#FF0000",
				GoalAmount: &planproto.Amount{
					Amount:       1000,
					Denomination: "USD",
				},
				Description: "Plan for a vacation",
				Icon:        "https://some-url.com/icon.png",
				UpdatedAt:   timestamppb.New(time.Date(1993, 7, 22, 7, 30, 0, 0, time.UTC)),
			}}}).
			Catch(plandomain.ErrPlanNotFound).
			Assert()
	})

	t.Run("updates a plan", func(t *testing.T) {
		onepiecetesting.NewTestCase(t, plandomain.Decider).
			Given(&planproto.Event{Event: &planproto.Event_PlanCreated{PlanCreated: &planproto.PlanCreated{
				PlanId: "d83a3744-0e53-4fb7-88f7-7ffc831f0090",
				Title:  "Vacation",
				Color:  "#FF0000",
				GoalAmount: &planproto.Amount{
					Amount:       1000,
					Denomination: "USD",
				},
				Description:      "Plan for a vacation",
				Icon:             "https://some-url.com/icon.png",
				CreatedAt:        timestamppb.New(time.Date(1993, 7, 22, 7, 30, 0, 0, time.UTC)),
				DepositAccountId: "583448c0-696f-4ce5-a4c0-785a3b5c1603",
			}}}).
			When(&planproto.Command{Command: &planproto.Command_UpdatePlan{UpdatePlan: &planproto.UpdatePlan{
				PlanId: "d83a3744-0e53-4fb7-88f7-7ffc831f0090",
				Title:  "Vacation",
				Color:  "#4A0336",
				GoalAmount: &planproto.Amount{
					Amount:       1000,
					Denomination: "USD",
				},
				Description: "Plan for a vacation",
				Icon:        "https://some-url.com/icon.png",
				UpdatedAt:   timestamppb.New(time.Date(1993, 7, 22, 7, 30, 0, 0, time.UTC)),
			}}}).
			Then(&planproto.Event{Event: &planproto.Event_PlanUpdated{PlanUpdated: &planproto.PlanUpdated{
				PlanId: "d83a3744-0e53-4fb7-88f7-7ffc831f0090",
				Title:  "Vacation",
				Color:  "#4A0336",
				GoalAmount: &planproto.Amount{
					Amount:       1000,
					Denomination: "USD",
				},
				Description: "Plan for a vacation",
				Icon:        "https://some-url.com/icon.png",
				UpdatedAt:   timestamppb.New(time.Date(1993, 7, 22, 7, 30, 0, 0, time.UTC)),
			}}}).Assert()
	})
}
