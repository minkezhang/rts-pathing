package cluster

import (
	"testing"

	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

        "github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
)

func (m *ClusterMap) Equal(other *ClusterMap) bool {
        return m.l == other.l && proto.Equal(m.d, other.d) && cmp.Equal(m.m, other.m)
}

func (c *Cluster) Equal(other *Cluster) bool {
        return proto.Equal(c.c, other.c)
}

func TestIsAdjacent(t *testing.T) {
	testConfigs := []struct {
		name string
		c1   *rtsspb.Coordinate
		c2   *rtsspb.Coordinate
		want bool
	}{
		{name: "IsAdjacent", c1: &rtsspb.Coordinate{X: 0, Y: 0}, c2: &rtsspb.Coordinate{X: 0, Y: 1}, want: true},
		{name: "IsSame", c1: &rtsspb.Coordinate{X: 0, Y: 0}, c2: &rtsspb.Coordinate{X: 0, Y: 0}, want: false},
		{name: "IsDiagonal", c1: &rtsspb.Coordinate{X: 0, Y: 0}, c2: &rtsspb.Coordinate{X: 1, Y: 1}, want: false},
		{name: "IsNotAdjacent", c1: &rtsspb.Coordinate{X: 0, Y: 0}, c2: &rtsspb.Coordinate{X: 100, Y: 100}, want: false},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if res := IsAdjacent(
				&Cluster{c: &rtsspb.Cluster{Coordinate: c.c1}},
				&Cluster{c: &rtsspb.Cluster{Coordinate: c.c2}}); res != c.want {
				t.Errorf("IsAdjacent((%v, %v), (%v, %v)) = %v, want = %v", c.c1.GetX(), c.c1.GetY(), c.c2.GetX(), c.c2.GetY(), res, c.want)
			}
		})
	}

}

func TestPartition(t *testing.T) {
	testConfigs := []struct {
		name string
		tileMapDimension int32
		tileDimension int32
		want []partitionInfo
		wantSuccess bool
	}{
		{name: "ZeroWidthMapTest", tileMapDimension: 0, tileDimension: 1, want: nil, wantSuccess: true},
		{name: "ZeroWidthMapZeroDimTest", tileMapDimension: 0, tileDimension: 0, want: nil, wantSuccess: false},
		{name: "SimplePartitionTest", tileMapDimension: 1, tileDimension: 1, want: []partitionInfo{
			{TileBoundary: 0, TileDimension: 1},
		}, wantSuccess: true},
		{name: "SimplePartitionMultipleTest", tileMapDimension: 2, tileDimension: 1, want: []partitionInfo{
			{TileBoundary: 0, TileDimension: 1},
			{TileBoundary: 1, TileDimension: 1},
		}, wantSuccess: true},
		{name: "PartialPartitionTest", tileMapDimension: 1, tileDimension: 2, want: []partitionInfo{
			{TileBoundary: 0, TileDimension: 1},
		}, wantSuccess: true},
		{name: "PartialPartitionMultipleTest", tileMapDimension: 3, tileDimension: 2, want: []partitionInfo{
			{TileBoundary: 0, TileDimension: 2},
			{TileBoundary: 2, TileDimension: 1},
		}, wantSuccess: true},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			ps, err := partition(c.tileMapDimension, c.tileDimension)
			if (err == nil) != c.wantSuccess {
				t.Fatalf("partition() = _, %v, want wantSuccess = %v", err, c.wantSuccess)
			}
			if err == nil && !cmp.Equal(ps, c.want) {
				t.Errorf("partition() = %v, want = %v", ps, c.want)
			}
		})
	}
}

func TestBuildCluster(t *testing.T) {
	testConfigs := []struct {
		name string
		tileMapDimension *rtsspb.Coordinate
		tileDimension *rtsspb.Coordinate
		want *ClusterMap
		wantSuccess bool
	}{
		{name: "ZeroWidthDimTest", tileMapDimension: &rtsspb.Coordinate{X: 1, Y: 1}, tileDimension: &rtsspb.Coordinate{X: 0, Y: 0}, want: nil, wantSuccess: false},
		{name: "ZeroXDimTest", tileMapDimension: &rtsspb.Coordinate{X: 1, Y: 1}, tileDimension: &rtsspb.Coordinate{X: 0, Y: 1}, want: nil, wantSuccess: false},
		{name: "ZeroYDimTest", tileMapDimension: &rtsspb.Coordinate{X: 1, Y: 1}, tileDimension: &rtsspb.Coordinate{X: 1, Y: 0}, want: nil, wantSuccess: false},
		{name: "ZeroWidthMapTest", tileMapDimension: &rtsspb.Coordinate{X: 0, Y: 0}, tileDimension: &rtsspb.Coordinate{X: 1, Y: 1}, want: &ClusterMap{
			l: 1, d: &rtsspb.Coordinate{X: 0, Y: 0}, m: nil}, wantSuccess: true},
		{name: "ZeroXMapTest", tileMapDimension: &rtsspb.Coordinate{X: 0, Y: 1}, tileDimension: &rtsspb.Coordinate{X: 1, Y: 1}, want: &ClusterMap{
			l: 1, d: &rtsspb.Coordinate{X: 0, Y: 0}, m: nil}, wantSuccess: true},
		{name: "ZeroYMapTest", tileMapDimension: &rtsspb.Coordinate{X: 1, Y: 0}, tileDimension: &rtsspb.Coordinate{X: 1, Y: 1}, want: &ClusterMap{
			l: 1, d: &rtsspb.Coordinate{X: 0, Y: 0}, m: nil}, wantSuccess: true},
		{name: "SimpleTest", tileMapDimension: &rtsspb.Coordinate{X: 1, Y: 1}, tileDimension: &rtsspb.Coordinate{X: 1, Y: 1}, want: &ClusterMap{
			l: 1, d: &rtsspb.Coordinate{X: 1, Y: 1}, m: map[int32]map[int32]*Cluster{
				0: map[int32]*Cluster{
					0: &Cluster{
						c: &rtsspb.Cluster{
							Coordinate: &rtsspb.Coordinate{X: 0, Y: 0},
							TileBoundary: &rtsspb.Coordinate{X: 0, Y: 0},
							TileDimension: &rtsspb.Coordinate{X: 1, Y: 1},
						},
					},
				},
			}}, wantSuccess: true},
		{name: "MultiplePartitionTest", tileMapDimension: &rtsspb.Coordinate{X: 2, Y: 3}, tileDimension: &rtsspb.Coordinate{X: 2, Y: 2}, want: &ClusterMap{
			l: 1, d: &rtsspb.Coordinate{X: 1, Y: 2}, m: map[int32]map[int32]*Cluster{
				0: map[int32]*Cluster{
					0: &Cluster{
						c: &rtsspb.Cluster{
							Coordinate: &rtsspb.Coordinate{X: 0, Y: 0},
							TileBoundary: &rtsspb.Coordinate{X: 0, Y: 0},
							TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
						},
					},
					1: &Cluster{
						c: &rtsspb.Cluster{
							Coordinate: &rtsspb.Coordinate{X: 0, Y: 1},
							TileBoundary: &rtsspb.Coordinate{X: 0, Y: 2},
							TileDimension: &rtsspb.Coordinate{X: 2, Y: 1},
						},
					},
				},
			}}, wantSuccess: true},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			m, err := BuildClusterMap(c.tileMapDimension, c.tileDimension, 1)
			if (err == nil) != c.wantSuccess {
				t.Fatalf("BuildClusterMap() = _, %v, want wantSuccess = %v", err, c.wantSuccess)
			}
			if err == nil && !cmp.Equal(m, c.want) {
				t.Errorf("BuildClusterMap() = %v, want = %v", m, c.want)
			}
		})
	}
}