// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package geomfn

import (
	"github.com/cockroachdb/cockroach/pkg/geo"
	"github.com/cockroachdb/cockroach/pkg/geo/geopb"
	"github.com/cockroachdb/cockroach/pkg/geo/geos"
)

// Polygonize returns a GeometryCollection containing the polygons
// formed by the constituent linework of a set of geometries
func Polygonize(
	g []geo.Geometry, ngeoms int,
) (geo.Geometry, error) {
	var geoms []geopb.EWKB
	for _, s := range g {
		geoms = append(geoms, s.EWKB())
	}
	paths, err := geos.Polygonize(geoms, len(geoms))
	if err != nil {
		return geo.Geometry{}, err
	}
	gm, err := geo.ParseGeometryFromEWKB(paths)
	if err != nil {
		return geo.Geometry{}, err
	}
	return gm, nil
}
