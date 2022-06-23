# is ec2

If you are running outside EC2 and don't have credentials, the AWS SDK will hang
for 5+ seconds attempting to connect to the EC2 metadata API. That's annoying.

This is a helper library to detect, quickly, if you are running in EC2 or not.
If you are not, then you can avoid including the EC2 metadata API in your
credential list.

There are no dependencies.

## Usage

```go

import "github.com/kevinburke/isec2"

func main() {
    yes, err := isec2.IsEC2(context.Background())
    if err == nil {
        fmt.Println("running on EC2:", yes)
    }
}
```

The main reason to run this is to avoid making endless calls to the EC2 metadata
API, 169.254.169.254, and taking 5+ seconds to time out when you are not running
in EC2.

To avoid using the EC2 metadata API, build a `credentials.Credentials` object as
follows, and then assign it to `aws.Config.Credentials`.

```go
import "github.com/kevinburke/isec2"
import "github.com/aws/aws-sdk-go/aws"
import "github.com/aws/aws-sdk-go/aws/credentials"
import "github.com/aws/aws-sdk-go/aws/defaults"

func main() {
	conf := defaults.Config()
	var creds *credentials.Credentials
	handlers := defaults.Handlers()
	isEC2, err := isec2.IsEC2(ctx)
	if err == nil && !isEC2 {
		// these are all of the defaults from the AWS SDK, except we skip the
		// AWS Metadata service API ("defaults.RemoteCredProvider")
		creds = credentials.NewCredentials(&credentials.ChainProvider{
			VerboseErrors: aws.BoolValue(conf.CredentialsChainVerboseErrors),
			Providers: []credentials.Provider{
				&credentials.EnvProvider{},
				&credentials.SharedCredentialsProvider{Filename: "", Profile: ""},
			},
		})
		conf.Credentials = creds
	}
}
```

## Errata

We will try to return an answer within a maximum of 50 milliseconds.

## Correctness

This library synthesizes [the advice from AWS][aws] and [the advice from
ServerFault][serverfault] to do the best possible EC2 detection.

[aws]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/identify_ec2_instances.html
[serverfault]: https://serverfault.com/a/903599
