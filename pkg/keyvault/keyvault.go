/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package keyvault

import (
	// "fmt"
	// "strings"

	"github.com/ritazh/keyvault-csi-driver/pkg/csi-common"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
)

type keyvault struct {
	driver   *csicommon.CSIDriver
	ns       *nodeServer
}

type keyvaultVolume struct {
	VolName string `json:"volName"`
	VolID   string `json:"volID"`
	VolSize int64  `json:"volSize"`
	VolPath string `json:"volPath"`
}

var keyvaultVolumes map[string]keyvaultVolume

var (
	keyvaultDriver *keyvault
	vendorVersion   = "0.0.2"
)

func init() {
	keyvaultVolumes = map[string]keyvaultVolume{}
}

func GetKeyvaultDriver() *keyvault {
	return &keyvault{}
}

func NewNodeServer(d *csicommon.CSIDriver) *nodeServer {
	return &nodeServer{
		DefaultNodeServer: csicommon.NewDefaultNodeServer(d),
	}
}

func (k *keyvault) Run(driverName, nodeID, endpoint string) {
	glog.Infof("Driver: %v ", driverName)
	glog.Infof("Version: %s", vendorVersion)

	// Initialize default library driver
	k.driver = csicommon.NewCSIDriver(driverName, vendorVersion, nodeID)
	if k.driver == nil {
		glog.Fatalln("Failed to initialize keyvault CSI Driver.")
	}
	k.driver.AddControllerServiceCapabilities(
		[]csi.ControllerServiceCapability_RPC_Type{
			csi.ControllerServiceCapability_RPC_PUBLISH_READONLY,
		})
	k.driver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{
		csi.VolumeCapability_AccessMode_SINGLE_NODE_READER_ONLY,
		csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY,
	})

	k.ns = NewNodeServer(k.driver)

	s := csicommon.NewNonBlockingGRPCServer()
	s.Start(endpoint, csicommon.NewDefaultIdentityServer(k.driver), csicommon.NewDefaultControllerServer(k.driver), k.ns)
	s.Wait()
}
