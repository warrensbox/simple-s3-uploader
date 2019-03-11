package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/viper"
)

//Constructor : struct for session
type Constructor struct {
	Session     *session.Session
	Directory   string
	Bucket      string
	IncludeBase bool
	ConfigFile  string
	Exclude     string
}

//NewConstructor :validate object
func NewConstructor(attr *Constructor) (*Constructor, error) {

	if attr.ConfigFile != "" {
		exists, _ := exists(attr.ConfigFile)
		if exists {
			fmt.Println("Reading from custom s3config file")
			dir, basename := filepath.Split(attr.ConfigFile)
			filename := strings.TrimSuffix(basename, filepath.Ext(basename))
			attr = configuration(attr, filename, dir)

		} else {
			fmt.Println("Cannot find config file")
		}
	} else {
		exists, _ := exists("./s3config.json")
		if exists {
			fmt.Println("Reading from local s3config file")
			attr = configuration(attr, "s3config", "./")
		}
	}

	return attr, nil
}

func configuration(attr *Constructor, filename string, dirpath string) *Constructor {

	viper.SetConfigType("json")
	viper.SetConfigName(filename)
	viper.AddConfigPath(dirpath)

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		fmt.Println(err)
		//panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	if attr.Exclude == "" {
		exclude := viper.Get("exclude")
		if exclude != nil {
			attr.Exclude = exclude.(string)
		}
	}

	if attr.Directory == "" {
		directory := viper.Get("source")
		if directory != nil {
			attr.Directory = directory.(string)
		} else {
			attr.Directory = "./"
		}
	}

	fmt.Printf("Directory: %s\n", attr.Directory)

	if attr.Bucket == "" {
		bucket := viper.Get("bucket")
		if bucket != nil {
			attr.Bucket = bucket.(string)
		} else {
			fmt.Println("You must provide a S3 bucket")
			os.Exit(1)
		}
	}

	accessKey := viper.Get("aws_access_key_id")
	secretAccessKey := viper.Get("aws_secret_access_key")
	region := viper.Get("aws_region")

	if accessKey != nil && secretAccessKey != nil && region != nil {
		sess, err := session.NewSession(&aws.Config{
			Region:      aws.String(region.(string)),
			Credentials: credentials.NewStaticCredentials(accessKey.(string), secretAccessKey.(string), ""),
		})

		if err != nil {
			fmt.Println("Unable to set ssm based on credentional provided")
			os.Exit(1)
		}

		attr.Session = sess
	}

	return attr

}
