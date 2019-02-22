package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	strr "strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fatih/color"
)

var options struct {
	costStartTime   string
	costEndTime     string
	costGranularity string
	groupBy         string
	version         bool
	help            bool
	listBucket      bool
	cost            bool
	totalSizeIn     string
	region          string
	unit            string
	filter          string
}

type bucketInfo struct {
	name                  string
	creationDate          time.Time
	lastModified          time.Time
	numberOfFiles         int64
	numberOfFilesStandard int64
	numberOfFilesIA       int64
	numberOfFilesRR       int64
	totalSizeInStandard   int64
	totalSizeInIA         int64
	totalSizeInRR         int64
	totalSize             int64
	region                string
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func main() {
	var Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s[post_count]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&options.costStartTime, "costST", "2019-02-15", "Start time for query cost and usage (Default: \"2019-02-15\"")
	flag.StringVar(&options.costEndTime, "costET", "2019-02-17", "End time for query cost and usage (Default: \"2019-02-17\"")
	flag.StringVar(&options.costGranularity, "granularity", "DAILY", "Granularity for query cost and usage (DAILY | MONTHLY | HOURLY)")
	flag.StringVar(&options.region, "region", "us-east-2", "Set Region for the output of results (defualt: us-east-2")
	flag.BoolVar(&options.cost, "cost", false, "Output costs and usage")
	flag.StringVar(&options.unit, "unit", "MB", "Unit for output total Size of objects in a bucket (default:GB) ")
	flag.BoolVar(&options.listBucket, "ls", false, "List Available Buckets")
	flag.BoolVar(&options.version, "version", false, "View the version of this application")
	flag.BoolVar(&options.help, "help", false, "Help Function")
	flag.StringVar(&options.groupBy, "groupBy", "Buckets", "Buckets|Region|Storage ")
	flag.StringVar(&options.filter, "filter", "s3://", " Filter the output (e.g.: s3://mybucket/Folder/SubFolder/log*)")

	flag.Parse()
	if options.version {
		fmt.Printf("%s", "goCmd installed version 0.0.1\n")
	}

	if options.help {
		Usage()
	}
	if options.listBucket {

		proccessedBucketInfoRegion := lsBuckets()
		showResult(proccessedBucketInfoRegion)

	}
	if options.cost {
		amount, unit := outputCost(options.costStartTime, options.costEndTime, options.costGranularity)
		red := color.New(color.FgRed)
		green := color.New(color.FgGreen)
		boldRed := red.Add(color.Bold)
		fmt.Println("--------------------------------------------")
		boldRed.Printf("Total Cost: ")
		green.Println(amount, " ", unit, "\n")
	}
}

// Check the unit option to print in TB, GB, MB, KB, or Byte
// Input is the option flag for unit (Default is GB)
// Output is (DividedBy, unitSize)
// DividedBy is an float64 that is the convert rate like 1024 to conver Byte to KB
// unitSize is an string shows the unit (GB,KB)
func getUnit(unit string) (float64, string) {
	// Initialize the default variables
	devideBy := 1.0
	unitSize := "Byte"

	switch unit {
	case "TB":
		devideBy = math.Pow(2.0, 40)
		unitSize = "TB"
	case "GB":
		devideBy = math.Pow(2.0, 30)
		unitSize = "GB"
	case "MB":
		devideBy = math.Pow(2.0, 20)
		unitSize = "MB"
	case "KB":
		devideBy = math.Pow(2.0, 10)
		unitSize = "KB"
	}

	return devideBy, unitSize
}

// Output the bucket infor
func showResult(proccessedBucketInfoRegion map[string][]bucketInfo) {

	// Iterate over gathered information from Buckets in different location
	if options.groupBy == "Storage" {
		devideBy, unitSize := getUnit(options.unit)
		red := color.New(color.FgRed)
		magenta := color.New(color.FgHiMagenta)
		boldRed := red.Add(color.Bold)
		yello := color.New(color.FgHiYellow)
		fmt.Println("##############################################")
		yello.Println("            Storage:  Standard")
		fmt.Println("##############################################")
		for _, valueBucketInfoList := range proccessedBucketInfoRegion {

			for _, elementBucketInformation := range valueBucketInfoList {
				boldRed.Print("Bucket Name: ")
				color.Cyan(elementBucketInformation.name)

				magenta.Print("Creation: ")
				fmt.Print(elementBucketInformation.creationDate)

				magenta.Print("    Last Modified: ")
				fmt.Println(elementBucketInformation.lastModified)

				magenta.Print("Number of files Standard: ")
				fmt.Print(elementBucketInformation.numberOfFilesStandard)

				magenta.Print("    Total Size Standard: ")
				fmt.Print(float64(elementBucketInformation.totalSizeInStandard)/devideBy, " ", unitSize)

				magenta.Print("     Region: ")
				fmt.Println(elementBucketInformation.region)

				fmt.Println("--------------------------------------------")

			}
		}
		fmt.Println("##############################################")
		yello.Println("            Storage:  IA")
		fmt.Println("##############################################")
		for _, valueBucketInfoList := range proccessedBucketInfoRegion {

			for _, elementBucketInformation := range valueBucketInfoList {

				boldRed.Print("Bucket Name: ")
				color.Cyan(elementBucketInformation.name)

				magenta.Print("Creation: ")
				fmt.Print(elementBucketInformation.creationDate)

				magenta.Print("    Last Modified: ")
				fmt.Println(elementBucketInformation.lastModified)

				magenta.Print("Number of files IA: ")
				fmt.Print(elementBucketInformation.numberOfFilesIA)

				magenta.Print("    Total Size IA: ")
				fmt.Print(float64(elementBucketInformation.totalSizeInIA)/devideBy, " ", unitSize)

				magenta.Print("     Region: ")
				fmt.Println(elementBucketInformation.region)

				fmt.Println("--------------------------------------------")

			}
		}
		fmt.Println("##############################################")
		yello.Println("            Storage:  RR")
		fmt.Println("##############################################")
		for _, valueBucketInfoList := range proccessedBucketInfoRegion {

			for _, elementBucketInformation := range valueBucketInfoList {
				boldRed.Print("Bucket Name: ")
				color.Cyan(elementBucketInformation.name)

				magenta.Print("Creation: ")
				fmt.Print(elementBucketInformation.creationDate)

				magenta.Print("    Last Modified: ")
				fmt.Println(elementBucketInformation.lastModified)

				magenta.Print("Number of files RR: ")
				fmt.Print(elementBucketInformation.numberOfFilesRR)

				magenta.Print("    Total Size RR: ")
				fmt.Print(float64(elementBucketInformation.totalSizeInRR)/devideBy, " ", unitSize)

				magenta.Print("     Region: ")
				fmt.Println(elementBucketInformation.region)

				fmt.Println("--------------------------------------------")
			}
		}
		return
	}
	for keyRegion, valueBucketInfoList := range proccessedBucketInfoRegion {

		// Check if the output should be grouped by Region
		if options.groupBy == "Region" {
			yello := color.New(color.FgHiYellow)
			fmt.Println("##############################################")
			yello.Println("            Region: ", keyRegion)
			fmt.Println("##############################################")
		}
		red := color.New(color.FgRed)
		magenta := color.New(color.FgHiMagenta)
		boldRed := red.Add(color.Bold)
		// Iterate over the Bucket List Information in a specific Region
		for _, elementBucketInformation := range valueBucketInfoList {
			// Define some color for the output

			boldRed.Print("Bucket Name: ")
			color.Cyan(elementBucketInformation.name)

			magenta.Print("Creation: ")
			fmt.Print(elementBucketInformation.creationDate)

			magenta.Print("    Last Modified: ")
			fmt.Println(elementBucketInformation.lastModified)

			magenta.Print("Number of Files: ")
			fmt.Print(elementBucketInformation.numberOfFiles)

			magenta.Print("    Total Size: ")
			devideBy, unitSize := getUnit(options.unit)
			fmt.Print(float64(elementBucketInformation.totalSize)/devideBy, " ", unitSize)

			magenta.Print("     Region: ")
			fmt.Println(elementBucketInformation.region)

			fmt.Println("--------------------------------------------")
		}
	}
}

func parsFilter(filter string) (string, string, string) {
	//filter := "s3://mybucket/Folder/SubFolder/log*"

	var fileName string = ""
	var subfolderName string = ""
	var bucketName string = ""
	//filter := "s3://"
	sizeFilter := len(filter)
	numberOfSlash := strr.Count(filter, "/")
	indexBucketName := strr.Index(filter, "//") + 2

	if numberOfSlash > 2 {
		indexEndBucketName := strr.Index(filter[indexBucketName:], "/")
		bucketName = filter[indexBucketName : indexEndBucketName+5]
		if !strr.HasSuffix(filter, "/") {
			indexFileName := strr.LastIndex(filter, "/") + 1
			fileName = filter[indexFileName:]
		}
		// It has subfolder
		if numberOfSlash > 3 {
			indexEndBucketName := strr.Index(filter[indexBucketName:], "/")

			subfolderName = filter[indexEndBucketName+6 : sizeFilter-len(fileName)]
		}
	}
	// len s3:// is 5 so we added 5
	return bucketName, subfolderName, fileName

}

// Output is a hashmap containing a list of bucket info in each region
// It uses the default existing user for Amazon S3 account
func lsBuckets() map[string][]bucketInfo {

	// Define the output of lsBuckets function
	proccessedBucketInfoRegion := make(map[string][]bucketInfo)

	// Create an AWS session
	sessionAWS, errSessionAWS := session.NewSession(&aws.Config{
		Region: aws.String("ca-central-1"),
	},
	)
	if errSessionAWS != nil {
		exitErrorf("Error in initializing the account, %v", errSessionAWS)
	}

	// A S3 inctance based on our account
	s3Connected := s3.New(sessionAWS)

	// List the bucket name in our account
	resultListBuckets, errListBuckets := s3Connected.ListBuckets(nil)

	// Check if the ListBuckets face any error
	if errListBuckets != nil {
		exitErrorf("Unable to list buckets, %v", errListBuckets)
	}

	bucketName, subfolderName, fileName := parsFilter(options.filter)
	indexStar := strr.Index(fileName, "*")

	// Iterate over the list of Buckets to retrieve information about their existing files
	for _, bucketStruct := range resultListBuckets.Buckets {

		if len(bucketName) > 0 && bucketName != *bucketStruct.Name {
			continue
		}
		// trunc shows there is more than 1000 objects in a bucket
		// It is initialized as true then it will be change by the IsTruncated variable
		trunc := true
		var maxObjects = 999
		var markerString = ""
		var totalSizeInStandard int64 = 0
		var totalSizeInIA int64 = 0
		var totalSizeInRR int64 = 0
		var totalSize int64 = 0

		numberOfFiles := 0
		numberOfFilesIA := 0
		numberOfFilesStandard := 0
		numberOfFilesRR := 0
		lastModifiedTime := *bucketStruct.CreationDate

		// Get the location of bucket in order to be able to access it
		inputGetBucketLocation := &s3.GetBucketLocationInput{
			Bucket: aws.String(*bucketStruct.Name),
		}
		bucketLocation, _ := s3Connected.GetBucketLocation(inputGetBucketLocation)

		// Check if list object should send more request (Objects are more than 1000)
		for trunc {

			inputListObject := &s3.ListObjectsInput{
				Bucket: aws.String(*bucketStruct.Name),
				Prefix: aws.String(subfolderName),
				//Delimiter: aws.String("/"),
				Marker: aws.String(markerString),
			}
			sessionWithLocation, err := session.NewSession(&aws.Config{
				Region: aws.String(*bucketLocation.LocationConstraint),
			},
			)
			s3ConnectionWithLocation := s3.New(sessionWithLocation)

			resultListBucketsWithLocation, errListObjects := s3ConnectionWithLocation.ListObjects(inputListObject)
			if errListObjects != nil {
				if aerr, ok := errListObjects.(awserr.Error); ok {
					switch aerr.Code() {
					case s3.ErrCodeNoSuchBucket:
						fmt.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
					default:
						fmt.Println(aerr.Error())
					}
				} else {
					// Print the error, cast err to awserr.Error to get the Code and
					// Message from an error.
					fmt.Println(err.Error())
				}

			}
			// Iterate over objects in the bucket to retrieve information
			for _, objectsInBuckets := range resultListBucketsWithLocation.Contents {
				if len(fileName) > 0 && indexStar > 0 && !strr.HasPrefix(*objectsInBuckets.Key, fileName[:strr.Index(fileName, "*")]) {
					// start with HasPrefix()
					continue
				} else if len(fileName) > 0 && indexStar == 0 && !strr.HasSuffix(*objectsInBuckets.Key, fileName[1:]) {
					// end with HasSuffix()
					continue
				}

				totalSize += *objectsInBuckets.Size

				if (*objectsInBuckets.LastModified).Unix() > lastModifiedTime.Unix() {
					lastModifiedTime = *objectsInBuckets.LastModified
				}
				numberOfFiles++

				// Different type of storage class
				if *objectsInBuckets.StorageClass == "STANDARD" {
					numberOfFilesStandard++
					totalSizeInStandard += *objectsInBuckets.Size
				} else if *objectsInBuckets.StorageClass == "IA" {
					numberOfFilesIA++
					totalSizeInIA += *objectsInBuckets.Size
				} else if *objectsInBuckets.StorageClass == "RR" {
					numberOfFilesRR++
					totalSizeInRR += *objectsInBuckets.Size
				}
			}
			if resultListBucketsWithLocation.IsTruncated == nil {
				break
			}
			// update trunc to see if there is more objects
			trunc = *resultListBucketsWithLocation.IsTruncated
			if trunc {
				markerString = *resultListBucketsWithLocation.Contents[maxObjects].Key
			}
		}

		// Add the new bucket information to hashmap of all buckets information
		mapKeyBucketLocation := *bucketLocation.LocationConstraint
		proccessedBucketInfoRegion[mapKeyBucketLocation] = append(proccessedBucketInfoRegion[mapKeyBucketLocation],
			bucketInfo{*bucketStruct.Name, *bucketStruct.CreationDate, lastModifiedTime, int64(numberOfFiles),
				int64(numberOfFilesStandard), int64(numberOfFilesIA), int64(numberOfFilesRR), totalSizeInStandard,
				totalSizeInIA, totalSizeInRR, totalSize, mapKeyBucketLocation})

	}

	return proccessedBucketInfoRegion
}

func outputCost(costStartTime string, costEndTime string, costGranularity string) (string, string) {
	sessionCost, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")},
	)
	if err != nil {
		exitErrorf("Unable to find, %v", err)
	}
	// Create a CostExplorer client with additional configuration
	var costAndUsage costexplorer.GetCostAndUsageInput
	var dateInt costexplorer.DateInterval
	dateInt.SetStart(costStartTime)
	dateInt.SetEnd(costEndTime)
	// Valid Values: DAILY | MONTHLY | HOURLY
	costAndUsage.SetGranularity(costGranularity)
	costAndUsage.SetTimePeriod(&dateInt)
	//Valid values are AmortizedCost, BlendedCost, NetAmortizedCost,
	//NetUnblendedCost, NormalizedUsageAmount, UnblendedCost, and UsageQuantity
	var requestedMetrics []string = []string{"AmortizedCost", "UsageQuantity"}

	costAndUsage.SetMetrics(aws.StringSlice(requestedMetrics))
	costExplorerConnection := costexplorer.New(sessionCost, aws.NewConfig().WithRegion("us-west-2"))
	// Example sending a request using the GetCostAndUsageRequest method.
	getCostAndUsage, resp := costExplorerConnection.GetCostAndUsageRequest(&costAndUsage)

	errQueryCostAndUsage := getCostAndUsage.Send()
	if errQueryCostAndUsage == nil { // resp is now filled

		costTotal := resp.ResultsByTime[0].Total["AmortizedCost"]

		return *costTotal.Amount, *costTotal.Unit
		//fmt.Println("Total Cost: ", *costTotal.Amount, " ", *costTotal.Unit)
	} else {
		fmt.Println(errQueryCostAndUsage)
	}
	return "", ""
}
