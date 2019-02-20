package main

import (
	"math"
	"testing"
)

type GetUnitResult struct {
	unit             string
	expectedDivideBy float64
	expectedUnitBy   string
}

var GetUnitResults = []GetUnitResult{
	{"Byte", 1.0, "Byte"},
	{"KB", math.Pow(2, 10), "KB"},
	{"MB", math.Pow(2, 20), "MB"},
	{"GB", math.Pow(2, 30), "GB"},
	{"TB", math.Pow(2, 40), "TB"},
}

func TestUnit(t *testing.T) {
	for _, test := range GetUnitResults {
		divideBy, unitBy := getUnit(test.unit)
		if divideBy != test.expectedDivideBy || unitBy != test.expectedUnitBy {
			t.Fatal("Expected Result for getUnit not Given")
		}
	}
}

type ParseFilterResult struct {
	filter             string
	expectedBucketName string
	expectedPathName   string
	expectedFileName   string
}

var ParseFilterResults = []ParseFilterResult{
	{"s3://mybucket/Folder/SubFolder/log*", "mybucket", "Folder/SubFolder/", "log*"},
	{"s3://hani-first/test01/subTest01/*.txt", "hani-first", "test01/subTest01/", "*.txt"},
	{"s3://hani-first/test01/*.txt", "hani-first", "test01/", "*.txt"},
	{"s3://hani-first/*.txt", "hani-first", "", "*.txt"},
	{"s3://hani-first/", "hani-first", "", ""},
}

func TestParsFilter(t *testing.T) {

	for _, test := range ParseFilterResults {
		bucketName, pathName, fileName := parsFilter(test.filter)
		if test.expectedBucketName != bucketName || test.expectedPathName != pathName ||
			test.expectedFileName != fileName {
			t.Fatal("Expected Result for getUnit not Given")
		}
	}
}
