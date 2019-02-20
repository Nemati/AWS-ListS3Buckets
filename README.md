# AWS-ListS3Buckets
 This project is a CLI based AWS S3 storage analysis tool. It uses [Amazon Go SDK](https://aws.amazon.com/sdk-for-go/). The SDK makes it easy to integrate the Go application with the Amazon S3 service.  
 ### How to Install
 We assume that you have an Amazon account and the user credential is setup on your machine. You could follow the steps written in [Configuring the AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html) to setup your account. 
 
 You need to install the following packages as well as golang on your machine:
 
 `go get -u github.com/aws/aws-sdk-go`
 
 `go get -u	github.com/fatih/color`
 
 ### How to Use
 Our tool supports different flags. 
 
 - ls: List Available Buckets and their objects
 - cost: Output cost based on Amazon cost Explorer
 - region: Set Region for the output of results (defualt: us-east-2")
 - costST: Start time for query cost and usage (Default: "2019-02-15")
 - costET: End time for query cost and usage (Default: "2019-02-17")
 - granularity: Set the granularity for query cost and usage (DAILY | MONTHLY | HOURLY)")
 - unit: set the unit for output total Size of objects in a bucket (default:GB)
 - groupBy: Group the output based on Buckets|Region|Storage
 - filter: Filter the output (e.g.: s3://mybucket/Folder/SubFolder/log*)
 - help: Help Function
 - version: Output the version of ListS3Buckets
 ### Example
 `go run main.go -ls -groupBy Region -filter  s3://hani-first/test01/subTest01/*.txt -cost`
 
 This command filters the output for hani-first bucket. It checks for all files with txt extension in test01/subTest01/ folder. It also shows the cost of the account. 
 
<a href="https://ibb.co/8x9b2kv"><img src="https://i.ibb.co/n1jBfpd/example1.png" alt="example1" border="0"></a> 

`go run main.go -ls -groupBy Storage`
 This command all available buckets and group them based on their storage class.

<a href="https://ibb.co/PxQN5Fh"><img src="https://i.ibb.co/nrCbD6c/example2.png" alt="example2" border="0"></a>
