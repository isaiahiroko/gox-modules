package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/origine-run/makr/pkg/job"
	"github.com/origine-run/makr/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// job
// job add/create
// job list
// job get
// job get
// job delete
// job delete all

func init() {
	for _, flag := range jobFlags {
		jobCmd.Flags().String(flag.name, "", flag.desc)
		// jobCmd.MarkFlagRequired(flag.name)
		viper.BindPFlag(flag.name, jobCmd.Flags().Lookup(flag.name))
	}

	RootCmd.AddCommand(jobCmd)
}

var (
	jobFlags = []flag{
		{name: "git-host", desc: "Git repository host"},
		{name: "git-username", desc: "Git repository username"},
		{name: "git-repo", desc: "Git repository name"},
		{name: "git-password", desc: "Git repository password (developer key)"},
		{name: "git-remote", desc: "Git repository remote name (defaults to origin)"},
		{name: "git-branch", desc: "Git repository branch"},
		{name: "docker-host", desc: "Docker registry host"},
		{name: "docker-username", desc: "Docker registry username"},
		{name: "docker-registry", desc: "Docker registry name"},
		{name: "docker-password", desc: "Docker registry password (developer key)"},
		{name: "image-version", desc: "Docker image version"},

		{name: "id", desc: "Job id"},
	}
)

var jobCmd = &cobra.Command{
	Use:       "job",
	Aliases:   []string{"r"},
	Short:     "Build container images from source codes",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"get-all", "get", "delete-all", "delete", "new"},
	Run: func(cmd *cobra.Command, args []string) {
		action := args[0]

		switch action {
		case "get-all":
			nJob := job.New()
			buf, err := nJob.GetAll()
			if err != nil {
				fmt.Println(err)
				return
			}

			var jobs []job.JobModel
			err = json.Unmarshal(buf, &jobs)
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, v := range jobs {
				fmt.Println(fmt.Sprintf("%s - %s/%s/%s:%s", v.Id, v.GitHost, v.GitUsername, v.GitRepo, v.ImageVersion))
			}
			break
		case "get":
			jobId, _ := cmd.Flags().GetString("id")

			nJob := job.New()
			buf, err := nJob.GetOne(jobId)
			if err != nil {
				fmt.Println(err)
				return
			}

			var v struct {
				Job  job.JobModel
				Logs []any
			}

			err = json.Unmarshal(buf, &v)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(fmt.Sprintf("%s - %s/%s/%s:%s", v.Job.Id, v.Job.GitHost, v.Job.GitUsername, v.Job.GitRepo, v.Job.ImageVersion))
			fmt.Println(v.Logs)
			for _, log := range v.Logs {
				fmt.Println(log)
			}
			break
		case "delete-all":
			nJob := job.New()
			err := nJob.DeleteAll()
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("OK")
			break
		case "delete":
			jobId, _ := cmd.Flags().GetString("id")

			nJob := job.New()
			err := nJob.DeleteOne(jobId)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("OK")
			break
		case "new":
			// obtain and set variables
			gHost, _ := cmd.Flags().GetString("git-host")
			gUsername, _ := cmd.Flags().GetString("git-username")
			gRepo, _ := cmd.Flags().GetString("git-repo")
			gPassword, _ := cmd.Flags().GetString("git-password")
			gPassword = utils.DecodeBase64(gPassword)
			gRemote, _ := cmd.Flags().GetString("git-remote")
			gBranch, _ := cmd.Flags().GetString("git-branch")

			dHost, _ := cmd.Flags().GetString("docker-host")
			dUsername, _ := cmd.Flags().GetString("docker-username")
			dRegistry, _ := cmd.Flags().GetString("docker-registry")
			dPassword, _ := cmd.Flags().GetString("docker-password")
			dPassword = utils.DecodeBase64(dPassword)

			iVersion, _ := cmd.Flags().GetString("image-version")

			id := utils.UUID()
			body := job.JobModel{
				Id:             id,
				GitHost:        gHost,
				GitUsername:    gUsername,
				GitRepo:        gRepo,
				GitPassword:    gPassword,
				GitRemote:      gRemote,
				GitBranch:      gBranch,
				DockerHost:     dHost,
				DockerUsername: dUsername,
				DockerRegistry: dRegistry,
				DockerPassword: dPassword,
				ImageVersion:   iVersion,
				Checksum:       utils.Checksum(fmt.Sprintf("%s/%s/%s:%s", gHost, gUsername, gRepo, iVersion)),
			}

			nJob := job.New()
			_, err := nJob.Add(body)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(fmt.Sprintf("%s - %s/%s/%s:%s", body.Id, body.GitHost, body.GitUsername, body.GitRepo, body.ImageVersion))
			break
		}
	},
}
