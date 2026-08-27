package v1alpha1

import (
	"testing"
)

func TestDeepCopy_PlatformService(t *testing.T) {
	ps := &PlatformService{
		Spec: PlatformServiceSpec{
			Team:             "my-team",
			Tier:             "backend",
			ExposeExternally: true,
		},
		Status: PlatformServiceStatus{
			Phase:     "Ready",
			Namespace: "svc-test",
		},
	}
	ps.Labels = map[string]string{"key": "value"}

	copy := ps.DeepCopy()
	if copy == nil {
		t.Fatal("expected non-nil copy")
	}
	if copy.Spec.Team != ps.Spec.Team {
		t.Errorf("expected team %s, got %s", ps.Spec.Team, copy.Spec.Team)
	}
	if copy.Status.Phase != ps.Status.Phase {
		t.Errorf("expected phase %s, got %s", ps.Status.Phase, copy.Status.Phase)
	}
	if copy.Labels["key"] != "value" {
		t.Errorf("expected labels to be copied")
	}

	obj := ps.DeepCopyObject()
	if obj == nil {
		t.Fatal("expected non-nil object copy")
	}
}

func TestDeepCopy_NilPlatformService(t *testing.T) {
	var ps *PlatformService
	if ps.DeepCopy() != nil {
		t.Error("expected nil for nil input")
	}
}

func TestDeepCopy_PlatformServiceList(t *testing.T) {
	psl := &PlatformServiceList{
		Items: []PlatformService{
			{Spec: PlatformServiceSpec{Team: "team1"}},
			{Spec: PlatformServiceSpec{Team: "team2"}},
		},
	}

	copy := psl.DeepCopy()
	if copy == nil {
		t.Fatal("expected non-nil copy")
	}
	if len(copy.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(copy.Items))
	}

	obj := psl.DeepCopyObject()
	if obj == nil {
		t.Fatal("expected non-nil object copy")
	}
}

func TestDeepCopy_NilPlatformServiceList(t *testing.T) {
	var psl *PlatformServiceList
	if psl.DeepCopy() != nil {
		t.Error("expected nil for nil input")
	}
}

func TestSchemeGroupVersion(t *testing.T) {
	gv := SchemeGroupVersion()
	if gv.Group != "platform.example.com" {
		t.Errorf("unexpected group: %s", gv.Group)
	}
	if gv.Version != "v1alpha1" {
		t.Errorf("unexpected version: %s", gv.Version)
	}
}

func TestResource(t *testing.T) {
	gr := Resource("platformservices")
	if gr.Resource != "platformservices" {
		t.Errorf("unexpected resource: %s", gr.Resource)
	}
}
