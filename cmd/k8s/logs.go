package k8s

import (
	"bytes"
	"context"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"shipyard/constants"
	"shipyard/display"
)

func NewLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logs",
		Short:   "Get logs from a service in an environment",
		GroupID: constants.GroupKubernetes,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("kubeconfig", cmd.Flags().Lookup("kubeconfig"))
			viper.BindPFlag("service", cmd.Flags().Lookup("service"))
			viper.BindPFlag("env", cmd.Flags().Lookup("env"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return handleLogsCmd()
		},
	}

	cmd.Flags().String("kubeconfig", "", "Path to kubeconfig")

	cmd.Flags().String("service", "", "Service name")
	cmd.MarkFlagRequired("service")

	cmd.Flags().String("env", "", "environment ID")
	cmd.MarkFlagRequired("env")

	return cmd
}

func handleLogsCmd() error {
	if err := SetKubeconfig(viper.GetString("env")); err != nil {
		return err
	}

	config, namespace, err := getConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	serviceName := viper.GetString("service")
	podName, err := getPodName(config, namespace, serviceName)
	if err != nil {
		return err
	}

	podLogOpts := corev1.PodLogOptions{}
	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, &podLogOpts)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return err
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return err
	}

	writer := display.NewSimpleDisplay()
	writer.Output(buf.String())

	return nil
}