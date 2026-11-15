package models

type ProcessingTask int32

const (
	ProcessingTaskThumbnails ProcessingTask = 1
	ProcessingTaskFrames     ProcessingTask = 2
	ProcessingTaskTranscode  ProcessingTask = 3
)
