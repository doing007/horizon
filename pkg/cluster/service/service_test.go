// Copyright © 2023 Horizoncd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/horizoncd/horizon/lib/orm"
	applicationmodels "github.com/horizoncd/horizon/pkg/application/models"
	applicationservice "github.com/horizoncd/horizon/pkg/application/service"
	clustermodels "github.com/horizoncd/horizon/pkg/cluster/models"
	groupmodels "github.com/horizoncd/horizon/pkg/group/models"
	groupservice "github.com/horizoncd/horizon/pkg/group/service"
	"github.com/horizoncd/horizon/pkg/param/managerparam"
	"github.com/stretchr/testify/assert"
)

var (
	// use tmp sqlite
	db, _   = orm.NewSqliteDB("")
	ctx     = context.TODO()
	manager = managerparam.InitManager(db)
)

// nolint
func init() {
	// create table
	err := db.AutoMigrate(&clustermodels.Cluster{}, &applicationmodels.Application{}, &groupmodels.Group{})
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
}

func TestServiceGetByID(t *testing.T) {
	group := &groupmodels.Group{
		Name:         "a",
		Path:         "a",
		TraversalIDs: "1",
	}
	db.Save(group)

	application := &applicationmodels.Application{
		Name:    "b",
		GroupID: group.ID,
	}
	db.Save(application)

	cluster := &clustermodels.Cluster{
		Name:          "c",
		ApplicationID: application.ID,
	}
	db.Save(cluster)

	t.Run("GetByID", func(t *testing.T) {
		s := service{
			applicationService: applicationservice.NewService(groupservice.NewService(manager), manager),
			clusterManager:     manager.ClusterMgr,
		}
		result, err := s.GetByID(ctx, application.ID)
		assert.Nil(t, err)
		assert.Equal(t, "/a/b/c", result.FullPath)
	})
}
