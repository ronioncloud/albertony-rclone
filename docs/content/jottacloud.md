---
title: "Jottacloud"
description: "Rclone docs for Jottacloud"
date: "2018-08-07"
---

<i class="fa fa-archive"></i> Jottacloud
-----------------------------------------

Paths are specified as `remote:path`

Paths may be as deep as required, eg `remote:directory/subdirectory`.

To configure Jottacloud you will need to enter your username and password and select a mountpoint.

Here is an example of how to make a remote called `remote`.  First run:

     rclone config

This will guide you through an interactive setup process:

```
No remotes found - make a new one
n) New remote
s) Set configuration password
q) Quit config
n/s/q> n
name> remote
Type of storage to configure.
Enter a string value. Press Enter for the default ("").
Choose a number from below, or type in your own value
[snip]
13 / JottaCloud
   \ "jottacloud"
[snip]
Storage> jottacloud
User Name
Enter a string value. Press Enter for the default ("").
user> user
Password.
y) Yes type in my own password
g) Generate random password
n) No leave this optional password blank
y/g/n> y
Enter the password:
password:
Confirm the password:
password:
The mountpoint to use.
Enter a string value. Press Enter for the default ("").
Choose a number from below, or type in your own value
 1 / Will be synced by the official client.
   \ "Sync"
 2 / Archive
   \ "Archive"
mountpoint> Archive
Remote config
--------------------
[remote]
type = jottacloud
user = user
pass = *** ENCRYPTED ***
mountpoint = Archive
--------------------
y) Yes this is OK
e) Edit this remote
d) Delete this remote
y/e/d> y
```
Once configured you can then use `rclone` like this,

List directories in top level of your Jottacloud

    rclone lsd remote:

List all the files in your Jottacloud

    rclone ls remote:

To copy a local directory to an Jottacloud directory called backup

    rclone copy /home/source remote:backup

### --fast-list ###

This remote supports `--fast-list` which allows you to use fewer
transactions in exchange for more memory. See the [rclone
docs](/docs/#fast-list) for more details.

Note that the implementation in Jottacloud always uses only a single
API request to get the entire list, so for large folders this could
lead to long wait time before the first results are shown.

### Modified time and hashes ###

Jottacloud allows modification times to be set on objects accurate to 1
second.  These will be used to detect whether objects need syncing or
not.

Jottacloud supports MD5 type hashes, so you can use the `--checksum`
flag.

Note that Jottacloud requires the MD5 hash before upload so if the
source does not have an MD5 checksum then the file will be cached
temporarily on disk (wherever the `TMPDIR` environment variable points
to) before it is uploaded.  Small files will be cached in memory - see
the `--jottacloud-md5-memory-limit` flag.

### Deleting files ###

By default rclone will send all files to the trash when deleting files.
To delete permanently use the `--jottacloud-hard-delete` option,
or set the equivalent environment variable.

The option `--jottacloud-trashed-files` can be set to list trashed files
in their original location.

Due to a lack of API documentation emptying the trash is currently
only possible via the Jottacloud website.

### Versions ###

Jottacloud supports file versioning. When rclone uploads a new version of a file it creates a new version of it. Currently rclone only supports retrieving the current version but older versions can be accessed via the Jottacloud Website.

### Quota information ###

To view your current quota you can use the `rclone about remote:`
command which will display your usage limit (unless it is unlimited)
and the current usage.

### Limitations ###

Note that Jottacloud is case insensitive so you can't have a file called
"Hello.doc" and one called "hello.doc".

There are quite a few characters that can't be in Jottacloud file names. Rclone will map these names to and from an identical looking unicode equivalent. For example if a file has a ? in it will be mapped to ？ instead.

Jottacloud only supports filenames up to 255 characters in length.

### Specific options ###

Here are the command line options specific to this cloud storage
system.

#### --jottacloud-md5-memory-limit SizeSuffix

Files bigger than this will be cached on disk to calculate the MD5 if
required. (default 10M)

#### --jottacloud-hard-delete ####

Controls whether files are sent to the trash or deleted
permanently. Defaults to false, namely sending files to the trash.
Use `--jottacloud-hard-delete=true` to delete files permanently instead.

#### --jottacloud-trashed-files ####

Advanced option to only show files that are in the trash. This will show
the trashed files in their original directory structure. Listed
directories may contain files that are not trashed, and therefore not listed.

When deleting trashed files listed with this option, they will be deleted
permanently, regardless of the `--jottacloud-use-trash` option.

When deleting any of the listed directories, they will be moved to
trash or deleted permanently according to the `--jottacloud-use-trash`
option.

Not that this option is not supported and will be ignored when used in
combination with the `--fast-list` option.

#### --jottacloud-incomplete-files ####

Advanced option to only show files that are incomplete. Listed
directories may contain files that are not incomplete, and therefore not listed.

Incomplete files may be in trash, use in combination with option `--jottacloud-trashed-files`
to list only such files.

If used in combination with `--jottacloud-corrupt-files` then both incomplete and corrupt files will be shown.

#### --jottacloud-corrupt-files ####

Advanced option to only show files that are corrupt. Listed
directories may contain files that are not incomplete, and therefore not listed.

Corrupt files may be in trash, use in combination with option `--jottacloud-trashed-files`
to list only such files.

If used in combination with `--jottacloud-incomplete-files` then both corrupt and incomplete files will be shown.

### Troubleshooting ###

Jottacloud exhibits some inconsistent behaviours regarding deleted files and folders which may cause Copy, Move and DirMove operations to previously deleted paths to fail. Emptying the trash should help in such cases.