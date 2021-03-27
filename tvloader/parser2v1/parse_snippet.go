// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v1

import (
	"fmt"
	"strconv"

	"github.com/spdx/tools-golang/spdx"
)

func (parser *tvParser2_1) parsePairFromSnippet2_1(tag string, value string) error {
	switch tag {
	// tag for creating new snippet section
	case "SnippetSPDXID":
		parser.snippet = &spdx.Snippet2_1{}
		eID, err := extractElementID(value)
		if err != nil {
			return err
		}
		// FIXME: how should we handle where not associated with current file?
		if parser.file != nil {
			if parser.file.Snippets == nil {
				parser.file.Snippets = map[spdx.ElementID]*spdx.Snippet2_1{}
			}
			parser.file.Snippets[eID] = parser.snippet
		}
		parser.snippet.SnippetSPDXIdentifier = eID
	// tag for creating new file section and going back to parsing File
	case "FileName":
		parser.st = psFile2_1
		parser.snippet = nil
		return parser.parsePairFromFile2_1(tag, value)
	// tag for creating new package section and going back to parsing Package
	case "PackageName":
		parser.st = psPackage2_1
		//check here whether the last file of the previous package contained the FileSpdxIdentifier
		if parser.file != nil && parser.file.FileSPDXIdentifier == spdx.ElementID("") {
			return fmt.Errorf("Invalid file without a file SPDX identifier")
		}
		parser.file = nil
		parser.snippet = nil
		return parser.parsePairFromPackage2_1(tag, value)
	// tag for going on to other license section
	case "LicenseID":
		parser.st = psOtherLicense2_1
		return parser.parsePairFromOtherLicense2_1(tag, value)
	// tags for snippet data
	case "SnippetFromFileSPDXID":
		deID, err := extractDocElementID(value)
		if err != nil {
			return err
		}
		parser.snippet.SnippetFromFileSPDXIdentifier = deID
	case "SnippetByteRange":
		byteStart, byteEnd, err := extractSubs(value)
		if err != nil {
			return err
		}
		bIntStart, err := strconv.Atoi(byteStart)
		if err != nil {
			return err
		}
		bIntEnd, err := strconv.Atoi(byteEnd)
		if err != nil {
			return err
		}
		parser.snippet.SnippetByteRangeStart = bIntStart
		parser.snippet.SnippetByteRangeEnd = bIntEnd
	case "SnippetLineRange":
		lineStart, lineEnd, err := extractSubs(value)
		if err != nil {
			return err
		}
		lInttStart, err := strconv.Atoi(lineStart)
		if err != nil {
			return err
		}
		lInttEnd, err := strconv.Atoi(lineEnd)
		if err != nil {
			return err
		}
		parser.snippet.SnippetLineRangeStart = lInttStart
		parser.snippet.SnippetLineRangeEnd = lInttEnd
	case "SnippetLicenseConcluded":
		parser.snippet.SnippetLicenseConcluded = value
	case "LicenseInfoInSnippet":
		parser.snippet.LicenseInfoInSnippet = append(parser.snippet.LicenseInfoInSnippet, value)
	case "SnippetLicenseComments":
		parser.snippet.SnippetLicenseComments = value
	case "SnippetCopyrightText":
		parser.snippet.SnippetCopyrightText = value
	case "SnippetComment":
		parser.snippet.SnippetComment = value
	case "SnippetName":
		parser.snippet.SnippetName = value
	// for relationship tags, pass along but don't change state
	case "Relationship":
		parser.rln = &spdx.Relationship2_1{}
		parser.doc.Relationships = append(parser.doc.Relationships, parser.rln)
		return parser.parsePairForRelationship2_1(tag, value)
	case "RelationshipComment":
		return parser.parsePairForRelationship2_1(tag, value)
	// for annotation tags, pass along but don't change state
	case "Annotator":
		parser.ann = &spdx.Annotation2_1{}
		parser.doc.Annotations = append(parser.doc.Annotations, parser.ann)
		return parser.parsePairForAnnotation2_1(tag, value)
	case "AnnotationDate":
		return parser.parsePairForAnnotation2_1(tag, value)
	case "AnnotationType":
		return parser.parsePairForAnnotation2_1(tag, value)
	case "SPDXREF":
		return parser.parsePairForAnnotation2_1(tag, value)
	case "AnnotationComment":
		return parser.parsePairForAnnotation2_1(tag, value)
	// tag for going on to review section (DEPRECATED)
	case "Reviewer":
		parser.st = psReview2_1
		return parser.parsePairFromReview2_1(tag, value)
	default:
		return fmt.Errorf("received unknown tag %v in Snippet section", tag)
	}

	return nil
}
