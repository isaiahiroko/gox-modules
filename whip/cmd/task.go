package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/origine-run/whip/pkg/task"
	"github.com/origine-run/whip/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// task
// task add/create
// task list
// task get
// task get
// task delete
// task delete all

func init() {
	for _, flag := range taskFlags {
		taskCmd.Flags().String(flag.name, "", flag.desc)
		// taskCmd.MarkFlagRequired(flag.name)
		viper.BindPFlag(flag.name, taskCmd.Flags().Lookup(flag.name))
	}

	RootCmd.AddCommand(taskCmd)
}

var (
	taskFlags = []flag{
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

		{name: "id", desc: "Task id"},
	}
)

var taskCmd = &cobra.Command{
	Use:       "task",
	Aliases:   []string{"r"},
	Short:     "Build container images from source codes",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"get-all", "get", "delete-all", "delete", "new"},
	Run: func(cmd *cobra.Command, args []string) {
		action := args[0]

		switch action {
		case "get-all":
			buf, err := task.GetAll()
			if err != nil {
				fmt.Println(err)
				return
			}

			var tasks []task.TaskPayload
			err = json.Unmarshal(buf, &tasks)
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, v := range tasks {
				fmt.Println(fmt.Sprintf("%s/%s", v.Source, v.Id))
			}
			break
		case "get":
			taskId, _ := cmd.Flags().GetString("id")

			buf, err := task.GetOne(taskId)
			if err != nil {
				fmt.Println(err)
				return
			}

			var v struct {
				Task task.TaskPayload
				Logs []any
			}

			err = json.Unmarshal(buf, &v)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(fmt.Sprintf("%s/%ss", v.Task.Source, v.Task.Id))
			fmt.Println(v.Logs)
			for _, log := range v.Logs {
				fmt.Println(log)
			}
			break
		case "delete-all":
			err := task.DeleteAll()
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("OK")
			break
		case "delete":
			taskId, _ := cmd.Flags().GetString("id")

			err := task.DeleteOne(taskId)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("OK")
			break
		case "new":
			// obtain and set variables
			source, _ := cmd.Flags().GetString("source")

			id := utils.UUID()
			body := task.TaskPayload{
				Id:       id,
				Checksum: utils.Checksum(fmt.Sprintf("%s/%s", source, id)),
			}

			_, err := task.Add(body)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(fmt.Sprintf("%s/%s", body.Source, body.Id))
			break
		}
	},
}
