package a2

// Province adjacency map.
//   - Keys are province names, all lowercase, spaces converted to underscores.
//   - Values are a slice of province names, that the given key province is adjacent to.
type AdjacenciesType map[string][]string

// The actual adjacency map, going from top-to-bottom, left-to-right.
// Provinces are ordered roughly by:
//   - land
//   - sea
//
// then by continent:
//   - north america
//   - south america
//   - europe
//   - africa
//   - asia
//
// then by ocean:
//   - arctic
//   - antarctic
//   - atlantic
//   - indian
//   - pacific
//
// Adjacencies are also ordered clockwise.
var Adjacencies = AdjacenciesType{
	Nome:    p(ChukchiSea, ArcticOcean, Koyukon, Yupik),
	Yupik:   p(Nome, Koyukon, Denaina, KvichakBay, ChukchiSea),
	Aleuts:  p(ChukchiSea, KvichakBay, NorthEquatorialCurrent),
	Koyukon: p(Nome, Yupik, Denaina, Tutchone, Tanana, Paulatuk, ArcticOcean),
}

func p(provinces ...string) []string {
	return provinces
}
