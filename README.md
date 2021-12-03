# papyrus
A Jenkins based API & Toolchain.

## How to use
First off, papyrus is only compatible with the linux platform.

### Requirements
To install papyrus, you will need to have the following:
* [Go 1.17](https://golang.org/doc/install)
* [A working Jenkins Server](https://jenkins.io/download/)
* [Git](https://git-scm.com/downloads)

### Installation
Once you have installed the required dependencies, 
you must clone the repository, this can be done with the command:
```shell
git clone https://github.com/PurpurMC/papyrus
```

Then `cd` into the directory and run the following command:
```
./build.sh
```

The output files of this command will be in the `out` directory.

To install the CLI, `cd` into the out directory and copy
the `papyrus` binary to the `/usr/local/bin` directory.

To make sure papyrus is working, you can test the CLI with
the `papyrus` command.

Then copy the rest of the files in the `out` directory to a
location of your choice (this will be the location of the webserver).

### Setup
To set up papyrus you will want to run `sudo papyrus setup`. This
will create the `/etc/papyrus.json` file and the `/srv/papyrus` directory.

You must make sure that the user who will run the CLI has
write permissions to the `/srv/papyrus` directory.

You will want to configure aspects of papyrus, so feel free to
explore the config file. Everything is pretty self-explanatory.

To run the webserver, use the `./papyrus web` command in
the location where your webserver is. **It is important you run
this command from the local binary, not the global one.**

### Using the CLI
To add a new build to the CLI, you can use the `papyrus add` command.

The command syntax is: `papyrus add [project] [version] [build] [file-path]`.
- The `project` is the name of the project you want to add.
- The `version` is the version of the project you want to add.
- The `build` is the build number of the project you want to add.
- The `file-path` is the path to the file you want to add
(this part replaces the `{file}` in `jenkins_file_path`).

An example of the command is:
```
papyrus add purpur 1.18 1429 build/libs/purpur.jar
```
