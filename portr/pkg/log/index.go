package log

import (
	"context"

	"github.com/origine-run/portr/pkg/cluster"
	coreV1 "k8s.io/api/core/v1"
)

type Log struct {
	cluster *cluster.Cluster
}

func (l *Log) Watch(namespace string, pod string, container string) {
	client := l.cluster.Client

	stream, err := client.CoreV1().Pods(namespace).GetLogs(
		pod,
		&coreV1.PodLogOptions{
			Container: container,
			Follow:    true,
			Previous:  true,
		},
	).Stream(context.TODO())
	if err != nil {
		return
	}
	defer stream.Close()

	for {
		// buf := make([]byte)
		// numBytes, err := stream.Read(buf)
		// if numBytes == 0 {
		// 	continue
		// }
		// if err == io.EOF {
		// 	time.Sleep(30 * time.Second)
		// 	continue
		// }
		// if err != nil {
		// 	return err
		// }
		// message := string(buf)
		// fmt.Print(message)
	}
}
