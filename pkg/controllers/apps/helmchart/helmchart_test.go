package helmchart

import (
	"context"
	"testing"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clienttesting "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"

	appsapi "github.com/clusternet/clusternet/pkg/apis/apps/v1alpha1"
	clusternetfake "github.com/clusternet/clusternet/pkg/generated/clientset/versioned/fake"
	informers "github.com/clusternet/clusternet/pkg/generated/informers/externalversions"
)

func TestShouldEnqueueHelmChartOnMetadataOnlyUpdateWhileStatusEmpty(t *testing.T) {
	oldChart := &appsapi.HelmChart{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-demo",
			Namespace: "default",
		},
		Spec: appsapi.HelmChartSpec{
			HelmOptions: appsapi.HelmOptions{
				Repository: "oci://registry/charts",
				Chart:      "nginx",
			},
			TargetNamespace: "default",
		},
	}
	newChart := oldChart.DeepCopy()
	newChart.Finalizers = []string{"apps.clusternet.io/finalizer"}

	if !shouldEnqueueHelmChart(oldChart, newChart) {
		t.Fatal("expected metadata-only update to be enqueued while status is empty")
	}
}

func TestShouldNotEnqueueHelmChartOnMetadataOnlyUpdateAfterVerification(t *testing.T) {
	oldChart := &appsapi.HelmChart{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-demo",
			Namespace: "default",
		},
		Spec: appsapi.HelmChartSpec{
			HelmOptions: appsapi.HelmOptions{
				Repository: "oci://registry/charts",
				Chart:      "nginx",
			},
			TargetNamespace: "default",
		},
		Status: appsapi.HelmChartStatus{
			Phase: appsapi.HelmChartFound,
		},
	}
	newChart := oldChart.DeepCopy()
	newChart.Finalizers = []string{"apps.clusternet.io/feed-protection"}

	if shouldEnqueueHelmChart(oldChart, newChart) {
		t.Fatal("expected metadata-only update to be skipped after verification")
	}
}

func TestUpdateChartStatusRetriesOnConflict(t *testing.T) {
	chart := &appsapi.HelmChart{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "nginx-ingress",
			Namespace:       "default",
			ResourceVersion: "1",
		},
	}
	controller, client := newTestController(t, chart)

	updateCalls := 0
	client.Fake.PrependReactor("update", "helmcharts", func(action clienttesting.Action) (bool, runtime.Object, error) {
		if action.GetSubresource() != "status" {
			return false, nil, nil
		}
		updateCalls++
		if updateCalls == 1 {
			return true, nil, apierrors.NewConflict(schema.GroupResource{Group: appsapi.SchemeGroupVersion.Group, Resource: "helmcharts"}, chart.Name, nil)
		}
		return false, nil, nil
	})

	err := controller.UpdateChartStatus(chart.DeepCopy(), &appsapi.HelmChartStatus{
		Phase: appsapi.HelmChartFound,
	})
	if err != nil {
		t.Fatalf("UpdateChartStatus returned error: %v", err)
	}
	if updateCalls != 2 {
		t.Fatalf("expected 2 status update attempts, got %d", updateCalls)
	}

	got, err := client.AppsV1alpha1().HelmCharts(chart.Namespace).Get(context.Background(), chart.Name, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("failed to get helmchart: %v", err)
	}
	if got.Status.Phase != appsapi.HelmChartFound {
		t.Fatalf("expected phase %q, got %q", appsapi.HelmChartFound, got.Status.Phase)
	}
}

func TestUpdateChartStatusReturnsNonConflictError(t *testing.T) {
	chart := &appsapi.HelmChart{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-ingress",
			Namespace: "default",
		},
	}
	controller, client := newTestController(t, chart)

	client.Fake.PrependReactor("update", "helmcharts", func(action clienttesting.Action) (bool, runtime.Object, error) {
		if action.GetSubresource() != "status" {
			return false, nil, nil
		}
		return true, nil, apierrors.NewForbidden(schema.GroupResource{Group: appsapi.SchemeGroupVersion.Group, Resource: "helmcharts"}, chart.Name, nil)
	})

	err := controller.UpdateChartStatus(chart.DeepCopy(), &appsapi.HelmChartStatus{
		Phase:  appsapi.HelmChartNotFound,
		Reason: "forbidden",
	})
	if err == nil {
		t.Fatal("expected UpdateChartStatus to return an error")
	}

	got, getErr := client.AppsV1alpha1().HelmCharts(chart.Namespace).Get(context.Background(), chart.Name, metav1.GetOptions{})
	if getErr != nil {
		t.Fatalf("failed to get helmchart: %v", getErr)
	}
	if got.Status.Phase != "" {
		t.Fatalf("expected empty phase after failed status update, got %q", got.Status.Phase)
	}
}

func newTestController(t *testing.T, chart *appsapi.HelmChart) (*Controller, *clusternetfake.Clientset) {
	t.Helper()

	client := clusternetfake.NewSimpleClientset(chart.DeepCopy())
	factory := informers.NewSharedInformerFactory(client, 0*time.Second)
	chartInformer := factory.Apps().V1alpha1().HelmCharts()

	if err := chartInformer.Informer().GetIndexer().Add(chart.DeepCopy()); err != nil {
		t.Fatalf("failed to seed helmchart informer: %v", err)
	}

	controller := &Controller{
		clusternetClient: client,
		helmChartLister:  chartInformer.Lister(),
		baseIndexer:      cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}),
		recorder:         record.NewFakeRecorder(10),
	}
	return controller, client
}
