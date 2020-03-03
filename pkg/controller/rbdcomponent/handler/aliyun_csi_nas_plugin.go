package handler

import (
	"context"

	rainbondv1alpha1 "github.com/goodrain/rainbond-operator/pkg/apis/rainbond/v1alpha1"
	"github.com/goodrain/rainbond-operator/pkg/util/commonutil"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type aliyunCsiNasPlugin struct {
	ctx    context.Context
	client client.Client

	component *rainbondv1alpha1.RbdComponent
	labels    map[string]string
}

var _ ComponentHandler = &aliyunCsiNasPlugin{}
var _ Replicaser = &aliyunCsiNasPlugin{}

// NewAliyunCSINasPlugin creates a new aliyun csi nas plugin handler.
func NewAliyunCSINasPlugin(ctx context.Context, client client.Client, component *rainbondv1alpha1.RbdComponent, cluster *rainbondv1alpha1.RainbondCluster) ComponentHandler {
	return &nfs{
		ctx:       ctx,
		client:    client,
		component: component,
		labels:    LabelsForRainbondComponent(component),
	}
}

func (h *aliyunCsiNasPlugin) Before() error {
	return nil
}

func (h *aliyunCsiNasPlugin) Resources() []interface{} {
	return nil
}

func (h *aliyunCsiNasPlugin) After() error {
	return nil
}

func (h *aliyunCsiNasPlugin) ListPods() ([]corev1.Pod, error) {
	return listPods(h.ctx, h.client, h.component.Namespace, h.labels)
}

func (h *aliyunCsiNasPlugin) Replicas() *int32 {
	nodeList := &corev1.NodeList{}
	if err := h.client.List(h.ctx, nodeList); err != nil {
		return nil
	}
	return commonutil.Int32(int32(len(nodeList.Items)))
}
