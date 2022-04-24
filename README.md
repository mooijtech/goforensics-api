<h1 align="center">
  <br>
  <a href="https://github.com/mooijtech/goforensics-api"><img src="https://i.imgur.com/kd7fwOf.png" alt="Go Forensics API" width="180"></a>
  <br>
  Go Forensics API
  <br>
</h1>

<h4 align="center">Open source forensic software to analyze digital evidence to be presented in court.</h4>

---

API which communicates with [Go Forensics Core](https://github.com/mooijtech/goforensics-core).

### Installation

MinIO and Tus must be in the same directory as the Go Forensics API.

```bash
# Download or clone the Go Forensics API
$ git clone https://github.com/mooijtech/goforensics-api

# Change directory
$ cd goforensics-api
```

### MinIO

The core uses [MinIO](https://min.io/) to store PST files and attachments.
Move the MinIO executable to the goforensics-api directory.

```bash
# Change directory
$ cd ~/path/to/goforensics-api

# Make MinIO executable
$ chmod +x ./minio

$ MINIO_ROOT_USER=admin MINIO_ROOT_PASSWORD=yourrootpassword ./minio server data/ --console-address ":9001"

# The console can be accessed via the browser at http://127.0.0.1:9001
# Create a new user via the console with an access key, secret key and the "readwrite" permission.
# The access key and secret key (environment variables) are required when starting the Go Forensics API.
```

### Tus

The dashboard uploads files to [Tus](https://github.com/tus/tusd) (resumable file uploads) which Tus uploads to MinIO.
Move the Tus executable to the goforensics-api.

```bash
# Change directory
$ cd ~/path/to/goforensics-api

# These environment variables are from setting up MinIO (in the Core).
export AWS_ACCESS_KEY_ID=yourMinIOaccesskey
export AWS_SECRET_ACCESS_KEY=yourMinIOsecretkey
export AWS_REGION=eu-west-1

# Replace the bucket name with the one you created in MinIO.
$ ./tusd -s3-endpoint http://127.0.0.1:9000 -s3-bucket BUCKET_NAME
```


### Ory Kratos

Ory Kratos is used for identity management (authentication).

```bash
$ bash <(curl https://raw.githubusercontent.com/ory/meta/master/install.sh) -d -b . kratos v0.9.0-alpha.3
$ sudo mv ./kratos /usr/local/bin/

# Edit kratos.yml to your SMTP provider (we use Postmark) and path to the outlook-mapper, user-identity-schema.
# Start Kratos 
$ kratos serve -c kratos-development.yml --watch-courier
```

### Go Forensics API

```bash
# cd ~/path/to/goforensics-api

# Export required environment variables.
$ export MINIO_BUCKET=yourMinIObucket
$ export MINIO_ENDPOINT=127.0.0.1:9000
$ export MINIO_ACCESS_KEY=yourMinIOaccesskey
$ export MINIO_SECRET_KEY=yourMinIOsecretkey
$ export MINIO_SECURE=false
$ export OUTLOOK_CLIENT_ID=yourOutlookClientID
$ export OUTLOOK_CLIENT_SECRET=yourOutlookClientSecret

# Start the API
$ go run cmd/api.go
```

### Libraries

- [goforensics-core](https://github.com/mooijtech/goforensics-core)
- [logrus](https://github.com/sirupsen/logrus)
- [Mux](https://github.com/gorilla/mux)
- [SCS](https://github.com/alexedwards/scs)
- [CORS](https://github.com/rs/cors)
- [SSE](https://github.com/r3labs/sse)
