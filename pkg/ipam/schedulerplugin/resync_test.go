package schedulerplugin

import (
	"net"
	"testing"

	"git.code.oa.com/gaiastack/galaxy/pkg/ipam/schedulerplugin/util"
)

func TestResyncAppNotExist(t *testing.T) {
	pod1 := createDeploymentPod("dp-xxx-yyy", "ns1", poolAnnotation("pool1"))
	pod2 := createDeploymentPod("dp2-aaa-bbb", "ns2", immutableAnnotation)
	fipPlugin, stopChan, _ := createPluginTestNodes(t)
	defer func() { stopChan <- struct{}{} }()
	var replicas int32 = 1
	fipPlugin.getDeployment = getDeploymentFunc(pod1.ObjectMeta, &replicas)
	pod1Key, pod2Key := util.FormatKey(pod1), util.FormatKey(pod2)

	if err := fipPlugin.ipam.AllocateSpecificIP(pod1Key.KeyInDB, net.ParseIP("10.49.27.205"), parseReleasePolicy(&pod1.ObjectMeta), ""); err != nil {
		t.Fatal(err)
	}
	if err := fipPlugin.ipam.AllocateSpecificIP(pod2Key.KeyInDB, net.ParseIP("10.49.27.216"), parseReleasePolicy(&pod2.ObjectMeta), ""); err != nil {
		t.Fatal(err)
	}
	if err := fipPlugin.resyncPod(fipPlugin.ipam); err != nil {
		t.Fatal(err)
	}
	if err := checkIPKey(fipPlugin.ipam, "10.49.27.205", pod1Key.PoolPrefix()); err != nil {
		t.Fatal(err)
	}
	if err := checkIPKey(fipPlugin.ipam, "10.49.27.216", ""); err != nil {
		t.Fatal(err)
	}
}