package v1

func NewGeoJsonMessageV1(feature1, feature2 Feature) (FeatureCollection, error) {
	return FeatureCollection{Features: []Feature{
		feature1,
		feature2,
	}}, nil
}

func NewGeoJsonMessageV2(features []Feature) (FeatureCollection, error) {
	resultFeatures := make([]Feature, 0, len(feature))
	resultFeatures = append(resultFeatures, features...)
	return FeatureCollection{Features: resultFeatures}, nil
}

func NewGeoJsonMessageV3(features ...Feature) (FeatureCollection, error) {
	return FeatureCollection{
		Type:     featureCollection,
		Features: features}, nil
}
