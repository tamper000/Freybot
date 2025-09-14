package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var AIRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "ai_requests_total",
	Help: "Total number of messages sent to AI models",
})

var PhotoRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "ai_photo_requests_total",
	Help: "Total number of messages containing photos sent to AI",
})

var VoiceRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "ai_voice_requests_total",
	Help: "Total number of voice messages sent to AI",
})

var ImagesGeneratedTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "ai_images_generated_total",
	Help: "Total number of images generated with AI",
})

var ImagesEditedTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "ai_images_edit_total",
	Help: "Total number of images edited with AI",
})

var ErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "ai_errors_total",
	Help: "Total number of errors occurred in the bot",
}, []string{"type"}) // stt, llm, db, format, genimage, photo, telegram, edit

var ModelUsageTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "ai_model_usage_total",
	Help: "Total usage count per AI model",
}, []string{"model"})
