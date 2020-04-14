package v1

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	coreV1 "k8s.io/api/core/v1"
)

func TestShadowDiff(t *testing.T) {
	tests := []struct {
		name      string
		desire    Desire
		report    Report
		wantDelta Desire
		wantErr   error
	}{
		{
			name:      "nil-1",
			desire:    Desire{},
			report:    nil,
			wantDelta: Desire{},
		},
		{
			name:      "0",
			desire:    Desire{},
			report:    Report{},
			wantDelta: Desire{},
		},
		{
			name:      "1",
			desire:    Desire{"name": "module", "version": "45"},
			report:    Report{"name": "module", "version": "43"},
			wantDelta: Desire{"version": "45"},
		},
		{
			name:      "2",
			desire:    Desire{"name": "module", "module": map[string]interface{}{"image": "test:v2"}},
			report:    Report{"name": "module", "module": map[string]interface{}{"image": "test:v1"}},
			wantDelta: Desire{"module": map[string]interface{}{"image": "test:v2"}},
		},
		{
			name:      "3",
			desire:    Desire{"module": map[string]interface{}{"image": "test:v2", "array": []interface{}{}}},
			report:    Report{"module": map[string]interface{}{"image": "test:v1", "object": map[string]interface{}{"attr": "value"}}},
			wantDelta: Desire{"module": map[string]interface{}{"image": "test:v2", "array": []interface{}{}}},
		},
		{
			name:      "6",
			desire:    Desire{"1": map[string]interface{}{"2": map[string]interface{}{"3": map[string]interface{}{"4": map[string]interface{}{"n": nil, "5": map[string]interface{}{"6": "x"}}}}}},
			report:    Report{"1": map[string]interface{}{"2": map[string]interface{}{"3": map[string]interface{}{"4": map[string]interface{}{"5": map[string]interface{}{"n": nil, "6": "y"}}}}}},
			wantDelta: Desire{"1": map[string]interface{}{"2": map[string]interface{}{"3": map[string]interface{}{"4": map[string]interface{}{"5": map[string]interface{}{"6": "x"}}}}}},
		},
		{
			name:      "apps",
			desire:    Desire{"apps": []interface{}{map[string]interface{}{"name": "a", "version": "1"}, map[string]interface{}{"name": "b", "version": "1"}}},
			report:    Report{"apps": []interface{}{map[string]interface{}{"name": "b", "version": "2"}, map[string]interface{}{"name": "c", "version": "2"}}},
			wantDelta: Desire{"apps": []interface{}{map[string]interface{}{"name": "a", "version": "1"}, map[string]interface{}{"name": "b", "version": "1"}}},
		},
		{
			name:      "apps-2",
			desire:    Desire{"apps": []interface{}{map[string]interface{}{"name": "a", "version": "1"}, map[string]interface{}{"name": "b", "version": "1"}}},
			report:    Report{"apps": nil},
			wantDelta: Desire{"apps": []interface{}{map[string]interface{}{"name": "a", "version": "1"}, map[string]interface{}{"name": "b", "version": "1"}}},
		},
		{
			name:      "apps-3",
			desire:    Desire{"apps": nil},
			report:    Report{"apps": []interface{}{map[string]interface{}{"name": "a", "version": "1"}, map[string]interface{}{"name": "b", "version": "1"}}},
			wantDelta: Desire{},
		},
		{
			name:      "apps-4",
			desire:    Desire{"apps": []interface{}{}},
			report:    Report{"apps": []interface{}{map[string]interface{}{"name": "a", "version": "1"}, map[string]interface{}{"name": "b", "version": "1"}}},
			wantDelta: Desire{"apps": []interface{}{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDelta, err := tt.desire.Diff(tt.report)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantDelta, gotDelta)
			assert.Equal(t, tt.desire.AppInfos(), gotDelta.AppInfos())
		})
	}
}

func TestShadowMerge(t *testing.T) {
	tests := []struct {
		name     string
		oldData  map[string]interface{}
		newData  map[string]interface{}
		wantData map[string]interface{}
		wantErr  error
	}{
		{
			name:     "nil-1",
			oldData:  map[string]interface{}{},
			newData:  nil,
			wantData: map[string]interface{}{},
		},
		{
			name:     "0",
			oldData:  map[string]interface{}{},
			newData:  map[string]interface{}{},
			wantData: map[string]interface{}{},
		},
		{
			name:     "1",
			oldData:  map[string]interface{}{"name": "module", "version": "45"},
			newData:  map[string]interface{}{"name": "module", "version": "43"},
			wantData: map[string]interface{}{"name": "module", "version": "43"},
		},
		{
			name:     "2",
			oldData:  map[string]interface{}{"name": "module", "module": map[string]interface{}{"image": "test:v2"}},
			newData:  map[string]interface{}{"name": "module", "module": map[string]interface{}{"image": "test:v1"}},
			wantData: map[string]interface{}{"name": "module", "module": map[string]interface{}{"image": "test:v1"}},
		},
		{
			name:     "3",
			oldData:  map[string]interface{}{"module": map[string]interface{}{"image": "test:v2", "array": []interface{}{}}},
			newData:  map[string]interface{}{"module": map[string]interface{}{"image": "test:v1", "object": map[string]interface{}{"attr": "value"}}},
			wantData: map[string]interface{}{"module": map[string]interface{}{"image": "test:v1", "array": []interface{}{}, "object": map[string]interface{}{"attr": "value"}}},
		},
		{
			name:     "err",
			oldData:  map[string]interface{}{"1": map[string]interface{}{"2": map[string]interface{}{"3": map[string]interface{}{"4": map[string]interface{}{"5": map[string]interface{}{"6": "y"}}}}}},
			newData:  map[string]interface{}{"1": map[string]interface{}{"2": map[string]interface{}{"3": map[string]interface{}{"4": map[string]interface{}{"n": nil, "5": map[string]interface{}{"n": nil, "6": "x"}}}}}},
			wantData: map[string]interface{}{"1": map[string]interface{}{"2": map[string]interface{}{"3": map[string]interface{}{"4": map[string]interface{}{"5": map[string]interface{}{"6": "y"}}}}}},
			wantErr:  ErrJSONLevelExceedsLimit,
		},
		{
			name:     "apps-1",
			oldData:  map[string]interface{}{"apps": []interface{}{map[string]interface{}{"name": "a", "version": "1"}, map[string]interface{}{"name": "b", "version": "1"}}},
			newData:  map[string]interface{}{"apps": []interface{}{map[string]interface{}{"name": "b", "version": "2"}, map[string]interface{}{"name": "c", "version": "2"}}},
			wantData: map[string]interface{}{"apps": []interface{}{map[string]interface{}{"name": "b", "version": "2"}, map[string]interface{}{"name": "c", "version": "2"}}},
		},
		{
			name:     "apps-2",
			oldData:  map[string]interface{}{"apps": []interface{}{map[string]interface{}{"name": "a", "version": "1"}, map[string]interface{}{"name": "b", "version": "1"}}},
			newData:  map[string]interface{}{"apps": nil},
			wantData: map[string]interface{}{"apps": nil},
		},
		{
			name:     "apps-3",
			oldData:  map[string]interface{}{"apps": []interface{}{map[string]interface{}{"name": "a", "version": "1"}, map[string]interface{}{"name": "b", "version": "1"}}},
			newData:  map[string]interface{}{"apps": []interface{}{}},
			wantData: map[string]interface{}{"apps": []interface{}{}},
		},
		{
			name:     "apps-4",
			oldData:  map[string]interface{}{"apps": nil},
			newData:  map[string]interface{}{"apps": []interface{}{map[string]interface{}{"name": "a", "version": "1"}, map[string]interface{}{"name": "b", "version": "1"}}},
			wantData: map[string]interface{}{"apps": []interface{}{map[string]interface{}{"name": "a", "version": "1"}, map[string]interface{}{"name": "b", "version": "1"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			or, nr := Report(tt.oldData), Report(tt.newData)
			err := or.Merge(nr)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, Report(tt.wantData), or)

			if tt.name == "err" {
				od, nd := Desire(tt.oldData), Desire(tt.newData)
				err = od.Merge(nd)
				assert.Equal(t, tt.wantErr, err)
				assert.Equal(t, Desire(tt.wantData), od)
			} else {
				assert.Equal(t, nr.AppInfos(), or.AppInfos())
			}
		})
	}
}

func TestDesireSysAppInfos(t *testing.T) {
	sysApps := Desire{
		"sysapps": []interface{}{
			map[string]interface{}{"name": "app1", "version": "1"},
			map[string]interface{}{"name": "app2", "version": "2"},
		},
	}

	expectApps := []AppInfo{
		{
			Name:    "app1",
			Version: "1",
		},
		{
			Name:    "app2",
			Version: "2",
		},
	}

	assert.Equal(t, expectApps, sysApps.SysAppInfos())
}

func TestAppInfos(t *testing.T) {
	assert.Nil(t, Report{}.AppInfos())
	assert.Nil(t, Report{}.SysAppInfos())
	assert.Nil(t, Report{"apps": nil}.AppInfos())
	assert.Nil(t, Report{"sysapps": nil}.SysAppInfos())
	assert.Nil(t, Report{"apps": []string{}}.AppInfos())
	assert.Nil(t, Report{"sysapps": []string{}}.SysAppInfos())
	assert.Equal(t, []AppInfo{}, Report{"apps": []AppInfo{}}.AppInfos())
	assert.Equal(t, []AppInfo{}, Report{"sysapps": []AppInfo{}}.SysAppInfos())
	assert.Equal(t, []AppInfo{}, Report{"apps": []interface{}{}}.AppInfos())
	assert.Equal(t, []AppInfo{}, Report{"sysapps": []interface{}{}}.SysAppInfos())

	expectApps := []AppInfo{
		{
			Name:    "app1",
			Version: "1",
		},
		{
			Name:    "app2",
			Version: "2",
		},
	}

	r := Report{
		"apps": []interface{}{
			map[string]interface{}{"name": "app1", "version": "1"},
			map[string]interface{}{"name": "app2", "version": "2"},
		},
	}
	assert.Equal(t, expectApps, r.AppInfos())

	r = Report{
		"sysapps": []AppInfo{
			{Name: "app1", Version: "1"},
			{Name: "app2", Version: "2"},
		},
	}

	assert.Equal(t, expectApps, r.SysAppInfos())
}

func TestTranslateNodeToNodeReportView(t *testing.T) {
	nodeData := `
{
		"namespace":"default",
		"name":"baetyl",
		"version":"v1",
		"createTime":"2020-04-11T00:21:35.588279937Z",
		"report":{"apps":[{"name":"function-testnode4","version":"1436880"},{"name":"core-testnode4","version":"1436874"},{"name":"timer","version":"1449120"}],"appstats":[{"name":"function-testnode4","version":"1436880","status":"Running","services":{"baetyl-function":{"name":"baetyl-function","container":{"name":"baetyl-function","id":"docker://ad4474b017ace7324884452ecbfaa7206b8201b33d5ce909d776f0887443b611"},"usage":{"cpu":"37310n","memory":"1728Ki"},"status":"Running","createTime":"2020-04-10T06:07:46Z"}}},{"name":"core-testnode4","version":"1436874","status":"Running","services":{"core-testnode4":{"name":"core-testnode4","container":{"name":"core-testnode4","id":"docker://65d7e9c5754a7702f28cb030202a639e729061f59e88d8d0afbb175b86da203b"},"usage":{"cpu":"8553885n","memory":"8708Ki"},"status":"Running","createTime":"2020-04-10T06:07:50Z"}}},{"name":"timer","version":"1449120","status":"Running","services":{"timer":{"name":"timer","container":{"name":"timer","id":"docker://3e468a0a55f0aa3dfab65cc74f3d7ee1e3b3e012803bca0179e9f7275d8afe89"},"status":"Running","cause":"Back-off restarting failed container\n","createTime":"2020-04-10T09:24:25Z"}}}],"node":{"hostname":"docker-desktop","address":"192.168.65.3","arch":"amd64","kernelVer":"4.19.76-linuxkit","os":"linux","containerRuntime":"docker://19.3.8","machineID":"301b485d-749e-4dc5-9bc3-319795e4021c","bootID":"aa584840-2149-44ac-acab-55ef67f69792","systemUUID":"df4949c4-0000-0000-b259-6cc320f45d7d","osImage":"Docker Desktop"},"nodestats":{"usage":{"cpu":"336037951n","memory":"1206552Ki"},"capacity":{"cpu":"2","memory":"4033160Ki"}},"time":"2020-04-11T00:21:35.588279937Z"},
		"desire":{"apps":"name","age":"12"}
}
`
	node := new(Node)
	err := json.Unmarshal([]byte(nodeData), node)
	assert.NoError(t, err)

	view := node.View()
	assert.NoError(t, err)
	assert.NotNil(t, view)
	assert.Equal(t, view.Namespace, "default")
	assert.Equal(t, view.Name, "baetyl")
	assert.Equal(t, view.Version, "v1")
	assert.NotNil(t, view.Report)
	assert.NotNil(t, view.Report.Apps)
	assert.NotNil(t, view.Report.NodeStatus)
	assert.NotNil(t, view.Report.Appstats)
	assert.NotNil(t, view.Desire)
	assert.Equal(t, view.Report.NodeStatus.Capacity[string(coreV1.ResourceMemory)], "4129955840")
	assert.Equal(t, view.Report.NodeStatus.Capacity[string(coreV1.ResourceCPU)], "2")
	assert.Equal(t, view.Report.NodeStatus.Usage[string(coreV1.ResourceMemory)], "1235509248")
	assert.Equal(t, view.Report.NodeStatus.Usage[string(coreV1.ResourceCPU)], "0.337")
	assert.Equal(t, view.Report.NodeStatus.Percent[string(coreV1.ResourceCPU)], "0.1685")
	assert.Equal(t, view.Report.NodeStatus.Percent[string(coreV1.ResourceMemory)], "0.2991579803429569")
}

func TestPopulateNodeStatus(t *testing.T) {
	report1 := &ReportView{
		NodeStatus: &NodeStatus{
			Usage: map[string]string{
				"cpu":    "1",
				"memory": "512Mi",
			},
			Capacity: map[string]string{
				"cpu":    "2",
				"memory": "1024Mi",
			},
		},
	}
	err := report1.NodeStatus.populateNodeStatus()
	assert.NoError(t, err)
	m1, err1 := translateQuantityToDecimal("1024Mi", false)
	assert.NoError(t, err1)
	s1 := strconv.FormatInt(m1, 10)
	m2, err2 := translateQuantityToDecimal("512Mi", false)
	assert.NoError(t, err2)
	s2 := strconv.FormatInt(m2, 10)
	assert.Equal(t, report1.NodeStatus.Capacity[string(coreV1.ResourceMemory)], s1)
	assert.Equal(t, report1.NodeStatus.Capacity[string(coreV1.ResourceCPU)], "2")
	assert.Equal(t, report1.NodeStatus.Usage[string(coreV1.ResourceMemory)], s2)
	assert.Equal(t, report1.NodeStatus.Usage[string(coreV1.ResourceCPU)], "1")
	assert.Equal(t, report1.NodeStatus.Percent[string(coreV1.ResourceMemory)], "0.5")
	assert.Equal(t, report1.NodeStatus.Percent[string(coreV1.ResourceCPU)], "0.5")

	report2 := &ReportView{
		NodeStatus: &NodeStatus{
			Usage: map[string]string{
				"cpu":    "500m",
				"memory": "512Mi",
			},
			Capacity: map[string]string{
				"cpu":    "2.0",
				"memory": "1024Mi",
			},
		},
	}
	err3 := report2.NodeStatus.populateNodeStatus()
	assert.NoError(t, err3)
	assert.Equal(t, report2.NodeStatus.Capacity[string(coreV1.ResourceCPU)], "2")
	assert.Equal(t, report2.NodeStatus.Usage[string(coreV1.ResourceCPU)], "0.5")
	assert.Equal(t, report2.NodeStatus.Percent[string(coreV1.ResourceMemory)], "0.5")
	assert.Equal(t, report2.NodeStatus.Percent[string(coreV1.ResourceCPU)], "0.25")

	report3 := &ReportView{
		NodeStatus: &NodeStatus{
			Usage: map[string]string{
				"cpu":    "0.5",
				"memory": "512a",
			},
			Capacity: map[string]string{
				"cpu":    "2.5",
				"memory": "1024a",
			},
		},
	}
	err4 := report3.NodeStatus.populateNodeStatus()
	assert.Error(t, err4)

	report4 := &ReportView{
		NodeStatus: &NodeStatus{
			Usage: map[string]string{
				"cpu":    "0.5a",
				"memory": "512Mi",
			},
			Capacity: map[string]string{
				"cpu":    "2.5s",
				"memory": "1024Mi",
			},
		},
	}
	err5 := report4.NodeStatus.populateNodeStatus()
	assert.Error(t, err5)

	report5 := &ReportView{
		NodeStatus: &NodeStatus{
			Usage: map[string]string{
				"cpu":    "0",
				"memory": "0",
			},
			Capacity: map[string]string{
				"cpu":    "0",
				"memory": "0",
			},
		},
	}
	err6 := report5.NodeStatus.populateNodeStatus()
	assert.NoError(t, err6)
	assert.Equal(t, report5.NodeStatus.Capacity[string(coreV1.ResourceCPU)], "0")
	assert.Equal(t, report5.NodeStatus.Usage[string(coreV1.ResourceCPU)], "0")
	assert.Equal(t, report5.NodeStatus.Percent[string(coreV1.ResourceMemory)], "0")
	assert.Equal(t, report5.NodeStatus.Percent[string(coreV1.ResourceCPU)], "0")
}