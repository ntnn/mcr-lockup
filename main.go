package mcrlockup

import (
	"context"
	"fmt"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	mctrl "sigs.k8s.io/multicluster-runtime"
	mcbuilder "sigs.k8s.io/multicluster-runtime/pkg/builder"
	mcreconcile "sigs.k8s.io/multicluster-runtime/pkg/reconcile"
	"sigs.k8s.io/multicluster-runtime/providers/file"
)

func Run(ctx context.Context, kubeconfigs []string) error {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if len(kubeconfigs) == 0 {
		return fmt.Errorf("at least one kubeconfig file must be specified")
	}

	provider, err := file.New(file.Options{
		KubeconfigFiles: kubeconfigs,
	})
	if err != nil {
		return err
	}

	kubeconfigBytes, err := os.ReadFile(kubeconfigs[0])
	if err != nil {
		return err
	}

	restConfig, err := clientcmd.RESTConfigFromKubeConfig(kubeconfigBytes)
	if err != nil {
		return err
	}

	mcmanager, err := mctrl.NewManager(restConfig, provider, mctrl.Options{})
	if err != nil {
		return err
	}

	recv := make(chan string, 10)

	if err := mcbuilder.ControllerManagedBy(mcmanager).
		For(&corev1.Namespace{}).
		Complete(mcreconcile.Func(func(ctx context.Context, req mctrl.Request) (mctrl.Result, error) {
			ctrl.Log.Info("reconciling namespace", "namespace", req.Namespace, "name", req.Name, "cluster", req.ClusterName)
			recv <- req.ClusterName
			return mctrl.Result{}, nil
		})); err != nil {
		return err
	}

	go func() {
		if err := mcmanager.Start(ctx); err != nil {
			ctrl.Log.Error(err, "multicluster manager exited with error")
		}
	}()

	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	recvClusters := map[string]struct{}{}

	for {
		select {
		case <-ctx.Done():
			ctrl.Log.Info("context cancelled before first reconcile event")
			return nil
		case clusterName := <-recv:
			recvClusters[clusterName] = struct{}{}
			ctrl.Log.Info("received reconcile event, shutting down")
			if len(recvClusters) == len(kubeconfigs) {
				continue
			}
			cancel()
		case <-timeoutCtx.Done():
			cancel()
			ctrl.Log.Info("timeout waiting for first reconcile event, shutting down")
		}
	}

	return nil
}
