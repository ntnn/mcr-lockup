# mcr lockup

Create a kind cluster:

    kind create cluster --name mcr-lockup --kubeconfig mcr-lockup.kubeconfig


Run the working version:

    ( cd cmd/working && go run . -kubeconfigs ../../mcr-lockup.kubeconfig )

The working version will shutdown again after a few seconds after the
first reconcile event.

```
2025-11-23T22:56:48+01:00       INFO    controller-runtime.metrics      Starting metrics server
2025-11-23T22:56:48+01:00       INFO    Starting Controller     {"controller": "namespace", "controllerGroup": "", "controllerKind": "Namespace"}
2025-11-23T22:56:48+01:00       INFO    Starting workers        {"controller": "namespace", "controllerGroup": "", "controllerKind": "Namespace", "worker count": 1}
2025-11-23T22:56:48+01:00       INFO    Starting EventSource    {"controller": "namespace", "controllerGroup": "", "controllerKind": "Namespace", "source": "func source: 0x1376180"}
2025-11-23T22:56:48+01:00       INFO    controller-runtime.metrics      Serving metrics server  {"bindAddress": ":8080", "secure": false}
2025-11-23T22:56:48+01:00       INFO    reconciling namespace   {"namespace": "", "name": "kube-node-lease", "cluster": "cl"}
2025-11-23T22:56:48+01:00       INFO    received first reconcile event, shutting down
2025-11-23T22:56:48+01:00       INFO    reconciling namespace   {"namespace": "", "name": "kube-public", "cluster": "cl"}
```

Run the broken manager:

    ( cd cmd/broken && go run . -kubeconfigs ../../mcr-lockup.kubeconfig )

The broken manager will shutdown after a timeout of 5s without receiving
any reconcile events.

```
2025-11-23T22:56:53+01:00       INFO    controller-runtime.metrics      Starting metrics server
2025-11-23T22:56:53+01:00       INFO    Starting Controller     {"controller": "namespace", "controllerGroup": "", "controllerKind": "Namespace"}
2025-11-23T22:56:53+01:00       INFO    Starting workers        {"controller": "namespace", "controllerGroup": "", "controllerKind": "Namespace", "worker count": 1}
2025-11-23T22:56:53+01:00       INFO    controller-runtime.metrics      Serving metrics server  {"bindAddress": ":8080", "secure": false}
2025-11-23T22:56:58+01:00       INFO    timeout waiting for first reconcile event, shutting down
```
