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

/*
	not all types are matched
	http://www.editeur.org/onix/2.1/02/reference/onix-international.dtd
	http://www.editeur.org/onix/2.1/03/reference/onix-international.dtd
*/
type (
	ProductIdentifier struct {
		ProductIDType int    `xml:"ProductIDType" sql:"bigint(14)"`
		IDTypeName    string `xml:"IDTypeName" sql:"varchar(255) NULL"`
		IDValue       string `xml:"IDValue" sql:"varchar(255) NULL"`
	}

	// // P.7 Authorship
	Contributor struct {
		SequenceNumber    int    `xml:"SequenceNumber" sql:"int(10) NULL"`
		ContributorRole   string `xml:"ContributorRole" sql:"varchar(20) NOT NULL"`
		PersonName        string `xml:"PersonName" sql:"varchar(255) NULL"`
		TitlesBeforeNames string `xml:"TitlesBeforeNames" sql:"varchar(100) NULL"`
		NamesBeforeKey    string `xml:"NamesBeforeKey" sql:"varchar(100) NULL"`
		KeyNames          string `xml:"KeyNames" sql:"varchar(100) NULL"`
		CorporateName     string `xml:"CorporateName" sql:"varchar(255) NULL"`
	}

	// // P.12 Subject
	Subject struct {
		SubjectSchemeIdentifier int    `xml:"SubjectSchemeIdentifier" sql:"int(10) NOT NULL"`
		SubjectCode             string `xml:"SubjectCode" sql:"varchar(100) NULL"`
		SubjectHeadingText      string `xml:"SubjectHeadingText" sql:"varchar(255) NULL"`
	}

	TextContent struct {
		TextType        int    `xml:"TextType" sql:"int(10) NOT NULL"`
		ContentAudience int    `xml:"ContentAudience" sql:"int(10) NOT NULL"`
		Text            string `xml:"Text" sql:"TEXT NOT NULL"`
		SourceTitle     string `xml:"SourceTitle" sql:"varchar(255) NULL"`
	}

	SupportingResource struct {
		ResourceContentType int `xml:"ResourceContentType" sql:"int(10) NOT NULL"`
		ContentAudience     int `xml:"ContentAudience" sql:"int(10) NOT NULL"`
		ResourceMode        int `xml:"ResourceMode" sql:"int(10) NOT NULL"`
		ResourceVersion     []ResourceVersion
	}

	//todo parse more for valid images
	ResourceVersion struct {
		ResourceLink           string `xml:"ResourceLink" sql:"varchar(255) NOT NULL"`
		ImageHeight            string `sql:"varchar(255) NOT NULL"`
		ResourceVersionFeature []ResourceVersionFeature
	}

	ResourceVersionFeature struct {
		ResourceVersionFeatureType int    `xml:"ResourceVersionFeatureType" sql:"int(10) NOT NULL"`
		FeatureValue               string `xml:"FeatureValue" sql:"varchar(255) NOT NULL"`
	}

	RelatedProduct struct {
		ProductRelationCode int    `xml:"ProductRelationCode" sql:"int(10) NOT NULL"`
		ProductIDType       int    `xml:"ProductIdentifier>ProductIDType" sql:"int(10) NOT NULL"`
		IDValue             string `xml:"ProductIdentifier>IDValue" sql:"bigint(15) NOT NULL"`
	}

	SupplyDetail struct {
		SupplierName        string `xml:"Supplier>SupplierName" sql:"varchar(255) NOT NULL"`
		ProductAvailability int    `xml:"ProductAvailability" sql:"int(10) NOT NULL DEFAULT 0"`
		Price               []Price
	}

	Price struct {
		SupplierName       string `sql:"varchar(255) NOT NULL"`
		PriceType          int    `xml:"PriceType" sql:"int(10) NOT NULL DEFAULT 0"`
		PriceQualifier     int    `xml:"PriceQualifier" sql:"int(10) NOT NULL DEFAULT 0"`
		PriceAmount        string `xml:"PriceAmount" sql:"decimal(10,2) NOT NULL DEFAULT 0"`
		CurrencyCode       string `xml:"CurrencyCode" sql:"varchar(10) NULL"`
		PriceConditionType int    `xml:"PriceCondition>PriceConditionType" sql:"int(10) NOT NULL DEFAULT 0"`
		Quantity					 int    `xml:"PriceCondition>PriceConditionQuantity>Quantity" sql:"int(5) NULL"` 
	}

	EpubUsageConstraint struct {
		EpubUsageType int `xml:"EpubUsageType" sql:"int(10) NOT NULL DEFAULT 0"`
		EpubUsageStatus int `xml:"EpubUsageStatus" sql:"int(10) NOT NULL DEFAULT 0"`
		Quantity int `xml:"EpubUsageLimit>Quantity" sql:"int(10) NOT NULL DEFAULT 0"`
		EpubUsageUnit int `xml:"EpubUsageLimit>EpubUsageUnit" sql:"int(10) NOT NULL DEFAULT 0"`
	}

	Product struct {
		RecordReference         string `xml:"RecordReference" sql:"varchar(255) NOT NULL"`
		NotificationType        int    `xml:"NotificationType" sql:"int(10) NOT NULL DEFAULT 0"`
		ProductIdentifier       []ProductIdentifier
		ProductForm             string             `xml:"DescriptiveDetail>ProductForm" sql:"varchar(20) NULL"`
		ProductFormDetail       string             `xml:"DescriptiveDetail>ProductFormDetail" sql:"varchar(20) NULL"`
		EpubTechnicalProtection int                `xml:"DescriptiveDetail>EpubTechnicalProtection" sql:"int(10) NULL"`
		TitleType               int                `xml:"DescriptiveDetail>TitleDetail>TitleType" sql:"int(10) NOT NULL"`
		TitleElementLevel       int                `xml:"DescriptiveDetail>TitleDetail>TitleElement>TitleElementLevel" sql:"int(10) NOT NULL"`
		TitleText               string             `xml:"DescriptiveDetail>TitleDetail>TitleElement>TitleText" sql:"varchar(255) NULL"`
		TitlePrefix             string             `xml:"DescriptiveDetail>TitleDetail>TitleElement>TitlePrefix" sql:"varchar(255) NULL"`
		TitleWithoutPrefix      string             `xml:"DescriptiveDetail>TitleDetail>TitleElement>TitleWithoutPrefix" sql:"varchar(255) NULL"`
		Subtitle                string             `xml:"DescriptiveDetail>TitleDetail>TitleElement>Subtitle" sql:"varchar(255) NULL"`
		Contributor             []Contributor      `xml:"DescriptiveDetail>Contributor"`
		Subject                 []Subject          `xml:"DescriptiveDetail>Subject"`
		
		TextContent             []TextContent      `xml:"CollateralDetail>TextContent"`
		SupportingResource      SupportingResource `xml:"CollateralDetail>SupportingResource"`
		// DescriptiveDetail
		// CollateralDetail
		ImprintName      				string           `xml:"PublishingDetail>Imprint>ImprintName" sql:"varchar(255) NULL"`
		PublisherName    				string           `xml:"PublishingDetail>Publisher>PublisherName" sql:"varchar(255) NULL"`
		PublishingStatus 				int              `xml:"PublishingDetail>PublishingStatus" sql:"int(10) NOT NULL DEFAULT 0"`
		PublishingDate   				string           `xml:"PublishingDetail>PublishingDate>Date" sql:"varchar(255) NULL"`
		RelatedProduct   				[]RelatedProduct `xml:"RelatedMaterial>RelatedProduct"`
		SupplyDetail     				[]SupplyDetail   `xml:"ProductSupply>SupplyDetail"`
		EditionNumber						int							 `xml:"DescriptiveDetail>EditionNumber" sql:"int(4) NULL"`
		EpubUsageConstraint			[]EpubUsageConstraint `xml:"DescriptiveDetail>EpubUsageConstraint"`
	}
)
