package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	session "github.com/aws/aws-sdk-go/aws/session"
	lib "github.com/warrensbox/s3upload/lib"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

/*
* Version 0.1.0
* Compatible with Mac OS X and linux OS ONLY
 */

/*** OPERATION WORKFLOW ***/
/*
* 1- SSM library gets AWS credentials from host machine
* 2- Parses command line arguments
* 3- Checks if default s3config file exist
* 4- Establishes S3 connection
* 5- Push files to S3
 */

var version = "0.1.0\n"

var (
	versionFlag  *bool
	helpFlag     *bool
	includeBase  *bool
	awsRegion    *string
	directory    *string
	bucket       *string
	addkey       *string
	configFile   *string
	excludeFiles *string
)

func init() {

	const (
		versionFlagDesc = "Displays the version of s3-pusher"
		skipBaseDesc    = "Skip base directory"
		awsRegionDesc   = "Provide AWS Region. Default is us-east-1"
		directoryDesc   = "Directory where files are stored. Default is current directory"
		bucketDesc      = "S3 bucket. Defaults are in config file"
		keyDesc         = "Append key to s3 bucket. For example: key/my.files"
		confDesc        = "S3 config info"
	)

	versionFlag = kingpin.Flag("version", versionFlagDesc).Short('v').Bool()
	includeBase = kingpin.Flag("ignoreBase", skipBaseDesc).Short('i').Bool()
	awsRegion = kingpin.Flag("region", awsRegionDesc).Short('r').String()
	directory = kingpin.Flag("directory", directoryDesc).Short('d').String()
	bucket = kingpin.Flag("bucket", bucketDesc).Short('b').String()
	addkey = kingpin.Flag("addkey", keyDesc).Short('k').String()
	configFile = kingpin.Flag("config", confDesc).Short('c').String()
	excludeFiles = kingpin.Flag("exclude", confDesc).Short('e').String()
}

func main() {

	kingpin.CommandLine.Interspersed(false)
	kingpin.Parse()

	config := &aws.Config{Region: aws.String(*awsRegion)}

	session := session.Must(session.NewSession(config))

	construct := &lib.Constructor{*directory, *bucket, *addkey, *includeBase, *configFile, *excludeFiles, session}
	profile := lib.NewConstructor(construct)

	err := profile.PushToS3()

	if err != nil {
		fmt.Println(err)
	}

}
