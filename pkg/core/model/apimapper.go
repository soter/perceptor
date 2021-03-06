/*
Copyright (C) 2018 Synopsys, Inc.

Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements. See the NOTICE file
distributed with this work for additional information
regarding copyright ownership. The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied. See the License for the
specific language governing permissions and limitations
under the License.
*/

package model

import (
	"github.com/blackducksoftware/perceptor/pkg/api"
	log "github.com/sirupsen/logrus"
)

// api -> model

func ApiImageToCoreImage(apiImage api.Image) *Image {
	return NewImage(apiImage.Name, DockerImageSha(apiImage.Sha))
}

func ApiContainerToCoreContainer(apiContainer api.Container) *Container {
	return NewContainer(*ApiImageToCoreImage(apiContainer.Image), apiContainer.Name)
}

func ApiPodToCorePod(apiPod api.Pod) *Pod {
	containers := []Container{}
	for _, apiContainer := range apiPod.Containers {
		containers = append(containers, *ApiContainerToCoreContainer(apiContainer))
	}
	return NewPod(apiPod.Name, apiPod.UID, apiPod.Namespace, containers)
}

// model -> api

func (model *Model) ScanResults() api.ScanResults {
	// pods
	pods := []api.ScannedPod{}
	for podName, pod := range model.Pods {
		podScan, err := model.ScanResultsForPod(podName)
		if err != nil {
			log.Errorf("unable to retrieve scan results for Pod %s: %s", podName, err.Error())
			continue
		}
		if podScan == nil {
			log.Debugf("image scans not complete for pod %s, skipping (pod info: %+v)", podName, pod)
			continue
		}
		pods = append(pods, api.ScannedPod{
			Namespace:        pod.Namespace,
			Name:             pod.Name,
			PolicyViolations: podScan.PolicyViolations,
			Vulnerabilities:  podScan.Vulnerabilities,
			OverallStatus:    podScan.OverallStatus})
	}

	// images
	images := []api.ScannedImage{}
	for sha, imageInfo := range model.Images {
		if imageInfo.ScanStatus != ScanStatusComplete {
			continue
		}
		if imageInfo.ScanResults == nil {
			log.Errorf("model inconsistency: found ScanStatusComplete for image %s, but nil ScanResults (imageInfo %+v)", sha, imageInfo)
			continue
		}
		image := imageInfo.Image()
		apiImage := api.ScannedImage{
			Name:             image.HumanReadableName(),
			Sha:              string(image.Sha),
			PolicyViolations: imageInfo.ScanResults.PolicyViolationCount(),
			Vulnerabilities:  imageInfo.ScanResults.VulnerabilityCount(),
			OverallStatus:    imageInfo.ScanResults.OverallStatus().String(),
			ComponentsURL:    imageInfo.ScanResults.ComponentsHref}
		images = append(images, apiImage)
	}

	return *api.NewScanResults(model.HubVersion, model.HubVersion, pods, images)
}
