package models

import "moviemicroservice.com/src/gen"

func MetadataToProto(m *MetaData) *gen.Metadata {
	return &gen.Metadata{
		Id:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Director:    m.Director,
	}
}

func MetadataFromProto(p *gen.Metadata) *MetaData {
	return &MetaData{
		ID:          p.Id,
		Title:       p.Title,
		Description: p.Description,
		Director:    p.Director,
	}
}
