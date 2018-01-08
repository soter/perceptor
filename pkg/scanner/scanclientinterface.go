package scanner

import (
	common "bitbucket.org/bdsengineering/perceptor/pkg/common"
)

type ScanClientInterface interface {
	Scan(job ScanJob) error
	FetchProject(projectName string) (*Project, error)
}

type ScanJob struct {
	ProjectName string
	Image       common.Image
}

func NewScanJob(projectName string, image common.Image) *ScanJob {
	return &ScanJob{
		ProjectName: projectName,
		Image:       image,
	}
}