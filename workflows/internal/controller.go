package internal

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/om-baji/kinetic/shared"
)

type WorkflowController struct {
	svc *WorkflowService
}

func NewWorkflowController(svc *WorkflowService) *WorkflowController {
	return &WorkflowController{svc: svc}
}

func (ctrl *WorkflowController) StatusHealth(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"health":    "OK!",
		"timestamp": time.Now(),
	})
}

func (ctrl *WorkflowController) PostWorkflow(ctx fiber.Ctx) error {
	var payload shared.Workflow
	if err := ctx.Bind().Body(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	record, err := ctrl.svc.Create(&payload)
	if err != nil {
		return toFiberError(err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(record)
}

func (ctrl *WorkflowController) FetchWorkflowById(ctx fiber.Ctx) error {
	id := ctx.Params("id", "")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "workflow id is required")
	}

	record, err := ctrl.svc.GetByID(id)
	if err != nil {
		return toFiberError(err)
	}

	return ctx.JSON(record)
}

func (ctrl *WorkflowController) DeleteWorkflow(ctx fiber.Ctx) error {
	id := ctx.Params("id", "")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "workflow id is required")
	}

	if err := ctrl.svc.Delete(id); err != nil {
		return toFiberError(err)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (ctrl *WorkflowController) PauseWorkflow(ctx fiber.Ctx) error {
	id := ctx.Params("id", "")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "workflow id is required")
	}

	record, err := ctrl.svc.Pause(id)
	if err != nil {
		return toFiberError(err)
	}

	return ctx.JSON(record)
}

func (ctrl *WorkflowController) ResumeWorkflow(ctx fiber.Ctx) error {
	id := ctx.Params("id", "")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "workflow id is required")
	}

	record, err := ctrl.svc.Resume(id)
	if err != nil {
		return toFiberError(err)
	}

	return ctx.JSON(record)
}

func (ctrl *WorkflowController) FetchWorkflowGraph(ctx fiber.Ctx) error {
	id := ctx.Params("id", "")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "workflow id is required")
	}

	graph, err := ctrl.svc.GetGraph(id)
	if err != nil {
		return toFiberError(err)
	}

	return ctx.JSON(graph)
}

func (ctrl *WorkflowController) FetchWorkflowLogs(ctx fiber.Ctx) error {
	id := ctx.Params("id", "")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "workflow id is required")
	}

	logs, err := ctrl.svc.GetLogs(id)
	if err != nil {
		return toFiberError(err)
	}

	return ctx.JSON(logs)
}

func toFiberError(err error) error {
	if apiErr, ok := err.(*shared.APIError); ok {
		return fiber.NewError(apiErr.Code, apiErr.Message)
	}
	return fiber.NewError(fiber.StatusInternalServerError, err.Error())
}
