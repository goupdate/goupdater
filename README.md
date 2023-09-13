## GoUpdater

GoUpdater is a fortified solution for Golang applications, enabling developers to manage updates across multiple projects on a single server. With enhanced security measures, including TLS certificate hash validation, DDOS protection, and an automatic ban mechanism, it ensures uninterrupted and secure software distribution.

## Features:

1. **Multi-Project Management**: Host multiple projects on a single server, streamlining update deployments for various apps simultaneously.

2. **Version Management**: Directly upload your binaries to specific branches using commands such as `-updater=upload-to-main` or `-updater=upload-to-test`.

3. **Secure Updates**: Application instances automatically seek server-side updates, fetched over a protected channel with file integrity verification.

4. **TLS Certificate Hash Validation**: Ensures that the server's TLS certificate is genuine, reducing risks associated with man-in-the-middle attacks.

5. **DDOS Protection & Auto-Ban**: Shields the server from DDOS attacks and bans suspicious activity, ensuring operational resilience.

6. **Easy Server Setup**: Launching the server application is a breeze. Simply specify the root directory, TLS certificate keys file, and project descriptions in the configuration.

7. **Efficient Storage**: Updates are stored in compressed formats, optimizing storage space.

8. **Version Archiving**: TODO: *Keep track of your app's evolution with the retention of past versions. By default, the last five versions are stored, but this is configurable.*

## How to Use:

### Server Configuration:

Set up the server with ease:
To set up the server:

1. Provide the necessary details in the JSON configuration file:

      ```json
	{
		"Listen": "0.0.0.0:1985", // ip:port to listen to
		"Storage": "./storage", // specifies the root directory for update storage
		"Projects": // an array that defines descriptions of the projects hosted:
	      [
	        {
			"Name": "PROJECT_NAME", 
			"Key": "KEY", // auth key to upload and check for updates
			"Branch": "main" // "test", etc. project branch for upload or chack for updates
		},
	        ...
	      ]
	}
      ```
2. Launch the server application.

### Client Integration:

To use GoUpdater within your application:

```go
var update = goupdater.New("server ip:port", "server auth key", "project name", "project branch");
```

### Update Check:

Verify server-hosted updates:

```go
if ok,err=update.Check(); ok {
    // Handle new updates.
    update.DownloadAndReplaceMe()
}
```

`update.Check()` will return `true` if a newer version is available.
`update.DownloadAndReplaceMe()` will download new file and replace current application in-place.

### Upload new version to server:

Upload this copy of application:
```bash
	>soft --goupload=main  --branch=main
	Upload done!
```

Now all other copies of "soft" for its "projectName" and "branchName" will be updated automatically to this version.

### Start server:

Just build, edit `server/app/config.json` and start `server/app`.

---

Empower your projects with GoUpdater, combining seamless versioning, security, and efficiency in a unified platform.
