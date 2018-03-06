/*
Copyright (C) 2018 Synopsys, Inc.

Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements. See the NOTICE file
distributed with this work for additional information
regarding copyright ownership. The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied. See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"testing"

	"github.com/spf13/viper"
)

func TestProto(t *testing.T) {
	viper.SetDefault("protoform", map[string]string{
		"HubUserPassword": "ASDF",
		"DryRun":          "true",
	})
	rcsArray := runProtoform("./")

	// Image facade needs to be privileged !
	if *rcsArray[2].Spec.Template.Spec.Containers[1].SecurityContext.Privileged == false {
		t.Log("%v %v", rcsArray[3].Spec.Template.Spec.Containers[0].Name, *rcsArray[3].Spec.Template.Spec.Containers[0].SecurityContext.Privileged)
		t.Fail()
	}

	// The scanner needs to be UNPRIVILEGED
	if *rcsArray[2].Spec.Template.Spec.Containers[0].SecurityContext.Privileged == true {
		t.Log("%v %v", rcsArray[3].Spec.Template.Spec.Containers[0].Name, *rcsArray[3].Spec.Template.Spec.Containers[0].SecurityContext.Privileged)
		t.Fail()
	}

	scanner_svc := rcsArray[2].Spec.Template.Spec.ServiceAccountName
	if scanner_svc == "" {
		t.Log("scanner svc ==> ( %v ) EMPTY !", scanner_svc)
		t.Fail()
	}

	s0 := rcsArray[2].Spec.Template.Spec.Containers[0].Name
	s := rcsArray[2].Spec.Template.Spec.Containers[0].VolumeMounts[1].Name
	if s != "var-images" {
		t.Log("%v %v", s0, s)
		t.Fail()
	}
}
