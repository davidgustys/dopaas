package do_state

import (
	"context"
	"errors"
	"fmt"
	"github.com/digitalocean/godo"
	"github.com/harshpreet93/dopaas/do_auth"
	"log"
	"os"
)

type dropletInfo struct {
	slug string
	dc string
	tags [] string
}

type projectState struct {
	numDroplets []dropletInfo
}

func getProject(projectId string) (*godo.Project, error) {
	fmt.Println("finding project")
	client := do_auth.Auth()
	ctx := context.Background()
	// create options. initially, these will be blank
	opt := &godo.ListOptions{}

	for {
		projects, resp, err := client.Projects.List(ctx, opt)

		if err != nil {
			return nil, err
		}

		for _, project := range projects {
			if project.ID == projectId {
				return &project, nil
			}
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	return nil, errors.New("cannot find project with ID " + projectId)
}


func extractProjectResourceInfo(project *godo.Project) (*projectState, error) {
	fmt.Println("extracting project resource info")
	client := do_auth.Auth()
	ctx := context.Background()
	opt := &godo.ListOptions{}

	for {
		projectResources, resp, err := client.Projects.ListResources(ctx, project.ID, opt)

		if err != nil {
			return nil, err
		}

		for _, projectResource := range projectResources {
			fmt.Println("project resource ", projectResource)
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	return nil, nil
}

func GetState(projectId string) (*projectState, error) {
	log.Println("getting current state of project ", projectId)

	project, err := getProject(projectId)

	if err != nil {
		log.Println("error getting project ", err)
		os.Exit(1)
	}

	currState, err := extractProjectResourceInfo(project)

	if err != nil {
		log.Println("error getting current state", err)
		os.Exit(1)
	}
	return currState, nil
}
