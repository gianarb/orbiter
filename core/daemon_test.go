package core

import "testing"

func TestNewCore(t *testing.T) {
	conf := map[string]AutoscalerConf{
		"first-scaler": AutoscalerConf{
			Provider:   "fake",
			Parameters: map[string]string{},
			Policies: map[string]PolicyConf{
				"frontend": PolicyConf{
					Up:   3,
					Down: 10,
				},
			},
		},
		"second-scaler": AutoscalerConf{
			Provider:   "fake",
			Parameters: map[string]string{},
			Policies: map[string]PolicyConf{
				"micro": PolicyConf{
					Up:   6,
					Down: 2,
				},
				"service": PolicyConf{
					Up:   3,
					Down: 1,
				},
			},
		},
	}
	core, err := NewCore(conf)
	if err != nil {
		t.Fatal(err)
	}
	if len(core.Autoscalers) != 3 {
		t.Fatalf("This core needs to have 2 autoscalers. Not %d", len(core.Autoscalers))
	}
}

func TestGetSingleAutoscaler(t *testing.T) {
	conf := map[string]AutoscalerConf{
		"second": AutoscalerConf{
			Provider:   "fake",
			Parameters: map[string]string{},
			Policies: map[string]PolicyConf{
				"micro": PolicyConf{
					Up:   6,
					Down: 2,
				},
				"service": PolicyConf{
					Up:   3,
					Down: 1,
				},
			},
		},
	}
	core, _ := NewCore(conf)
	_, ok := core.Autoscalers["second/micro"]
	if ok == false {
		t.Fatal("micro exist")
	}
}

func TestNewCoreWithUnsupportedProvider(t *testing.T) {
	conf := map[string]AutoscalerConf{
		"second-scaler": AutoscalerConf{
			Provider:   "fake",
			Parameters: map[string]string{},
			Policies: map[string]PolicyConf{
				"micro": PolicyConf{
					Up:   6,
					Down: 2,
				},
				"service": PolicyConf{
					Up:   3,
					Down: 1,
				},
			},
		},
		"first-scaler": AutoscalerConf{
			Provider:   "lalala",
			Parameters: map[string]string{},
			Policies: map[string]PolicyConf{
				"frontend": PolicyConf{
					Up:   3,
					Down: 10,
				},
			},
		},
	}
	_, err := NewCore(conf)
	if err.Error() != "lalala not supported." {
		t.Fatal(err)
	}
}
