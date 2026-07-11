package internal

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(app *fiber.App, ctrl *WorkflowController) {
	app.Get("/health", ctrl.StatusHealth)
	app.Post("/workflow", ctrl.PostWorkflow)
	app.Get("/workflow/:id", ctrl.FetchWorkflowById)
	app.Delete("/workflow/:id", ctrl.DeleteWorkflow)
	app.Post("/workflow/:id/pause", ctrl.PauseWorkflow)
	app.Post("/workflow/:id/resume", ctrl.ResumeWorkflow)
	app.Get("/workflow/:id/graph", ctrl.FetchWorkflowGraph)
	app.Get("/workflow/:id/logs", ctrl.FetchWorkflowLogs)
}
