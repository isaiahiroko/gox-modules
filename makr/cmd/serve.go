package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/origine-run/makr/pkg/job"
	"github.com/origine-run/makr/pkg/server"
	"github.com/origine-run/makr/pkg/utils"
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
	endpoints = []server.Endpoint{
		{
			Method: "GET",
			Path:   "/",
			Handler: func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
				rw.Write([]byte("Makr v0.0.0"))
			},
		},
		{
			Method: "POST",
			Path:   "/job",
			Handler: func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
				// obtain and set variables
				var body job.JobModel

				err := json.NewDecoder(r.Body).Decode(&body)
				if err != nil {
					http.Error(rw, err.Error(), http.StatusBadRequest)
					return
				}

				body.GitPassword = utils.DecodeBase64(body.GitPassword)
				body.DockerPassword = utils.DecodeBase64(body.DockerPassword)
				body.Id = utils.UUID()
				body.Checksum = utils.Checksum(fmt.Sprintf("%s/%s/%s:%s", body.GitHost, body.GitUsername, body.GitRepo, body.ImageVersion))

				nJob := job.New()
				buf, err := nJob.Add(body)
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}

				rw.Write(buf)
			},
		},
		{
			Method: "GET",
			Path:   "/job",
			Handler: func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
				nJob := job.New()
				buf, err := nJob.GetAll()
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}
				rw.Write(buf)
			},
		},
		{
			Method: "GET",
			Path:   "/job/:id",
			Handler: func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
				jobId := params.ByName("id")

				nJob := job.New()
				buf, err := nJob.GetOne(jobId)
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}

				rw.Write(buf)
			},
		},
		{
			Method: "DELETE",
			Path:   "/job",
			Handler: func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
				nJob := job.New()
				err := nJob.DeleteAll()
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}

				rw.Write([]byte("Ok"))
			},
		},
		{
			Method: "DELETE",
			Path:   "/job/:id",
			Handler: func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
				jobId := params.ByName("id")

				nJob := job.New()
				err := nJob.DeleteOne(jobId)
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
	Short:   "Expose makr via HTTP",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.SetVersionTemplate("v0.0.0")
		port, _ := cmd.Flags().GetString("port")

		s := server.New()
		s.Start(":"+port, endpoints)
	},
}
