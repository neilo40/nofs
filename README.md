# NOFS - Network Object FileSystem
This exposes a filesystem across the network in S3 protocol, like NFS but using objects.  On the server, specify root directory to expose.  Every dir with becomes an individual bucket.  Within that "bucket" files are objects, subdirs can be traversed with slashes in the object name

Any S3 compatible client can be used to list buckets & objects, and perform CRUD operations on objects 