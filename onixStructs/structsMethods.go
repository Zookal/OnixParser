/*
	Copyright (C) 2014  Cyrill AT Schumacher dot fm

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.

    Contribute @ https://github.com/SchumacherFM/OnixParser
*/
package onixStructs

import (
	"bytes"
	"sort"
	"strconv"
	// "strings"
)

var (
	internalByteCounter     = make(map[string]int)
	currentWriteToTableName = make(map[string]string)
)

func getCurrentWriteToTableName(tableName string) string {
	tn, isSet := currentWriteToTableName[tableName]
	if false == isSet {
		return tableName
	}
	return tn
}

// get around of mysql max allowed packet which is hardcoded in the mysql driver at 8MB :-(
func countByte(tableName string, bytes int) {
	counted, isSet := internalByteCounter[tableName]
	if false == isSet {
		internalByteCounter[tableName] = 0
	}
	internalByteCounter[tableName] = bytes + counted
}

func moreThanMySqlMaxAllowedPacket(tableName string) bool {
	if internalByteCounter[tableName] > appConfig.MaxPacketSize {
		internalByteCounter[tableName] = 0
		return true
	}
	return false
}

func writeOneElementToFile(anyStruct interface{}, args map[int]string) (int, error) {
	tableName := appConfig.GetNameOfStruct(anyStruct)
	mapLen := len(args) - 1
	var buffer bytes.Buffer

	// important to keep the correct order of the map
	keys := []int{}
	for k := range args {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		buffer.WriteByte(appConfig.Csv.Enclosure)
		buffer.WriteString(args[k])
		buffer.WriteByte(appConfig.Csv.Enclosure)
		if k < mapLen {
			buffer.WriteByte(appConfig.Csv.Delimiter)
		}
	}
	buffer.WriteByte(appConfig.Csv.LineEnding)
	countByte(tableName, buffer.Len())

	if true == moreThanMySqlMaxAllowedPacket(tableName) {
		// create new file
		nextTableName := appConfig.GetNextTableName(tableName)
		currentWriteToTableName[tableName] = nextTableName
	}
	writeToTn := getCurrentWriteToTableName(tableName)
	return appConfig.WriteBytes(writeToTn, buffer.Bytes())
}

func (p *Product) Xml2CsvRoot() {
	_, writeErr := writeOneElementToFile(p, map[int]string{
		0:  p.RecordReference,
		1:  p.RecordReference,
		2:  strconv.Itoa(p.NotificationType),
		3:  p.ProductForm,
		4:  p.ProductFormDetail,
		5:  strconv.Itoa(p.EpubTechnicalProtection),
		6:  strconv.Itoa(p.TitleType),
		7:  strconv.Itoa(p.TitleElementLevel),
		8:  p.TitleText,
		9:  p.TitlePrefix,
		10: p.TitleWithoutPrefix,
		11: p.Subtitle,
		12: p.ImprintName,
		13: p.PublisherName,
		14: strconv.Itoa(p.PublishingStatus),
		15: p.PublishingDate,
		// 7:  strconv.Itoa(p.TitleElementLevel),
	})
	appConfig.HandleErr(writeErr)
}

func (p *ProductIdentifier) Xml2Csv(id string) {
	if p.ProductIDType > 0 {
		_, writeErr := writeOneElementToFile(p, map[int]string{
			0: id,
			1: strconv.Itoa(p.ProductIDType),
			2: p.IDTypeName,
			3: p.IDValue,
		})
		appConfig.HandleErr(writeErr)
	}
}

func (r *RelatedProduct) Xml2Csv(id string) {
	_, writeErr := writeOneElementToFile(r, map[int]string{
		0: id,
		1: strconv.Itoa(r.ProductRelationCode),
		2: strconv.Itoa(r.ProductIDType),
		3: r.IDValue,
	})
	appConfig.HandleErr(writeErr)
}

func (c *Contributor) Xml2Csv(id string) {
	_, writeErr := writeOneElementToFile(c, map[int]string{
		0: id,
		1: strconv.Itoa(c.SequenceNumber),
		2: c.ContributorRole,
		3: c.PersonName,
		4: c.TitlesBeforeNames,
		5: c.NamesBeforeKey,
		6: c.KeyNames,
	})
	appConfig.HandleErr(writeErr)
}

func (t *TextContent) Xml2Csv(id string) {
	if t.TextType > 0 {
		_, writeErr := writeOneElementToFile(t, map[int]string{
			0: id,
			1: strconv.Itoa(t.TextType),
			2: strconv.Itoa(t.ContentAudience),
			3: t.Text,
			4: t.SourceTitle,
		})
		appConfig.HandleErr(writeErr)
	}
}

func (s *SupportingResource) Xml2Csv(id string) {
	if s.ResourceMode == 3 { //still image only
		_, writeErr := writeOneElementToFile(s, map[int]string{
			0: id,
			1: strconv.Itoa(s.ResourceContentType),
			2: strconv.Itoa(s.ContentAudience),
			3: strconv.Itoa(s.ResourceMode),
		})
		appConfig.HandleErr(writeErr)
		if len(s.ResourceVersion) > 0 {
			for _, sResourceVersion := range s.ResourceVersion {
				sResourceVersion.Xml2Csv(id)
			}
		}
	}
}

func (r *ResourceVersion) Xml2Csv(id string) {
	var FeatureValue string
	if len(r.ResourceVersionFeature) > 0 {
		for _, rResourceVersionFeature := range r.ResourceVersionFeature {
			if rResourceVersionFeature.ResourceVersionFeatureType == 2 { // 2 mean image height
				FeatureValue = rResourceVersionFeature.FeatureValue
			}
		}
	}
	_, writeErr := writeOneElementToFile(r, map[int]string{
		0: id,
		1: r.ResourceLink,
		2: FeatureValue,
	})
	appConfig.HandleErr(writeErr)
}

func (s *Subject) Xml2Csv(id string) {
	_, writeErr := writeOneElementToFile(s, map[int]string{
		0: id,
		1: strconv.Itoa(s.SubjectSchemeIdentifier),
		2: s.SubjectCode,
		3: s.SubjectHeadingText,
	})
	appConfig.HandleErr(writeErr)
}

func (s *SupplyDetail) Xml2Csv(id string) {
	_, writeErr := writeOneElementToFile(s, map[int]string{
		0: id,
		1: s.SupplierName,
		2: strconv.Itoa(s.ProductAvailability),
	})
	appConfig.HandleErr(writeErr)
	if len(s.Price) > 0 {
		for _, sPrice := range s.Price {
			sPrice.Xml2Csv(id, s.SupplierName)
		}
	}
}

func (p *Price) Xml2Csv(id string, supplierName string) {
	_, writeErr := writeOneElementToFile(p, map[int]string{
		0: id,
		1: supplierName,
		2: strconv.Itoa(p.PriceType),
		3: strconv.Itoa(p.PriceQualifier),
		4: p.PriceAmount,
		5: p.CurrencyCode,
		6: strconv.Itoa(p.PriceConditionType),
	})
	appConfig.HandleErr(writeErr)
}
