#!/usr/bin/env bash

#cd s3upload
#go build -v -o s3uploadBuild s3upload.go
#zip s3upload.zip s3uploadBuild
#rm s3uploadBuild
#aws lambda update-function-code \
#  --region eu-west-1 \
#  --function-name s3upload \
#  --zip-file fileb://s3upload.zip
#cd ..


cd dynamoUpload
go build -v -o dynamoUploadBuild dynamoUpload.go
zip dynamoUpload.zip dynamoUploadBuild
rm dynamoUploadBuild
aws lambda update-function-code \
  --region eu-west-1 \
  --function-name dynamoUpload \
  --zip-file fileb://dynamoUpload.zip
cd ..