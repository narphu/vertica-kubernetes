/*
 (c) Copyright [2021-2022] Micro Focus or one of its affiliates.
 Licensed under the Apache License, Version 2.0 (the "License");
 You may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	vapi "github.com/vertica/vertica-kubernetes/api/v1beta1"
	"github.com/vertica/vertica-kubernetes/pkg/cmds"
	"github.com/vertica/vertica-kubernetes/pkg/events"
	"github.com/vertica/vertica-kubernetes/pkg/names"
	"github.com/vertica/vertica-kubernetes/pkg/status"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	ConsolePodName       = "verticadb-operator-console-0"
	ConsoleContainerName = "console"
)

// AgentReconciler will ensure the agent is running
type MCImportDBReconciler struct {
	VRec    *VerticaDBReconciler
	Log     logr.Logger
	Vdb     *vapi.VerticaDB // Vdb is the CRD we are acting on.
	PRunner cmds.PodRunner
	PFacts  *PodFacts
}

// MakeMCImportDBReconcoler will build a MCImportDBReconciler object
func MakeMCImportDBReconciler(vdbrecon *VerticaDBReconciler, log logr.Logger,
	vdb *vapi.VerticaDB, prunner cmds.PodRunner, pfacts *PodFacts) ReconcileActor {
	return &MCImportDBReconciler{VRec: vdbrecon, Log: log, Vdb: vdb, PRunner: prunner, PFacts: pfacts}
}

// Reconcile will ensure the agent is running and start it if it isn't
func (m *MCImportDBReconciler) Reconcile(ctx context.Context, req *ctrl.Request) (ctrl.Result, error) {
	if m.Vdb.IsMCImportComplete() {
		m.Log.Info("MC import has already happened")
		return ctrl.Result{}, nil
	}
	if consolePodRunning, err := m.isConsolePodRunning(ctx); err != nil {
		return ctrl.Result{}, err
	} else if !consolePodRunning {
		m.Log.Info("Console pod isn't running, so skipping the import of the database")
		return ctrl.Result{}, nil
	}

	if err := m.PFacts.Collect(ctx, m.Vdb); err != nil {
		return ctrl.Result{}, err
	}

	if err := m.runImportDBCommand(ctx); err != nil {
		return ctrl.Result{}, err
	}

	// SPILLY -rename a to m
	err := status.UpdateCondition(ctx, m.VRec.Client, m.Vdb,
		vapi.VerticaDBCondition{Type: vapi.MCImportCompleted, Status: corev1.ConditionTrue},
	)
	return ctrl.Result{}, err
}

func (m *MCImportDBReconciler) isConsolePodRunning(ctx context.Context) (bool, error) {
	pod := &corev1.Pod{}
	nm := names.GenNamespacedName(m.Vdb, ConsolePodName)
	if err := m.VRec.Client.Get(ctx, nm, pod); err != nil && !errors.IsNotFound(err) {
		return false, fmt.Errorf("could not fetch console pod %w", err)
	}
	return pod.Status.Phase == corev1.PodRunning, nil
}

// startAgentInPod will start the agent in the given pod.
func (m *MCImportDBReconciler) runImportDBCommand(ctx context.Context) error {
	node1IP, err := m.getIPForNode1()
	if err != nil {
		return err
	}
	passwd, err := GetSuperuserPasswordForClient(ctx, m.VRec.Client, m.Vdb, m.Log)
	if err != nil {
		return err
	}

	m.VRec.EVRec.Eventf(m.Vdb, corev1.EventTypeNormal, events.MCImport,
		"Importing the database into the management console")

	cmd := []string{
		"java",
		"-jar", "/jars/MCClient.jar",
		"-import_database",
		"-api_key", APIKey,
		"-db_admin_user", "dbadmin",
		"-db_admin_pwd", passwd,
		"-db_name", m.Vdb.Spec.DBName,
		"-debug",
		"-vertica_node1_ip", node1IP,
		"-mc_url", m.getMCUrl(),
		// SPILLY - need to get password from somewhere
		"-mc_dbadmin_pwd", "dbadmin",
		"-mc_dbadmin_user", "dbadmin",
	}
	m.Log.Info("Importing DB into MC", "cmd", cmd)
	pn := names.GenNamespacedName(m.Vdb, ConsolePodName)
	_, _, err = m.PRunner.ExecInPod(ctx, pn, ConsoleContainerName, cmd...)
	return err
}

func (m *MCImportDBReconciler) getIPForNode1() (string, error) {
	pf, ok := m.PFacts.findRunningPod()
	if !ok {
		return "", fmt.Errorf("could not find running pod")
	}
	return pf.podIP, nil
}

func (m *MCImportDBReconciler) getMCUrl() string {
	return fmt.Sprintf("https://%s.verticadb-operator-console.%s.svc.cluster.local:5450/webui",
		ConsolePodName, m.Vdb.Namespace)
}
