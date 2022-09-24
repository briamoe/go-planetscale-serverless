module github.com/briamoe/go-planetscale-serverless/examples/aws

go 1.18

require github.com/aws/aws-lambda-go v1.34.1 // indirect

require (
  github.com/briamoe/go-planetscale-serverless v0.1.0
)

replace github.com/briamoe/go-planetscale-serverless => "../../"