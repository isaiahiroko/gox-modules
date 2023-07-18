package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/origine-run/whip/pkg/server"
	"github.com/origine-run/whip/pkg/task"
	"github.com/origine-run/whip/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	for _, flag := range serveFlags {
		serveCmd.Flags().String(flag.name, "", flag.desc)
		serveCmd.MarkFlagRequired(flag.name)
		viper.BindPFlag(flag.name, serveCmd.Flags().Lookup(flag.name))
	}

	RootCmd.AddCommand(serveCmd)
}

var (
	serveFlags = []flag{
		{name: "port", desc: "Server listening port"},
	}
	// endpoints = []server.Endpoint{
	// 	{
	// 		Method: "GET",
	// 		Path:   "/",
	// 		Handler: func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// 			httpTask := &tasks.HttpTask{}
	// 			httpTask.Queue()
	// 			httpTask.Args()
	// 			httpTask.Args()
	// 			httpTask.Args()
	// 			queue := queue.NewQueue("http-tasks")
	// 			worker.Add(httpTask)
	// 			rw.Write([]byte("Whip v0.0.0"))
	// 		},
	// 	},
	// }
	endpoints = []server.Endpoint{
		{
			Method: "GET",
			Path:   "/",
			Handler: func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
				rw.Write([]byte("Whip v0.0.0"))
			},
		},
		{
			Method: "POST",
			Path:   "/task",
			Handler: func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
				// obtain and set variables
				var body task.TaskPayload

				err := json.NewDecoder(r.Body).Decode(&body)
				if err != nil {
					http.Error(rw, err.Error(), http.StatusBadRequest)
					return
				}

				body.Id = utils.UUID()
				body.Checksum = utils.Checksum(fmt.Sprintf("%s/%s", body.Source, body.Id))

				buf, err := task.Add(body)
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}

				rw.Write(buf)
			},
		},
		{
			Method: "GET",
			Path:   "/task",
			Handler: func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
				buf, err := task.GetAll()
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}
				rw.Write(buf)
			},
		},
		{
			Method: "GET",
			Path:   "/task/:id",
			Handler: func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
				taskId := params.ByName("id")

				buf, err := task.GetOne(taskId)
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}

				rw.Write(buf)
			},
		},
		{
			Method: "DELETE",
			Path:   "/task",
			Handler: func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
				err := task.DeleteAll()
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}

				rw.Write([]byte("Ok"))
			},
		},
		{
			Method: "DELETE",
			Path:   "/task/:id",
			Handler: func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
				taskId := params.ByName("id")

				err := task.DeleteOne(taskId)
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}

				rw.Write([]byte("OK"))
			},
		},
	}
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"s"},
	Short:   "Expose whip via HTTP",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.SetVersionTemplate("v0.0.0")
		port, _ := cmd.Flags().GetString("port")

		s := server.New()
		s.Start(":"+port, endpoints)
	},
}
