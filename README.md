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

## Installation

### MinIO

The core uses [MinIO](https://min.io/) to store PST files and attachments.

```bash
# Make MinIO executable
$ chmod +x ./minio

$ MINIO_ROOT_USER=admin MINIO_ROOT_PASSWORD=yourrootpassword ./minio server data/ --console-address ":9001"

# The console can be accessed via the browser at http://127.0.0.1:9001
# Create a new user via the console with an access key, secret key and the "readwrite" permission.
# The access key and secret key are required when starting the Go Forensics API (in the configuration file).
```

### Tus

The dashboard uploads files to [Tus](https://github.com/tus/tusd) (resumable file uploads) which Tus uploads to MinIO.

```bash
# These environment variables are from setting up MinIO.
export AWS_ACCESS_KEY_ID=yourMinIOaccesskey
export AWS_SECRET_ACCESS_KEY=yourMinIOsecretkey
export AWS_REGION=eu-west-1

# Replace the bucket name with the one you created in MinIO.
$ ./tusd -s3-endpoint http://127.0.0.1:9000 -s3-bucket BUCKET_NAME
```


### Ory Kratos

[Ory Kratos](https://www.ory.sh/kratos/) is used for identity management (authentication).

```bash
# Edit kratos.yml:
# - Your SMTP provider (we use Postmark)
# - Path to outlook-mapper and user-identity-schema
#
# Start Kratos 
$ kratos serve -c kratos.yml --watch-courier
```

### Go Forensics API

**Required configuration** before starting the API can be found in [goforensics.yml](https://github.com/mooijtech/goforensics-api/blob/main/goforensics.yml).

```bash
# Change directory
$ cd ~/path/to/goforensics-api

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
