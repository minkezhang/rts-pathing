package tileastar

import (
	"math"
	"testing"

	"github.com/downflux/game/map/utils"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	gdpb "github.com/downflux/game/api/data_go_proto"
	mcpb "github.com/downflux/game/map/api/constants_go_proto"
	mdpb "github.com/downflux/game/map/api/data_go_proto"
	tile "github.com/downflux/game/map/map"
)

var (
	/**
	 * Y = 0 -
	 *   X = 0
	 */
	trivialOpenMap = &mdpb.TileMap{
		Dimension: &gdpb.Coordinate{X: 1, Y: 1},
		Tiles: []*mdpb.Tile{
			{
				Coordinate:  &gdpb.Coordinate{X: 0, Y: 0},
				TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS,
			},
		},
		TerrainCosts: []*mdpb.TerrainCost{},
	}

	/**
	 * Y = 0 W
	 *   X = 0
	 */
	trivialClosedMap = &mdpb.TileMap{
		Dimension: &gdpb.Coordinate{X: 1, Y: 1},
		Tiles: []*mdpb.Tile{
			{
				Coordinate:  &gdpb.Coordinate{X: 0, Y: 0},
				TerrainType: mcpb.TerrainType_TERRAIN_TYPE_BLOCKED,
			},
		},
		TerrainCosts: []*mdpb.TerrainCost{
			{TerrainType: mcpb.TerrainType_TERRAIN_TYPE_BLOCKED, Cost: math.Inf(0)},
		},
	}

	/**
	 *       W
	 * Y = 0 -
	 *   X = 0
	 */
	trivialSemiOpenMap = &mdpb.TileMap{
		Dimension: &gdpb.Coordinate{X: 1, Y: 2},
		Tiles: []*mdpb.Tile{
			{
				Coordinate:  &gdpb.Coordinate{X: 0, Y: 0},
				TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS,
			},
			{
				Coordinate:  &gdpb.Coordinate{X: 0, Y: 1},
				TerrainType: mcpb.TerrainType_TERRAIN_TYPE_BLOCKED,
			},
		},
		TerrainCosts: []*mdpb.TerrainCost{
			{TerrainType: mcpb.TerrainType_TERRAIN_TYPE_BLOCKED, Cost: math.Inf(0)},
		},
	}

	/**
	 *       -
	 *       W
	 * Y = 0 -
	 *   X = 0
	 */
	impassableMap = &mdpb.TileMap{
		Dimension: &gdpb.Coordinate{X: 1, Y: 3},
		Tiles: []*mdpb.Tile{
			{
				Coordinate:  &gdpb.Coordinate{X: 0, Y: 0},
				TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS,
			},
			{
				Coordinate:  &gdpb.Coordinate{X: 0, Y: 1},
				TerrainType: mcpb.TerrainType_TERRAIN_TYPE_BLOCKED,
			},
			{
				Coordinate:  &gdpb.Coordinate{X: 0, Y: 2},
				TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS,
			},
		},
		TerrainCosts: []*mdpb.TerrainCost{
			{TerrainType: mcpb.TerrainType_TERRAIN_TYPE_BLOCKED, Cost: math.Inf(0)},
		},
	}

	/**
	 *       - - -
	 *       W W W
	 * Y = 0 - - -
	 *   X = 0
	 */
	passableMap = &mdpb.TileMap{
		Dimension: &gdpb.Coordinate{X: 3, Y: 3},
		Tiles: []*mdpb.Tile{
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 2, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 1}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_BLOCKED},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 1}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_BLOCKED},
			{Coordinate: &gdpb.Coordinate{X: 2, Y: 1}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_BLOCKED},
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 2}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 2}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 2, Y: 2}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
		},
		TerrainCosts: []*mdpb.TerrainCost{
			{TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS, Cost: 1},
		},
	}

	/**
	 *       - - -
	 * Y = 0 - W -
	 *   X = 0
	 */
	blockedRowMap = &mdpb.TileMap{
		Dimension: &gdpb.Coordinate{X: 3, Y: 2},
		Tiles: []*mdpb.Tile{
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_BLOCKED},
			{Coordinate: &gdpb.Coordinate{X: 2, Y: 0}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 0, Y: 1}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 1, Y: 1}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
			{Coordinate: &gdpb.Coordinate{X: 2, Y: 1}, TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS},
		},
		TerrainCosts: []*mdpb.TerrainCost{
			{TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS, Cost: 1},
			{TerrainType: mcpb.TerrainType_TERRAIN_TYPE_BLOCKED, Cost: math.Inf(0)},
		},
	}
)

type aStarResult struct {
	path []*gdpb.Coordinate
	cost float64
}

func TestAStarSearchError(t *testing.T) {
	testConfigs := []struct {
		name                string
		m                   *mdpb.TileMap
		src, dest           *gdpb.Coordinate
		boundary, dimension *gdpb.Coordinate
	}{
		{
			name:      "SourceOutOfBounds",
			m:         trivialOpenMap,
			src:       &gdpb.Coordinate{X: 1, Y: 1},
			dest:      &gdpb.Coordinate{X: 0, Y: 0},
			dimension: trivialOpenMap.GetDimension(),
		},
		{
			name:      "DestinationOutOfBounds",
			m:         trivialOpenMap,
			src:       &gdpb.Coordinate{X: 0, Y: 0},
			dest:      &gdpb.Coordinate{X: 1, Y: 1},
			dimension: trivialOpenMap.GetDimension(),
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			tm, err := tile.ImportMap(c.m)
			if err != nil {
				t.Fatalf("ImportMap() = %v, want = nil", err)
			}

			if _, _, err = Path(tm, utils.MC(c.src), utils.MC(c.dest), c.boundary, c.dimension); err == nil {
				t.Fatal("Path() = nil, want a non-nil error")
			}
		})
	}
}

func TestAStarSearch(t *testing.T) {
	testConfigs := []struct {
		name                string
		m                   *mdpb.TileMap
		src, dest           *gdpb.Coordinate
		boundary, dimension *gdpb.Coordinate
		want                aStarResult
	}{
		{
			name:      "TrivialOpenMap",
			m:         trivialOpenMap,
			src:       &gdpb.Coordinate{X: 0, Y: 0},
			dest:      &gdpb.Coordinate{X: 0, Y: 0},
			dimension: trivialOpenMap.GetDimension(),
			want: aStarResult{
				path: []*gdpb.Coordinate{{X: 0, Y: 0}},
			},
		},
		{name: "TrivialClosedMap", m: trivialClosedMap, src: &gdpb.Coordinate{X: 0, Y: 0}, dest: &gdpb.Coordinate{X: 0, Y: 0}, dimension: trivialClosedMap.GetDimension(), want: aStarResult{
			path: nil,
			cost: math.Inf(0),
		}},
		{name: "BlockedSource", m: trivialSemiOpenMap, src: &gdpb.Coordinate{X: 0, Y: 1}, dest: &gdpb.Coordinate{X: 0, Y: 0}, dimension: trivialSemiOpenMap.GetDimension(), want: aStarResult{
			path: nil,
			cost: math.Inf(0),
		}},
		{name: "BlockedDestination", m: trivialSemiOpenMap, src: &gdpb.Coordinate{X: 0, Y: 0}, dest: &gdpb.Coordinate{X: 0, Y: 1}, dimension: trivialSemiOpenMap.GetDimension(), want: aStarResult{
			path: nil,
			cost: math.Inf(0),
		}},
		{name: "ImpassableMap", m: impassableMap, src: &gdpb.Coordinate{X: 0, Y: 0}, dest: &gdpb.Coordinate{X: 0, Y: 2}, dimension: impassableMap.GetDimension(), want: aStarResult{
			path: nil,
			cost: math.Inf(0),
		}},
		{
			name:      "SimpleSearch",
			m:         passableMap,
			src:       &gdpb.Coordinate{X: 0, Y: 0},
			dest:      &gdpb.Coordinate{X: 2, Y: 0},
			dimension: passableMap.GetDimension(),
			want: aStarResult{
				path: []*gdpb.Coordinate{
					{X: 0, Y: 0},
					{X: 1, Y: 0},
					{X: 2, Y: 0},
				},
				cost: 2,
			},
		},
		{
			name:      "SameSourceDestination",
			m:         passableMap,
			src:       &gdpb.Coordinate{X: 0, Y: 0},
			dest:      &gdpb.Coordinate{X: 0, Y: 0},
			dimension: passableMap.GetDimension(),
			want: aStarResult{
				path: []*gdpb.Coordinate{
					{X: 0, Y: 0},
				},
				cost: 0,
			},
		},
		{
			name:      "BlockedScopeSearch",
			m:         blockedRowMap,
			src:       &gdpb.Coordinate{X: 0, Y: 0},
			dest:      &gdpb.Coordinate{X: 2, Y: 0},
			dimension: &gdpb.Coordinate{X: 3, Y: 1},
			want: aStarResult{
				path: nil,
				cost: math.Inf(0),
			},
		},
		{
			name:      "ExpandedScopeSearch",
			m:         blockedRowMap,
			src:       &gdpb.Coordinate{X: 0, Y: 0},
			dest:      &gdpb.Coordinate{X: 2, Y: 0},
			dimension: blockedRowMap.GetDimension(),
			want: aStarResult{
				path: []*gdpb.Coordinate{
					{X: 0, Y: 0},
					{X: 0, Y: 1},
					{X: 1, Y: 1},
					{X: 2, Y: 1},
					{X: 2, Y: 0},
				},
				cost: 4,
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			tm, err := tile.ImportMap(c.m)
			if err != nil {
				t.Fatalf("ImportMap() = %v, want = nil", err)
			}

			tiles, cost, err := Path(tm, utils.MC(c.src), utils.MC(c.dest), c.boundary, c.dimension)
			if err != nil {
				t.Fatalf("Path() = %v, want = nil", err)
			}

			var path []*gdpb.Coordinate
			for _, t := range tiles {
				path = append(path, t.Val.GetCoordinate())
			}

			got := aStarResult{
				path: path,
				cost: cost,
			}
			if diff := cmp.Diff(c.want, got, cmp.AllowUnexported(aStarResult{}), protocmp.Transform()); diff != "" {
				t.Errorf("Path() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
