package server

import "net/http"

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	handler := newHandlers()

	mux.HandleFunc("GET /healthz", handler.healthz)
	mux.HandleFunc("POST /chat", handler.chat)
	mux.HandleFunc("POST /analyze-pipeline", handler.analyzePipeline)
	mux.HandleFunc("POST /deployment-summary", handler.deploymentSummary)
	mux.HandleFunc("POST /root-cause", handler.rootCause)
	mux.HandleFunc("POST /optimize", handler.optimize)
	mux.HandleFunc("POST /security-review", handler.securityReview)
	mux.HandleFunc("POST /explain-workflow", handler.explainWorkflow)
	mux.HandleFunc("POST /generate-pipeline", handler.generatePipeline)
	mux.HandleFunc("POST /generate-terraform", handler.generateTerraform)
	mux.HandleFunc("POST /generate-kubernetes", handler.generateKubernetes)
	mux.HandleFunc("POST /generate-video-plan", handler.generateVideoPlan)
	mux.HandleFunc("POST /upload-video", handler.uploadVideo)

	return mux
}
